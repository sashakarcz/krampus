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
import {
  ThumbUp,
  ThumbDown,
  Add as AddIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import apiClient from '../api/client';
import { useAuth } from '../contexts/AuthContext';

const Proposals = () => {
  const { isAdmin } = useAuth();
  const [proposals, setProposals] = useState([]);
  const [loading, setLoading] = useState(true);
  const [openDialog, setOpenDialog] = useState(false);
  const [formData, setFormData] = useState({
    identifier: '',
    rule_type: 'BINARY',
    proposed_policy: 'BLOCKLIST',
    custom_message: '',
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const fetchProposals = async () => {
    try {
      const response = await apiClient.get('/api/proposals');
      setProposals(response.data);
    } catch (error) {
      console.error('Failed to fetch proposals:', error);
      setError('Failed to fetch proposals');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProposals();
  }, []);

  const handleVote = async (proposalId, voteType) => {
    try {
      await apiClient.post(`/api/proposals/${proposalId}/vote`, { vote_type: voteType });
      setSuccess('Vote submitted successfully!');
      fetchProposals();
      setTimeout(() => setSuccess(''), 3000);
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to submit vote');
      setTimeout(() => setError(''), 5000);
    }
  };

  const handleApprove = async (proposalId, policy) => {
    try {
      await apiClient.post(`/api/proposals/${proposalId}/approve`, { policy });
      setSuccess('Proposal approved!');
      fetchProposals();
      setTimeout(() => setSuccess(''), 3000);
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to approve proposal');
      setTimeout(() => setError(''), 5000);
    }
  };

  const handleCreate = async () => {
    try {
      await apiClient.post('/api/proposals', formData);
      setSuccess('Proposal created successfully!');
      setOpenDialog(false);
      setFormData({
        identifier: '',
        rule_type: 'BINARY',
        proposed_policy: 'BLOCKLIST',
        custom_message: '',
      });
      fetchProposals();
      setTimeout(() => setSuccess(''), 3000);
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to create proposal');
      setTimeout(() => setError(''), 5000);
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'PENDING':
        return 'warning';
      case 'APPROVED':
        return 'success';
      case 'REJECTED':
        return 'error';
      default:
        return 'default';
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">Proposals</Typography>
        <Box>
          <IconButton onClick={fetchProposals} sx={{ mr: 1 }}>
            <RefreshIcon />
          </IconButton>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => setOpenDialog(true)}
          >
            New Proposal
          </Button>
        </Box>
      </Box>

      {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      {loading ? (
        <LinearProgress />
      ) : (
        <Grid container spacing={2}>
          {proposals.map((proposal) => (
            <Grid item xs={12} key={proposal.id}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', mb: 2 }}>
                    <Box sx={{ flex: 1 }}>
                      <Typography variant="h6" sx={{ wordBreak: 'break-all' }}>
                        {proposal.identifier}
                      </Typography>
                      <Box sx={{ mt: 1 }}>
                        <Chip label={proposal.rule_type} size="small" sx={{ mr: 1 }} />
                        <Chip
                          label={proposal.proposed_policy}
                          size="small"
                          color={proposal.proposed_policy === 'BLOCKLIST' ? 'error' : 'success'}
                          sx={{ mr: 1 }}
                        />
                        <Chip label={proposal.status} size="small" color={getStatusColor(proposal.status)} />
                      </Box>
                      {proposal.custom_message && (
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                          {proposal.custom_message}
                        </Typography>
                      )}
                      <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                        Created by {proposal.creator_username} on {new Date(proposal.created_at).toLocaleString()}
                      </Typography>
                    </Box>
                    <Box sx={{ ml: 2, minWidth: 200 }}>
                      <Typography variant="body2" gutterBottom>
                        Votes:
                      </Typography>
                      <Box sx={{ mb: 1 }}>
                        <Typography variant="caption">Allowlist: {proposal.allowlist_votes}</Typography>
                        <LinearProgress
                          variant="determinate"
                          value={(proposal.allowlist_votes / 3) * 100}
                          sx={{ height: 8, borderRadius: 1, backgroundColor: '#e0e0e0' }}
                          color="success"
                        />
                      </Box>
                      <Box>
                        <Typography variant="caption">Blocklist: {proposal.blocklist_votes}</Typography>
                        <LinearProgress
                          variant="determinate"
                          value={(proposal.blocklist_votes / 3) * 100}
                          sx={{ height: 8, borderRadius: 1, backgroundColor: '#e0e0e0' }}
                          color="error"
                        />
                      </Box>
                      {proposal.status === 'PENDING' && (
                        <Box sx={{ mt: 2, display: 'flex', gap: 1 }}>
                          <Button
                            size="small"
                            variant="outlined"
                            color="success"
                            startIcon={<ThumbUp />}
                            onClick={() => handleVote(proposal.id, 'ALLOWLIST')}
                          >
                            Allowlist
                          </Button>
                          <Button
                            size="small"
                            variant="outlined"
                            color="error"
                            startIcon={<ThumbDown />}
                            onClick={() => handleVote(proposal.id, 'BLOCKLIST')}
                          >
                            Blocklist
                          </Button>
                        </Box>
                      )}
                      {isAdmin() && proposal.status === 'PENDING' && (
                        <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
                          <Button
                            size="small"
                            variant="contained"
                            color="success"
                            onClick={() => handleApprove(proposal.id, 'ALLOWLIST')}
                          >
                            Approve Allow
                          </Button>
                          <Button
                            size="small"
                            variant="contained"
                            color="error"
                            onClick={() => handleApprove(proposal.id, 'BLOCKLIST')}
                          >
                            Approve Block
                          </Button>
                        </Box>
                      )}
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create New Proposal</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            label="Identifier (SHA256 hash, cert hash, etc.)"
            value={formData.identifier}
            onChange={(e) => setFormData({ ...formData, identifier: e.target.value })}
            margin="normal"
          />
          <TextField
            fullWidth
            select
            label="Rule Type"
            value={formData.rule_type}
            onChange={(e) => setFormData({ ...formData, rule_type: e.target.value })}
            margin="normal"
          >
            <MenuItem value="BINARY">Binary</MenuItem>
            <MenuItem value="CERTIFICATE">Certificate</MenuItem>
            <MenuItem value="SIGNINGID">Signing ID</MenuItem>
            <MenuItem value="TEAMID">Team ID</MenuItem>
            <MenuItem value="CDHASH">CD Hash</MenuItem>
          </TextField>
          <TextField
            fullWidth
            select
            label="Proposed Policy"
            value={formData.proposed_policy}
            onChange={(e) => setFormData({ ...formData, proposed_policy: e.target.value })}
            margin="normal"
          >
            <MenuItem value="ALLOWLIST">Allowlist</MenuItem>
            <MenuItem value="BLOCKLIST">Blocklist</MenuItem>
          </TextField>
          <TextField
            fullWidth
            label="Custom Message (optional)"
            value={formData.custom_message}
            onChange={(e) => setFormData({ ...formData, custom_message: e.target.value })}
            margin="normal"
            multiline
            rows={3}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button onClick={handleCreate} variant="contained">
            Create
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Proposals;
