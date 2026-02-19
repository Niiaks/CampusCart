"use client";

import { useState } from "react";
import { ImageIcon, Upload } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";

export default function BrandAppearance() {
  const [profilePreview, setProfilePreview] = useState<string | null>(null);
  const [bannerPreview, setBannerPreview] = useState<string | null>(null);

  const handleFileChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    setter: (url: string | null) => void,
  ) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      toast.error("Please select an image file");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.error("Image must be under 5MB");
      return;
    }

    const url = URL.createObjectURL(file);
    setter(url);
  };

  const handleSave = () => {
    // TODO: upload images and call API
    toast.info("Appearance update is not yet implemented");
  };

  return (
    <div>
      <h2 className="text-lg font-semibold text-foreground">Appearance</h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Customize your brand&apos;s profile picture and banner
      </p>

      <div className="mt-6 space-y-8">
        {/* Profile picture */}
        <div>
          <p className="mb-3 text-sm font-medium text-foreground">
            Profile picture
          </p>
          <div className="flex items-center gap-4">
            <div className="relative flex size-20 items-center justify-center overflow-hidden rounded-xl border-2 border-dashed border-border bg-muted">
              {profilePreview ? (
                <img
                  src={profilePreview}
                  alt="Profile preview"
                  className="size-full object-cover"
                />
              ) : (
                <ImageIcon className="size-8 text-muted-foreground" />
              )}
            </div>
            <div>
              <label className="inline-flex cursor-pointer items-center gap-2 rounded-lg border border-border bg-background px-3 py-2 text-sm font-medium text-foreground transition-colors hover:bg-muted">
                <Upload className="size-4" />
                Upload
                <input
                  type="file"
                  accept="image/*"
                  className="hidden"
                  onChange={(e) => handleFileChange(e, setProfilePreview)}
                />
              </label>
              <p className="mt-1 text-xs text-muted-foreground">
                Recommended: 200×200px, max 5MB
              </p>
            </div>
          </div>
        </div>

        {/* Banner */}
        <div>
          <p className="mb-3 text-sm font-medium text-foreground">
            Banner image
          </p>
          <div className="relative flex h-36 w-full items-center justify-center overflow-hidden rounded-xl border-2 border-dashed border-border bg-muted">
            {bannerPreview ? (
              <img
                src={bannerPreview}
                alt="Banner preview"
                className="size-full object-cover"
              />
            ) : (
              <div className="flex flex-col items-center gap-2 text-muted-foreground">
                <ImageIcon className="size-10" />
                <span className="text-sm">No banner uploaded</span>
              </div>
            )}
          </div>
          <label className="mt-3 inline-flex cursor-pointer items-center gap-2 rounded-lg border border-border bg-background px-3 py-2 text-sm font-medium text-foreground transition-colors hover:bg-muted">
            <Upload className="size-4" />
            Upload Banner
            <input
              type="file"
              accept="image/*"
              className="hidden"
              onChange={(e) => handleFileChange(e, setBannerPreview)}
            />
          </label>
          <p className="mt-1 text-xs text-muted-foreground">
            Recommended: 1200×300px, max 5MB
          </p>
        </div>

        <Button
          onClick={handleSave}
          className="bg-brand text-brand-foreground hover:bg-brand-hover"
        >
          Save Appearance
        </Button>
      </div>
    </div>
  );
}
