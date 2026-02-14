import ListingCard, { dummyListings } from "../../ui/ListingCard";

const SimilarListings = () => {
  return (
    <section className="mt-10 border-t border-border/40 pt-8">
      <h2 className="text-xl font-semibold tracking-tight text-foreground sm:text-2xl">
        Similar Listings
      </h2>
      <div className="mt-4 grid grid-cols-2 gap-3 sm:gap-4 md:grid-cols-3 lg:grid-cols-4">
        {dummyListings.slice(0, 4).map((listing) => (
          <ListingCard key={listing.id} listing={listing} />
        ))}
      </div>
    </section>
  );
};

export default SimilarListings;
