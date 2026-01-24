export type Difficulty = 'easy' | 'medium' | 'hard'
export type SubmissionStatus = 'Pending' | 'Processing' | 'Accepted' | 'Wrong Answer' | 'Time Limit Exceeded' | 'Memory Limit Exceeded' | 'Runtime Error' | 'Compilation Error' | 'Internal Error'

export interface Problem {
    id: number
    title: string
    slug: string
    description: string
    difficulty: Difficulty
    time_limit: number
    memory_limit: number
    input_format: string
    output_format: string
    constraints: string
    acceptance_rate: number
    total_submissions: number
    total_accepted: number
    user_status?: 'solved' | 'attempted' | 'unsolved'
    creator?: User
    tags?: Tag[]
    categories?: Category[]
}


export interface Boilerplate {
    id: number
    problem_id: number
    language_id: number
    stub_code: string
}

export interface ProblemResponse {
    id: number;
    title: string;
    slug: string;
    description: string;
    difficulty: "easy" | "medium" | "hard" | string;
    time_limit: number;
    memory_limit: number;
    visibility: "public" | "private" | string;
    is_active: boolean;
    input_format?: string;
    output_format?: string;
    constraints?: string;

    acceptance_rate: number;
    total_submissions: number;
    total_accepted: number;

    created_by: number;
    creator: Creator;

    test_cases: TestCase[];
    boilerplates: Boilerplate[];

    created_at: string; // ISO date
    updated_at: string; // ISO date
}

export interface Creator {
    id: number;
    email: string;
    username: string;
    role: "admin" | "user" | string;

    is_active: boolean;
    email_verified: boolean;
    is_bot: boolean;

    xp: number;
    level: number;

    created_at: string;
    updated_at: string;
}

export interface TestCase {
    id: number;
    problem_id: number;

    input: string;
    expected_output: string;

    is_sample: boolean;
    validation_config: unknown | null;

    order_index: number;
    created_at: string;
}

export interface Boilerplate {
    id: number;
    problem_id: number;
    language_id: number;

    stub_code: string;
    test_harness_template: string | null;

    created_at: string;
    updated_at: string;

    language: Language;
}


export interface Tag {
    id: number
    name: string
    slug: string
}

export interface Category {
    id: number
    name: string
    slug: string
}

export interface User {
    id: number
    username: string
    email?: string
    role?: string
    created_at?: string
}

export interface Language {
    id: number
    language_id: string
    name: string
    version: string
    extension: string
    default_template: string
}

export interface ProblemLanguage {
    problem_id: number
    language_id: number
    language_name: string
    language_version: string
    function_code: string
    main_code: string
    language: Language
}

export interface TestCase {
    id: number
    problem_id: number
    input: string
    expected_output: string
    is_sample: boolean
}

export interface Submission {
    id: number
    problem_id: number
    language_id: number
    language: Language
    problem?: Problem
    status: SubmissionStatus
    code: string
    function_code: string
    error_message?: string
    runtime: number
    memory: number
    passed_test_cases: number
    total_test_cases: number
    created_at: string
    is_run_only?: boolean
    test_case_results?: TestCaseResult[]
}

export interface TestCaseResult {
    input: string
    expected_output: string
    actual_output: string
    status: string
    is_sample: boolean
}

export interface ListProblemsRequest {
    page?: number
    limit?: number
    difficulty?: string
    search?: string
    tags?: string[]
    categories?: string[]
}
