import Link from "next/link";
import ListingCard, { dummyListings, type Listing } from "../../ui/ListingCard";

interface ListingGridProps {
  listings?: Listing[];
}

const ListingGrid = ({ listings = dummyListings }: ListingGridProps) => {
  return (
    <div className="mt-4 grid grid-cols-2 gap-3 sm:gap-4 md:grid-cols-3 lg:grid-cols-4">
      {listings.map((listing) => (
        <ListingCard key={listing.id} listing={listing} />
      ))}
    </div>
  );
};

export default ListingGrid;
