import { MessageCircle } from "lucide-react";

import { Button } from "@/components/ui/button";

import Navbar from "../components/layout/Navbar";
import Link from "next/link";

const DashboardLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <>
      <Navbar />
      <div className="pt-24">{children}</div>
      <Link href={"/chat"}>
        <Button
          size="icon"
          aria-label="Open messages"
          className="fixed bottom-4 right-4 z-40 size-12 rounded-full bg-brand text-brand-foreground shadow-2xl shadow-brand/30 hover:bg-brand-hover sm:bottom-6 sm:right-6 sm:size-14"
        >
          <MessageCircle className="size-5 sm:size-6" />
        </Button>
      </Link>
    </>
  );
};

export default DashboardLayout;
