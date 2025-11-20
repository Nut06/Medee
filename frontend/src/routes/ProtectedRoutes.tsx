import { useUserStore } from "@/stores/userStore";
import { Navigate, Outlet } from "react-router-dom";

export default function ProtectedRoute(){
    const {user} = useUserStore();

    if (!user) return <Navigate to="/login" replace/> 

    return <Outlet/>;
};