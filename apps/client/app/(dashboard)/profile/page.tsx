"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  User,
  Phone,
  MessageSquareOff,
  MessageCircleOff,
  Bell,
  Lock,
  Trash2,
  LogOut,
  ChevronLeft,
} from "lucide-react";
import { toast } from "sonner";

import { useAuth, useLogout } from "@/hooks/useAuth";

import PersonalDetails from "./sections/PersonalDetails";
import ChangeNumber from "./sections/ChangeNumber";
import DisableChat from "./sections/DisableChat";
import DisableFeedbacks from "./sections/DisableFeedbacks";
import ManageNotifications from "./sections/ManageNotifications";
import ChangePassword from "./sections/ChangePassword";
import DeleteAccount from "./sections/DeleteAccount";

const sidebarItems = [
  { key: "personal", label: "Personal Details", icon: User },
  { key: "phone", label: "Change Number", icon: Phone },
  { key: "chat", label: "Disable Chat", icon: MessageSquareOff },
  { key: "feedback", label: "Disable Feedbacks", icon: MessageCircleOff },
  { key: "notifications", label: "Manage Notifications", icon: Bell },
  { key: "password", label: "Change Password", icon: Lock },
  { key: "delete", label: "Delete Account", icon: Trash2, destructive: true },
] as const;

type SectionKey = (typeof sidebarItems)[number]["key"];

const sectionComponents: Record<SectionKey, React.ComponentType> = {
  personal: PersonalDetails,
  phone: ChangeNumber,
  chat: DisableChat,
  feedback: DisableFeedbacks,
  notifications: ManageNotifications,
  password: ChangePassword,
  delete: DeleteAccount,
};

export default function ProfilePage() {
  const [activeSection, setActiveSection] = useState<SectionKey>("personal");
  const { data: user, isLoading } = useAuth();
  const { mutate: logout } = useLogout();
  const router = useRouter();

  const handleLogout = () => {
    logout(undefined, {
      onSuccess: () => {
        toast.success("Signed out successfully");
        router.push("/");
      },
    });
  };

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
        <h1 className="text-2xl font-bold text-foreground">Settings</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage your account preferences
        </p>
      </div>

      <div className="flex flex-col gap-6 lg:flex-row">
        {/* Sidebar */}
        <aside className="w-full shrink-0 lg:w-64">
          <nav className="flex flex-col gap-1 rounded-xl border border-border/60 bg-card p-2">
            {sidebarItems.map(({ key, label, icon: Icon, ...rest }) => {
              const isDestructive = "destructive" in rest && rest.destructive;
              const isActive = activeSection === key;

              return (
                <button
                  key={key}
                  onClick={() => setActiveSection(key)}
                  className={`flex items-center gap-3 rounded-lg px-3 py-2.5 text-left text-sm font-medium transition-colors ${
                    isActive
                      ? "bg-brand text-brand-foreground"
                      : isDestructive
                        ? "text-destructive hover:bg-destructive/10"
                        : "text-foreground hover:bg-muted"
                  }`}
                >
                  <Icon className="size-4 shrink-0" />
                  {label}
                </button>
              );
            })}

            <div className="my-1 border-t border-border/60" />

            <button
              onClick={handleLogout}
              className="flex items-center gap-3 rounded-lg px-3 py-2.5 text-left text-sm font-medium text-foreground transition-colors hover:bg-muted"
            >
              <LogOut className="size-4 shrink-0" />
              Logout
            </button>
          </nav>
        </aside>

        {/* Content */}
        <section className="min-w-0 flex-1">
          {/* Mobile back button */}
          <button
            onClick={() => setActiveSection("personal")}
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
