import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useUserStore } from "@/stores/userStore";
import { deleteAvatar, updateUser, uploadAvatar } from "@/services/userService";
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
      const updatedUser = await updateUser(user.id, data);
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
      const updatedUser = await uploadAvatar(user.id, file);
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
      const updatedUser = await deleteAvatar(user.id);
      setUser(updatedUser); // Update store with new avatar URL
      toast.success("Avatar deleted successfully");
    } catch (error) {
      console.error(error);
      toast.error("Failed to delete avatar");
    } finally {
      setIsLoading(false);
    }
  };

  return {
    form,
    onSubmit,
    onUploadAvatar,
    onDeleteAvatar,
    isLoading,
    user,
  };
};
