import { api } from "@/lib/api";
import type { User } from "@/utils/types/user.type";
import type {
  UpdateUserRequest,
  UpdateUserResponse,
} from "@/utils/types/user.type";

export const getUser = async ():Promise<User> => {
  const res = await api.get("/user");
  const data = (await res.data) as User;
  return data;
};

export const updateUser = async (input: UpdateUserRequest):Promise<User> => {
  const res = await api.put<User>(`/user`, input);
  const data = (await res.data) as User;
  return data;
};

export const uploadAvatar = async (file: File):Promise<User> => {
  const formData = new FormData();
  formData.append("avatar", file);
  const res = await api.put<User>(`/user/avatar`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  const data = (await res.data) as User;
  return data;
};

export const deleteAvatar = async ():Promise<User> => {
  const res = await api.delete<User>(`/user/avatar`);
  const data = (await res.data) as User;
  return data;
};