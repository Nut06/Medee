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
  company?: Company;
  companies?: Company[]; // List of companies the user belongs to
  projects?: Project[]; // List of projects the user belongs to
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
  user: {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
  };
  companies?: Company[];
}

export interface RegisterRequest {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  user: {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
  };
  companies?: Company[];
}

export interface UpdateUserRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  phoneNumber?: string;
  AvatarURL?: string;
  company?: Company;
  skills?: Skill[];
}

export interface UpdateUserResponse {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
}

export interface Skill {
  name?: string;
}

export interface PortfolioItem {
  id?: string;
  title: string;
  description: string;
  imageURL?: string;
  githubURL?: string;
  demoURL?: string;
}

export interface WorkExperience {
  id?: string;
  position: string;
  companyName: string;
  startDate?: string; // ISO Date string
  endDate?: string; // ISO Date string
  description?: string;
}

export interface Education {
  id?: string;
  instituteName: string; // Simplified for frontend, backend might need ID
  degree: string;
  fieldOfStudy: string; // Simplified
  graduationYear?: number;
}

// For API Requests
export interface AddExperienceRequest {
  position: string;
  companyName: string;
  startDate?: string;
  endDate?: string;
  description?: string;
}

export interface AddEducationRequest {
  instituteName: string;
  degree: string;
  fieldOfStudy: string;
  graduationYear?: number;
}

export interface AddProjectRequest {
  title: string;
  description: string;
  imageURL?: string;
  githubURL?: string;
  demoURL?: string;
}
