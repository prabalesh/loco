import React from 'react';
import { ProblemCreationForm } from '../components/v2/ProblemCreationForm';
import { Box, Container, IconButton, Typography, Stack } from '@mui/material';
import { ArrowBack as ArrowBackIcon } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

const ProblemCreateV2: React.FC = () => {
    const navigate = useNavigate();

    return (
        <Box sx={{ py: 4 }}>
            <Container maxWidth="lg">
                <Box sx={{ mb: 3 }}>
                    <Stack direction="row" spacing={2} alignItems="center">
                        <IconButton onClick={() => navigate('/problems')} size="small">
                            <ArrowBackIcon />
                        </IconButton>
                        <Typography variant="h5" fontWeight="bold">
                            Create New Problem (Automated V2)
                        </Typography>
                    </Stack>
                </Box>
                <ProblemCreationForm />
            </Container>
        </Box>
    );
};

export default ProblemCreateV2;
