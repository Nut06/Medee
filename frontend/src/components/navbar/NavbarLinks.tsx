import { Link } from "react-router-dom";
import { NAV_LINKS } from "@/constants/navigation";

export function NavbarLinks() {
  return (
    <div className="hidden md:flex md:items-center md:gap-8">
      {NAV_LINKS.map((link) => (
        <Link
          key={link.href}
          to={link.href}
          className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
        >
          {link.label}
        </Link>
      ))}
    </div>
  );
}
