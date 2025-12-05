import { NavbarBrand } from "./navbar/NavbarBrand";
import { NavbarLinks } from "./navbar/NavbarLinks";
import { NavbarActions } from "./navbar/NavbarActions";
import { MobileMenu } from "./navbar/MobileMenu";

export function Navbar() {
  return (
    <nav className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Logo & Brand */}
          <NavbarBrand />

          {/* Desktop Navigation Links */}
          <NavbarLinks />

          {/* Desktop Auth Buttons */}
          <NavbarActions />

          {/* Mobile Menu */}
          <MobileMenu />
        </div>
      </div>
    </nav>
  );
}
