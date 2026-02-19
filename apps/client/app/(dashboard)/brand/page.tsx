"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  Store,
  FileText,
  ImageIcon,
  Package,
  Plus,
  BarChart3,
  ChevronLeft,
} from "lucide-react";

import { useAuth } from "@/hooks/useAuth";

import BrandOverview from "./sections/BrandOverview";
import EditBrand from "./sections/EditBrand";
import BrandAppearance from "./sections/BrandAppearance";
import MyListings from "./sections/MyListings";
import AddListing from "./sections/AddListing";
import BrandAnalytics from "./sections/BrandAnalytics";

const sidebarItems = [
  { key: "overview", label: "Overview", icon: Store },
  { key: "edit", label: "Edit Brand", icon: FileText },
  { key: "appearance", label: "Appearance", icon: ImageIcon },
  { key: "listings", label: "My Listings", icon: Package },
  { key: "add-listing", label: "Add Listing", icon: Plus, accent: true },
  { key: "analytics", label: "Analytics", icon: BarChart3 },
] as const;

type SectionKey = (typeof sidebarItems)[number]["key"];

const sectionComponents: Record<SectionKey, React.ComponentType> = {
  overview: BrandOverview,
  edit: EditBrand,
  appearance: BrandAppearance,
  listings: MyListings,
  "add-listing": AddListing,
  analytics: BrandAnalytics,
};

export default function BrandPage() {
  const [activeSection, setActiveSection] = useState<SectionKey>("overview");
  const { data: user, isLoading } = useAuth();
  const router = useRouter();

  if (isLoading) {
    return (
      <main className="mx-auto flex min-h-[60vh] max-w-7xl items-center justify-center px-4 py-8">
        <div className="size-8 animate-spin rounded-full border-4 border-brand border-t-transparent" />
      </main>
    );
  }

  if (!user) {
    router.push("/login");
    return null;
  }

  const ActiveComponent = sectionComponents[activeSection];

  return (
    <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-foreground">My Brand</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage your storefront and listings
        </p>
      </div>

      <div className="flex flex-col gap-6 lg:flex-row">
        {/* Sidebar */}
        <aside className="w-full shrink-0 lg:w-64">
          <nav className="flex flex-col gap-1 rounded-xl border border-border/60 bg-card p-2">
            {sidebarItems.map(({ key, label, icon: Icon, ...rest }) => {
              const isAccent = "accent" in rest && rest.accent;
              const isActive = activeSection === key;

              return (
                <button
                  key={key}
                  onClick={() => setActiveSection(key)}
                  className={`flex items-center gap-3 rounded-lg px-3 py-2.5 text-left text-sm font-medium transition-colors ${
                    isActive
                      ? "bg-brand text-brand-foreground"
                      : isAccent
                        ? "text-brand hover:bg-brand-muted"
                        : "text-foreground hover:bg-muted"
                  }`}
                >
                  <Icon className="size-4 shrink-0" />
                  {label}
                </button>
              );
            })}
          </nav>
        </aside>

        {/* Content */}
        <section className="min-w-0 flex-1">
          <button
            onClick={() => setActiveSection("overview")}
            className="mb-4 flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground lg:hidden"
          >
            <ChevronLeft className="size-4" />
            Back to menu
          </button>

          <div className="rounded-xl border border-border/60 bg-card p-6">
            <ActiveComponent />
          </div>
        </section>
      </div>
    </main>
  );
}
