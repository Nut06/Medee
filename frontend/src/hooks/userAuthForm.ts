import { loginLocal, register } from "@/services/auth.service";
import { useUserStore } from "@/stores/userStore";

import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  User,
} from "@/utils/types/user.type";
import { useReducer } from "react";
import axios from "axios";

type AuthFormState = {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  error: string;
  loading: boolean;
};

type Action =
  | { type: "SET_FIELD"; key: keyof AuthFormState; val: string | boolean }
  | { type: "RESET" };

const initialState: AuthFormState = {
  firstName: "",
  lastName: "",
  email: "",
  password: "",
  error: "",
  loading: false,
};

function reducer(state: AuthFormState, action: Action): AuthFormState {
  switch (action.type) {
    case "SET_FIELD":
      return { ...state, [action.key]: action.val };
    case "RESET":
      return initialState;
    default:
      return state;
  }
}

export function useAuthForm() {
  const [state, dispatch] = useReducer(reducer, initialState);
  const { setUser } = useUserStore();

  const setField = (key: keyof AuthFormState, val: string | boolean) => {
    dispatch({ type: "SET_FIELD", key, val });
  };

  const resetForm = () => dispatch({ type: "RESET" });

  const registerUser = async (input: RegisterRequest): Promise<void> => {
    try {
      setField("loading", true);
      setField("error", "");

      const response = await register(input);
      if (response.user) {
        setUser(response.user);
      }
    } catch (error: unknown) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.message || "Registration failed";
        setField("error", message);
      } else if (error instanceof Error) {
        setField("error", error.message);
      } else {
        setField("error", "An unknown error occurred");
      }
    } finally {
      setField("loading", false);
    }
  };

  const loginUser = async (input: LoginRequest): Promise<boolean> => {
    try {
      setField("loading", true);
      setField("error", "");

      const loginRes: LoginResponse = await loginLocal(input);
      const { user: rawUser, companies = [] } = loginRes;
      const user = rawUser as User;
      if (companies != null) {
        user.companies = companies;
      }

      setUser(user);
      resetForm();
      return true; // ← บอกว่า success
    } catch (error: unknown) {
      setField("error", "Invalid username or password");
      return false; // ← บอกว่า fail
    } finally {
      setField("loading", false);
    }
  };

  return {
    ...state,
    setField,
    resetForm,
    loginUser,
    registerUser,
  };
}
