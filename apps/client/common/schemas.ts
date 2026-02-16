import { z } from "zod";

export const loginSchema = z.object({
  email: z.email("Enter a valid email address"),
  password: z
    .string()
    .min(1, "Password is required")
    .min(8, "Password must be at least 8 characters"),
});

export const registerSchema = z.object({
  username: z
    .string()
    .min(1, "Username is required")
    .min(3, "Username must be at least 3 characters"),
  email: z.email("Enter a valid email address"),
  password: z
    .string()
    .min(1, "Password is required")
    .min(8, "Password must be at least 8 characters"),
  phone: z
    .string()
    .min(1, "Phone number is required")
    .regex(/^\d{10,15}$/, "Enter a valid phone number"),
});

export const verifyEmailSchema = z.object({
  email: z.email("Enter a valid email address"),
  code: z
    .string()
    .min(1, "Verification code is required")
    .regex(/^\d{6}$/, "Code must be exactly 6 digits"),
});

export type LoginFormData = z.infer<typeof loginSchema>;
export type RegisterFormData = z.infer<typeof registerSchema>;
export type VerifyEmailFormData = z.infer<typeof verifyEmailSchema>;

export const createBrandSchema = z.object({
  name: z
    .string()
    .min(1, "Brand name is required")
    .min(2, "Brand name must be at least 2 characters")
    .max(50, "Brand name must be at most 50 characters"),
  description: z
    .string()
    .min(1, "Description is required")
    .min(10, "Description must be at least 10 characters")
    .max(500, "Description must be at most 500 characters"),
});

export type CreateBrandFormData = z.infer<typeof createBrandSchema>;

export const createListingSchema = z.object({
  title: z
    .string()
    .min(1, "Title is required")
    .min(3, "Title must be at least 3 characters")
    .max(120, "Title must be at most 120 characters"),
  description: z
    .string()
    .min(1, "Description is required")
    .min(10, "Description must be at least 10 characters")
    .max(2000, "Description must be at most 2000 characters"),
  price: z
    .number({ error: "Enter a valid price" })
    .min(0, "Price cannot be negative"),
  condition: z.enum(["new", "used", "second-hand"], {
    error: "Select a condition",
  }),
  category_id: z.string().min(1, "Category is required"),
  brand: z.string().optional(),
  model: z.string().optional(),
  size: z.string().optional(),
  storage_size: z.string().optional(),
  color: z.string().optional(),
});

export type CreateListingFormData = z.infer<typeof createListingSchema>;
