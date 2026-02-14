"use client";

import { Bell, Bookmark, Search, UserRound } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";

const Navbar = () => {
  return (
    <header className="w-full border-b border-border/40 bg-brand text-brand-foreground">
      <nav className="mx-auto flex max-w-7xl flex-wrap items-center gap-4 px-4 py-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-3">
          <div className="flex size-11 items-center justify-center rounded-full bg-brand-light/60 text-lg font-semibold uppercase text-brand-foreground">
            CC
          </div>
          <div className="flex flex-col leading-tight">
            <span className="text-sm font-semibold tracking-[0.2em] text-brand-foreground/70">
              campusCart
            </span>
            <Button
              variant="ghost"
              size="xs"
              className="h-5 w-fit gap-1 rounded-full bg-brand-foreground/10 px-2 text-xs font-medium text-brand-foreground hover:bg-brand-foreground/15"
            >
              Legon
            </Button>
          </div>
        </div>

        <div className="order-3 w-full flex-1 md:order-0 md:flex">
          <label className="relative flex w-full items-center">
            <Search className="pointer-events-none absolute left-4 size-4 text-brand/80" />
            <Input
              placeholder="Search textbooks, electronics, dorm essentials..."
              className="h-11 w-full rounded-2xl border border-transparent bg-white pl-12 text-base font-semibold text-brand shadow-lg shadow-black/5 placeholder:text-brand/70 focus-visible:border-brand focus-visible:ring-brand/30"
            />
          </label>
        </div>

        <div className="ml-auto flex items-center gap-2">
          <Button className="hidden bg-white text-brand hover:bg-white/90 md:inline-flex">
            Sell on campusCart
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="text-brand-foreground hover:bg-brand-foreground/20"
            aria-label="Notifications"
          >
            <Bell className="size-5" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="text-brand-foreground hover:bg-brand-foreground/20"
            aria-label="Saved items"
          >
            <Bookmark className="size-5" />
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="rounded-full border border-brand-foreground/40 text-brand-foreground hover:bg-brand-foreground/15"
                aria-label="Account menu"
              >
                <UserRound className="size-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel>Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Profile</DropdownMenuItem>
              <DropdownMenuItem>Orders</DropdownMenuItem>
              <DropdownMenuItem>Saved listings</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Sign out</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </nav>
    </header>
  );
};

export default Navbar;
