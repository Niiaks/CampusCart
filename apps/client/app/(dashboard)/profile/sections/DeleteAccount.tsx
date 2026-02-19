"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { AlertTriangle, Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export default function DeleteAccount() {
  const [confirmation, setConfirmation] = useState("");
  const [isPending, setIsPending] = useState(false);
  const router = useRouter();

  const canDelete = confirmation === "DELETE";

  const handleDelete = async () => {
    if (!canDelete) return;
    setIsPending(true);
    try {
      // TODO: call API to delete account
      toast.info("Account deletion is not yet implemented");
    } finally {
      setIsPending(false);
    }
  };

  return (
    <div>
      <h2 className="text-lg font-semibold text-destructive">Delete Account</h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Permanently delete your account and all associated data
      </p>

      <div className="mt-6 rounded-lg border border-destructive/30 bg-destructive/5 p-4">
        <div className="flex items-start gap-3">
          <AlertTriangle className="mt-0.5 size-5 shrink-0 text-destructive" />
          <div className="space-y-2">
            <p className="text-sm font-medium text-foreground">
              This action is irreversible
            </p>
            <ul className="list-inside list-disc space-y-1 text-xs text-muted-foreground">
              <li>All your listings will be removed</li>
              <li>Your messages and conversations will be deleted</li>
              <li>Your profile and reviews will be permanently erased</li>
              <li>You will not be able to recover your account</li>
            </ul>
          </div>
        </div>
      </div>

      <div className="mt-6 max-w-sm space-y-4">
        <div className="space-y-1.5">
          <label
            htmlFor="confirmation"
            className="text-sm font-medium text-foreground"
          >
            Type <span className="font-bold text-destructive">DELETE</span> to
            confirm
          </label>
          <Input
            id="confirmation"
            placeholder="DELETE"
            value={confirmation}
            onChange={(e) => setConfirmation(e.target.value)}
          />
        </div>

        <Button
          variant="destructive"
          disabled={!canDelete || isPending}
          onClick={handleDelete}
        >
          {isPending ? (
            <Loader2 className="animate-spin" />
          ) : (
            "Permanently Delete Account"
          )}
        </Button>
      </div>
    </div>
  );
}
