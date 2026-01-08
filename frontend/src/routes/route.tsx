import { Routes, Route, Navigate } from "react-router-dom";
import RegisterPage from "@/views/RegisterPage";
import LoginPage from "@/views/LoginPage";
import LandingPage from "@/views/LandingPage";
import ProtectedRoute from "./ProtectedRoutes";
import GuestRoute from "./GuestRoute";
import RequireRole from "./RequireRole";
import CompanyPage from "@/views/CompanyPage";
import ProfilePage from "@/views/ProfilePage";
import EditProfilePage from "@/views/EditProfilePage";

export const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />

      <Route element={<GuestRoute />}>
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/login" element={<LoginPage />} />
      </Route>

      {/* Protected routes (require authentication) */}
      <Route element={<ProtectedRoute />}>
        <Route element={<RequireRole role="company" />}>
          <Route path="/company" element={<CompanyPage />} />
        </Route>
        <Route element={<RequireRole role="user" />}>
          <Route path="/user" element={<ProfilePage />} />
          {/* <Route path="/user/edit-profile" element={<EditProfilePage />} /> */}
        </Route>
      </Route>

      {/* 404 - Redirect to login */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
};
