package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type TableHandler struct {
	Client *minio.Client
	Bucket string
}

func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	obj, err := h.Client.GetObject(r.Context(), h.Bucket, "table.json", minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer obj.Close()
	data, err := io.ReadAll(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *TableHandler) PutTable(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	_, err = h.Client.PutObject(r.Context(), h.Bucket, "table.json",
		bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
