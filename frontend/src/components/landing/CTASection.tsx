import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";

export default function CTASection() {
  return (
    <section className="py-20">
      <div className="container mx-auto px-4">
        <div className="mx-auto max-w-4xl rounded-2xl bg-primary p-12 text-center text-primary-foreground">
          <h2 className="mb-4 text-4xl font-bold">
            พร้อมเริ่มต้นเส้นทางของคุณแล้วหรือยัง?
          </h2>
          <p className="mb-8 text-lg opacity-90">
            ไม่ว่าคุณจะกำลังมองหาโอกาสในการเริ่มต้นอาชีพ
            หรือกำลังมองหาผู้มีความสามารถ เรามีทุกอย่างที่คุณต้องการ
          </p>
          <Button size="lg" variant="secondary" asChild>
            <Link to="/register">เริ่มต้นเลยตอนนี้</Link>
          </Button>
        </div>
      </div>
    </section>
  );
}
