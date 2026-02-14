"use client";

import { useRef } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

import { Button } from "@/components/ui/button";

import ListingCard, { dummyListings, type Listing } from "../../ui/ListingCard";

interface ListingCarouselProps {
  listings?: Listing[];
}

const ListingCarousel = ({
  listings = dummyListings,
}: ListingCarouselProps) => {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: "left" | "right") => {
    if (!scrollRef.current) return;
    const amount = scrollRef.current.clientWidth * 0.6;
    scrollRef.current.scrollBy({
      left: direction === "left" ? -amount : amount,
      behavior: "smooth",
    });
  };

  return (
    <div className="group relative mt-4">
      {/* Left arrow */}
      <Button
        variant="outline"
        size="icon"
        onClick={() => scroll("left")}
        className="absolute -left-3 top-1/2 z-10 size-9 -translate-y-1/2 rounded-full border bg-white/90 shadow-md opacity-0 transition-opacity group-hover:opacity-100"
        aria-label="Scroll left"
      >
        <ChevronLeft className="size-4" />
      </Button>

      <div
        ref={scrollRef}
        className="scrollbar-none flex gap-3 overflow-x-auto scroll-smooth sm:gap-4"
      >
        {listings.map((listing) => (
          <div
            key={listing.id}
            className="w-40 shrink-0 sm:w-48 md:w-52 lg:w-56"
          >
            <ListingCard listing={listing} />
          </div>
        ))}
      </div>

      {/* Right arrow */}
      <Button
        variant="outline"
        size="icon"
        onClick={() => scroll("right")}
        className="absolute -right-3 top-1/2 z-10 size-9 -translate-y-1/2 rounded-full border bg-white/90 shadow-md opacity-0 transition-opacity group-hover:opacity-100"
        aria-label="Scroll right"
      >
        <ChevronRight className="size-4" />
      </Button>
    </div>
  );
};

export default ListingCarousel;
