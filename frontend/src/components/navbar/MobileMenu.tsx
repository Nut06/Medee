import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu } from "lucide-react";
import { useState } from "react";
import { useUserStore } from "@/stores/userStore";
import { NAV_LINKS } from "@/constants/navigation";
import { NavbarBrand } from "./NavbarBrand";

export function MobileMenu() {
  const [isOpen, setIsOpen] = useState(false);
  const { isAuth } = useUserStore();

  const closeMenu = () => setIsOpen(false);

  return (
    <Sheet open={isOpen} onOpenChange={setIsOpen}>
      <SheetTrigger asChild className="md:hidden">
        <Button variant="ghost" size="icon">
          <Menu className="h-6 w-6" />
          <span className="sr-only">Toggle menu</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="right" className="w-[300px] sm:w-[400px]">
        <div className="flex flex-col gap-6 py-6">
          {/* Mobile Logo */}
          <div onClick={closeMenu}>
            <NavbarBrand />
          </div>

          {/* Mobile Navigation Links */}
          <nav className="flex flex-col gap-4">
            {NAV_LINKS.map((link) => (
              <Link
                key={link.href}
                to={link.href}
                className="text-lg font-medium text-foreground transition-colors hover:text-primary"
                onClick={closeMenu}
              >
                {link.label}
              </Link>
            ))}
          </nav>

          {/* Mobile Auth Buttons */}
          <div className="flex flex-col gap-3 pt-4 border-t">
            {isAuth ? (
              <Button asChild onClick={closeMenu}>
                <Link to="/user">แดชบอร์ด</Link>
              </Button>
            ) : (
              <>
                <Button variant="outline" asChild onClick={closeMenu}>
                  <Link to="/login">เข้าสู่ระบบ</Link>
                </Button>
                <Button asChild onClick={closeMenu}>
                  <Link to="/register">สมัครสมาชิก</Link>
                </Button>
              </>
            )}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
