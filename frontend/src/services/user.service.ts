import { api } from "@/lib/api";
import type {
  User,
  UpdateUserRequest,
  AddExperienceRequest,
  AddEducationRequest,
  Skill,
  AddProjectRequest,
} from "@/utils/types/user.type";

export const getUser = async (): Promise<User> => {
  const res = await api.get("/user");
  const data = (await res.data) as User;
  return data;
};

export const updateUser = async (input: UpdateUserRequest): Promise<User> => {
  const res = await api.put<User>(`/user`, input);
  const data = (await res.data) as User;
  return data;
};

export const uploadAvatar = async (file: File): Promise<User> => {
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

export const deleteAvatar = async (): Promise<User> => {
  const res = await api.delete<User>(`/user/avatar`);
  const data = (await res.data) as User;
  return data;
};

export const uploadResume = async (file: File): Promise<User> => {
  const formData = new FormData();
  formData.append("resume", file);
  const res = await api.put<User>(`/user/resume`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  const data = (await res.data) as User;
  return data;
};

export const deleteResume = async (): Promise<User> => {
  const res = await api.delete<User>(`/user/resume`);
  const data = (await res.data) as User;
  return data;
};

// Candidate Features

export const getFullProfile = async (): Promise<User> => {
  const res = await api.get("/user/profile");
  return res.data as User;
};

// Experience
export const addExperience = async (
  input: AddExperienceRequest
): Promise<User> => {
  const res = await api.post("/user/experience", input);
  return res.data as User;
};

export const updateExperience = async (
  id: string,
  input: AddExperienceRequest
): Promise<User> => {
  const res = await api.put(`/user/experience/${id}`, input);
  return res.data as User;
};

export const deleteExperience = async (id: string): Promise<User> => {
  const res = await api.delete(`/user/experience/${id}`);
  return res.data as User;
};

// Education
export const addEducation = async (
  input: AddEducationRequest
): Promise<User> => {
  const res = await api.post("/user/education", input);
  return res.data as User;
};

export const updateEducation = async (
  id: string,
  input: AddEducationRequest
): Promise<User> => {
  const res = await api.put(`/user/education/${id}`, input);
  return res.data as User;
};

export const deleteEducation = async (id: string): Promise<User> => {
  const res = await api.delete(`/user/education/${id}`);
  return res.data as User;
};

// Skills
export const updateSkills = async (skills: Skill[]): Promise<User> => {
  const { data } = await api.put("/user/skills", { skills });
  return data as User;
};

// Projects
export const addProject = async (input: AddProjectRequest): Promise<User> => {
  const res = await api.post("/user/project", input);
  return res.data as User;
};

export const updateProject = async (
  id: string,
  input: AddProjectRequest
): Promise<User> => {
  const res = await api.put(`/user/project/${id}`, input);
  return res.data as User;
};

export const deleteProject = async (id: string): Promise<User> => {
  const res = await api.delete(`/user/project/${id}`);
  return res.data as User;
};
