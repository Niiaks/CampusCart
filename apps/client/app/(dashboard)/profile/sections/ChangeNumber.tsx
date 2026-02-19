"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuth } from "@/hooks/useAuth";

const changeNumberSchema = z.object({
  phone: z
    .string()
    .min(1, "Phone number is required")
    .regex(/^\d{10,15}$/, "Enter a valid phone number"),
});

type ChangeNumberData = z.infer<typeof changeNumberSchema>;

export default function ChangeNumber() {
  const { data: user } = useAuth();
  const [isPending, setIsPending] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ChangeNumberData>({
    resolver: zodResolver(changeNumberSchema),
    defaultValues: {
      phone: user?.phone ?? "",
    },
  });

  const onSubmit = async (data: ChangeNumberData) => {
    setIsPending(true);
    try {
      // TODO: call API to update phone number
      toast.info("Phone number update is not yet implemented");
    } finally {
      setIsPending(false);
    }
  };

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">Change Number</h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Update the phone number linked to your account
      </p>

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="mt-6 max-w-sm space-y-4"
      >
        <div className="space-y-1.5">
          <label
            htmlFor="phone"
            className="text-sm font-medium text-foreground"
          >
            New phone number
          </label>
          <Input
            id="phone"
            type="tel"
            placeholder="0201234567"
            {...register("phone")}
          />
          {errors.phone && (
            <p className="text-xs text-destructive">{errors.phone.message}</p>
          )}
        </div>

        <Button
          type="submit"
          disabled={isPending}
          className="bg-brand text-brand-foreground hover:bg-brand-hover"
        >
          {isPending ? <Loader2 className="animate-spin" /> : "Update Number"}
        </Button>
      </form>
    </div>
  );
}
