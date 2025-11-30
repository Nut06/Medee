import { useUserStore } from "@/stores/userStore";
import { Navigate, Outlet } from "react-router-dom";

export default function ProtectedRoute() {
  const { isAuth } = useUserStore();

  if (!isAuth)
    return <Navigate to="/login" replace />;
  return <Outlet />;
}