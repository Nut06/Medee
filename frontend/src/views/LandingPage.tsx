import { Navbar } from "@/components/Navbar";
import HeroSection from "@/components/landing/HeroSection";
import FeaturesSection from "@/components/landing/FeaturesSection";
import HowItWorkSection from "@/components/landing/HowItWorkSection";
import TestimonialSection from "@/components/landing/TestimonialSection";
import CTASection from "@/components/landing/CTASection";

export default function LandingPage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <main className="flex-1">
        <HeroSection />
        <FeaturesSection />
        <HowItWorkSection />
        <TestimonialSection />
        <CTASection />
      </main>
    </div>
  );
}
