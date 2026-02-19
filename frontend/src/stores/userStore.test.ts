import { describe, it, expect, vi, beforeEach } from "vitest";
import { useUserStore } from "./userStore";
import type { User } from "@/utils/types/user.type";

// Mock the auth service
vi.mock("@/services/auth.service", () => ({
  logout: vi.fn().mockResolvedValue(undefined),
}));

// ============================================================
// userStore Tests
// ============================================================

describe("useUserStore", () => {
  beforeEach(() => {
    // Reset store state before each test
    useUserStore.setState({ user: {} as User, isAuth: false });
  });

  // ============================================================
  // setUser Tests
  // ============================================================

  describe("setUser", () => {
    it("should set user with full payload", () => {
      const user = {
        firstName: "John",
        lastName: "Doe",
        email: "john@example.com",
      } as User;

      useUserStore.getState().setUser(user);

      const state = useUserStore.getState();
      expect(state.user.firstName).toBe("John");
      expect(state.user.lastName).toBe("Doe");
      expect(state.isAuth).toBe(true);
    });

    it("should set user with key-value pair", () => {
      // First set a user
      useUserStore.getState().setUser({ firstName: "John" } as User);

      // Then update a specific key
      useUserStore.getState().setUser(null, "firstName", "Jane");

      expect(useUserStore.getState().user.firstName).toBe("Jane");
    });

    it("should clear user when payload is null", () => {
      // Set initial user
      useUserStore.getState().setUser({ firstName: "John" } as User);
      expect(useUserStore.getState().isAuth).toBe(true);

      // Clear
      useUserStore.getState().setUser(null);

      const state = useUserStore.getState();
      expect(state.isAuth).toBe(false);
    });

    it("should merge partial user data", () => {
      useUserStore
        .getState()
        .setUser({ firstName: "John", email: "john@test.com" } as User);
      useUserStore.getState().setUser({ lastName: "Doe" } as Partial<User>);

      const state = useUserStore.getState();
      expect(state.user.firstName).toBe("John");
      expect(state.user.lastName).toBe("Doe");
      expect(state.user.email).toBe("john@test.com");
    });
  });

  // ============================================================
  // setAuth Tests
  // ============================================================

  describe("setAuth", () => {
    it("should set isAuth to true", () => {
      useUserStore.getState().setAuth(true);
      expect(useUserStore.getState().isAuth).toBe(true);
    });

    it("should set isAuth to false", () => {
      useUserStore.getState().setAuth(true);
      useUserStore.getState().setAuth(false);
      expect(useUserStore.getState().isAuth).toBe(false);
    });
  });

  // ============================================================
  // clearAuth Tests
  // ============================================================

  describe("clearAuth", () => {
    it("should clear user and set isAuth to false", async () => {
      // Setup authenticated state
      useUserStore.getState().setUser({ firstName: "John" } as User);
      expect(useUserStore.getState().isAuth).toBe(true);

      // Clear auth
      await useUserStore.getState().clearAuth();

      const state = useUserStore.getState();
      expect(state.isAuth).toBe(false);
    });
  });
});
