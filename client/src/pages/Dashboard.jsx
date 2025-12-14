import { useEffect, useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Grid,
  Typography,
  CircularProgress,
} from '@mui/material';
import {
  HowToVote as VoteIcon,
  Rule as RuleIcon,
  Computer as ComputerIcon,
  People as PeopleIcon,
} from '@mui/icons-material';
import apiClient from '../api/client';
import { useAuth } from '../contexts/AuthContext';

const StatCard = ({ title, value, icon, color }) => (
  <Card>
    <CardContent>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <Box
          sx={{
            backgroundColor: color,
            borderRadius: '50%',
            p: 1,
            mr: 2,
            display: 'flex',
          }}
        >
          {icon}
        </Box>
        <Typography variant="h6">{title}</Typography>
      </Box>
      <Typography variant="h3">{value !== null ? value : '...'}</Typography>
    </CardContent>
  </Card>
);

const Dashboard = () => {
  const { user, isAdmin } = useAuth();
  const [stats, setStats] = useState({
    proposals: null,
    rules: null,
    machines: null,
    users: null,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [proposalsRes, rulesRes, machinesRes] = await Promise.all([
          apiClient.get('/api/proposals'),
          apiClient.get('/api/rules'),
          apiClient.get('/api/machines'),
        ]);

        const newStats = {
          proposals: proposalsRes.data.length,
          rules: rulesRes.data.length,
          machines: machinesRes.data.length,
        };

        if (isAdmin()) {
          const usersRes = await apiClient.get('/api/users');
          newStats.users = usersRes.data.length;
        }

        setStats(newStats);
      } catch (error) {
        console.error('Failed to fetch stats:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, [isAdmin]);

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Welcome, {user?.username}!
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
        Santa Sync Server Dashboard
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Proposals"
            value={stats.proposals}
            icon={<VoteIcon sx={{ color: 'white' }} />}
            color="#1976d2"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Active Rules"
            value={stats.rules}
            icon={<RuleIcon sx={{ color: 'white' }} />}
            color="#2e7d32"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Machines"
            value={stats.machines}
            icon={<ComputerIcon sx={{ color: 'white' }} />}
            color="#ed6c02"
          />
        </Grid>
        {isAdmin() && (
          <Grid item xs={12} sm={6} md={3}>
            <StatCard
              title="Users"
              value={stats.users}
              icon={<PeopleIcon sx={{ color: 'white' }} />}
              color="#9c27b0"
            />
          </Grid>
        )}
      </Grid>
    </Box>
  );
};

export default Dashboard;
