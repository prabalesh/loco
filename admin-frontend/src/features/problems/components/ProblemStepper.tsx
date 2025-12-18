import {Box, Step, StepLabel, Stepper} from "@mui/material"

const steps = ['Metadata', 'Test Cases', 'Languages', 'Validate', 'Publish']

export type ProblemStepperStep = 1 | 2 | 3| 4

interface ProblemStepperProps {
    currentStep: ProblemStepperStep;
    model: 'validate' | 'publish'
}

export const ProblemStepper = ({currentStep, model} : ProblemStepperProps) => {
    let activeStep = currentStep - 1;

    if(currentStep == 4 && model === 'publish') {
        activeStep = 4
    }

    return (
        <Box sx={{width: '100%', mb: 3}}>
            <Stepper activeStep={activeStep} alternativeLabel>
                {steps.map(step => (<Step key={step}><StepLabel>{step}</StepLabel></Step>))}
            </Stepper>
        </Box>
    )
}