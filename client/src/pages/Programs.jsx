import { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Chip,
} from '@mui/material';
import { getPrograms } from '../api/events';

const Programs = () => {
  const [programs, setPrograms] = useState([]);

  useEffect(() => {
    loadPrograms();
  }, []);

  const loadPrograms = async () => {
    try {
      const response = await getPrograms();
      setPrograms(response.data.programs || []);
    } catch (error) {
      console.error('Failed to load programs:', error);
    }
  };

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 3 }}>
        Installed Programs & Binaries
      </Typography>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Program</TableCell>
              <TableCell>Bundle</TableCell>
              <TableCell>Certificate</TableCell>
              <TableCell>Executions</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Last Seen</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {programs.map((program) => (
              <TableRow key={program.file_hash}>
                <TableCell>
                  <Typography variant="body2">
                    {program.file_path || 'Unknown'}
                  </Typography>
                  <Typography variant="caption" color="text.secondary" sx={{ fontFamily: 'monospace' }}>
                    {program.file_hash.substring(0, 16)}...
                  </Typography>
                </TableCell>
                <TableCell>
                  {program.bundle_name && (
                    <>
                      <Typography variant="body2">{program.bundle_name}</Typography>
                      <Typography variant="caption" color="text.secondary">
                        {program.bundle_id}
                      </Typography>
                    </>
                  )}
                </TableCell>
                <TableCell>
                  {program.cert_cn && (
                    <Typography variant="body2">{program.cert_cn}</Typography>
                  )}
                  {program.team_id && (
                    <Typography variant="caption" color="text.secondary">
                      Team: {program.team_id}
                    </Typography>
                  )}
                </TableCell>
                <TableCell>{program.execution_count}</TableCell>
                <TableCell>
                  <Box sx={{ display: 'flex', gap: 1 }}>
                    {program.allow_count > 0 && (
                      <Chip
                        label={`${program.allow_count} Allowed`}
                        color="success"
                        size="small"
                        variant="outlined"
                      />
                    )}
                    {program.block_count > 0 && (
                      <Chip
                        label={`${program.block_count} Blocked`}
                        color="error"
                        size="small"
                        variant="outlined"
                      />
                    )}
                  </Box>
                </TableCell>
                <TableCell>
                  {new Date(program.last_seen).toLocaleString()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {programs.length === 0 && (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="body1" color="text.secondary">
            No program execution data yet
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default Programs;
