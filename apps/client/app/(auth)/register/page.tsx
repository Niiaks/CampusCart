"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { registerSchema, type RegisterFormData } from "@/common/schemas";
import { useRegister } from "@/hooks/useAuth";
import type { APIError } from "@/common/types";

export default function RegisterPage() {
  const router = useRouter();
  const { mutate: signup, isPending } = useRegister();

  const {
    register,
    handleSubmit,
    setError,
    getValues,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  const onSubmit = (data: RegisterFormData) => {
    signup(data, {
      onSuccess: () => {
        toast.success(
          "Account created! Check your email for a verification code.",
        );
        router.push(
          `/verify-email?email=${encodeURIComponent(getValues("email"))}`,
        );
      },
      onError: (err: unknown) => {
        const apiError = err as APIError;

        if (apiError.errors?.length) {
          apiError.errors.forEach(({ field, error }) => {
            setError(field as keyof RegisterFormData, { message: error });
          });
          return;
        }

        toast.error(apiError.message || "Something went wrong");
      },
    });
  };

  return (
    <div className="w-full max-w-md space-y-8">
      <div className="text-center">
        <h1 className="text-2xl font-bold tracking-tight">Create an account</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Join CampusCart and start trading with fellow students
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div className="space-y-1.5">
          <label htmlFor="username" className="text-sm font-medium">
            Username
          </label>
          <Input
            id="username"
            placeholder="kwame_asante"
            autoComplete="username"
            {...register("username")}
          />
          {errors.username && (
            <p className="text-xs text-destructive">
              {errors.username.message}
            </p>
          )}
        </div>

        <div className="space-y-1.5">
          <label htmlFor="email" className="text-sm font-medium">
            Student email
          </label>
          <Input
            id="email"
            type="email"
            placeholder="you@st.ug.edu.gh"
            autoComplete="email"
            {...register("email")}
          />
          {errors.email && (
            <p className="text-xs text-destructive">{errors.email.message}</p>
          )}
        </div>

        <div className="space-y-1.5">
          <label htmlFor="phone" className="text-sm font-medium">
            Phone number
          </label>
          <Input
            id="phone"
            type="tel"
            placeholder="0241234567"
            autoComplete="tel"
            {...register("phone")}
          />
          {errors.phone && (
            <p className="text-xs text-destructive">{errors.phone.message}</p>
          )}
        </div>

        <div className="space-y-1.5">
          <label htmlFor="password" className="text-sm font-medium">
            Password
          </label>
          <Input
            id="password"
            type="password"
            placeholder="********"
            autoComplete="new-password"
            {...register("password")}
          />
          {errors.password && (
            <p className="text-xs text-destructive">
              {errors.password.message}
            </p>
          )}
        </div>

        <Button
          type="submit"
          disabled={isPending}
          className="w-full bg-brand text-brand-foreground hover:bg-brand-hover"
        >
          {isPending ? <Loader2 className="animate-spin" /> : "Create account"}
        </Button>
      </form>

      <p className="text-center text-sm text-muted-foreground">
        Already have an account?{" "}
        <Link href="/login" className="font-medium text-brand hover:underline">
          Sign in
        </Link>
      </p>
    </div>
  );
}
