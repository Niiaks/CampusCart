"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { verifyEmailSchema, type VerifyEmailFormData } from "@/common/schemas";
import { useVerifyEmail } from "@/hooks/useAuth";

export default function VerifyEmailPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const emailFromQuery = searchParams.get("email") ?? "";
  const { mutate: verify, isPending } = useVerifyEmail();

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<VerifyEmailFormData>({
    resolver: zodResolver(verifyEmailSchema),
    defaultValues: {
      email: emailFromQuery,
    },
  });

  const onSubmit = (data: VerifyEmailFormData) => {
    verify(data, {
      onSuccess: () => {
        toast.success("Email verified! Welcome to CampusCart.");
        router.push("/");
      },
      onError: (err: unknown) => {
        const apiError = err as APIError;

        if (apiError.errors?.length) {
          apiError.errors.forEach(({ field, error }) => {
            setError(field as keyof VerifyEmailFormData, { message: error });
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
        <h1 className="text-2xl font-bold tracking-tight">Verify your email</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Enter the 6-digit code sent to your student email
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div className="space-y-1.5">
          <label htmlFor="email" className="text-sm font-medium">
            Email
          </label>
          <Input
            id="email"
            type="email"
            placeholder="you@st.ug.edu.gh"
            autoComplete="email"
            disabled={!!emailFromQuery}
            {...register("email")}
          />
          {errors.email && (
            <p className="text-xs text-destructive">{errors.email.message}</p>
          )}
        </div>

        <div className="space-y-1.5">
          <label htmlFor="code" className="text-sm font-medium">
            Verification code
          </label>
          <Input
            id="code"
            placeholder="000000"
            maxLength={6}
            autoComplete="one-time-code"
            className="text-center text-lg tracking-[0.5em]"
            {...register("code")}
          />
          {errors.code && (
            <p className="text-xs text-destructive">{errors.code.message}</p>
          )}
        </div>

        <Button
          type="submit"
          disabled={isPending}
          className="w-full bg-brand text-brand-foreground hover:bg-brand-hover"
        >
          {isPending ? <Loader2 className="animate-spin" /> : "Verify email"}
        </Button>
      </form>
    </div>
  );
}
