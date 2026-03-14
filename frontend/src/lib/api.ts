import axios, {
  AxiosError,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";
import { useUserStore } from "@/stores/userStore";

export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:3000";

interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

// Instance สำหรับ request ธรรมดา (มี interceptor)
export const api = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

const apiWithoutInterceptors = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

let isRefreshing = false;
api.interceptors.request.use(
  (request) => {
    const accessToken = localStorage.getItem("accessToken");
    if (accessToken) {
      request.headers.Authorization = `Bearer ${accessToken}`;
    }
    return request;
  },
  (error) => {
    console.error("Request interceptor error:", error);
    return Promise.reject(error);
  },
);

api.interceptors.response.use(
  (response) => response as AxiosResponse,
  async (error: AxiosError) => {
    const originalRequest: CustomAxiosRequestConfig | undefined = error.config;

    if (error.response?.status !== 401 || !originalRequest) {
      return Promise.reject(error);
    }

    if (isRefreshing) {
      return new Promise((token) => {
        originalRequest.headers.Authorization = `Bearer ${token}`;
        return api(originalRequest);
      }).catch((err) => {
        return Promise.reject(err);
      });
    }

    if (originalRequest._retry) {
      useUserStore.getState().clearAuth();
      return Promise.reject(error); // ← ไม่ hard redirect
    }

    // เริ่มการ refresh token
    isRefreshing = true;
    originalRequest._retry = true;

    try {
      const { data } = await apiWithoutInterceptors.post(
        "/auth/refresh",
        {},
        { withCredentials: true },
      );

      const { Token: accessToken } = data;

      useUserStore.getState().setAccessToken(accessToken);

      originalRequest.headers.Authorization = `Bearer ${accessToken}`;

      return api(originalRequest);
    } catch (refreshError) {
      console.log("Token refresh failed:", refreshError);

      useUserStore.getState().clearAuth();
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false; 
    }
  },
);
