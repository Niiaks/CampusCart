import { ArrowRight } from "lucide-react";

interface SectionProps {
  children: React.ReactNode;
  LeadTitle: string;
}

const Section = ({ children, LeadTitle }: SectionProps) => {
  return (
    <section className="w-full mb-4">
      <div className="mx-auto flex w-full max-w-7xl flex-col py-8">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <h2 className="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl">
            {LeadTitle}
          </h2>
          <button className="inline-flex items-center gap-1 text-sm font-medium text-brand hover:text-brand-hover">
            <span>See all</span>
            <ArrowRight size={18} />
          </button>
        </div>
        <div className="w-full">{children}</div>
      </div>
    </section>
  );
};

export default Section;
