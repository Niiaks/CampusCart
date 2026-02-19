"use client";

import { useState } from "react";
import { MessageCircleOff, MessageCircle } from "lucide-react";

import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function DisableFeedbacks() {
  const [feedbackEnabled, setFeedbackEnabled] = useState(true);

  const toggleFeedback = () => {
    // TODO: call API to toggle feedback
    setFeedbackEnabled((prev) => !prev);
    toast.success(feedbackEnabled ? "Feedbacks disabled" : "Feedbacks enabled");
  };

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">
        Feedback Settings
      </h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Control whether other users can leave feedback on your profile
      </p>

      <div className="mt-6 flex items-start gap-4 rounded-lg border border-border/60 p-4">
        <div
          className={`flex size-10 shrink-0 items-center justify-center rounded-lg ${
            feedbackEnabled ? "bg-success/10" : "bg-destructive/10"
          }`}
        >
          {feedbackEnabled ? (
            <MessageCircle className="size-5 text-success" />
          ) : (
            <MessageCircleOff className="size-5 text-destructive" />
          )}
        </div>

        <div className="flex-1">
          <p className="text-sm font-medium text-foreground">
            Feedbacks are{" "}
            <span
              className={feedbackEnabled ? "text-success" : "text-destructive"}
            >
              {feedbackEnabled ? "enabled" : "disabled"}
            </span>
          </p>
          <p className="mt-1 text-xs text-muted-foreground">
            {feedbackEnabled
              ? "Buyers and sellers can leave reviews and feedback on your profile."
              : "No one can leave feedback on your profile. Existing feedback remains visible."}
          </p>
        </div>

        <Button
          variant={feedbackEnabled ? "destructive" : "default"}
          size="sm"
          onClick={toggleFeedback}
          className={
            !feedbackEnabled
              ? "bg-brand text-brand-foreground hover:bg-brand-hover"
              : ""
          }
        >
          {feedbackEnabled ? "Disable" : "Enable"}
        </Button>
      </div>
    </div>
  );
}
