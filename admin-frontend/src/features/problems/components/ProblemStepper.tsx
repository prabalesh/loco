import { Box, Step, StepButton, Stepper } from "@mui/material"
import { Link } from "react-router-dom";

export type ProblemStepperStep = 1 | 2 | 3 | 4

interface ProblemStepperProps {
    currentStep: ProblemStepperStep;
    model: 'validate' | 'publish'
    problemId: number | string
}

export const ProblemStepper = ({ currentStep, model, problemId }: ProblemStepperProps) => {
    const problemSteps = [{
        label: 'Metadata',
        href: `/problems/edit/${problemId}`,
    }, {
        label: 'Test Cases',
        href: `/problems/${problemId}/testcases`,
    }, {
        label: 'Languages',
        href: `/problems/${problemId}/languages`
    }, {
        label: 'Validate',
        href: `/problems/${problemId}/validate`
    }, {
        label: 'Publish',
        href: `/problems/${problemId}/validate`
    }]
    let activeStep = currentStep - 1;


    if (currentStep == 4 && model === 'publish') {
        activeStep = 4
    }

    return (
        <Box sx={{ width: "100%", mb: 3 }}>
            <Stepper activeStep={activeStep} alternativeLabel>
                {problemSteps.map((step) => (
                    <Step key={step.label}>
                        <StepButton
                            component={Link}
                            to={step.href}
                        >
                            {step.label}
                        </StepButton>
                    </Step>
                ))}
            </Stepper>
        </Box>
    )
}