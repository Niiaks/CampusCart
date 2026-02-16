"use client";

import { Package, MoreVertical, Eye, Pencil, Trash2 } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";

// Placeholder listings — replace with API data
const dummyListings = [
  {
    id: "1",
    title: "MacBook Pro M2 2023",
    price: 5500,
    condition: "used" as const,
    image_url: "/placeholder.svg",
    created_at: "2026-01-20T10:00:00Z",
  },
  {
    id: "2",
    title: "Calculus Early Transcendentals 8th Ed",
    price: 120,
    condition: "second-hand" as const,
    image_url: "/placeholder.svg",
    created_at: "2026-02-01T14:30:00Z",
  },
  {
    id: "3",
    title: "JBL Flip 6 Bluetooth Speaker",
    price: 450,
    condition: "new" as const,
    image_url: "/placeholder.svg",
    created_at: "2026-02-10T09:15:00Z",
  },
];

const conditionBadge: Record<string, string> = {
  new: "bg-success/10 text-success",
  used: "bg-warning/10 text-warning",
  "second-hand": "bg-info/10 text-info",
};

export default function MyListings() {
  // TODO: fetch listings from API
  const listings = dummyListings;

  return (
    <div>
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold text-foreground">My Listings</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            {listings.length} listing{listings.length !== 1 && "s"}
          </p>
        </div>
      </div>

      {listings.length === 0 ? (
        <div className="mt-12 flex flex-col items-center gap-3 text-center">
          <div className="flex size-14 items-center justify-center rounded-full bg-muted">
            <Package className="size-7 text-muted-foreground" />
          </div>
          <p className="text-sm font-medium text-foreground">No listings yet</p>
          <p className="text-xs text-muted-foreground">
            Add your first listing to start selling
          </p>
        </div>
      ) : (
        <div className="mt-6 space-y-3">
          {listings.map((listing) => (
            <div
              key={listing.id}
              className="flex items-center gap-4 rounded-lg border border-border/60 p-3 transition-colors hover:bg-muted/40"
            >
              {/* Thumbnail */}
              <div className="size-16 shrink-0 overflow-hidden rounded-lg bg-muted">
                <img
                  src={listing.image_url}
                  alt={listing.title}
                  className="size-full object-cover"
                />
              </div>

              {/* Info */}
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-medium text-foreground">
                  {listing.title}
                </p>
                <div className="mt-1 flex flex-wrap items-center gap-2">
                  <span className="text-sm font-bold text-foreground">
                    GH₵ {listing.price.toLocaleString()}
                  </span>
                  <span
                    className={`rounded-full px-2 py-0.5 text-xs font-medium capitalize ${conditionBadge[listing.condition] ?? ""}`}
                  >
                    {listing.condition}
                  </span>
                </div>
                <p className="mt-1 text-xs text-muted-foreground">
                  Listed{" "}
                  {new Date(listing.created_at).toLocaleDateString("en-GB", {
                    day: "numeric",
                    month: "short",
                    year: "numeric",
                  })}
                </p>
              </div>

              {/* Actions */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    className="shrink-0"
                    aria-label="Listing actions"
                  >
                    <MoreVertical className="size-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-40">
                  <DropdownMenuItem>
                    <Eye className="size-4" />
                    View
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <Pencil className="size-4" />
                    Edit
                  </DropdownMenuItem>
                  <DropdownMenuItem className="text-destructive focus:text-destructive">
                    <Trash2 className="size-4" />
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
