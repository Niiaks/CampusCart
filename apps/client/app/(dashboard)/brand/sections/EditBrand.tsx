"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { createBrandSchema, type CreateBrandFormData } from "@/common/schemas";

export default function EditBrand() {
  const [isPending, setIsPending] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateBrandFormData>({
    resolver: zodResolver(createBrandSchema),
    defaultValues: {
      name: "",
      description: "",
    },
  });

  const onSubmit = async (data: CreateBrandFormData) => {
    setIsPending(true);
    try {
      // TODO: call API to create/update brand
      toast.info("Brand update is not yet implemented");
    } finally {
      setIsPending(false);
    }
  };

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">Edit Brand</h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Update your brand name and description
      </p>

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="mt-6 max-w-lg space-y-4"
      >
        <div className="space-y-1.5">
          <label htmlFor="name" className="text-sm font-medium text-foreground">
            Brand name
          </label>
          <Input
            id="name"
            placeholder="e.g. TechDeals Legon"
            {...register("name")}
          />
          {errors.name && (
            <p className="text-xs text-destructive">{errors.name.message}</p>
          )}
        </div>

        <div className="space-y-1.5">
          <label
            htmlFor="description"
            className="text-sm font-medium text-foreground"
          >
            Description
          </label>
          <textarea
            id="description"
            rows={4}
            placeholder="Tell buyers what your brand is about..."
            className="flex w-full rounded-xl border border-input bg-background/80 px-4 py-2 text-sm text-foreground shadow-sm transition-all placeholder:text-muted-foreground/70 focus-visible:border-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/30 disabled:cursor-not-allowed disabled:opacity-50"
            {...register("description")}
          />
          {errors.description && (
            <p className="text-xs text-destructive">
              {errors.description.message}
            </p>
          )}
        </div>

        <Button
          type="submit"
          disabled={isPending}
          className="bg-brand text-brand-foreground hover:bg-brand-hover"
        >
          {isPending ? <Loader2 className="animate-spin" /> : "Save Changes"}
        </Button>
      </form>
    </div>
  );
}
