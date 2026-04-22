import type { User } from "@/utils/types/user.type";
import { create } from "zustand";

interface UserState {
  user: User;
  setUser: (
    user: User | Partial<User> | null,
    key?: keyof User,
    value?: unknown,
  ) => void;
  isAuth: boolean;
  setAuth: (b: boolean) => void;
  clearAuth: () => Promise<void>;
}

export const useUserStore = create<UserState>((set) => ({
  user: {},

  setUser: (payload, key, value) => {
    if (key) {
      set((state) =>
        state.user
          ? { user: { ...state.user, [key]: value as User[keyof User] } }
          : state,
      );
      return;
    }

    if (!payload) {
      set({ user: {}, isAuth: false });
      return;
    }

    set((state) => ({
      user: { ...(state.user ?? {}), ...(payload as Partial<User>) },
      isAuth: true,
    }));
  },

  isAuth: false,
  setAuth: (b: boolean) => {
    set({
      isAuth: b,
    });
  },
  clearAuth: async () => {
    localStorage.removeItem("accessToken");
    set({
      user: {},
      isAuth: false,
    });
  },
}));
