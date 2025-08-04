import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import { secretsApi, Secret, CreateSecretRequest } from '../services/api';

export function Secrets() {
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [loading, setLoading] = useState(true);
  const [createDialog, setCreateDialog] = useState(false);
  const [editDialog, setEditDialog] = useState<{ open: boolean; secret?: Secret }>({ open: false });
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; secret?: Secret }>({ open: false });
  const [formData, setFormData] = useState<CreateSecretRequest>({ name: '', value: '' });
  const [editValue, setEditValue] = useState('');

  useEffect(() => {
    fetchSecrets();
  }, []);

  const fetchSecrets = async () => {
    try {
      setLoading(true);
      const response = await secretsApi.list();
      setSecrets(response.secrets);
    } catch (error) {
      console.error('Error fetching secrets:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateSecret = async () => {
    try {
      await secretsApi.create(formData);
      setCreateDialog(false);
      setFormData({ name: '', value: '' });
      fetchSecrets();
    } catch (error) {
      console.error('Error creating secret:', error);
    }
  };

  const handleUpdateSecret = async () => {
    if (!editDialog.secret) return;
    
    try {
      await secretsApi.update(editDialog.secret.id, editValue);
      setEditDialog({ open: false });
      setEditValue('');
      fetchSecrets();
    } catch (error) {
      console.error('Error updating secret:', error);
    }
  };

  const handleDeleteSecret = async () => {
    if (!deleteDialog.secret) return;
    
    try {
      await secretsApi.delete(deleteDialog.secret.id);
      setDeleteDialog({ open: false });
      fetchSecrets();
    } catch (error) {
      console.error('Error deleting secret:', error);
    }
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Secrets</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setCreateDialog(true)}
        >
          Create Secret
        </Button>
      </Box>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>ID</TableCell>
                <TableCell>Created</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {secrets.map((secret) => (
                <TableRow key={secret.id}>
                  <TableCell>
                    <Typography variant="subtitle2">{secret.name}</Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                      {secret.id}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    {new Date(secret.created_at).toLocaleDateString()}
                  </TableCell>
                  <TableCell>
                    <IconButton
                      onClick={() => {
                        setEditDialog({ open: true, secret });
                        setEditValue('');
                      }}
                      size="small"
                      title="Edit Secret"
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => setDeleteDialog({ open: true, secret })}
                      size="small"
                      title="Delete Secret"
                      color="error"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Create Secret Dialog */}
      <Dialog
        open={createDialog}
        onClose={() => setCreateDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Create New Secret</DialogTitle>
        <DialogContent>
          <TextField
            label="Name"
            fullWidth
            margin="normal"
            value={formData.name}
            onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
          />
          <TextField
            label="Value"
            fullWidth
            margin="normal"
            type="password"
            value={formData.value}
            onChange={(e) => setFormData(prev => ({ ...prev, value: e.target.value }))}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialog(false)}>Cancel</Button>
          <Button
            onClick={handleCreateSecret}
            variant="contained"
            disabled={!formData.name.trim() || !formData.value.trim()}
          >
            Create
          </Button>
        </DialogActions>
      </Dialog>

      {/* Edit Secret Dialog */}
      <Dialog
        open={editDialog.open}
        onClose={() => setEditDialog({ open: false })}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Update Secret: {editDialog.secret?.name}</DialogTitle>
        <DialogContent>
          <TextField
            label="New Value"
            fullWidth
            margin="normal"
            type="password"
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            placeholder="Enter new secret value"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialog({ open: false })}>Cancel</Button>
          <Button
            onClick={handleUpdateSecret}
            variant="contained"
            disabled={!editValue.trim()}
          >
            Update
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialog.open}
        onClose={() => setDeleteDialog({ open: false })}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Delete Secret</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete the secret "{deleteDialog.secret?.name}"?
            This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialog({ open: false })}>Cancel</Button>
          <Button
            onClick={handleDeleteSecret}
            variant="contained"
            color="error"
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
```
