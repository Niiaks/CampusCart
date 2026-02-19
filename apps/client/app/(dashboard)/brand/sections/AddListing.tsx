"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Loader2,
  Upload,
  X,
  ImageIcon,
  Video,
  PackagePlus,
  Play,
} from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  createListingSchema,
  type CreateListingFormData,
} from "@/common/schemas";

// TODO: fetch categories from API
const dummyCategories = [
  { category_id: "1", name: "Electronics" },
  { category_id: "2", name: "Textbooks" },
  { category_id: "3", name: "Dorm Essentials" },
  { category_id: "4", name: "Clothing" },
  { category_id: "5", name: "Accessories" },
  { category_id: "6", name: "Sports" },
];

const conditions = [
  { value: "new", label: "New", description: "Brand new, unused" },
  { value: "used", label: "Used", description: "Previously owned, good shape" },
  {
    value: "second-hand",
    label: "Second-hand",
    description: "Shows signs of wear",
  },
] as const;

export default function AddListing() {
  const [isPending, setIsPending] = useState(false);
  const [images, setImages] = useState<{ file: File; preview: string }[]>([]);
  const [videos, setVideos] = useState<{ file: File; preview: string }[]>([]);

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors },
  } = useForm<CreateListingFormData>({
    resolver: zodResolver(createListingSchema),
    defaultValues: {
      title: "",
      description: "",
      price: 0,
      condition: undefined,
      category_id: "",
      brand: "",
      model: "",
      size: "",
      storage_size: "",
      color: "",
    },
  });

  const selectedCondition = watch("condition");

  const handleImageAdd = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files ?? []);
    if (images.length + files.length > 8) {
      toast.error("Maximum 8 images allowed");
      return;
    }

    const newImages = files
      .filter((f) => {
        if (!f.type.startsWith("image/")) {
          toast.error(`${f.name} is not an image`);
          return false;
        }
        if (f.size > 5 * 1024 * 1024) {
          toast.error(`${f.name} exceeds 5MB`);
          return false;
        }
        return true;
      })
      .map((file) => ({ file, preview: URL.createObjectURL(file) }));

    setImages((prev) => [...prev, ...newImages]);
    e.target.value = "";
  };

  const removeImage = (index: number) => {
    setImages((prev) => {
      URL.revokeObjectURL(prev[index].preview);
      return prev.filter((_, i) => i !== index);
    });
  };

  const handleVideoAdd = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files ?? []);
    if (videos.length + files.length > 2) {
      toast.error("Maximum 2 videos allowed");
      return;
    }

    const newVideos = files
      .filter((f) => {
        if (!f.type.startsWith("video/")) {
          toast.error(`${f.name} is not a video`);
          return false;
        }
        if (f.size > 50 * 1024 * 1024) {
          toast.error(`${f.name} exceeds 50MB`);
          return false;
        }
        return true;
      })
      .map((file) => ({ file, preview: URL.createObjectURL(file) }));

    setVideos((prev) => [...prev, ...newVideos]);
    e.target.value = "";
  };

  const removeVideo = (index: number) => {
    setVideos((prev) => {
      URL.revokeObjectURL(prev[index].preview);
      return prev.filter((_, i) => i !== index);
    });
  };

  const onSubmit = async (data: CreateListingFormData) => {
    if (images.length === 0) {
      toast.error("Please add at least one image");
      return;
    }

    setIsPending(true);
    try {
      // TODO: upload images & videos, get URLs, then call API with full payload
      const payload = {
        ...data,
        image_url: images.map((img) => img.preview), // replace with real URLs after upload
        video_url: videos.map((vid) => vid.preview), // replace with real URLs after upload
      };
      console.log("Listing payload:", payload);
      toast.info("Listing creation is not yet implemented");
      // reset();
      // setImages([]);
    } finally {
      setIsPending(false);
    }
  };

  return (
    <div>
      <div className="flex items-center gap-3">
        <div className="flex size-10 items-center justify-center rounded-lg bg-brand-muted">
          <PackagePlus className="size-5 text-brand" />
        </div>
        <div>
          <h2 className="text-lg font-semibold text-foreground">
            Add New Listing
          </h2>
          <p className="text-sm text-muted-foreground">
            Fill in the details to list your item for sale
          </p>
        </div>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="mt-8 space-y-8">
        {/* ─── Images ─── */}
        <fieldset>
          <legend className="mb-3 text-sm font-semibold text-foreground">
            Photos <span className="text-muted-foreground">(max 8)</span>
          </legend>
          <div className="flex flex-wrap gap-3">
            {images.map((img, i) => (
              <div
                key={i}
                className="group relative size-24 overflow-hidden rounded-lg border border-border"
              >
                <img
                  src={img.preview}
                  alt={`Upload ${i + 1}`}
                  className="size-full object-cover"
                />
                <button
                  type="button"
                  onClick={() => removeImage(i)}
                  className="absolute right-1 top-1 flex size-5 items-center justify-center rounded-full bg-black/60 text-white opacity-0 transition-opacity group-hover:opacity-100"
                >
                  <X className="size-3" />
                </button>
              </div>
            ))}

            {images.length < 8 && (
              <label className="flex size-24 cursor-pointer flex-col items-center justify-center gap-1 rounded-lg border-2 border-dashed border-border bg-muted/50 text-muted-foreground transition-colors hover:border-brand hover:text-brand">
                <Upload className="size-5" />
                <span className="text-[10px] font-medium">Add</span>
                <input
                  type="file"
                  accept="image/*"
                  multiple
                  className="hidden"
                  onChange={handleImageAdd}
                />
              </label>
            )}
          </div>
        </fieldset>

        {/* ─── Videos ─── */}
        <fieldset>
          <legend className="mb-3 text-sm font-semibold text-foreground">
            Videos{" "}
            <span className="text-muted-foreground">
              (optional, max 2, up to 50MB each)
            </span>
          </legend>
          <div className="flex flex-wrap gap-3">
            {videos.map((vid, i) => (
              <div
                key={i}
                className="group relative h-24 w-36 overflow-hidden rounded-lg border border-border bg-black"
              >
                <video
                  src={vid.preview}
                  className="size-full object-cover"
                  muted
                />
                <div className="absolute inset-0 flex items-center justify-center">
                  <Play className="size-6 text-white/80" />
                </div>
                <button
                  type="button"
                  onClick={() => removeVideo(i)}
                  className="absolute right-1 top-1 flex size-5 items-center justify-center rounded-full bg-black/60 text-white opacity-0 transition-opacity group-hover:opacity-100"
                >
                  <X className="size-3" />
                </button>
              </div>
            ))}

            {videos.length < 2 && (
              <label className="flex h-24 w-36 cursor-pointer flex-col items-center justify-center gap-1 rounded-lg border-2 border-dashed border-border bg-muted/50 text-muted-foreground transition-colors hover:border-brand hover:text-brand">
                <Video className="size-5" />
                <span className="text-[10px] font-medium">Add Video</span>
                <input
                  type="file"
                  accept="video/*"
                  className="hidden"
                  onChange={handleVideoAdd}
                />
              </label>
            )}
          </div>
        </fieldset>

        {/* ─── Basic info ─── */}
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-1.5 sm:col-span-2">
            <label
              htmlFor="title"
              className="text-sm font-medium text-foreground"
            >
              Title
            </label>
            <Input
              id="title"
              placeholder="e.g. MacBook Pro M2 – 256GB"
              {...register("title")}
            />
            {errors.title && (
              <p className="text-xs text-destructive">{errors.title.message}</p>
            )}
          </div>

          <div className="space-y-1.5 sm:col-span-2">
            <label
              htmlFor="description"
              className="text-sm font-medium text-foreground"
            >
              Description
            </label>
            <textarea
              id="description"
              rows={4}
              placeholder="Describe your item — condition details, what's included, etc."
              className="flex w-full rounded-xl border border-input bg-background/80 px-4 py-2 text-sm text-foreground shadow-sm transition-all placeholder:text-muted-foreground/70 focus-visible:border-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/30 disabled:cursor-not-allowed disabled:opacity-50"
              {...register("description")}
            />
            {errors.description && (
              <p className="text-xs text-destructive">
                {errors.description.message}
              </p>
            )}
          </div>

          <div className="space-y-1.5">
            <label
              htmlFor="price"
              className="text-sm font-medium text-foreground"
            >
              Price (GH₵)
            </label>
            <Input
              id="price"
              type="number"
              min={0}
              step={1}
              placeholder="0"
              {...register("price", { valueAsNumber: true })}
            />
            {errors.price && (
              <p className="text-xs text-destructive">{errors.price.message}</p>
            )}
          </div>

          <div className="space-y-1.5">
            <label
              htmlFor="category_id"
              className="text-sm font-medium text-foreground"
            >
              Category
            </label>
            <select
              id="category_id"
              className="flex h-11 w-full rounded-xl border border-input bg-background/80 px-4 py-2 text-sm text-foreground shadow-sm transition-all focus-visible:border-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/30"
              {...register("category_id")}
            >
              <option value="">Select a category</option>
              {dummyCategories.map((cat) => (
                <option key={cat.category_id} value={cat.category_id}>
                  {cat.name}
                </option>
              ))}
            </select>
            {errors.category_id && (
              <p className="text-xs text-destructive">
                {errors.category_id.message}
              </p>
            )}
          </div>
        </div>

        {/* ─── Condition ─── */}
        <fieldset>
          <legend className="mb-3 text-sm font-semibold text-foreground">
            Condition
          </legend>
          <div className="grid gap-3 sm:grid-cols-3">
            {conditions.map(({ value, label, description }) => (
              <button
                key={value}
                type="button"
                onClick={() =>
                  setValue("condition", value, { shouldValidate: true })
                }
                className={`rounded-lg border p-3 text-left transition-colors ${
                  selectedCondition === value
                    ? "border-brand bg-brand-muted ring-2 ring-brand/30"
                    : "border-border hover:border-brand/40"
                }`}
              >
                <p className="text-sm font-medium text-foreground">{label}</p>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  {description}
                </p>
              </button>
            ))}
          </div>
          {errors.condition && (
            <p className="mt-1 text-xs text-destructive">
              {errors.condition.message}
            </p>
          )}
        </fieldset>

        {/* ─── Optional specs ─── */}
        <fieldset>
          <legend className="mb-3 text-sm font-semibold text-foreground">
            Specifications{" "}
            <span className="font-normal text-muted-foreground">
              (optional)
            </span>
          </legend>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div className="space-y-1.5">
              <label
                htmlFor="brand"
                className="text-xs font-medium text-muted-foreground"
              >
                Brand
              </label>
              <Input
                id="brand"
                placeholder="e.g. Apple"
                {...register("brand")}
              />
            </div>
            <div className="space-y-1.5">
              <label
                htmlFor="model"
                className="text-xs font-medium text-muted-foreground"
              >
                Model
              </label>
              <Input
                id="model"
                placeholder="e.g. MacBook Pro 14″"
                {...register("model")}
              />
            </div>
            <div className="space-y-1.5">
              <label
                htmlFor="size"
                className="text-xs font-medium text-muted-foreground"
              >
                Size
              </label>
              <Input
                id="size"
                placeholder="e.g. XL, 42, 14-inch"
                {...register("size")}
              />
            </div>
            <div className="space-y-1.5">
              <label
                htmlFor="storage_size"
                className="text-xs font-medium text-muted-foreground"
              >
                Storage
              </label>
              <Input
                id="storage_size"
                placeholder="e.g. 256GB"
                {...register("storage_size")}
              />
            </div>
            <div className="space-y-1.5">
              <label
                htmlFor="color"
                className="text-xs font-medium text-muted-foreground"
              >
                Color
              </label>
              <Input
                id="color"
                placeholder="e.g. Space Gray"
                {...register("color")}
              />
            </div>
          </div>
        </fieldset>

        {/* ─── Submit ─── */}
        <div className="flex items-center gap-3 border-t border-border/60 pt-6">
          <Button
            type="submit"
            disabled={isPending}
            className="bg-brand text-brand-foreground hover:bg-brand-hover"
          >
            {isPending ? (
              <Loader2 className="animate-spin" />
            ) : (
              "Publish Listing"
            )}
          </Button>
          <Button
            type="button"
            variant="outline"
            onClick={() => {
              reset();
              setImages([]);
              setVideos([]);
            }}
          >
            Clear
          </Button>
        </div>
      </form>
    </div>
  );
}
