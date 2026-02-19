"use client";

import { useState } from "react";
import { toast } from "sonner";

const notificationOptions = [
  {
    key: "order_updates",
    label: "Order updates",
    description: "Get notified about order status changes",
  },
  {
    key: "messages",
    label: "New messages",
    description: "Receive alerts when someone sends you a message",
  },
  {
    key: "promotions",
    label: "Promotions & offers",
    description: "Deals, discounts, and featured listings",
  },
  {
    key: "price_drops",
    label: "Price drops",
    description: "Get notified when items in your saved list drop in price",
  },
  {
    key: "security",
    label: "Security alerts",
    description: "Login attempts and password changes",
  },
] as const;

type NotificationKey = (typeof notificationOptions)[number]["key"];

export default function ManageNotifications() {
  const [enabled, setEnabled] = useState<Record<NotificationKey, boolean>>({
    order_updates: true,
    messages: true,
    promotions: false,
    price_drops: true,
    security: true,
  });

  const toggle = (key: NotificationKey) => {
    setEnabled((prev) => {
      const next = { ...prev, [key]: !prev[key] };
      // TODO: call API to update notification preferences
      toast.success(
        next[key] ? "Notification enabled" : "Notification disabled",
      );
      return next;
    });
  };

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">
        Manage Notifications
      </h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Choose which notifications you want to receive
      </p>

      <div className="mt-6 space-y-3">
        {notificationOptions.map(({ key, label, description }) => (
          <div
            key={key}
            className="flex items-center justify-between gap-4 rounded-lg border border-border/60 p-4"
          >
            <div>
              <p className="text-sm font-medium text-foreground">{label}</p>
              <p className="mt-0.5 text-xs text-muted-foreground">
                {description}
              </p>
            </div>

            {/* Toggle switch */}
            <button
              role="switch"
              aria-checked={enabled[key]}
              onClick={() => toggle(key)}
              className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full transition-colors ${
                enabled[key] ? "bg-brand" : "bg-muted"
              }`}
            >
              <span
                className={`inline-block size-4 rounded-full bg-white shadow-sm transition-transform ${
                  enabled[key] ? "translate-x-6" : "translate-x-1"
                }`}
              />
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
