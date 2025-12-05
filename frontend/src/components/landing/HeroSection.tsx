import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";

export default function HeroSection() {
  return (
    <section className="relative overflow-hidden bg-background py-20 sm:py-32">
      <div className="container mx-auto px-4">
        <div className="grid gap-12 lg:grid-cols-2 lg:gap-8 items-center">
          {/* Left: Text Content */}
          <div className="flex flex-col justify-center space-y-8">
            <div className="space-y-4">
              <h1 className="text-4xl font-bold tracking-tight sm:text-5xl xl:text-6xl">
                สร้างประสบการณ์จริง{" "}
                <span className="text-primary">จากทุกที่ทุกเวลา</span>
              </h1>
              <p className="text-lg text-muted-foreground sm:text-xl max-w-2xl">
                เชื่อมต่อกับบริษัทชั้นนำและสร้างประสบการณ์การทำงานผ่านการฝึกงานออนไลน์
                สร้างพอร์ตโฟลิโอและเริ่มต้นอาชีพของคุณ
              </p>
            </div>

            {/* CTA Buttons */}
            <div className="flex flex-col gap-4 sm:flex-row">
              <Button size="lg" asChild className="text-base">
                <Link to="/register">สมัครสมาชิกสำหรับผู้สมัคร</Link>
              </Button>
              <Button size="lg" variant="outline" asChild className="text-base">
                <Link to="/register?type=company">สมัครสมาชิกสำหรับบริษัท</Link>
              </Button>
            </div>

            {/* Stats or Trust Indicators */}
            <div className="flex flex-wrap gap-8 pt-4">
              <div>
                <p className="text-3xl font-bold text-primary">500+</p>
                <p className="text-sm text-muted-foreground">
                  โปรเจกต์ที่เปิดรับ
                </p>
              </div>
              <div>
                <p className="text-3xl font-bold text-primary">1,000+</p>
                <p className="text-sm text-muted-foreground">นักศึกษา</p>
              </div>
              <div>
                <p className="text-3xl font-bold text-primary">100+</p>
                <p className="text-sm text-muted-foreground">บริษัทพันธมิตร</p>
              </div>
            </div>
          </div>

          {/* Right: Hero Image */}
          <div className="relative lg:order-last">
            <div className="relative aspect-[4/3] overflow-hidden rounded-2xl bg-muted shadow-2xl">
              <img
                src="https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=800&auto=format&fit=crop&q=80"
                alt="Team collaboration"
                className="h-full w-full object-cover"
              />
              {/* Overlay gradient */}
              <div className="absolute inset-0 bg-gradient-to-tr from-primary/20 to-transparent" />
            </div>

            {/* Floating card (optional decoration) */}
            <div className="absolute -bottom-6 -left-6 hidden rounded-lg bg-background p-4 shadow-lg sm:block">
              <div className="flex items-center gap-3">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                  <svg
                    className="h-6 w-6 text-primary"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
                <div>
                  <p className="text-sm font-semibold">
                    บริษัทที่ได้รับการยืนยัน
                  </p>
                  <p className="text-xs text-muted-foreground">
                    ได้รับความไว้วางใจจากผู้นำในอุตสาหกรรม
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Background decoration */}
      <div className="absolute inset-0 -z-10 overflow-hidden">
        <div className="absolute left-1/2 top-0 -translate-x-1/2 -translate-y-1/2">
          <div className="h-[600px] w-[600px] rounded-full bg-primary/5 blur-3xl" />
        </div>
      </div>
    </section>
  );
}
