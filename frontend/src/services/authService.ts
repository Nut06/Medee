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
  const res = await api.post<RegisterResponse>("/auth/login", input);
  const data = (await res.data) as RegisterResponse;
  return data;
};

export const loginLocal = async (
  input: LoginRequest
): Promise<LoginResponse> => {
  const res = await api.post<LoginResponse>("/auth/login", input);
  const data = (await res.data) as LoginResponse;
  return data;
};
