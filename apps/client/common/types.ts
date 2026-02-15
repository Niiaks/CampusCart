interface LoginRequest {
  email: string;
  password: string;
}

interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  phone: string;
}

interface VerifyEmailRequest {
  email: string;
  code: string;
}
interface APIError {
  code: string;
  message: string;
  status: number;
  override: boolean;
  errors?: { field: string; error: string }[];
  action?: { type: string; message: string; value: string };
}

interface User {
  id: string;
  username: string;
  role: string;
  email: string;
  phone: string;
  email_verified: boolean;
  is_active: boolean;
  last_active: Date;
  created_at: Date;
}
