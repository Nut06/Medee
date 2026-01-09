import { api } from "@/lib/api";
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
} from "@/utils/types/user.type";

export const register = async (
  input: RegisterRequest
): Promise<RegisterResponse> => {
  const { data } = await api.post<RegisterResponse>("/auth/register", input);
  return data as RegisterResponse;
};

export const loginLocal = async (
  input: LoginRequest
): Promise<LoginResponse> => {
  const { data } = await api.post<LoginResponse>("/auth/login", input);
  localStorage.setItem("accessToken", data.accessToken);
  return data as LoginResponse;
};

export const logout = async (): Promise<void> => {
  await api.post("/auth/logout");
};
