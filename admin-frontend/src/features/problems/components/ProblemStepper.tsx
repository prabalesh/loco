import {Box, Step, StepLabel, Stepper} from "@mui/material"
import { PROBLEM_STEPS } from "../../../config/constant";

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
                {PROBLEM_STEPS.map(step => (<Step key={step.label}><StepLabel>{step.label}</StepLabel></Step>))}
            </Stepper>
        </Box>
    )
}