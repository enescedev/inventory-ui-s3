package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/xuri/excelize/v2"
)

// TableHandler manages table routes backed by S3
type TableHandler struct {
	Client *minio.Client
	Bucket string
}

// GetTabs lists available directories under the bucket
func (h *TableHandler) GetTabs(w http.ResponseWriter, r *http.Request) {
	tabs := map[string]struct{}{}
	for object := range h.Client.ListObjects(r.Context(), h.Bucket, minio.ListObjectsOptions{Recursive: true}) {
		if object.Err != nil {
			http.Error(w, object.Err.Error(), http.StatusInternalServerError)
			return
		}
		parts := strings.SplitN(object.Key, "/", 2)
		if len(parts) > 1 {
			tabs[parts[0]] = struct{}{}
		}
	}
	list := make([]string, 0, len(tabs))
	for k := range tabs {
		list = append(list, k)
	}
	sort.Strings(list)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// GetTable fetches the latest file for a given tab and returns JSON
func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	tab := mux.Vars(r)["tab"]
	prefix := tab + "/"
	var latest minio.ObjectInfo
	found := false
	for obj := range h.Client.ListObjects(r.Context(), h.Bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if obj.Err != nil {
			http.Error(w, obj.Err.Error(), http.StatusInternalServerError)
			return
		}
		if !strings.HasSuffix(obj.Key, ".json") && !strings.HasSuffix(obj.Key, ".xlsx") {
			continue
		}
		if !found || obj.LastModified.After(latest.LastModified) {
			latest = obj
			found = true
		}
	}
	if !found {
		http.Error(w, "tab or file not found", http.StatusNotFound)
		return
	}

	object, err := h.Client.GetObject(r.Context(), h.Bucket, latest.Key, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch filepath.Ext(latest.Key) {
	case ".json":
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	case ".xlsx":
		f, err := excelize.OpenReader(bytes.NewReader(data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sheet := f.GetSheetName(0)
		rows, err := f.GetRows(sheet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(rows) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("[]"))
			return
		}
		headers := rows[0]
		result := make([]map[string]string, 0, len(rows)-1)
		for _, row := range rows[1:] {
			entry := map[string]string{}
			for i, cell := range row {
				if i < len(headers) {
					entry[headers[i]] = cell
				}
			}
			result = append(result, entry)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	default:
		http.Error(w, "unsupported file type", http.StatusUnsupportedMediaType)
	}
}

// PutTable uploads JSON data for a given tab
func (h *TableHandler) PutTable(w http.ResponseWriter, r *http.Request) {
	tab := mux.Vars(r)["tab"]
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	key := tab + "/" + tab + ".json"
	_, err = h.Client.PutObject(r.Context(), h.Bucket, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
