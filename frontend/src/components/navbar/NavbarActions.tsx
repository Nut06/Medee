import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useUserStore } from "@/stores/userStore";

export function NavbarActions() {
  const { isAuth } = useUserStore();

  return (
    <div className="hidden md:flex md:items-center md:gap-3">
      {isAuth ? (
        <Button asChild>
          <Link to="/user">Dashboard</Link>
        </Button>
      ) : (
        <>
          <Button variant="ghost" asChild>
            <Link to="/login">Log In</Link>
          </Button>
          <Button asChild>
            <Link to="/register">Sign Up</Link>
          </Button>
        </>
      )}
    </div>
  );
}
