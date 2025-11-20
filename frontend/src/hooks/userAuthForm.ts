import { loginLocal, register } from "@/services/authService";
import { useUserStore } from "@/stores/userStore";
import type { LoginRequest, RegisterRequest } from "@/utils/types/user.type";
import { useReducer } from "react";

type AuthFormState = {
  name: string;
  email: string;
  password: string;
  error: string;
  loading: boolean;
};

type Action =
  | { type: "SET_FIELD"; key: keyof AuthFormState; val: string | boolean }
  | { type: "RESET" };

const initialState: AuthFormState = {
  name: "",
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
      const data = await register(input);
      setUser(data);
    } catch (error: unknown) {
      if (error instanceof Error) {
        setField("error", error.message);
      } else {
        setField("error", "An unknown error occurred");
      }
    } finally {
      setField("loading", true);
      resetForm();
    }
  };

  const loginUser = async (input: LoginRequest): Promise<void> => {
    try {
      setField("loading", true);
      const data = await loginLocal(input);
      setUser(data);
    } catch (error: unknown) {
      if (error instanceof Error) {
        setField("error", error.message);
      } else {
        setField("error", "An unknown error occurred");
      }
    } finally {
      setField("loading", false);
      resetForm();
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
