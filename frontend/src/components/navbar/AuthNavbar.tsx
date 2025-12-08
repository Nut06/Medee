import { NavbarBrand } from "./NavbarBrand";
import { AuthNavbarLinks } from "./AuthNavbarLinks";
import { SearchBar } from "./SearchBar";
import { NotificationBell } from "./NotificationBell";
import { UserMenu } from "./UserMenu";

export function AuthNavbar() {
  return (
    <nav className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center gap-4">
          {/* Logo & Brand */}
          <NavbarBrand />

          {/* Desktop Navigation Links */}
          <AuthNavbarLinks />

          {/* Search Bar */}
          <SearchBar />

          {/* Notification Bell */}
          <NotificationBell />

          {/* User Menu */}
          <UserMenu />
        </div>
      </div>
    </nav>
  );
}
