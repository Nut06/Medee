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
  getFullProfile,
} from "@/services/userService";
import { toast } from "sonner";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import type { Skill } from "@/utils/types/user.type";

const profileSchema = z.object({
  firstName: z
    .string()
    .min(2, "First name must be at least 2 characters")
    .optional(),
  lastName: z
    .string()
    .min(2, "Last name must be at least 2 characters")
    .optional(),
  email: z.email("Invalid email address").optional(),
  bio: z.string().optional(),
});



export type ProfileFormValues = z.infer<typeof profileSchema>;

export const useProfile = () => {
  const { user, setUser } = useUserStore();
  const queryClient = useQueryClient();

  // Fetch Full Profile
  const { data: profileData, isLoading: isProfileLoading } = useQuery({
    queryKey: ["profile"],
    queryFn: getFullProfile,
    enabled: !!user?.id, // Only fetch if user is logged in
  });

  // Sync React Query data with Zustand Store
  useEffect(() => {
    if (profileData) {
      setUser(profileData);
    }
  }, [profileData, setUser]);

  const form = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      firstName: user?.firstName || "",
      lastName: user?.lastName || "",
      email: user?.email || "",
      bio: user?.bio || "",
    },
    values: {
      // Update form when user data changes
      firstName: user?.firstName || "",
      lastName: user?.lastName || "",
      email: user?.email || "",
      bio: user?.bio || "",
    },
  });

  // Helper to invalidate queries and update store
  const onSuccess = (data: any, message: string) => {
    queryClient.invalidateQueries({ queryKey: ["profile"] });
    setUser(data);
    toast.success(message);
  };
  
  const onError = (error: any, message: string) => {
    console.error(error);
    toast.error(message);
  };

  // Mutations
  const updateProfileMutation = useMutation({
    mutationFn: updateUser,
    onSuccess: (data) => onSuccess(data, "Profile updated successfully"),
    onError: (error) => onError(error, "Failed to update profile"),
  });

  const uploadAvatarMutation = useMutation({
    mutationFn: uploadAvatar,
    onSuccess: (data) => onSuccess(data, "Avatar updated successfully"),
    onError: (error) => onError(error, "Failed to upload avatar"),
  });

  const deleteAvatarMutation = useMutation({
    mutationFn: deleteAvatar,
    onSuccess: (data) => onSuccess(data, "Avatar deleted successfully"),
    onError: (error) => onError(error, "Failed to delete avatar"),
  });

  // Experience Mutations
  const addExperienceMutation = useMutation({
    mutationFn: addExperience,
    onSuccess: (data) => onSuccess(data, "Experience added"),
    onError: (error) => onError(error, "Failed to add experience"),
  });

  const updateExperienceMutation = useMutation({
    mutationFn: (vars: { id: string; data: any }) =>
      updateExperience(vars.id, vars.data),
    onSuccess: (data) => onSuccess(data, "Experience updated"),
    onError: (error) => onError(error, "Failed to update experience"),
  });

  const deleteExperienceMutation = useMutation({
    mutationFn: deleteExperience,
    onSuccess: (data) => onSuccess(data, "Experience deleted"),
    onError: (error) => onError(error, "Failed to delete experience"),
  });

  // Education Mutations
  const addEducationMutation = useMutation({
    mutationFn: addEducation,
    onSuccess: (data) => onSuccess(data, "Education added"),
    onError: (error) => onError(error, "Failed to add education"),
  });

  const updateEducationMutation = useMutation({
    mutationFn: (vars: { id: string; data: any }) =>
      updateEducation(vars.id, vars.data),
    onSuccess: (data) => onSuccess(data, "Education updated"),
    onError: (error) => onError(error, "Failed to update education"),
  });

  const deleteEducationMutation = useMutation({
    mutationFn: deleteEducation,
    onSuccess: (data) => onSuccess(data, "Education deleted"),
    onError: (error) => onError(error, "Failed to delete education"),
  });

  // Skills Mutations
  const updateSkillsMutation = useMutation({
    mutationFn: updateSkills,
    onSuccess: (data) => onSuccess(data, "Skills updated"),
    onError: (error) => onError(error, "Failed to update skills"),
  });

  // Projects Mutations
  const addProjectMutation = useMutation({
    mutationFn: addProject,
    onSuccess: (data) => onSuccess(data, "Project added"),
    onError: (error) => onError(error, "Failed to add project"),
  });

  const updateProjectMutation = useMutation({
    mutationFn: (vars: { id: string; data: any }) =>
      updateProject(vars.id, vars.data),
    onSuccess: (data) => onSuccess(data, "Project updated"),
    onError: (error) => onError(error, "Failed to update project"),
  });

  const deleteProjectMutation = useMutation({
    mutationFn: deleteProject,
    onSuccess: (data) => onSuccess(data, "Project deleted"),
    onError: (error) => onError(error, "Failed to delete project"),
  });

  // Handlers (Wrappers to match previous API)
  const onSubmit = (data: ProfileFormValues) => {
    updateProfileMutation.mutate(data);
    
    setUser(data);
  }
  const onUploadAvatar = (file: File) => uploadAvatarMutation.mutate(file);
  const onDeleteAvatar = () => deleteAvatarMutation.mutate();

  const onAddExperience = (data: any) => addExperienceMutation.mutate(data);
  const onUpdateExperience = (id: string, data: any) =>
    updateExperienceMutation.mutate({ id, data });
  const onDeleteExperience = (id: string) =>
    deleteExperienceMutation.mutate(id);

  const onAddEducation = (data: any) => addEducationMutation.mutate(data);
  const onUpdateEducation = (id: string, data: any) =>
    updateEducationMutation.mutate({ id, data });
  const onDeleteEducation = (id: string) => deleteEducationMutation.mutate(id);

  const onUpdateSkills = (skills: any[]) => updateSkillsMutation.mutate(skills);

  const onAddProject = (data: any) => addProjectMutation.mutate(data);
  const onUpdateProject = (id: string, data: any) =>
    updateProjectMutation.mutate({ id, data });
  const onDeleteProject = (id: string) => deleteProjectMutation.mutate(id);

  const isLoading =
    isProfileLoading ||
    updateProfileMutation.isPending ||
    uploadAvatarMutation.isPending ||
    deleteAvatarMutation.isPending ||
    addExperienceMutation.isPending ||
    updateExperienceMutation.isPending ||
    deleteExperienceMutation.isPending ||
    addEducationMutation.isPending ||
    updateEducationMutation.isPending ||
    deleteEducationMutation.isPending ||
    updateSkillsMutation.isPending ||
    addProjectMutation.isPending ||
    updateProjectMutation.isPending ||
    deleteProjectMutation.isPending;

  const isSaving =
    updateProfileMutation.isPending ||
    uploadAvatarMutation.isPending ||
    deleteAvatarMutation.isPending ||
    addExperienceMutation.isPending ||
    updateExperienceMutation.isPending ||
    deleteExperienceMutation.isPending ||
    addEducationMutation.isPending ||
    updateEducationMutation.isPending ||
    deleteEducationMutation.isPending ||
    updateSkillsMutation.isPending ||
    addProjectMutation.isPending ||
    updateProjectMutation.isPending ||
    deleteProjectMutation.isPending;

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
    isLoading, // Kept for backward compatibility if needed elsewhere
    isSaving, // New specific state for save buttons
    user,
  };
};
