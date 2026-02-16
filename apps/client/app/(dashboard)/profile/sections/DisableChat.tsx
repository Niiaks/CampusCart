"use client";

import { useState } from "react";
import { MessageSquareOff, MessageSquare } from "lucide-react";

import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function DisableChat() {
  const [chatEnabled, setChatEnabled] = useState(true);

  const toggleChat = () => {
    // TODO: call API to toggle chat
    setChatEnabled((prev) => !prev);
    toast.success(chatEnabled ? "Chat disabled" : "Chat enabled");
  };

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">Chat Settings</h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Control whether other users can message you
      </p>

      <div className="mt-6 flex items-start gap-4 rounded-lg border border-border/60 p-4">
        <div
          className={`flex size-10 shrink-0 items-center justify-center rounded-lg ${
            chatEnabled ? "bg-success/10" : "bg-destructive/10"
          }`}
        >
          {chatEnabled ? (
            <MessageSquare className="size-5 text-success" />
          ) : (
            <MessageSquareOff className="size-5 text-destructive" />
          )}
        </div>

        <div className="flex-1">
          <p className="text-sm font-medium text-foreground">
            Direct messages are{" "}
            <span className={chatEnabled ? "text-success" : "text-destructive"}>
              {chatEnabled ? "enabled" : "disabled"}
            </span>
          </p>
          <p className="mt-1 text-xs text-muted-foreground">
            {chatEnabled
              ? "Other users can send you messages about your listings."
              : "No one can send you direct messages. You won't receive inquiries about your listings."}
          </p>
        </div>

        <Button
          variant={chatEnabled ? "destructive" : "default"}
          size="sm"
          onClick={toggleChat}
          className={
            !chatEnabled
              ? "bg-brand text-brand-foreground hover:bg-brand-hover"
              : ""
          }
        >
          {chatEnabled ? "Disable" : "Enable"}
        </Button>
      </div>
    </div>
  );
}
