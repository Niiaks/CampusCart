import { Bookmark, MapPin } from "lucide-react";

import { Button } from "@/components/ui/button";

export interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  image: string;
  seller: {
    name: string;
    avatar: string;
  };
  location: string;
  condition: string;
  timeAgo: string;
}

export const dummyListings: Listing[] = [
  {
    id: 1,
    title: "MacBook Air M2 — 256GB, Midnight",
    description:
      "Barely used MacBook Air M2, comes with charger and original box. Perfect for students.",
    price: 4500,
    image: "https://picsum.photos/seed/macbook/400/400",
    seller: { name: "Kwame A.", avatar: "https://i.pravatar.cc/40?img=1" },
    location: "Legon Hall",
    condition: "Like New",
    timeAgo: "2h ago",
  },
  {
    id: 2,
    title: "Organic Chemistry Textbook (7th Ed.)",
    description:
      "Clean copy with no highlights. Covers all chapters for CHEM 101 & 102.",
    price: 120,
    image: "https://picsum.photos/seed/chembook/400/400",
    seller: { name: "Ama D.", avatar: "https://i.pravatar.cc/40?img=5" },
    location: "Akuafo Hall",
    condition: "Good",
    timeAgo: "5h ago",
  },
  {
    id: 3,
    title: "JBL Flip 6 Bluetooth Speaker",
    description:
      "Sealed in box, never opened. Waterproof portable speaker with amazing bass.",
    price: 650,
    image: "https://picsum.photos/seed/jbl/400/400",
    seller: { name: "Kofi M.", avatar: "https://i.pravatar.cc/40?img=3" },
    location: "TF Hostel",
    condition: "Brand New",
    timeAgo: "1d ago",
  },
  {
    id: 4,
    title: 'Samsung Monitor 24" Curved',
    description:
      "Great for coding and streaming. Slight scratch on the bezel, screen is flawless.",
    price: 1800,
    image: "https://picsum.photos/seed/monitor/400/400",
    seller: { name: "Efua K.", avatar: "https://i.pravatar.cc/40?img=9" },
    location: "Pentagon",
    condition: "Good",
    timeAgo: "3h ago",
  },
  {
    id: 5,
    title: "Desk Lamp with USB Charging",
    description:
      "Adjustable LED desk lamp with built-in USB port. Three brightness levels.",
    price: 85,
    image: "https://picsum.photos/seed/lamp/400/400",
    seller: { name: "Yaw B.", avatar: "https://i.pravatar.cc/40?img=7" },
    location: "Mensah Sarbah",
    condition: "Like New",
    timeAgo: "6h ago",
  },
  {
    id: 6,
    title: "Nike Air Force 1 — Size 43",
    description:
      "Worn a few times, still in decent shape. White colourway, size EU 43.",
    price: 350,
    image: "https://picsum.photos/seed/nike/400/400",
    seller: { name: "Adwoa S.", avatar: "https://i.pravatar.cc/40?img=10" },
    location: "Volta Hall",
    condition: "Used — Fair",
    timeAgo: "12h ago",
  },
  {
    id: 7,
    title: "HP Scientific Calculator",
    description:
      "HP 35s scientific calculator, works perfectly. Ideal for engineering courses.",
    price: 95,
    image: "https://picsum.photos/seed/calc/400/400",
    seller: { name: "Nana P.", avatar: "https://i.pravatar.cc/40?img=11" },
    location: "Legon Hall",
    condition: "Good",
    timeAgo: "1d ago",
  },
  {
    id: 8,
    title: "Mini Fridge — 50L Silver",
    description:
      "Compact 50-litre fridge, runs quietly. Used for one semester only.",
    price: 900,
    image: "https://picsum.photos/seed/fridge/400/400",
    seller: { name: "Akua N.", avatar: "https://i.pravatar.cc/40?img=16" },
    location: "ISH",
    condition: "Like New",
    timeAgo: "4h ago",
  },
];

const conditionColor: Record<string, string> = {
  "Brand New": "bg-emerald-100 text-emerald-700",
  "Like New": "bg-sky-100 text-sky-700",
  Good: "bg-amber-100 text-amber-700",
  "Used — Fair": "bg-orange-100 text-orange-700",
};

const ListingCard = ({ listing }: { listing: Listing }) => {
  return (
    <article className="group flex w-full flex-col overflow-hidden rounded-2xl border border-border/60 bg-card shadow-sm transition-shadow hover:shadow-lg">
      {/* Image */}
      <div className="relative aspect-square w-full overflow-hidden bg-muted">
        <img
          src={listing.image}
          alt={listing.title}
          className="size-full object-cover transition-transform duration-300 group-hover:scale-105"
        />

        {/* Condition badge */}
        <span
          className={`absolute left-2.5 top-2.5 rounded-full px-2.5 py-0.5 text-[11px] font-semibold ${conditionColor[listing.condition] ?? "bg-muted text-muted-foreground"}`}
        >
          {listing.condition}
        </span>

        {/* Save button */}
        <Button
          variant="ghost"
          size="icon-sm"
          className="absolute right-2.5 top-2.5 rounded-full bg-white/80 text-foreground/70 backdrop-blur hover:bg-white hover:text-brand"
          aria-label="Save listing"
        >
          <Bookmark className="size-4" />
        </Button>
      </div>

      {/* Details */}
      <div className="flex flex-1 flex-col gap-2 p-3">
        {/* Price & Title */}
        <div className="space-y-0.5">
          <p className="text-base font-bold text-foreground">
            GH₵ {listing.price.toLocaleString()}
          </p>
          <h3 className="line-clamp-2 text-sm font-medium leading-snug text-foreground/80">
            {listing.title}
          </h3>
        </div>

        {/* Description */}
        <p className="line-clamp-2 text-xs leading-relaxed text-muted-foreground">
          {listing.description}
        </p>

        {/* Location */}
        <div className="flex items-center gap-1 text-xs text-muted-foreground">
          <MapPin className="size-3 shrink-0" />
          <span className="truncate">{listing.location}</span>
        </div>

        {/* Seller & Time */}
        <div className="mt-auto flex items-center gap-2 border-t border-border/40 pt-2">
          <img
            src={listing.seller.avatar}
            alt={listing.seller.name}
            className="size-6 rounded-full object-cover ring-1 ring-border"
          />
          <span className="truncate text-xs font-medium text-foreground/70">
            {listing.seller.name}
          </span>
          <span className="ml-auto shrink-0 text-[11px] text-muted-foreground">
            {listing.timeAgo}
          </span>
        </div>
      </div>
    </article>
  );
};

export default ListingCard;
