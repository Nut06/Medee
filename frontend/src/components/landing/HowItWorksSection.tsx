import { HOW_IT_WORKS_STEPS } from "@/constants/steps";
import { StepCard } from "./StepCard";

export default function HowItWorksSection() {
  return (
    <section className="py-20">
      <div className="container mx-auto px-4">
        {/* Heading */}
        <h2 className="mb-12 text-center text-4xl font-bold">วิธีการใช้งาน</h2>

        {/* Steps */}
        <div className="mx-auto max-w-3xl space-y-8">
          {HOW_IT_WORKS_STEPS.map((step) => (
            <StepCard key={step.number} step={step} />
          ))}
        </div>
      </div>
    </section>
  );
}
