"use client";

import { useAuth } from "@/hooks/useAuth";
import { Mail, Phone, Calendar, Shield } from "lucide-react";

export default function PersonalDetails() {
  const { data: user } = useAuth();

  if (!user) return null;

  const details = [
    { label: "Username", value: user.username, icon: Shield },
    { label: "Email", value: user.email, icon: Mail },
    { label: "Phone", value: user.phone || "Not set", icon: Phone },
    {
      label: "Joined",
      value: new Date(user.created_at).toLocaleDateString("en-GB", {
        day: "numeric",
        month: "long",
        year: "numeric",
      }),
      icon: Calendar,
    },
  ];

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">
        Personal Details
      </h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Your account information
      </p>

      <div className="mt-6 space-y-4">
        {/* Avatar */}
        <div className="flex items-center gap-4">
          <div className="flex size-16 items-center justify-center rounded-full bg-brand text-2xl font-bold text-brand-foreground">
            {user.username.charAt(0).toUpperCase()}
          </div>
          <div>
            <p className="font-medium text-foreground">{user.username}</p>
            <span
              className={`inline-block rounded-full px-2 py-0.5 text-xs font-medium ${
                user.email_verified
                  ? "bg-success/10 text-success"
                  : "bg-warning/10 text-warning"
              }`}
            >
              {user.email_verified ? "Verified" : "Unverified"}
            </span>
          </div>
        </div>

        {/* Info grid */}
        <div className="mt-6 grid gap-4 sm:grid-cols-2">
          {details.map(({ label, value, icon: Icon }) => (
            <div
              key={label}
              className="flex items-start gap-3 rounded-lg border border-border/60 p-4"
            >
              <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-brand-muted">
                <Icon className="size-4 text-brand" />
              </div>
              <div>
                <p className="text-xs font-medium text-muted-foreground">
                  {label}
                </p>
                <p className="mt-0.5 text-sm font-medium text-foreground">
                  {value}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
