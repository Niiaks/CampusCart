"use client";

import { Store, Package, Star, TrendingUp } from "lucide-react";

export default function BrandOverview() {
  // TODO: fetch brand data from API
  const brand = {
    name: "My Store",
    description: "Welcome to my campus store!",
    profile_url: "",
    banner_url: "",
    created_at: new Date().toISOString(),
    totalListings: 0,
    totalSales: 0,
    avgRating: 0,
  };

  const stats = [
    {
      label: "Total Listings",
      value: brand.totalListings,
      icon: Package,
      color: "bg-brand-muted text-brand",
    },
    {
      label: "Total Sales",
      value: brand.totalSales,
      icon: TrendingUp,
      color: "bg-success/10 text-success",
    },
    {
      label: "Avg Rating",
      value: brand.avgRating > 0 ? brand.avgRating.toFixed(1) : "N/A",
      icon: Star,
      color: "bg-warning/10 text-warning",
    },
  ];

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">Brand Overview</h2>
      <p className="mt-1 text-sm text-muted-foreground">
        A snapshot of your storefront
      </p>

      {/* Brand card */}
      <div className="mt-6 flex items-center gap-4">
        <div className="flex size-16 items-center justify-center rounded-xl bg-brand text-2xl font-bold text-brand-foreground">
          {brand.name.charAt(0).toUpperCase()}
        </div>
        <div>
          <p className="text-lg font-semibold text-foreground">{brand.name}</p>
          <p className="text-sm text-muted-foreground">{brand.description}</p>
          <p className="mt-1 text-xs text-muted-foreground">
            Since{" "}
            {new Date(brand.created_at).toLocaleDateString("en-GB", {
              month: "long",
              year: "numeric",
            })}
          </p>
        </div>
      </div>

      {/* Stats grid */}
      <div className="mt-8 grid gap-4 sm:grid-cols-3">
        {stats.map(({ label, value, icon: Icon, color }) => (
          <div
            key={label}
            className="flex items-center gap-3 rounded-lg border border-border/60 p-4"
          >
            <div
              className={`flex size-10 shrink-0 items-center justify-center rounded-lg ${color}`}
            >
              <Icon className="size-5" />
            </div>
            <div>
              <p className="text-xs font-medium text-muted-foreground">
                {label}
              </p>
              <p className="text-xl font-bold text-foreground">{value}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
