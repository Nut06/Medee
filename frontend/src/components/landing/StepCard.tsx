import type { Step } from "@/utils/types/landing.type";

interface StepCardProps {
  step: Step;
}

export function StepCard({ step }: StepCardProps) {
  return (
    <div className="flex gap-6">
      {/* Number Badge */}
      <div className="flex-shrink-0">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary text-primary-foreground text-lg font-bold">
          {step.number}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1">
        <h3 className="mb-2 text-xl font-semibold">{step.title}</h3>
        <p className="text-muted-foreground">{step.description}</p>
      </div>
    </div>
  );
}
