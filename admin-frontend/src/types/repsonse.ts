export interface Pagination {
  page: number;
  limit: number;
  totalItems: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: Pagination;
}

export interface Response<T> {
  data: T;
}

export interface SimpleResponse {
  message: string;
}

export interface ApiError {
  code: string;
  message: string;
  field?: string;
  details?: Record<string, string[]>;
}

export interface ErrorResponse {
  error: {
    code: string;
    message: string;
    errors?: ApiError[];
  };
}

export interface SuccessResponse<T> {
  data: T;
  message?: string;
}

export type ApiResponse<T> = SuccessResponse<T> | ErrorResponse;

export interface CreateResponse<T> {
  data: T;
  message: string;
}

export interface DeleteResponse {
  message: string;
  deletedId: number | string;
}
