import {Routes, Route, Navigate} from "react-router-dom"
import {RegisterPage} from "@/views/RegisterPage"
import {LoginPage} from "@/views/LoginPage"

export const AppRoutes = () => {
    return(
        <Routes>
            <Route path="/" element={<Navigate to="/register" replace/>}></Route>
            {/* public route */}
            <Route path="/register" element={<RegisterPage />}></Route>
            <Route path="/login" element={<LoginPage />}></Route>
            {/* <Route path="/auth/callback" element={<AuthCallback />}></Route> */}

            {/* protected route */}
            <Route element={<ProtectedRoute />}>

            </Route>
        </Routes>
    );
};