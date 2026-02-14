export interface DetailedListing {
  id: number;
  title: string;
  description: string;
  price: number;
  media: MediaItem[];
  condition: string;
  category: string;
  location: string;
  timeAgo: string;
  specs: Record<string, string>;
  seller: Seller;
}

export interface MediaItem {
  id: number;
  type: "image" | "video";
  src: string;
  alt: string;
}

export interface Seller {
  id: number;
  name: string;
  avatar: string;
  rating: number;
  totalReviews: number;
  responseRate: string;
  responseTime: string;
  joined: string;
  totalListings: number;
  verified: boolean;
}

export interface Review {
  id: number;
  author: string;
  avatar: string;
  rating: number;
  comment: string;
  timeAgo: string;
}

export const dummyListing: DetailedListing = {
  id: 1,
  title: "MacBook Air M2 — 256GB, Midnight",
  description:
    "Barely used MacBook Air M2 in Midnight colour. Purchased in September 2025 from the Apple Store. Comes with the original charger, box, and documentation. Battery cycle count is under 50. No scratches, dents, or marks anywhere — screen is pristine. Runs perfectly for coding, design work, and everyday use. Selling because I upgraded to a Pro. Price is slightly negotiable for serious buyers.",
  price: 4500,
  media: [
    {
      id: 1,
      type: "image",
      src: "https://picsum.photos/seed/mac1/800/600",
      alt: "MacBook Air front view",
    },
    {
      id: 2,
      type: "image",
      src: "https://picsum.photos/seed/mac2/800/600",
      alt: "MacBook Air side angle",
    },
    {
      id: 3,
      type: "video",
      src: "https://www.w3schools.com/html/mov_bbb.mp4",
      alt: "MacBook Air quick demo",
    },
    {
      id: 4,
      type: "image",
      src: "https://picsum.photos/seed/mac4/800/600",
      alt: "MacBook Air keyboard",
    },
    {
      id: 5,
      type: "image",
      src: "https://picsum.photos/seed/mac5/800/600",
      alt: "MacBook Air with charger",
    },
  ],
  condition: "Like New",
  category: "Electronics",
  location: "Legon Hall",
  timeAgo: "2h ago",
  specs: {
    Brand: "Apple",
    Model: "MacBook Air M2 2024",
    Color: "Midnight",
    Storage: "256GB SSD",
    RAM: "8GB",
    "Screen Size": '13.6"',
    "Battery Health": "98%",
  },
  seller: {
    id: 1,
    name: "Kwame Asante",
    avatar: "https://i.pravatar.cc/80?img=1",
    rating: 4.8,
    totalReviews: 24,
    responseRate: "98%",
    responseTime: "~15 min",
    joined: "Sep 2024",
    totalListings: 12,
    verified: true,
  },
};

export const dummyReviews: Review[] = [
  {
    id: 1,
    author: "Ama D.",
    avatar: "https://i.pravatar.cc/40?img=5",
    rating: 5,
    comment:
      "Great seller! Item was exactly as described. Very responsive and friendly.",
    timeAgo: "1w ago",
  },
  {
    id: 2,
    author: "Kofi M.",
    avatar: "https://i.pravatar.cc/40?img=3",
    rating: 4,
    comment:
      "Smooth transaction. Delivery was quick. Would definitely buy from again.",
    timeAgo: "2w ago",
  },
  {
    id: 3,
    author: "Efua K.",
    avatar: "https://i.pravatar.cc/40?img=9",
    rating: 5,
    comment:
      "Legit seller, no issues at all. Product was in perfect condition.",
    timeAgo: "1mo ago",
  },
  {
    id: 4,
    author: "Yaw B.",
    avatar: "https://i.pravatar.cc/40?img=7",
    rating: 4,
    comment:
      "Good communication, item was clean. Slightly late for meetup but overall good.",
    timeAgo: "1mo ago",
  },
];
