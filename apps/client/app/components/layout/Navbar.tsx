"use client";

import Link from "next/link";
import {
  Bell,
  Bookmark,
  LogIn,
  LogOut,
  Search,
  UserPlus,
  UserRound,
} from "lucide-react";

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
import { useAuth, useLogout } from "@/hooks/useAuth";

const Navbar = () => {
  const { data: user } = useAuth();
  const { mutate: logout } = useLogout();
  return (
    <header className="fixed inset-x-0 top-0 z-50 w-full border-b border-border/40 bg-brand/95 text-brand-foreground backdrop-blur">
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

        <div className="order-3 basis-full md:order-0 md:basis-auto md:flex-1">
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
              {user ? (
                <>
                  <DropdownMenuLabel>{user.username}</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem>My Brand</DropdownMenuItem>
                  <DropdownMenuItem>Feedback</DropdownMenuItem>
                  <DropdownMenuItem>Profile</DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => logout()}>
                    <LogOut className="size-4" />
                    Sign out
                  </DropdownMenuItem>
                </>
              ) : (
                <>
                  <DropdownMenuLabel>Get started</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem asChild>
                    <Link href="/login">
                      <LogIn className="size-4" />
                      Sign in
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem asChild>
                    <Link href="/register">
                      <UserPlus className="size-4" />
                      Sign up
                    </Link>
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </nav>
    </header>
  );
};

export default Navbar;
