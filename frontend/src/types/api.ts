export interface ApiError {
  error: {
    code: string;
    message: string;
    details?: { field: string; message: string }[];
  };
}

export interface ListMeta {
  total: number;
  page: number;
  limit: number;
}

export interface ListResponse<T> {
  data: T[];
  meta: ListMeta;
}
