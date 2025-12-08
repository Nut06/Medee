import { NavLink } from "react-router-dom";
import { AUTH_NAV_LINKS } from "@/constants/auth-navigation";
import { cn } from "@/lib/utils";

export function AuthNavbarLinks() {
  return (
    <nav className="hidden md:flex md:items-center md:gap-6 md:flex-1 md:justify-center">
      {AUTH_NAV_LINKS.map((link) => (
        <NavLink
          key={link.href}
          to={link.href}
          className={({ isActive }) =>
            cn(
              "text-sm font-medium transition-colors hover:text-primary",
              isActive ? "text-primary" : "text-muted-foreground"
            )
          }
        >
          {link.label}
        </NavLink>
      ))}
    </nav>
  );
}
