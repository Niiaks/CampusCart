import { Bookmark, MapPin, Share2, Tag } from "lucide-react";

import { Button } from "@/components/ui/button";

import type { DetailedListing } from "./types";

interface ProductDetailsProps {
  listing: DetailedListing;
}

const conditionColor: Record<string, string> = {
  "Brand New": "bg-emerald-100 text-emerald-700",
  "Like New": "bg-sky-100 text-sky-700",
  Good: "bg-amber-100 text-amber-700",
  "Used — Fair": "bg-orange-100 text-orange-700",
};

const ProductDetails = ({ listing }: ProductDetailsProps) => {
  return (
    <div className="flex flex-col gap-5">
      {/* Header: price, title, badges */}
      <div className="space-y-2">
        <div className="flex items-start justify-between gap-3">
          <h1 className="text-2xl font-bold tracking-tight text-foreground sm:text-3xl">
            GH₵ {listing.price.toLocaleString()}
          </h1>
          <div className="flex shrink-0 gap-1.5">
            <Button
              variant="outline"
              size="icon-sm"
              className="rounded-full"
              aria-label="Share"
            >
              <Share2 className="size-4" />
            </Button>
            <Button
              variant="outline"
              size="icon-sm"
              className="rounded-full"
              aria-label="Save"
            >
              <Bookmark className="size-4" />
            </Button>
          </div>
        </div>
        <h2 className="text-lg font-semibold text-foreground/90">
          {listing.title}
        </h2>
        <div className="flex flex-wrap items-center gap-2 text-sm">
          <span
            className={`rounded-full px-2.5 py-0.5 text-xs font-semibold ${conditionColor[listing.condition] ?? "bg-muted text-muted-foreground"}`}
          >
            {listing.condition}
          </span>
          <span className="flex items-center gap-1 text-muted-foreground">
            <Tag className="size-3" />
            {listing.category}
          </span>
          <span className="flex items-center gap-1 text-muted-foreground">
            <MapPin className="size-3" />
            {listing.location}
          </span>
          <span className="text-xs text-muted-foreground">
            · Posted {listing.timeAgo}
          </span>
        </div>
      </div>

      {/* Description */}
      <div className="space-y-2">
        <h3 className="text-sm font-semibold uppercase tracking-wider text-foreground/70">
          Description
        </h3>
        <p className="text-sm leading-relaxed text-foreground/80">
          {listing.description}
        </p>
      </div>

      {/* Specs table */}
      <div className="space-y-2">
        <h3 className="text-sm font-semibold uppercase tracking-wider text-foreground/70">
          Product Details
        </h3>
        <div className="overflow-hidden rounded-xl border border-border/60">
          {Object.entries(listing.specs).map(([key, value], idx) => (
            <div
              key={key}
              className={`flex items-center justify-between px-4 py-2.5 text-sm ${
                idx % 2 === 0 ? "bg-muted/40" : "bg-card"
              }`}
            >
              <span className="font-medium text-muted-foreground">{key}</span>
              <span className="font-semibold text-foreground">{value}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ProductDetails;
