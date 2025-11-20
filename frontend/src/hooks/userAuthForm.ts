import { loginLocal, register } from "@/services/authService";
import { useUserStore } from "@/stores/userStore";
import type { LoginRequest, RegisterRequest } from "@/utils/types/user.type";
import { useState } from "react";

export function useAuthForm() {
  const [name, setName] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const { setUser } = useUserStore();

  const registerUser = async (input: RegisterRequest): Promise<void> => {
    try {
      setLoading(true);
      const data = await register(input);
      setUser(data);
    } catch (error: unknown) {
      if (error instanceof Error) {
        setError(error.message);
      } else {
        setError("An unknown error occurred");
      }
    } finally {
      setLoading(!loading);
    }
  };

  const loginUser = async (input: LoginRequest): Promise<void> => {
    try {
      setLoading(true);
      const data = await loginLocal(input);
      setUser(data);
    } catch (error: unknown) {
      if (error instanceof Error) {
        setError(error.message);
      } else {
        setError("An unknown error occurred");
      }
    } finally {
      setLoading(!loading);
    }
  };
  return {
    name,
    loginUser,
    registerUser,
    setName,
    password,
    setPassword,
    error,
    setError,
  };
}
