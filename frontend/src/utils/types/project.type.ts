type SubmitStatus = "submitted" | "approved" | "rejected";

export interface Project {
  name: string;
  detail: string;
  createdBy: string;
  createdAt: string;
}

export interface Submission {
  id: string;
  projectId: string;
  userId: string;
  fileUrl: string;
  status: SubmitStatus;
  submittedAt: string;
}

export interface Portfolio {
  id: string;
  userId: string;
  projecId: string;
  submissionId: string;
  addedAt: string;
}
