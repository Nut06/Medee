import "./App.css";
import { AppRoutes } from "./routes/route";
import { useEffect } from "react";
import { useUserStore } from "./stores/userStore";
import { getUser } from "./services/user.service";

function App() {
  const { setUser, setAuth } = useUserStore();

  useEffect(() => {
    const checkAuth = async () => {
      const accessToken = localStorage.getItem("accessToken");
      if (!accessToken) {
        setAuth(false);
        return;
      }

      const user = await getUser();
      if (user) {
        setUser(user);
        setAuth(true);
      }
    };
    checkAuth();
  }, [ setUser, setAuth]);

  return (
    <>
      <AppRoutes />
    </>
  );
}

export default App;
