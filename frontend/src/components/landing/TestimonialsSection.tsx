import { TESTIMONIALS } from "@/constants/testimonials";
import { TestimonialCard } from "./TestimonialCard";

export default function TestimonialsSection() {
  return (
    <section className="bg-muted/50 py-20">
      <div className="container mx-auto px-4">
        {/* Heading */}
        <div className="mb-12 text-center">
          <h2 className="mb-4 text-4xl font-bold">
            ได้รับความไว้วางใจจากนักศึกษาและบริษัท
          </h2>
          <p className="mx-auto max-w-2xl text-lg text-muted-foreground">
            เรียนรู้จากประสบการณ์ของผู้ใช้งานของเราว่าพวกเขาได้รับอะไรจากแพลตฟอร์ม
          </p>
        </div>

        {/* Testimonial Cards Grid */}
        <div className="mx-auto grid max-w-5xl gap-8 md:grid-cols-2">
          {TESTIMONIALS.map((testimonial, index) => (
            <TestimonialCard key={index} testimonial={testimonial} />
          ))}
        </div>
      </div>
    </section>
  );
}
