import { useUserStore } from "@/stores/userStore";
import type { Role } from "@/utils/types/user.type";
import { Navigate, Outlet, useLocation } from "react-router-dom";

type Props = {
  role: Role;
};

export default function RequireRole({ role }: Props) {
  const { user } = useUserStore();
  const location = useLocation();

  if (!user) return <Navigate to="/login" replace />;

  const hasCompany = user.companies && user.companies.length > 0;

  // If route requires "company" role but user has no companies -> Redirect to User Profile
  if (role === "company" && !hasCompany) {
    return <Navigate to="/user" replace state={{ from: location }} />;
  }

  // If route requires "user" (candidate) role but user belongs to a company -> Redirect to Company Dashboard
  // (Optional: depending on if we want to restrict company members from accessing candidate pages)
  if (role === "user" && hasCompany) {
    return <Navigate to="/company" replace state={{ from: location }} />;
  }

  return <Outlet />;
}
