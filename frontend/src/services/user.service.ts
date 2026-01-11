import { api } from "@/lib/api";
import type {
  User,
  UpdateUserRequest,
  AddExperienceRequest,
  AddEducationRequest,
  Skill,
  AddProjectRequest,
} from "@/utils/types/user.type";
import { uploadService } from "./upload.service";

export const getUser = async (): Promise<User> => {
  const { data } = await api.get("/user");
  return data as User;
};

export const updateUser = async (input: UpdateUserRequest): Promise<User> => {
  const { data } = await api.put<User>(`/user`, input);
  return data as User;
};

export const uploadAvatar = async (file: File): Promise<User> => {
  const publicUrl = await uploadService.uploadImage(file);

  const { data } = await api.put<User>("/user/avatar", { avatarURL: publicUrl });
  
  return data as User;
};

export const deleteAvatar = async (): Promise<User> => {
  const { data } = await api.delete<User>(`/user/avatar`);
  return data as User;
};

export const uploadResume = async (file: File): Promise<User> => {
  const formData = new FormData();
  formData.append("resume", file);
  const { data } = await api.put<User>(`/user/resume`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  return data as User;
};

export const deleteResume = async (): Promise<User> => {
  const { data } = await api.delete<User>(`/user/resume`);
  return data as User;
};

// Candidate Features

export const getFullProfile = async (): Promise<User> => {
  const { data } = await api.get("/user/profile");
  return data as User;
};

// Experience
export const addExperience = async (
  input: AddExperienceRequest
): Promise<User> => {
  const { data } = await api.post("/user/experience", input);
  return data as User;
};

export const updateExperience = async (
  id: string,
  input: AddExperienceRequest
): Promise<User> => {
  const { data } = await api.put(`/user/experience/${id}`, input);
  return data as User;
};

export const deleteExperience = async (id: string): Promise<User> => {
  const { data } = await api.delete(`/user/experience/${id}`);
  return data as User;
};

// Education
export const addEducation = async (
  input: AddEducationRequest
): Promise<User> => {
  const { data } = await api.post("/user/education", input);
  return data as User;
};

export const updateEducation = async (
  id: string,
  input: AddEducationRequest
): Promise<User> => {
  const { data } = await api.put(`/user/education/${id}`, input);
  return data as User;
};

export const deleteEducation = async (id: string): Promise<User> => {
  const { data } = await api.delete(`/user/education/${id}`);
  return data as User;
};

// Skills
export const updateSkills = async (skills: Skill[]): Promise<User> => {
  const { data } = await api.put("/user/skills", { skills });
  return data as User;
};

// Projects
export const addProject = async (input: AddProjectRequest): Promise<User> => {
  const { data } = await api.post("/user/project", input);
  return data as User;
};

export const updateProject = async (
  id: string,
  input: AddProjectRequest
): Promise<User> => {
  const { data } = await api.put(`/user/project/${id}`, input);
  return data as User;
};

export const deleteProject = async (id: string): Promise<User> => {
  const { data } = await api.delete(`/user/project/${id}`);
  return data as User;
};
