import { Facebook, Instagram, Mail, Twitter } from "lucide-react";

const footerLinks = {
  marketplace: [
    { label: "Browse Listings", href: "#" },
    { label: "Sell an Item", href: "#" },
    { label: "Categories", href: "#" },
    { label: "Trending", href: "#" },
  ],
  support: [
    { label: "Help Center", href: "#" },
    { label: "Safety Tips", href: "#" },
    { label: "Contact Us", href: "#" },
    { label: "Report a Listing", href: "#" },
  ],
  company: [
    { label: "About Us", href: "#" },
    { label: "Blog", href: "#" },
    { label: "Careers", href: "#" },
    { label: "Press", href: "#" },
  ],
  legal: [
    { label: "Terms of Service", href: "#" },
    { label: "Privacy Policy", href: "#" },
    { label: "Cookie Policy", href: "#" },
  ],
};

const socials = [
  { icon: Twitter, href: "#", label: "Twitter" },
  { icon: Instagram, href: "#", label: "Instagram" },
  { icon: Facebook, href: "#", label: "Facebook" },
  { icon: Mail, href: "#", label: "Email" },
];

const Footer = () => {
  return (
    <footer className="border-t border-border/60 bg-brand text-brand-foreground">
      {/* Main grid */}
      <div className="mx-auto grid max-w-7xl gap-10 px-4 py-12 sm:grid-cols-2 sm:px-6 lg:grid-cols-5 lg:px-8">
        {/* Brand column */}
        <div className="lg:col-span-2">
          <div className="flex items-center gap-3">
            <div className="flex size-10 items-center justify-center rounded-full bg-brand-light/60 text-base font-semibold uppercase">
              CC
            </div>
            <span className="text-lg font-semibold tracking-wider">
              campusCart
            </span>
          </div>
          <p className="mt-4 max-w-xs text-sm leading-relaxed text-brand-foreground/70">
            The student marketplace for your campus. Buy and sell textbooks,
            electronics, dorm essentials and more — safely and locally.
          </p>
          <div className="mt-6 flex gap-3">
            {socials.map((s) => (
              <a
                key={s.label}
                href={s.href}
                aria-label={s.label}
                className="flex size-9 items-center justify-center rounded-full bg-brand-foreground/10 text-brand-foreground/70 transition-colors hover:bg-brand-foreground/20 hover:text-brand-foreground"
              >
                <s.icon className="size-4" />
              </a>
            ))}
          </div>
        </div>

        {/* Link columns */}
        <LinkColumn title="Marketplace" links={footerLinks.marketplace} />
        <LinkColumn title="Support" links={footerLinks.support} />
        <LinkColumn title="Company" links={footerLinks.company} />
      </div>

      {/* Bottom bar */}
      <div className="border-t border-brand-foreground/10">
        <div className="mx-auto flex max-w-7xl flex-col items-center justify-between gap-4 px-4 py-5 sm:flex-row sm:px-6 lg:px-8">
          <p className="text-xs text-brand-foreground/50">
            &copy; {new Date().getFullYear()} campusCart. All rights reserved.
          </p>
          <div className="flex gap-6">
            {footerLinks.legal.map((link) => (
              <a
                key={link.label}
                href={link.href}
                className="text-xs text-brand-foreground/50 no-underline transition-colors hover:text-brand-foreground/80"
              >
                {link.label}
              </a>
            ))}
          </div>
        </div>
      </div>
    </footer>
  );
};

const LinkColumn = ({
  title,
  links,
}: {
  title: string;
  links: { label: string; href: string }[];
}) => (
  <div>
    <h4 className="text-sm font-semibold uppercase tracking-wider text-brand-foreground/90">
      {title}
    </h4>
    <ul className="mt-4 space-y-2.5">
      {links.map((link) => (
        <li key={link.label}>
          <a
            href={link.href}
            className="text-sm text-brand-foreground/60 no-underline transition-colors hover:text-brand-foreground"
          >
            {link.label}
          </a>
        </li>
      ))}
    </ul>
  </div>
);

export default Footer;
