import { useProfile } from "@/hooks/useProfile";
import ExperienceSection from "@/components/profile/ExperienceSection";
import AboutSection from "@/components/profile/AboutSection";
import SkillsSectionCard from "@/components/profile/SkillsSection";
import ProjectsSection from "@/components/profile/ProjectsSection";
import ProfileCompletenessCard from "@/components/profile/ProfileCompletenessCard";
import AvatarUpload from "@/components/profile/AvatarUpload";
import PersonalInfoCard from "@/components/profile/PersonalInfoSection";
import EducationSectionCard from "@/components/profile/EducationSection";
import ResumeSectionCard from "@/components/profile/ResumeSection";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Linkedin, Github, Globe, Share2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { getShareableLink, copyToClipboard } from "@/utils/profileUtils";
import { AuthNavbar } from "@/components/navbar/AuthNavbar";

export default function ProfilePage() {
  const { user, isLoading } = useProfile();

  const handleShareProfile = async () => {
    if (user.id) {
      try {
        const link = getShareableLink(user.id);
        await copyToClipboard(link);
        toast.success("Profile link copied to clipboard!");
      } catch (error) {
        toast.error("Failed to copy link");
      }
    }
  };

  return (
    <div className="container mx-auto py-10 space-y-8">
      <AuthNavbar />
      <div>
        <h1 className="text-3xl font-bold">Profile</h1>
        <p className="text-muted-foreground">
          จัดการหน้า Profile และ experience ของคุณได้เลย
        </p>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading profile...</p>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Left Column: User Info */}
          <div className="md:col-span-1 space-y-6">
            <div className="bg-card rounded-lg border p-6 shadow-sm sticky top-6">
              <div className="flex flex-col items-center space-y-4">
                <AvatarUpload
                  avatarURL={user.AvatarURL}
                  firstName={user.firstName}
                  lastName={user.lastName}
                />
                <div className="text-center w-full">
                  <h2 className="text-xl font-semibold">
                    {user.firstName} {user.lastName}
                  </h2>
                  {user.tagline && (
                    <p className="text-sm text-muted-foreground mt-1">
                      {user.tagline}
                    </p>
                  )}
                  <p className="text-sm text-muted-foreground mt-1">
                    {user.email}
                  </p>

                  <div className="flex gap-3 justify-center mt-3">
                    {user.linkedInURL && (
                      <a
                        href={user.linkedInURL}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-muted-foreground hover:text-primary transition-colors"
                      >
                        <Linkedin className="h-5 w-5" />
                      </a>
                    )}
                    {user.githubURL && (
                      <a
                        href={user.githubURL}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-muted-foreground hover:text-primary transition-colors"
                      >
                        <Github className="h-5 w-5" />
                      </a>
                    )}
                    {user.websiteURL && (
                      <a
                        href={user.websiteURL}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-muted-foreground hover:text-primary transition-colors"
                      >
                        <Globe className="h-5 w-5" />
                      </a>
                    )}
                  </div>

                  <div className="pt-4 w-full">
                    <Button
                      variant="outline"
                      className="w-full"
                      onClick={handleShareProfile}
                    >
                      <Share2 className="mr-2 h-4 w-4" />
                      Share Profile
                    </Button>
                  </div>
                </div>
              </div>
            </div>

            {/* Profile Completeness Card */}
            <ProfileCompletenessCard />
          </div>

          {/* Right Column: Details with Tabs */}
          <div className="md:col-span-2">
            <Tabs defaultValue="about" className="w-full">
              <TabsList className="grid w-full grid-cols-4 mb-6">
                <TabsTrigger value="about">About</TabsTrigger>
                <TabsTrigger value="experience">Experience</TabsTrigger>
                <TabsTrigger value="skills">Skills</TabsTrigger>
                <TabsTrigger value="projects">Projects</TabsTrigger>
              </TabsList>

              <TabsContent value="about" className="space-y-6">
                <PersonalInfoCard />
                <AboutSection />
                <ResumeSectionCard />
              </TabsContent>

              <TabsContent value="experience" className="space-y-6">
                <ExperienceSection />
                <EducationSectionCard />
              </TabsContent>

              <TabsContent value="skills" className="space-y-6">
                <SkillsSectionCard />
              </TabsContent>

              <TabsContent value="projects" className="space-y-6">
                <ProjectsSection />
              </TabsContent>
            </Tabs>
          </div>
        </div>
      )}
    </div>
  );
}
