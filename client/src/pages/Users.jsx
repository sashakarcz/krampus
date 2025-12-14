import { useEffect, useState } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Grid,
  LinearProgress,
  TextField,
  Typography,
  MenuItem,
  Alert,
  IconButton,
} from '@mui/material';
import { Refresh as RefreshIcon, Edit as EditIcon } from '@mui/icons-material';
import apiClient from '../api/client';

const Users = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openEdit, setOpenEdit] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);
  const [newRole, setNewRole] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const fetchUsers = async () => {
    try {
      const response = await apiClient.get('/api/users');
      setUsers(response.data);
    } catch (error) {
      console.error('Failed to fetch users:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleEditRole = async () => {
    try {
      await apiClient.put(`/api/users/${selectedUser.id}`, { role: newRole });
      setSuccess('User role updated successfully!');
      setOpenEdit(false);
      fetchUsers();
      setTimeout(() => setSuccess(''), 3000);
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to update user role');
      setTimeout(() => setError(''), 5000);
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">Users</Typography>
        <IconButton onClick={fetchUsers}>
          <RefreshIcon />
        </IconButton>
      </Box>

      {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      {loading ? (
        <LinearProgress />
      ) : (
        <Grid container spacing={2}>
          {users.map((user) => (
            <Grid item xs={12} md={6} key={user.id}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                    <Box sx={{ flex: 1 }}>
                      <Typography variant="h6">{user.username}</Typography>
                      <Box sx={{ mt: 1 }}>
                        <Chip
                          label={user.role}
                          size="small"
                          color={user.role === 'ADMIN' ? 'primary' : 'default'}
                          sx={{ mr: 1 }}
                        />
                      </Box>
                      {user.email && (
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                          {user.email}
                        </Typography>
                      )}
                      <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                        Created: {new Date(user.created_at).toLocaleString()}
                      </Typography>
                      {user.last_login && (
                        <Typography variant="caption" color="text.secondary">
                          Last login: {new Date(user.last_login).toLocaleString()}
                        </Typography>
                      )}
                    </Box>
                    <IconButton
                      onClick={() => {
                        setSelectedUser(user);
                        setNewRole(user.role);
                        setOpenEdit(true);
                      }}
                    >
                      <EditIcon />
                    </IconButton>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      <Dialog open={openEdit} onClose={() => setOpenEdit(false)}>
        <DialogTitle>Edit User Role</DialogTitle>
        <DialogContent>
          <Typography variant="body2" sx={{ mb: 2 }}>
            User: {selectedUser?.username}
          </Typography>
          <TextField
            fullWidth
            select
            label="Role"
            value={newRole}
            onChange={(e) => setNewRole(e.target.value)}
            margin="normal"
          >
            <MenuItem value="USER">User</MenuItem>
            <MenuItem value="ADMIN">Admin</MenuItem>
          </TextField>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenEdit(false)}>Cancel</Button>
          <Button onClick={handleEditRole} variant="contained">
            Update
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Users;
