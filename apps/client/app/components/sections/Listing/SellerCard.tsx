"use client";

import { useState } from "react";
import {
  BadgeCheck,
  ChevronDown,
  ChevronUp,
  Clock,
  ListChecks,
  MessageSquare,
  Star,
} from "lucide-react";

import { Button } from "@/components/ui/button";

import type { Review, Seller } from "./types";

interface SellerCardProps {
  seller: Seller;
  reviews: Review[];
}

const SellerCard = ({ seller, reviews }: SellerCardProps) => {
  const [showReviews, setShowReviews] = useState(false);

  return (
    <div className="flex flex-col gap-4 rounded-2xl border border-border/60 bg-card p-4 shadow-sm">
      {/* Seller header */}
      <div className="flex items-center gap-3">
        <img
          src={seller.avatar}
          alt={seller.name}
          className="size-14 rounded-full object-cover ring-2 ring-border"
        />
        <div className="flex-1">
          <div className="flex items-center gap-1.5">
            <span className="font-semibold text-foreground">{seller.name}</span>
            {seller.verified && <BadgeCheck className="size-4 text-brand" />}
          </div>
          <div className="flex items-center gap-1 text-sm text-muted-foreground">
            <Star className="size-3.5 fill-amber-400 text-amber-400" />
            <span className="font-medium text-foreground">{seller.rating}</span>
            <span>({seller.totalReviews} reviews)</span>
          </div>
          <p className="text-xs text-muted-foreground">
            Member since {seller.joined}
          </p>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-2 rounded-xl bg-muted/50 p-3">
        <Stat
          icon={<MessageSquare className="size-3.5" />}
          label="Response"
          value={seller.responseRate}
        />
        <Stat
          icon={<Clock className="size-3.5" />}
          label="Reply time"
          value={seller.responseTime}
        />
        <Stat
          icon={<ListChecks className="size-3.5" />}
          label="Listings"
          value={String(seller.totalListings)}
        />
      </div>

      {/* View profile */}
      <Button variant="outline" className="w-full rounded-xl">
        View Profile
      </Button>

      {/* Reviews toggle */}
      <button
        onClick={() => setShowReviews(!showReviews)}
        className="flex items-center justify-between text-sm font-medium text-brand hover:text-brand-hover"
      >
        <span>Seller Reviews ({seller.totalReviews})</span>
        {showReviews ? (
          <ChevronUp className="size-4" />
        ) : (
          <ChevronDown className="size-4" />
        )}
      </button>

      {/* Reviews list */}
      {showReviews && (
        <div className="flex flex-col gap-3">
          {reviews.map((review) => (
            <ReviewItem key={review.id} review={review} />
          ))}
          {seller.totalReviews > reviews.length && (
            <button className="text-xs font-medium text-brand hover:text-brand-hover">
              See all {seller.totalReviews} reviews →
            </button>
          )}
        </div>
      )}
    </div>
  );
};

const Stat = ({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode;
  label: string;
  value: string;
}) => (
  <div className="flex flex-col items-center gap-0.5 text-center">
    <div className="flex items-center gap-1 text-muted-foreground">{icon}</div>
    <span className="text-xs font-bold text-foreground">{value}</span>
    <span className="text-[10px] text-muted-foreground">{label}</span>
  </div>
);

const ReviewItem = ({ review }: { review: Review }) => (
  <div className="flex gap-3 rounded-xl bg-muted/30 p-3">
    <img
      src={review.avatar}
      alt={review.author}
      className="size-8 shrink-0 rounded-full object-cover ring-1 ring-border"
    />
    <div className="flex-1 space-y-1">
      <div className="flex items-center justify-between">
        <span className="text-xs font-semibold text-foreground">
          {review.author}
        </span>
        <span className="text-[10px] text-muted-foreground">
          {review.timeAgo}
        </span>
      </div>
      <div className="flex gap-0.5">
        {Array.from({ length: 5 }).map((_, i) => (
          <Star
            key={i}
            className={`size-3 ${
              i < review.rating
                ? "fill-amber-400 text-amber-400"
                : "text-border"
            }`}
          />
        ))}
      </div>
      <p className="text-xs leading-relaxed text-foreground/70">
        {review.comment}
      </p>
    </div>
  </div>
);

export default SellerCard;
