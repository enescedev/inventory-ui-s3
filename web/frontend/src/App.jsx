import React from 'react';
import { AppBar, Tabs, Tab, Box, Container } from '@mui/material';
import TabPanel from './components/TabPanel';

const tabNames = ['ceph', 'kubernetes', 'openshift'];

export default function App() {
  const [value, setValue] = React.useState(0);

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
