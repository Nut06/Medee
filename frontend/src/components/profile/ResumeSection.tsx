import {
  Upload,
  FileText,
  X,
  Download,
  Trash2,
  CheckCircle,
} from "lucide-react";
import { useState, useRef, type ChangeEvent } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useProfile } from "@/hooks/useProfile";

// Hook for resume logic
function useResumeUpload() {
  const { user, onUploadResume, onDeleteResume, isSaving } = useProfile();
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      // Check file type
      const validTypes = [
        "application/pdf",
        "application/msword",
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      ];
      if (!validTypes.includes(file.type)) {
        alert("Please upload a PDF or DOC file");
        return;
      }
      // Check file size (max 5MB)
      if (file.size > 5 * 1024 * 1024) {
        alert("File size must be less than 5MB");
        return;
      }
      setSelectedFile(file);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      fileInputRef.current!.files = e.dataTransfer.files;
      handleFileChange({ target: { files: e.dataTransfer.files } } as any);
    }
  };

  const removeFile = () => {
    setSelectedFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleUpload = () => {
    if (selectedFile) {
      onUploadResume(selectedFile);
      setSelectedFile(null);
    }
  };

  const handleDelete = () => {
    if (confirm("Are you sure you want to delete your resume?")) {
      onDeleteResume();
    }
  };

  return {
    user,
    selectedFile,
    fileInputRef,
    isSaving,
    handleFileChange,
    handleDragOver,
    handleDrop,
    removeFile,
    handleUpload,
    handleDelete,
    hasResume: !!user.resumeURL,
  };
}

// File Upload Zone Component
function FileUploadZone({
  fileInputRef,
  onFileChange,
  onDragOver,
  onDrop,
}: {
  fileInputRef: React.RefObject<HTMLInputElement | null>;
  onFileChange: (e: ChangeEvent<HTMLInputElement>) => void;
  onDragOver: (e: React.DragEvent) => void;
  onDrop: (e: React.DragEvent) => void;
}) {
  return (
    <div
      className="border-2 border-dashed rounded-lg p-8 text-center cursor-pointer hover:border-primary transition-colors"
      onDragOver={onDragOver}
      onDrop={onDrop}
      onClick={() => fileInputRef.current?.click()}
    >
      <Upload className="h-10 w-10 mx-auto text-muted-foreground mb-3" />
      <p className="text-sm font-medium mb-1">
        <span className="text-primary hover:underline">Upload a file</span> or
        drag and drop
      </p>
      <p className="text-xs text-muted-foreground">PDF, DOC, DOCX up to 5MB</p>
      <input
        ref={fileInputRef}
        type="file"
        className="hidden"
        accept=".pdf,.doc,.docx"
        onChange={onFileChange}
      />
    </div>
  );
}

// Selected File Preview Component
function SelectedFilePreview({
  file,
  onUpload,
  onRemove,
  isSaving,
}: {
  file: File;
  onUpload: () => void;
  onRemove: () => void;
  isSaving: boolean;
}) {
  return (
    <div className="p-4 border rounded-lg bg-card">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <FileText className="h-8 w-8 text-primary" />
          <div>
            <p className="font-medium text-sm">{file.name}</p>
            <p className="text-xs text-muted-foreground">
              {(file.size / 1024).toFixed(2)} KB
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={onUpload} disabled={isSaving}>
            {isSaving ? "Uploading..." : "Upload"}
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={onRemove}
            disabled={isSaving}
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}

// Current Resume Display Component
function CurrentResumeDisplay({
  resumeURL,
  onDelete,
  isSaving,
}: {
  resumeURL: string;
  onDelete: () => void;
  isSaving: boolean;
}) {
  return (
    <div className="p-4 border rounded-lg bg-green-50 dark:bg-green-950/20 border-green-200 dark:border-green-800">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <CheckCircle className="h-6 w-6 text-green-600" />
          <div>
            <p className="font-medium text-sm text-green-800 dark:text-green-200">
              Resume uploaded
            </p>
            <p className="text-xs text-green-600 dark:text-green-400">
              Your resume is ready for employers
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <a href={resumeURL} target="_blank" rel="noopener noreferrer">
            <Button variant="outline" size="sm">
              <Download className="mr-2 h-4 w-4" />
              Download
            </Button>
          </a>
          <Button
            variant="ghost"
            size="icon"
            onClick={onDelete}
            disabled={isSaving}
            className="text-destructive hover:text-destructive"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}

// Standalone Section for EditProfilePage
export function ResumeSection() {
  const {
    user,
    selectedFile,
    fileInputRef,
    isSaving,
    handleFileChange,
    handleDragOver,
    handleDrop,
    removeFile,
    handleUpload,
    handleDelete,
    hasResume,
  } = useResumeUpload();

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Upload your Resume (PDF, DOCX)
        </p>
        {hasResume && (
          <div className="flex gap-2">
            <a href={user.resumeURL} target="_blank" rel="noopener noreferrer">
              <Button variant="outline" size="sm">
                <Download className="mr-2 h-4 w-4" />
                Download Current Resume
              </Button>
            </a>
            <Button
              variant="destructive"
              size="sm"
              onClick={handleDelete}
              disabled={isSaving}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        )}
      </div>

      {!selectedFile ? (
        <FileUploadZone
          fileInputRef={fileInputRef}
          onFileChange={handleFileChange}
          onDragOver={handleDragOver}
          onDrop={handleDrop}
        />
      ) : (
        <SelectedFilePreview
          file={selectedFile}
          onUpload={handleUpload}
          onRemove={removeFile}
          isSaving={isSaving}
        />
      )}
    </div>
  );
}

// Card version for ProfilePage (LinkedIn-style inline editing)
export default function ResumeSectionCard() {
  const {
    user,
    selectedFile,
    fileInputRef,
    isSaving,
    handleFileChange,
    handleDragOver,
    handleDrop,
    removeFile,
    handleUpload,
    handleDelete,
    hasResume,
  } = useResumeUpload();

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="flex items-center gap-2">
          <FileText className="h-5 w-5" />
          Resume / CV
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Show current resume status */}
        {hasResume && (
          <CurrentResumeDisplay
            resumeURL={user.resumeURL!}
            onDelete={handleDelete}
            isSaving={isSaving}
          />
        )}

        {/* Upload new resume */}
        {!selectedFile ? (
          <FileUploadZone
            fileInputRef={fileInputRef}
            onFileChange={handleFileChange}
            onDragOver={handleDragOver}
            onDrop={handleDrop}
          />
        ) : (
          <SelectedFilePreview
            file={selectedFile}
            onUpload={handleUpload}
            onRemove={removeFile}
            isSaving={isSaving}
          />
        )}
      </CardContent>
    </Card>
  );
}
