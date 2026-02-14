"use client";

import { useRef } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

import { Button } from "@/components/ui/button";

const dummyCategories = [
  { id: 1, name: "Textbooks", color: "bg-amber-100" },
  { id: 2, name: "Electronics", color: "bg-sky-100" },
  { id: 3, name: "Dorm Essentials", color: "bg-emerald-100" },
  { id: 4, name: "Clothing", color: "bg-rose-100" },
  { id: 5, name: "Stationery", color: "bg-violet-100" },
  { id: 6, name: "Furniture", color: "bg-orange-100" },
  { id: 7, name: "Sports", color: "bg-teal-100" },
  { id: 8, name: "Kitchen", color: "bg-pink-100" },
  { id: 9, name: "Beauty", color: "bg-fuchsia-100" },
  { id: 10, name: "Accessories", color: "bg-lime-100" },
];

const Categories = () => {
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
        className="scrollbar-none flex gap-4 overflow-x-auto scroll-smooth"
      >
        {dummyCategories.map((category) => (
          <button
            key={category.id}
            className="flex shrink-0 flex-col items-center gap-2"
          >
            <div
              className={`${category.color} aspect-square w-20 rounded-2xl transition-transform hover:scale-105 sm:w-24 md:w-28`}
            />
            <span className="max-w-20 truncate text-xs font-medium text-foreground/80 sm:max-w-24 sm:text-sm">
              {category.name}
            </span>
          </button>
        ))}
      </div>

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

export default Categories;
