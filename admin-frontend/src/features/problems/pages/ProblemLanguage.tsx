import { useParams } from "react-router-dom";

export default function ProblemLanguage() {
    const { problemId } = useParams<{ problemId: string }>();
    return <>ProblemLanguage Page</>
}