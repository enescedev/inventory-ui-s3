import React, { useEffect, useState } from 'react';
import { AppBar, Tabs, Tab, Box, Container } from '@mui/material';
import TabPanel from './components/TabPanel';

const API_BASE_URL =
  process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

export default function App() {
  const [value, setValue] = useState(0);
  const [tabNames, setTabNames] = useState([]);

  useEffect(() => {
    fetch(`${API_BASE_URL}/tabs`)
      .then((res) => res.json())
      .then((data) => setTabNames(data))
      .catch(() => setTabNames([]));
  }, []);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4 }}>
      <AppBar position="static">
        <Tabs value={value} onChange={handleChange} textColor="inherit" indicatorColor="secondary">
          {tabNames.map((name) => (
            <Tab label={name} key={name} />
          ))}
        </Tabs>
      </AppBar>
      {tabNames.map((name, idx) => (
        <TabPanel key={name} value={value} index={idx} name={name} />
      ))}
    </Container>
  );
}
