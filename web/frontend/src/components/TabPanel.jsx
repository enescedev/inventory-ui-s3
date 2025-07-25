import React from 'react';
import Box from '@mui/material/Box';
import TableView from './TableView';

export default function TabPanel({ value, index, name }) {
  return (
    <div role="tabpanel" hidden={value !== index}">
      {value === index && (
        <Box sx={{ p: 3 }}>
          <TableView tableName={name} />
        </Box>
      )}
    </div>
  );
}
