import axios, { AxiosError } from 'axios';
import type { ApiError } from '../types/api';

const client = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response interceptor to format errors
client.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiError>) => {
    if (error.response && error.response.data && error.response.data.error) {
      // Backend structured error
      return Promise.reject(error.response.data.error);
    }
    // Fallback for network errors
    return Promise.reject({
      code: 'NETWORK_ERROR',
      message: error.message || 'Koneksi ke server gagal',
    });
  }
);

export default client;
