import type { FooterLinks } from "@/utils/types/landing.type";

export const FOOTER_LINKS: FooterLinks = {
  company: [
    { label: "เกี่ยวกับเรา", href: "/about" },
    { label: "ร่วมงานกับเรา", href: "/careers" },
    { label: "บล็อก", href: "/blog" },
  ],
  support: [
    { label: "ติดต่อเรา", href: "/contact" },
    { label: "คำถามที่พบบ่อย", href: "/faq" },
  ],
  legal: [
    { label: "นโยบายความเป็นส่วนตัว", href: "/privacy" },
    { label: "ข้อกำหนดการใช้งาน", href: "/terms" },
  ],
};
