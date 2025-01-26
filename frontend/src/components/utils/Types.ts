// Type definitions for the frontend

// Basic types
export interface ApiResponse<T> {
  code: number;
  msg: string;
  data: T;
}