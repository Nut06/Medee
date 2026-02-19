import axios, { type AxiosResponse } from "axios";

export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:3000";

export const api = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

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
    Promise.reject(error);
  },
);

api.interceptors.response.use(
  (response) => response as AxiosResponse,
  async (error) => {
    const originalRequest = error.config; 

    if (
      error.response.status === 401 &&
      !originalRequest._retry &&
      !originalRequest.url.includes("/auth/refresh")
    ) {
      originalRequest._retry = true;
      try {
        const response = await api.post("/auth/refresh");
        const { accessToken } = response.data; // Fixed: response.data.token based on backend struct
        localStorage.setItem("accessToken", accessToken);
        originalRequest.headers.Authorization = `Bearer ${accessToken}`;
        return api(originalRequest);
      } catch (refreshError) {
        console.error("Token refresh failed:", refreshError);
        localStorage.removeItem("accessToken");
        window.location.href = "/auth"; // Redirect to login page
        return Promise.reject(refreshError);
      }
    }

    const msg = error.response?.data?.message || "Unexpected error";
    return Promise.reject(new Error(msg));
  },
);
