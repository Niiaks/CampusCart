"use client";

import { useState } from "react";
import { ChevronLeft, ChevronRight, Play } from "lucide-react";

import { Button } from "@/components/ui/button";

import type { MediaItem } from "./types";

interface MediaGalleryProps {
  media: MediaItem[];
}

const MediaGallery = ({ media }: MediaGalleryProps) => {
  const [activeIndex, setActiveIndex] = useState(0);
  const activeItem = media[activeIndex];

  const prev = () =>
    setActiveIndex((i) => (i === 0 ? media.length - 1 : i - 1));
  const next = () =>
    setActiveIndex((i) => (i === media.length - 1 ? 0 : i + 1));

  return (
    <div className="flex flex-col gap-3">
      {/* Main display */}
      <div className="group relative aspect-4/3 w-full overflow-hidden rounded-2xl bg-muted">
        {activeItem.type === "video" ? (
          <video
            src={activeItem.src}
            controls
            className="size-full object-cover"
            poster={media.find((m) => m.type === "image")?.src}
          />
        ) : (
          <img
            src={activeItem.src}
            alt={activeItem.alt}
            className="size-full object-cover transition-transform duration-300"
          />
        )}

        {/* Navigation arrows */}
        {media.length > 1 && (
          <>
            <Button
              variant="ghost"
              size="icon"
              onClick={prev}
              className="absolute left-3 top-1/2 z-10 size-10 -translate-y-1/2 rounded-full bg-black/40 text-white opacity-0 backdrop-blur transition-opacity hover:bg-black/60 group-hover:opacity-100"
              aria-label="Previous"
            >
              <ChevronLeft className="size-5" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={next}
              className="absolute right-3 top-1/2 z-10 size-10 -translate-y-1/2 rounded-full bg-black/40 text-white opacity-0 backdrop-blur transition-opacity hover:bg-black/60 group-hover:opacity-100"
              aria-label="Next"
            >
              <ChevronRight className="size-5" />
            </Button>
          </>
        )}

        {/* Counter */}
        <span className="absolute bottom-3 right-3 rounded-full bg-black/50 px-3 py-1 text-xs font-medium text-white backdrop-blur">
          {activeIndex + 1} / {media.length}
        </span>
      </div>

      {/* Thumbnails */}
      <div className="scrollbar-none flex gap-2 overflow-x-auto">
        {media.map((item, idx) => (
          <button
            key={item.id}
            onClick={() => setActiveIndex(idx)}
            className={`relative size-16 shrink-0 overflow-hidden rounded-xl border-2 transition-all sm:size-20 ${
              idx === activeIndex
                ? "border-brand ring-2 ring-brand/30"
                : "border-transparent opacity-60 hover:opacity-100"
            }`}
          >
            {item.type === "video" ? (
              <div className="flex size-full items-center justify-center bg-muted">
                <Play className="size-5 text-muted-foreground" />
              </div>
            ) : (
              <img
                src={item.src}
                alt={item.alt}
                className="size-full object-cover"
              />
            )}
          </button>
        ))}
      </div>
    </div>
  );
};

export default MediaGallery;
