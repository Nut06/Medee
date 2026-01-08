import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useUserStore } from "@/stores/userStore";
import {
  deleteAvatar,
  updateUser,
  uploadAvatar,
  uploadResume,
  deleteResume,
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
} from "@/services/user.service";
import { toast } from "sonner";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";

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
      firstName: user.firstName || "",
      lastName: user.lastName || "",
      email: user.email || "",
      bio: user.bio || "",
    },
    values: {
      // Update form when user data changes
      firstName: user.firstName || "",
      lastName: user.lastName || "",
      email: user.email || "",
      bio: user.bio || "",
    },
  });

  // Helper to invalidate queries and refetch
  const onSuccess = async (message: string) => {
    queryClient.invalidateQueries({ queryKey: ["profile"] });
    toast.success(message);
  };

  const onError = (error: any, message: string) => {
    console.error(error);
    toast.error(message);
  };

  // ========== Factory Functions for DRY Optimistic Updates ==========
  type ProfileField = "experiences" | "educations" | "skills" | "projects";

  // Factory for DELETE mutations
  const createDeleteMutation = <T>(
    mutationFn: (id: string) => Promise<T>,
    field: ProfileField,
    successMsg: string,
    errorMsg: string
  ) =>
    useMutation({
      mutationFn,
      onMutate: async (id: string) => {
        await queryClient.cancelQueries({ queryKey: ["profile"] });
        const previousProfile = queryClient.getQueryData(["profile"]);

        const updater = (old: any) => ({
          ...old,
          [field]: old?.[field]?.filter((item: any) => item.id !== id) || [],
        });

        queryClient.setQueryData(["profile"], updater);
        setUser(updater(user));

        return { previousProfile };
      },
      onError: (error, _id, context: any) => {
        if (context?.previousProfile) {
          queryClient.setQueryData(["profile"], context.previousProfile);
          setUser(context.previousProfile);
        }
        onError(error, errorMsg);
      },
      onSuccess: () => onSuccess(successMsg),
    });

  // Factory for ADD mutations
  const createAddMutation = <TInput, TOutput>(
    mutationFn: (data: TInput) => Promise<TOutput>,
    field: ProfileField,
    successMsg: string,
    errorMsg: string
  ) =>
    useMutation({
      mutationFn,
      onMutate: async (newItem: TInput) => {
        await queryClient.cancelQueries({ queryKey: ["profile"] });
        const previousProfile = queryClient.getQueryData(["profile"]);

        const tempItem = { ...newItem, id: `temp-${Date.now()}` };

        const updater = (old: any) => ({
          ...old,
          [field]: [...(old?.[field] || []), tempItem],
        });

        queryClient.setQueryData(["profile"], updater);
        setUser(updater(user));

        return { previousProfile };
      },
      onError: (error, _item, context: any) => {
        if (context?.previousProfile) {
          queryClient.setQueryData(["profile"], context.previousProfile);
          setUser(context.previousProfile);
        }
        onError(error, errorMsg);
      },
      onSuccess: () => onSuccess(successMsg),
    });

  // Factory for UPDATE mutations
  const createUpdateMutation = <TInput, TOutput>(
    mutationFn: (vars: { id: string; data: TInput }) => Promise<TOutput>,
    field: ProfileField,
    successMsg: string,
    errorMsg: string
  ) =>
    useMutation({
      mutationFn,
      onMutate: async (vars: { id: string; data: TInput }) => {
        await queryClient.cancelQueries({ queryKey: ["profile"] });
        const previousProfile = queryClient.getQueryData(["profile"]);

        const updater = (old: any) => ({
          ...old,
          [field]:
            old?.[field]?.map((item: any) =>
              item.id === vars.id ? { ...item, ...vars.data } : item
            ) || [],
        });

        queryClient.setQueryData(["profile"], updater);
        setUser(updater(user));

        return { previousProfile };
      },
      onError: (error, _vars, context: any) => {
        if (context?.previousProfile) {
          queryClient.setQueryData(["profile"], context.previousProfile);
          setUser(context.previousProfile);
        }
        onError(error, errorMsg);
      },
      onSuccess: () => onSuccess(successMsg),
    });
  // ===================================================================

  // Mutations
  const updateProfileMutation = useMutation({
    mutationFn: updateUser,
    onSuccess: () => onSuccess("Profile updated successfully"),
    onError: (error) => onError(error, "Failed to update profile"),
  });

  const uploadAvatarMutation = useMutation({
    mutationFn: uploadAvatar,
    onSuccess: () => onSuccess("Avatar updated successfully"),
    onError: (error) => onError(error, "Failed to upload avatar"),
  });

  const deleteAvatarMutation = useMutation({
    mutationFn: deleteAvatar,
    onSuccess: () => onSuccess("Avatar deleted successfully"),
    onError: (error) => onError(error, "Failed to delete avatar"),
  });

  const uploadResumeMutation = useMutation({
    mutationFn: uploadResume,
    onSuccess: () => onSuccess("Resume uploaded successfully"),
    onError: (error) => onError(error, "Failed to upload resume"),
  });

  const deleteResumeMutation = useMutation({
    mutationFn: deleteResume,
    onSuccess: () => onSuccess("Resume deleted successfully"),
    onError: (error) => onError(error, "Failed to delete resume"),
  });

  // Experience Mutations
  const addExperienceMutation = createAddMutation(
    addExperience,
    "experiences",
    "Experience added",
    "Failed to add experience"
  );
  const updateExperienceMutation = createUpdateMutation(
    (vars: { id: string; data: any }) => updateExperience(vars.id, vars.data),
    "experiences",
    "Experience updated",
    "Failed to update experience"
  );
  const deleteExperienceMutation = createDeleteMutation(
    deleteExperience,
    "experiences",
    "Experience deleted",
    "Failed to delete experience"
  );

  // Education Mutations
  const addEducationMutation = createAddMutation(
    addEducation,
    "educations",
    "Education added",
    "Failed to add education"
  );
  const updateEducationMutation = createUpdateMutation(
    (vars: { id: string; data: any }) => updateEducation(vars.id, vars.data),
    "educations",
    "Education updated",
    "Failed to update education"
  );
  const deleteEducationMutation = createDeleteMutation(
    deleteEducation,
    "educations",
    "Education deleted",
    "Failed to delete education"
  );

  // Skills Mutations (replace entire array)
  const updateSkillsMutation = useMutation({
    mutationFn: updateSkills,
    onMutate: async (newSkills) => {
      await queryClient.cancelQueries({ queryKey: ["profile"] });
      const previousProfile = queryClient.getQueryData(["profile"]);

      const updater = (old: any) => ({ ...old, skills: newSkills });
      queryClient.setQueryData(["profile"], updater);
      setUser(updater(user));

      return { previousProfile };
    },
    onError: (error, _newSkills, context: any) => {
      if (context?.previousProfile) {
        queryClient.setQueryData(["profile"], context.previousProfile);
        setUser(context.previousProfile);
      }
      onError(error, "Failed to update skills");
    },
    onSuccess: () => onSuccess("Skills updated"),
  });

  // Projects Mutations
  const addProjectMutation = createAddMutation(
    addProject,
    "projects",
    "Project added",
    "Failed to add project"
  );
  const updateProjectMutation = createUpdateMutation(
    (vars: { id: string; data: any }) => updateProject(vars.id, vars.data),
    "projects",
    "Project updated",
    "Failed to update project"
  );
  const deleteProjectMutation = createDeleteMutation(
    deleteProject,
    "projects",
    "Project deleted",
    "Failed to delete project"
  );

  // Handlers (Wrappers to match previous API)
  const onSubmit = (data: ProfileFormValues) => {
    updateProfileMutation.mutate(data);
    setUser(data);
  };

  const onUploadAvatar = (file: File) => uploadAvatarMutation.mutate(file);
  const onDeleteAvatar = () => deleteAvatarMutation.mutate();

  const onUploadResume = (file: File) => uploadResumeMutation.mutate(file);
  const onDeleteResume = () => deleteResumeMutation.mutate();

  const onAddExperience = async (data: any) => {
    await addExperienceMutation.mutateAsync(data);
  };
  const onUpdateExperience = async (id: string, data: any) => {
    await updateExperienceMutation.mutateAsync({ id, data });
  };

  const onDeleteExperience = async (id: string) => {
    await deleteExperienceMutation.mutateAsync(id);
  };

  const onAddEducation = async (data: any) => {
    await addEducationMutation.mutateAsync(data);
  };
  const onUpdateEducation = async (id: string, data: any) => {
    await updateEducationMutation.mutateAsync({ id, data });
  };
  const onDeleteEducation = async (id: string) => {
    await deleteEducationMutation.mutateAsync(id);
  };

  const onUpdateSkills = async (skills: any[]) => {
    await updateSkillsMutation.mutateAsync(skills);
  };

  const onAddProject = async (data: any) => {
    await addProjectMutation.mutateAsync(data);
  };

  const onUpdateProject = async (id: string, data: any) => {
    await updateProjectMutation.mutateAsync({ id, data });
  };

  const onDeleteProject = async (id: string) => {
    await deleteProjectMutation.mutateAsync(id);
  };

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
    updateProjectMutation.isPending ||
    deleteProjectMutation.isPending ||
    uploadResumeMutation.isPending ||
    deleteResumeMutation.isPending;

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
    onUploadResume,
    onDeleteResume,
    isLoading, // Kept for backward compatibility if needed elsewhere
    isSaving, // New specific state for save buttons
    user,
  };
};
