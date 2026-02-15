import { MessageCircle } from "lucide-react";

import { Button } from "@/components/ui/button";

import Footer from "../components/layout/Footer";
import Navbar from "../components/layout/Navbar";

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <>
      <Navbar />
      <div className="pt-24">{children}</div>
      <Footer />
      <Button
        size="icon"
        aria-label="Open messages"
        className="fixed bottom-4 right-4 z-40 size-12 rounded-full bg-brand text-brand-foreground shadow-2xl shadow-brand/30 hover:bg-brand-hover sm:bottom-6 sm:right-6 sm:size-14"
      >
        <MessageCircle className="size-5 sm:size-6" />
      </Button>
    </>
  );
};

export default RootLayout;
