import React, { useEffect, useState, useCallback } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { Box, Button, TextField, CircularProgress, Alert } from '@mui/material';

const API_BASE = process.env.REACT_APP_API_URL || '';

export default function TableView({ tableName }) {
  const [rows, setRows] = useState([]);
  const [columns, setColumns] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [filterText, setFilterText] = useState('');

  const fetchData = useCallback(() => {
    setLoading(true);
    fetch(`${API_BASE}/api/table/${tableName}`)
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch data');
        return res.json();
      })
      .then((data) => {
        let rows = [];
        let cols = [];
        if (Array.isArray(data)) {
          cols = Object.keys(data[0] || {}).map((key) => ({
            field: key,
            headerName: key,
            flex: 1,
            editable: true,
          }));
          rows = data.map((row, idx) => ({ id: idx, ...row }));
        } else if (data && data.rows && data.columns) {
          cols = data.columns.map((c) => ({ ...c, editable: true }));
          rows = data.rows.map((r, idx) => ({ id: idx, ...r }));
        }
        setRows(rows);
        setColumns(cols);
        setError(null);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [tableName]);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 60000);
    return () => clearInterval(interval);
  }, [fetchData]);

  const handleSave = () => {
    const output = rows.map(({ id, ...rest }) => rest);
    fetch(`${API_BASE}/api/table/${tableName}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(output),
    }).catch((err) => setError(err.message));
  };

  const processRowUpdate = (newRow) => {
    setRows((prev) => prev.map((r) => (r.id === newRow.id ? newRow : r)));
    return newRow;
  };

  const handleFilterChange = (e) => setFilterText(e.target.value);

  const filteredRows = rows.filter((row) =>
    Object.values(row).some((val) =>
      String(val).toLowerCase().includes(filterText.toLowerCase())
    )
  );

  if (loading) return <CircularProgress />;
  if (error) return <Alert severity="error">{error}</Alert>;

  return (
    <Box>
      <Box mb={2} sx={{ display: 'flex', gap: 2 }}>
        <TextField label="Search" value={filterText} onChange={handleFilterChange} />
        <Button variant="contained" onClick={handleSave}>Save</Button>
      </Box>
      <div style={{ height: 400, width: '100%' }}>
        <DataGrid
          rows={filteredRows}
          columns={columns}
          pagination
          pageSizeOptions={[5, 10, 20]}
          processRowUpdate={processRowUpdate}
          experimentalFeatures={{ newEditingApi: true }}
        />
      </div>
    </Box>
  );
}
