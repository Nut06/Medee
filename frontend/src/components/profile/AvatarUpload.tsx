import { useState, useRef } from "react";
import { Camera, Trash2 } from "lucide-react";
import { useProfile } from "@/hooks/useProfile";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface AvatarUploadProps {
  avatarURL?: string;
  firstName?: string;
  lastName?: string;
  className?: string;
}

export default function AvatarUpload({
  avatarURL,
  firstName,
  lastName,
  className = "",
}: AvatarUploadProps) {
  const { onUploadAvatar, onDeleteAvatar, isSaving } = useProfile();
  const [isHovered, setIsHovered] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const initials = `${firstName?.[0] || ""}${lastName?.[0] || ""}`;

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      // Validate file type
      if (!file.type.startsWith("image/")) {
        alert("Please select an image file");
        return;
      }
      // Validate file size (max 5MB)
      if (file.size > 5 * 1024 * 1024) {
        alert("File size must be less than 5MB");
        return;
      }
      await onUploadAvatar(file);
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const handleDelete = async () => {
    await onDeleteAvatar();
    setShowDeleteDialog(false);
  };

  return (
    <>
      <div
        className={`relative ${className}`}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <div className="h-24 w-24 rounded-full bg-muted flex items-center justify-center text-2xl font-bold overflow-hidden relative group">
          {avatarURL ? (
            <img
              src={avatarURL}
              alt="Avatar"
              className="h-full w-full object-cover"
            />
          ) : (
            <>{initials}</>
          )}

          {/* Overlay on hover */}
          <div
            className={`absolute inset-0 bg-black/50 flex items-center justify-center gap-2 transition-opacity ${
              isHovered ? "opacity-100" : "opacity-0"
            }`}
          >
            <button
              onClick={handleUploadClick}
              disabled={isSaving}
              className="p-2 bg-white/90 rounded-full hover:bg-white transition-colors"
              aria-label="Upload avatar"
            >
              <Camera className="h-4 w-4 text-gray-700" />
            </button>
            {avatarURL && (
              <button
                onClick={() => setShowDeleteDialog(true)}
                disabled={isSaving}
                className="p-2 bg-white/90 rounded-full hover:bg-white transition-colors"
                aria-label="Delete avatar"
              >
                <Trash2 className="h-4 w-4 text-red-600" />
              </button>
            )}
          </div>
        </div>

        {/* Camera icon indicator */}
        <div className="absolute bottom-0 right-0 bg-primary rounded-full p-1.5 border-2 border-background">
          <Camera className="h-3 w-3 text-primary-foreground" />
        </div>

        {/* Hidden file input */}
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileChange}
          className="hidden"
        />
      </div>

      {/* Delete confirmation dialog */}
      <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Avatar</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete your profile picture? This action
              cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setShowDeleteDialog(false)}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={isSaving}
            >
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
