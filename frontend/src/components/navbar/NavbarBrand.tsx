import { Link } from "react-router-dom";
import { Logo } from "@/components/Logo";
import { BRAND_NAME } from "@/constants/navigation";

export function NavbarBrand() {
  return (
    <Link
      to="/"
      className="flex items-center gap-2 hover:opacity-80 transition-opacity"
    >
      <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary">
        <Logo className="h-5 w-5 text-primary-foreground" />
      </div>
      <span className="text-xl font-bold">{BRAND_NAME}</span>
    </Link>
  );
}
