import { useUserStore } from "@/stores/userStore";
import { Navigate, Outlet, useLocation } from "react-router-dom";

export default function GuestRoute() {
  const { isAuth } = useUserStore();
  const location = useLocation();

  if (isAuth) {
    // Get the intended destination from location state, or use default
    const from = location.state?.from?.pathname || "/user";

    // Redirect authenticated users away from guest-only pagesP
    return <Navigate to={from} state={{ from: location }} replace />;
  }

  // Allow access to guest-only routes
  return <Outlet />;
}
