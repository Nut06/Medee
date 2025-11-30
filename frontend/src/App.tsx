import "./App.css";
import { AppRoutes } from "./routes/route";
import { useEffect } from "react";
import { useUserStore } from "./stores/userStore";
import { getUser } from "./services/userService";

function App() {
  const { setUser, setAuth } = useUserStore();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const user = await getUser();
        if (user) {
          setUser(user);
          setAuth(true);
        }
      } catch (error) {
        console.log("Not authenticated");
        setAuth(false);
        setUser(null);
      }
    };
    checkAuth();
  }, [setUser, setAuth]);

  return (
    <>
      <AppRoutes />
    </>
  );
}

export default App;
