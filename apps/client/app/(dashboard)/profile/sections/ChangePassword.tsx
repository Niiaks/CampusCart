"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const changePasswordSchema = z
  .object({
    currentPassword: z.string().min(1, "Current password is required"),
    newPassword: z.string().min(8, "Password must be at least 8 characters"),
    confirmPassword: z.string().min(1, "Please confirm your new password"),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

type ChangePasswordData = z.infer<typeof changePasswordSchema>;

export default function ChangePassword() {
  const [isPending, setIsPending] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ChangePasswordData>({
    resolver: zodResolver(changePasswordSchema),
  });

  const onSubmit = async (data: ChangePasswordData) => {
    setIsPending(true);
    try {
      // TODO: call API to change password
      toast.info("Password change is not yet implemented");
      reset();
    } finally {
      setIsPending(false);
    }
  };

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">Change Password</h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Update your password to keep your account secure
      </p>

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="mt-6 max-w-sm space-y-4"
      >
        <div className="space-y-1.5">
          <label
            htmlFor="currentPassword"
            className="text-sm font-medium text-foreground"
          >
            Current password
          </label>
          <Input
            id="currentPassword"
            type="password"
            placeholder="••••••••"
            autoComplete="current-password"
            {...register("currentPassword")}
          />
          {errors.currentPassword && (
            <p className="text-xs text-destructive">
              {errors.currentPassword.message}
            </p>
          )}
        </div>

        <div className="space-y-1.5">
          <label
            htmlFor="newPassword"
            className="text-sm font-medium text-foreground"
          >
            New password
          </label>
          <Input
            id="newPassword"
            type="password"
            placeholder="••••••••"
            autoComplete="new-password"
            {...register("newPassword")}
          />
          {errors.newPassword && (
            <p className="text-xs text-destructive">
              {errors.newPassword.message}
            </p>
          )}
        </div>

        <div className="space-y-1.5">
          <label
            htmlFor="confirmPassword"
            className="text-sm font-medium text-foreground"
          >
            Confirm new password
          </label>
          <Input
            id="confirmPassword"
            type="password"
            placeholder="••••••••"
            autoComplete="new-password"
            {...register("confirmPassword")}
          />
          {errors.confirmPassword && (
            <p className="text-xs text-destructive">
              {errors.confirmPassword.message}
            </p>
          )}
        </div>

        <Button
          type="submit"
          disabled={isPending}
          className="bg-brand text-brand-foreground hover:bg-brand-hover"
        >
          {isPending ? <Loader2 className="animate-spin" /> : "Change Password"}
        </Button>
      </form>
    </div>
  );
}
