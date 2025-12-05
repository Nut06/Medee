import { Link } from "react-router-dom";
import { Separator } from "@/components/ui/separator";
import { Logo } from "@/components/Logo";
import { FOOTER_LINKS } from "@/constants/footer";
import { BRAND_NAME } from "@/constants/navigation";

export function Footer() {
  return (
    <footer className="border-t bg-background py-12">
      <div className="container mx-auto px-4">
        {/* Footer Links Grid */}
        <div className="mb-8 grid gap-8 md:grid-cols-4">
          {/* Company */}
          <div>
            <h3 className="mb-4 font-semibold">บริษัท</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              {FOOTER_LINKS.company.map((link) => (
                <li key={link.href}>
                  <Link
                    to={link.href}
                    className="transition-colors hover:text-primary"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Support */}
          <div>
            <h3 className="mb-4 font-semibold">ช่วยเหลือ</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              {FOOTER_LINKS.support.map((link) => (
                <li key={link.href}>
                  <Link
                    to={link.href}
                    className="transition-colors hover:text-primary"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Legal */}
          <div>
            <h3 className="mb-4 font-semibold">กฎหมาย</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              {FOOTER_LINKS.legal.map((link) => (
                <li key={link.href}>
                  <Link
                    to={link.href}
                    className="transition-colors hover:text-primary"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Logo */}
          <div>
            <Link to="/" className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary">
                <Logo className="h-5 w-5 text-primary-foreground" />
              </div>
              <span className="font-bold">{BRAND_NAME}</span>
            </Link>
          </div>
        </div>

        <Separator className="mb-8" />

        {/* Copyright */}
        <div className="text-center text-sm text-muted-foreground">
          © {new Date().getFullYear()} {BRAND_NAME}. สงวนลิขสิทธิ์.
        </div>
      </div>
    </footer>
  );
}
