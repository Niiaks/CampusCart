"use client";

import { useAuth } from "@/hooks/useAuth";
import Categories from "../components/sections/Home/categories";
import ListingCarousel from "../components/sections/Home/ListingCarousel";
import ListingGrid from "../components/sections/Home/ListingGrid";
import Section from "../components/sections/Home/section";

const Home = () => {
  const { data: user } = useAuth();
  console.log("the user is", user);
  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-2 px-4 py-4 sm:px-6 lg:px-8">
      <Section LeadTitle="Trending Categories">
        <Categories />
      </Section>

      <Section LeadTitle="Recommended For You">
        <ListingGrid />
      </Section>

      <Section LeadTitle="Top Trending">
        <ListingCarousel />
      </Section>

      <Section LeadTitle="Top Seller">
        <ListingGrid />
      </Section>

      <Section LeadTitle="New In">
        <ListingCarousel />
      </Section>
    </main>
  );
};

export default Home;
