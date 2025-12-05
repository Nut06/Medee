export interface Feature {
  icon: React.ComponentType<{ className?: string }>;
  title: string;
  description: string;
}

export interface Step {
  number: number;
  title: string;
  description: string;
}

export interface Testimonial {
  name: string;
  role: string;
  avatar: string;
  content: string;
}

export interface FooterLink {
  label: string;
  href: string;
}

export interface FooterLinks {
  company: FooterLink[];
  support: FooterLink[];
  legal: FooterLink[];
}
