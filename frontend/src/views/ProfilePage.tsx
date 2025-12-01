import { useUserStore } from "@/stores/userStore";
import ExperienceSection from "@/components/profile/ExperienceSection";
import AboutSection from "@/components/profile/AboutSection";
import SkillsSection from "@/components/profile/SkillsSection";
import ProjectsSection from "@/components/profile/ProjectsSection";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Linkedin, Github, Globe } from "lucide-react";
import { EditProfileDialog } from "@/components/profile/EditProfileDialog";

export default function ProfilePage() {
  const { user } = useUserStore();

  return (
    <div className="container mx-auto py-10 space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Profile</h1>
        <p className="text-muted-foreground">
          จัดการหน้า Profile และ experience ของคุณได้เลย
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Left Column: User Info */}
        <div className="md:col-span-1 space-y-6">
          <div className="bg-card rounded-lg border p-6 shadow-sm sticky top-6">
            <div className="flex flex-col items-center space-y-4">
              <div className="h-24 w-24 rounded-full bg-muted flex items-center justify-center text-2xl font-bold overflow-hidden">
                {user?.AvatarURL ? (
                  <img
                    src={user.AvatarURL}
                    alt="Avatar"
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <>
                    {user?.firstName?.[0]}
                    {user?.lastName?.[0]}
                  </>
                )}
              </div>
              <div className="text-center">
                <h2 className="text-xl font-semibold">
                  {user?.firstName} {user?.lastName}
                </h2>
                <p className="text-sm text-muted-foreground">{user?.email}</p>

                <div className="flex gap-3 justify-center mt-2">
                  {user?.linkedInURL && (
                    <a
                      href={user.linkedInURL}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-muted-foreground hover:text-primary"
                    >
                      <Linkedin className="h-5 w-5" />
                    </a>
                  )}
                  {user?.githubURL && (
                    <a
                      href={user.githubURL}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-muted-foreground hover:text-primary"
                    >
                      <Github className="h-5 w-5" />
                    </a>
                  )}
                  {user?.websiteURL && (
                    <a
                      href={user.websiteURL}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-muted-foreground hover:text-primary"
                    >
                      <Globe className="h-5 w-5" />
                    </a>
                  )}
                </div>

                <div className="pt-4 w-full">
                  <EditProfileDialog />
                </div>
              </div>
            </div>
          </div>
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
              <AboutSection />
            </TabsContent>

            <TabsContent value="experience" className="space-y-6">
              <ExperienceSection />
            </TabsContent>

            <TabsContent value="skills" className="space-y-6">
              <SkillsSection />
            </TabsContent>

            <TabsContent value="projects" className="space-y-6">
              <ProjectsSection />
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
