import { api } from "@/lib/api";
import type { User } from "@/utils/types/user.type";
import type { UpdateUserRequest, UpdateUserResponse } from "@/utils/types/user.type";

export const getUser = async () => {
  const res = await api.get("/user");
  const data = (await res.data) as User;
  return data;
};

export const updateUser = async (id: string, input: UpdateUserRequest) => {
  const res = await api.put(`/user/${id}`, input);
  const data = (await res.data) as UpdateUserResponse;
  return data;
};
