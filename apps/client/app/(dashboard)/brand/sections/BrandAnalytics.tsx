"use client";

import { Eye, ShoppingBag, TrendingUp, Star } from "lucide-react";

const stats = [
  { label: "Views this month", value: "0", icon: Eye, change: "–" },
  { label: "Items sold", value: "0", icon: ShoppingBag, change: "–" },
  { label: "Revenue (GH₵)", value: "0", icon: TrendingUp, change: "–" },
  { label: "Avg rating", value: "N/A", icon: Star, change: "–" },
];

export default function BrandAnalytics() {
  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">Analytics</h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Track your brand&apos;s performance
      </p>

      <div className="mt-6 grid gap-4 sm:grid-cols-2">
        {stats.map(({ label, value, icon: Icon, change }) => (
          <div
            key={label}
            className="flex items-center gap-4 rounded-lg border border-border/60 p-4"
          >
            <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-brand-muted">
              <Icon className="size-5 text-brand" />
            </div>
            <div className="flex-1">
              <p className="text-xs font-medium text-muted-foreground">
                {label}
              </p>
              <p className="text-xl font-bold text-foreground">{value}</p>
            </div>
            <span className="rounded-full bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground">
              {change}
            </span>
          </div>
        ))}
      </div>

      <div className="mt-8 flex flex-col items-center gap-2 rounded-lg border border-dashed border-border/60 py-12 text-center">
        <TrendingUp className="size-10 text-muted-foreground/50" />
        <p className="text-sm font-medium text-foreground">
          Analytics coming soon
        </p>
        <p className="text-xs text-muted-foreground">
          Start selling to see your performance data here
        </p>
      </div>
    </div>
  );
}
