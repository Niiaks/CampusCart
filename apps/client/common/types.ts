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

export interface Brand {
  id: string;
  seller_id: string;
  name: string;
  description: string;
  profile_url?: string;
  banner_url?: string;
  created_at: string;
  updated_at: string;
}

export interface Listing {
  id: string;
  brand_id: string;
  category_id: string;
  title: string;
  description: string;
  price: number;
  image_url: string[];
  video_url?: string[];
  condition: "new" | "used" | "second-hand";
  is_promoted: boolean;
  is_discounted: boolean;
  discount_percentage: number;
  brand?: string;
  model?: string;
  size?: string;
  storage_size?: string;
  color?: string;
  created_at: string;
  updated_at: string;
}

export interface Category {
  category_id: string;
  name: string;
  image_url: string;
}

export interface CreateBrandRequest {
  name: string;
  description: string;
  profile_url?: string;
  banner_url?: string;
}

export interface CreateListingRequest {
  brand_id: string;
  category_id: string;
  title: string;
  description: string;
  price: number;
  image_url: string[];
  video_url?: string[];
  condition: "new" | "used" | "second-hand";
  brand?: string;
  model?: string;
  size?: string;
  storage_size?: string;
  color?: string;
}
