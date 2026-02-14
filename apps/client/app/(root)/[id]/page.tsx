import MediaGallery from "../../components/sections/Listing/MediaGallery";
import ProductDetails from "../../components/sections/Listing/ProductDetails";
import QuickMessages from "../../components/sections/Listing/QuickMessages";
import SellerCard from "../../components/sections/Listing/SellerCard";
import SimilarListings from "../../components/sections/Listing/SimilarListings";
import {
  dummyListing,
  dummyReviews,
} from "../../components/sections/Listing/types";

const ListingPage = () => {
  const listing = dummyListing;
  const reviews = dummyReviews;

  return (
    <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
      <div className="flex flex-col gap-8 lg:flex-row lg:gap-10">
        {/* ─── Left Column: Media + Details ─── */}
        <div className="flex flex-1 flex-col gap-6 lg:max-w-[58%]">
          <MediaGallery media={listing.media} />
          <ProductDetails listing={listing} />
        </div>

        {/* ─── Right Column: Actions + Seller ─── */}
        <div className="flex flex-col gap-4 lg:w-85 lg:shrink-0">
          {/* Sticky sidebar on desktop */}
          <div className="flex flex-col gap-4 lg:sticky lg:top-28">
            <QuickMessages />
            <SellerCard seller={listing.seller} reviews={reviews} />
          </div>
        </div>
      </div>

      {/* ─── Similar Listings ─── */}
      <SimilarListings />
    </main>
  );
};

export default ListingPage;
