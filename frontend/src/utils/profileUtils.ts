import type { User } from "./types/user.type";

/**
 * Calculate profile completeness percentage based on filled fields
 * @param user - User object to calculate completeness for
 * @returns Percentage (0-100) of profile completion
 */
export const calculateProfileCompleteness = (user: User | null): number => {
  if (!user) return 0;

  // Define required fields for a complete profile
  const requiredFields = [
    "firstName",
    "lastName",
    "email",
    "bio",
    "tagline",
    "avatarURL",
    "phoneNumber",
    "linkedInURL",
    "githubURL",
    "websiteURL",
  ] as const;

  // Count filled basic fields
  let filledCount = 0;
  requiredFields.forEach((field) => {
    if (user[field] && user[field] !== "") {
      filledCount++;
    }
  });

  // Add weight for array fields (skills, experiences, projects)
  const hasSkills = user.skills && user.skills.length > 0;
  const hasExperiences = user.experiences && user.experiences.length > 0;
  const hasProjects = user.projects && user.projects.length > 0;

  // Each array field counts as 1 field
  const totalFields = requiredFields.length + 3; // +3 for skills, experiences, projects
  const totalFilled =
    filledCount +
    (hasSkills ? 1 : 0) +
    (hasExperiences ? 1 : 0) +
    (hasProjects ? 1 : 0);

  return Math.round((totalFilled / totalFields) * 100);
};

/**
 * Get list of missing profile fields
 * @param user - User object to check
 * @returns Array of missing field names (human-readable)
 */
export const getMissingFields = (user: User | null): string[] => {
  if (!user) return [];

  const missingFields: string[] = [];

  const fieldLabels: Record<string, string> = {
    firstName: "First Name",
    lastName: "Last Name",
    email: "Email",
    bio: "Bio",
    tagline: "Tagline",
    avatarURL: "Profile Picture",
    phoneNumber: "Phone Number",
    linkedInURL: "LinkedIn URL",
    githubURL: "GitHub URL",
    websiteURL: "Website URL",
  };

  Object.entries(fieldLabels).forEach(([field, label]) => {
    if (!user[field as keyof User] || user[field as keyof User] === "") {
      missingFields.push(label);
    }
  });

  if (!user.skills || user.skills.length === 0) {
    missingFields.push("Skills");
  }
  if (!user.experiences || user.experiences.length === 0) {
    missingFields.push("Work Experience");
  }
  if (!user.projects || user.projects.length === 0) {
    missingFields.push("Projects");
  }

  return missingFields;
};

/**
 * Generate shareable profile link
 * @param userId - User ID to generate link for
 * @returns Shareable profile URL
 */
export const getShareableLink = (userId: string): string => {
  const baseUrl = window.location.origin;
  return `${baseUrl}/profile/${userId}`;
};

/**
 * Copy text to clipboard
 * @param text - Text to copy
 * @returns Promise that resolves when text is copied
 */
export const copyToClipboard = async (text: string): Promise<void> => {
  if (navigator.clipboard && navigator.clipboard.writeText) {
    await navigator.clipboard.writeText(text);
  } else {
    // Fallback for older browsers
    const textArea = document.createElement("textarea");
    textArea.value = text;
    textArea.style.position = "fixed";
    textArea.style.left = "-999999px";
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    try {
      document.execCommand("copy");
    } finally {
      document.body.removeChild(textArea);
    }
  }
};

/**
 * Download resume from URL
 * @param resumeURL - URL of the resume to download
 * @param fileName - Optional custom file name
 */
export const downloadResume = async (
  resumeURL: string,
  fileName: string = "resume.pdf"
): Promise<void> => {
  try {
    const response = await fetch(resumeURL);
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = fileName;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error("Failed to download resume:", error);
    throw error;
  }
};

export const oneHour = 1000 * 60 * 60;
export const oneDay = 1000 * 60 * 60 * 24;
