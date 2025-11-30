import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useUserStore } from "@/stores/userStore";
import {
  deleteAvatar,
  updateUser,
  uploadAvatar,
  addExperience,
  updateExperience,
  deleteExperience,
  addEducation,
  updateEducation,
  deleteEducation,
  updateSkills,
  addProject,
  updateProject,
  deleteProject,
} from "@/services/userService";
import { toast } from "sonner";

const profileSchema = z.object({
  firstName: z.string().min(2, "First name must be at least 2 characters"),
  lastName: z.string().min(2, "Last name must be at least 2 characters"),
  email: z.string().email("Invalid email address"),
  bio: z.string().optional(),
});

export type ProfileFormValues = z.infer<typeof profileSchema>;

export const useProfile = () => {
  const { user, setUser } = useUserStore();
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      firstName: user?.firstName || "",
      lastName: user?.lastName || "",
      email: user?.email || "",
      bio: user?.bio || "",
    },
  });

  const onSubmit = async (data: ProfileFormValues) => {
    if (!user?.id) return;
    setIsLoading(true);
    try {
      const updatedUser = await updateUser(data);
      setUser(updatedUser);
      toast.success("Profile updated successfully");
    } catch (error) {
      console.error(error);
      toast.error("Failed to update profile");
    } finally {
      setIsLoading(false);
    }
  };

  const onUploadAvatar = async (file: File) => {
    if (!user?.id) return;
    setIsLoading(true);
    try {
      const updatedUser = await uploadAvatar(file);
      setUser(updatedUser); // Update store with new avatar URL
      toast.success("Avatar updated successfully");
    } catch (error) {
      console.error(error);
      toast.error("Failed to upload avatar");
    } finally {
      setIsLoading(false);
    }
  };
  const onDeleteAvatar = async () => {
    if (!user?.id) return;
    setIsLoading(true);
    try {
      const updatedUser = await deleteAvatar();
      setUser(updatedUser); // Update store with new avatar URL
      toast.success("Avatar deleted successfully");
    } catch (error) {
      console.error(error);
      toast.error("Failed to delete avatar");
    } finally {
      setIsLoading(false);
    }
  };

  // Candidate Features Handlers

  const onAddExperience = async (data: any) => {
    setIsLoading(true);
    try {
      const updatedUser = await addExperience(data);
      setUser(updatedUser);
      toast.success("Experience added");
    } catch (error) {
      toast.error("Failed to add experience");
    } finally {
      setIsLoading(false);
    }
  };

  const onUpdateExperience = async (id: string, data: any) => {
    setIsLoading(true);
    try {
      const updatedUser = await updateExperience(id, data);
      setUser(updatedUser);
      toast.success("Experience updated");
    } catch (error) {
      toast.error("Failed to update experience");
    } finally {
      setIsLoading(false);
    }
  };

  const onDeleteExperience = async (id: string) => {
    setIsLoading(true);
    try {
      const updatedUser = await deleteExperience(id);
      setUser(updatedUser);
      toast.success("Experience deleted");
    } catch (error) {
      toast.error("Failed to delete experience");
    } finally {
      setIsLoading(false);
    }
  };

  const onAddEducation = async (data: any) => {
    setIsLoading(true);
    try {
      const updatedUser = await addEducation(data);
      setUser(updatedUser);
      toast.success("Education added");
    } catch (error) {
      toast.error("Failed to add education");
    } finally {
      setIsLoading(false);
    }
  };

  const onUpdateEducation = async (id: string, data: any) => {
    setIsLoading(true);
    try {
      const updatedUser = await updateEducation(id, data);
      setUser(updatedUser);
      toast.success("Education updated");
    } catch (error) {
      toast.error("Failed to update education");
    } finally {
      setIsLoading(false);
    }
  };

  const onDeleteEducation = async (id: string) => {
    setIsLoading(true);
    try {
      const updatedUser = await deleteEducation(id);
      setUser(updatedUser);
      toast.success("Education deleted");
    } catch (error) {
      toast.error("Failed to delete education");
    } finally {
      setIsLoading(false);
    }
  };

  const onUpdateSkills = async (skills: any[]) => {
    setIsLoading(true);
    try {
      const updatedUser = await updateSkills(skills);
      setUser(updatedUser);
      toast.success("Skills updated");
    } catch (error) {
      toast.error("Failed to update skills");
    } finally {
      setIsLoading(false);
    }
  };

  const onAddProject = async (data: any) => {
    setIsLoading(true);
    try {
      const updatedUser = await addProject(data);
      setUser(updatedUser);
      toast.success("Project added");
    } catch (error) {
      toast.error("Failed to add project");
    } finally {
      setIsLoading(false);
    }
  };

  const onUpdateProject = async (id: string, data: any) => {
    setIsLoading(true);
    try {
      const updatedUser = await updateProject(id, data);
      setUser(updatedUser);
      toast.success("Project updated");
    } catch (error) {
      toast.error("Failed to update project");
    } finally {
      setIsLoading(false);
    }
  };

  const onDeleteProject = async (id: string) => {
    setIsLoading(true);
    try {
      const updatedUser = await deleteProject(id);
      setUser(updatedUser);
      toast.success("Project deleted");
    } catch (error) {
      toast.error("Failed to delete project");
    } finally {
      setIsLoading(false);
    }
  };

  return {
    form,
    onSubmit,
    onUploadAvatar,
    onDeleteAvatar,
    onAddExperience,
    onUpdateExperience,
    onDeleteExperience,
    onAddEducation,
    onUpdateEducation,
    onDeleteEducation,
    onUpdateSkills,
    onAddProject,
    onUpdateProject,
    onDeleteProject,
    isLoading,
    user,
  };
};
