"use client";

import { useState } from "react";
import { Send } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const quickMessages = [
  "Is this still available?",
  "Can you do a discount?",
  "Where can we meet?",
  "Can I see more photos?",
  "What's the lowest price?",
];

const QuickMessages = () => {
  const [message, setMessage] = useState("");

  return (
    <div className="flex flex-col gap-3 rounded-2xl border border-border/60 bg-card p-4 shadow-sm">
      <h3 className="text-sm font-semibold text-foreground">
        Message the seller
      </h3>

      {/* Quick replies */}
      <div className="flex flex-wrap gap-2">
        {quickMessages.map((msg) => (
          <button
            key={msg}
            onClick={() => setMessage(msg)}
            className="rounded-full border border-border/60 bg-muted/50 px-3 py-1.5 text-xs font-medium text-foreground/80 transition-colors hover:border-brand/40 hover:bg-brand-muted hover:text-brand"
          >
            {msg}
          </button>
        ))}
      </div>

      {/* Input */}
      <div className="flex gap-2">
        <Input
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Type a message..."
          className="h-10 flex-1 rounded-xl text-sm"
        />
        <Button
          size="icon"
          className="size-10 shrink-0 rounded-xl"
          aria-label="Send message"
        >
          <Send className="size-4" />
        </Button>
      </div>
    </div>
  );
};

export default QuickMessages;
