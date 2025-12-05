import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useUserStore } from "@/stores/userStore";

export function NavbarActions() {
  const { isAuth } = useUserStore();

  return (
    <div className="hidden md:flex md:items-center md:gap-3">
      {isAuth ? (
        <Button asChild>
          <Link to="/user">แดชบอร์ด</Link>
        </Button>
      ) : (
        <>
          <Button variant="ghost" asChild>
            <Link to="/login">เข้าสู่ระบบ</Link>
          </Button>
          <Button asChild>
            <Link to="/register">สมัครสมาชิก</Link>
          </Button>
        </>
      )}
    </div>
  );
}
