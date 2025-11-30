import { useUserStore } from "@/stores/userStore";
import type { Role } from "@/utils/types/user.type";
import { Navigate, Outlet, useLocation } from "react-router-dom";

type Props = {
  role: Role;
};

export default function RequireRole({ role }: Props) {
  const { user } = useUserStore();
  const location = useLocation();

  if (!user || !user.company) return <Navigate to="/login" replace />;

  if ( user && !user.company)
    return <Navigate to="/user" replace state={{ from: location }} />;

  return <Outlet />;
}
