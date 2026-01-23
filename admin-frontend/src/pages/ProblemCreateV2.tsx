import React from 'react';
import { ProblemWizard } from '../features/problems/components/wizard/ProblemWizard';
import { Box } from '@mui/material';

const ProblemCreateV2: React.FC = () => {
    return (
        <Box>
            <ProblemWizard />
        </Box>
    );
};

export default ProblemCreateV2;
