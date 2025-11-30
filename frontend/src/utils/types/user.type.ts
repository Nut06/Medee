import type { Project } from "./project.type";

export type Role = "user" | "company";

export interface User {
  id?: string;
  firstName?: string;
  lastName?: string;
  email?: string;
  bio?: string;
  phoneNumber?: string;
  AvatarURL?: string;
  role?: Role;
  projects?: Project[];
  skills?: Skill[];
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface Company {
  id?: string;
  name?: string;
  description?: string;
}

export interface LoginResponse {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  role: string;
}

export interface RegisterRequest {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  role: string;
}

export interface RegisterResponse {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  role: string;
}

export interface UpdateUserRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  phoneNumber?: string;
  AvatarURL?: string;
  role?: Role;
  skills?: Skill[];
}

export interface UpdateUserResponse {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  role: string;
}

export interface Skill {
  name?: string;
}