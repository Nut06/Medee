import { api } from "@/lib/api";
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
} from "@/utils/types/user.type";

export const register = async (
  input: RegisterRequest,
): Promise<RegisterResponse> => {
  try {
    const { data } = await api.post<LoginResponse>("/auth/register", input);
    // Store accessToken (refreshToken stays in HttpOnly cookie)
    return data as RegisterResponse;
  } catch (error) {
    console.error("Register error:", error);
    throw error; // Throw for caller to handle
  }
};

export const initializeCsrf = async (): Promise<void> => {
  try {
    await api.get("/auth/csrf-token");
  } catch (error) {
    console.error("CSRF initialization error:", error);
    throw error;
  }
};

export const loginLocal = async (
  input: LoginRequest,
): Promise<LoginResponse> => {
  try {
    const { data } = await api.post<LoginResponse>("/auth/login", input);
    // Store accessToken (refreshToken stays in HttpOnly cookie)
    return data as LoginResponse;
  } catch (error) {
    console.error("Login error:", error);
    throw error; // Throw for caller to handle
  }
};

export const logout = async (): Promise<void> => {
  try {
    await api.post("/auth/logout");
  } catch (error) {
    console.error("Logout request failed:", error);
    // Still clear auth even if logout endpoint fails
  }
};
