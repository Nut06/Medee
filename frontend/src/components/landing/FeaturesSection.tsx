import { FEATURES } from "@/constants/features";
import { FeatureCard } from "./FeatureCard";

export default function FeaturesSection() {
  return (
    <section className="bg-muted/50 py-20">
      <div className="container mx-auto px-4">
        {/* Heading */}
        <div className="mb-12 text-center">
          <h2 className="mb-4 text-4xl font-bold">ปลดล็อกศักยภาพของคุณ</h2>
          <p className="mx-auto max-w-2xl text-lg text-muted-foreground">
            แพลตฟอร์มของเราให้เครื่องมือและโอกาสที่คุณต้องการ
            เพื่อเริ่มต้นอาชีพและพัฒนาทักษะของคุณ
          </p>
        </div>

        {/* Feature Cards Grid */}
        <div className="grid gap-8 md:grid-cols-3">
          {FEATURES.map((feature, index) => (
            <FeatureCard key={index} feature={feature} />
          ))}
        </div>
      </div>
    </section>
  );
}
