export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  phone: string;
}

export interface VerifyEmailRequest {
  email: string;
  code: string;
}

export interface APIError {
  code: string;
  message: string;
  status: number;
  override: boolean;
  errors?: { field: string; error: string }[];
  action?: { type: string; message: string; value: string };
}

export interface User {
  id: string;
  username: string;
  role: string;
  email: string;
  phone: string;
  email_verified: boolean;
  is_active: boolean;
  last_active: string;
  created_at: string;
}
