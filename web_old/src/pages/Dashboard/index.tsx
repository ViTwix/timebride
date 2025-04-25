import React from 'react';
import { Grid, Paper, Typography, Box } from '@mui/material';
import { styled } from '@mui/material/styles';

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#1A2027' : '#fff',
  ...theme.typography.body2,
  padding: theme.spacing(2),
  textAlign: 'center',
  color: theme.palette.text.secondary,
}));

const Dashboard: React.FC = () => {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Typography variant="h4" component="h1" gutterBottom>
            Dashboard
          </Typography>
        </Grid>
        <Grid item xs={12} md={6} lg={3}>
          <Item>
            <Typography variant="h6">Total Bookings</Typography>
            <Typography variant="h3">0</Typography>
          </Item>
        </Grid>
        <Grid item xs={12} md={6} lg={3}>
          <Item>
            <Typography variant="h6">Upcoming Events</Typography>
            <Typography variant="h3">0</Typography>
          </Item>
        </Grid>
        <Grid item xs={12} md={6} lg={3}>
          <Item>
            <Typography variant="h6">Templates</Typography>
            <Typography variant="h3">0</Typography>
          </Item>
        </Grid>
        <Grid item xs={12} md={6} lg={3}>
          <Item>
            <Typography variant="h6">Files</Typography>
            <Typography variant="h3">0</Typography>
          </Item>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard; 