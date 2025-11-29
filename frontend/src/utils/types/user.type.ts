import type { Project } from "./project.type";

export type Role = "candidate" | "company";

export interface User {
  id?: string;
  name?: string;
  email?: string;
  phoneNumber?: string;
  AvatarURL?: string;
  role?: Role;
  projects?: Project[];
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  id: string;
  name: string;
  email: string;
  role: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  role: string;
}

export interface RegisterResponse {
  id: string;
  name: string;
  email: string;
  role: string;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  phoneNumber?: string;
  AvatarURL?: string;
  role?: Role;
}

export interface UpdateUserResponse {
  id: string;
  name: string;
  email: string;
  role: string;
}