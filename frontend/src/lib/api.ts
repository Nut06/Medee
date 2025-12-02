import axios from "axios";

export const API_URL = import.meta.env.VITE_API_URL || "http://localhost:5000";

export const api = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Prevent infinite loops
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Attempt to refresh token
        // We don't need to capture the response here because the cookie is set automatically by the browser
        // But if we want to update the user store with the new user info, we can.
        // For now, let's just ensure the refresh succeeds.
        await api.post("/auth/refresh");

        // Retry the original request
        return api(originalRequest);
      } catch (refreshError) {
        // Refresh failed (e.g., refresh token expired or invalid)
        // Clear auth state and redirect to login
        // We can't use hooks here, so we might need to access the store directly or dispatch an event.
        // For now, let's just let the error propagate or redirect.
        console.error("Token refresh failed:", refreshError);
        // window.location.href = '/login'; // Optional: Force redirect
        return Promise.reject(refreshError);
      }
    }

    const msg = error.response?.data?.message || "Unexpected error";
    return Promise.reject(new Error(msg));
  }
);
