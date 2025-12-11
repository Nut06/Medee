import { ArrowLeft, Eye } from "lucide-react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { PersonalInfoSection } from "@/components/profile/PersonalInfoSection";
import { SkillsSection } from "@/components/profile/SkillsSection";
import { WorkExperienceSection } from "@/components/profile/WorkExperienceSection";
import { EducationSection } from "@/components/profile/EducationSection";
import { PortfolioSection } from "@/components/profile/PortfolioSection";
import { ResumeSection } from "@/components/profile/ResumeSection";

export default function EditProfilePage() {
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b sticky top-0 z-10">
        <div className="max-w-5xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link to="/user">
              <Button variant="ghost" size="icon">
                <ArrowLeft className="h-5 w-5" />
              </Button>
            </Link>
            <div>
              <h1 className="text-2xl font-bold">Edit Profile</h1>
              <p className="text-sm text-muted-foreground">
                Update your profile to keep it fresh and appealing to companies.
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Link to="/user">
              <Button variant="outline" size="sm">
                <Eye className="mr-2 h-4 w-4" />
                Preview Profile
              </Button>
            </Link>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-5xl mx-auto px-4 py-8">
        <Accordion
          type="multiple"
          defaultValue={["personal", "skills"]}
          className="space-y-4"
        >
          {/* Personal Information */}
          <AccordionItem
            value="personal"
            className="bg-white rounded-lg border"
          >
            <AccordionTrigger className="px-6 py-4 hover:no-underline">
              <h2 className="text-lg font-semibold">Personal Information</h2>
            </AccordionTrigger>
            <AccordionContent>
              <Separator />
              <div className="px-6 py-4">
                <PersonalInfoSection />
              </div>
            </AccordionContent>
          </AccordionItem>

          {/* Resume/CV */}
          <AccordionItem value="resume" className="bg-white rounded-lg border">
            <AccordionTrigger className="px-6 py-4 hover:no-underline">
              <h2 className="text-lg font-semibold">Resume/CV</h2>
            </AccordionTrigger>
            <AccordionContent>
              <Separator />
              <div className="px-6 py-4">
                <ResumeSection />
              </div>
            </AccordionContent>
          </AccordionItem>

          {/* Work Experience */}
          <AccordionItem
            value="experience"
            className="bg-white rounded-lg border"
          >
            <AccordionTrigger className="px-6 py-4 hover:no-underline">
              <h2 className="text-lg font-semibold">Work Experience</h2>
            </AccordionTrigger>
            <AccordionContent>
              <Separator />
              <div className="px-6 py-4">
                <WorkExperienceSection />
              </div>
            </AccordionContent>
          </AccordionItem>

          {/* Education */}
          <AccordionItem
            value="education"
            className="bg-white rounded-lg border"
          >
            <AccordionTrigger className="px-6 py-4 hover:no-underline">
              <h2 className="text-lg font-semibold">Education</h2>
            </AccordionTrigger>
            <AccordionContent>
              <Separator />
              <div className="px-6 py-4">
                <EducationSection />
              </div>
            </AccordionContent>
          </AccordionItem>

          {/* Skills */}
          <AccordionItem value="skills" className="bg-white rounded-lg border">
            <AccordionTrigger className="px-6 py-4 hover:no-underline">
              <h2 className="text-lg font-semibold">Skills</h2>
            </AccordionTrigger>
            <AccordionContent>
              <Separator />
              <div className="px-6 py-4">
                <SkillsSection />
              </div>
            </AccordionContent>
          </AccordionItem>

          {/* Portfolio/Projects */}
          <AccordionItem
            value="portfolio"
            className="bg-white rounded-lg border"
          >
            <AccordionTrigger className="px-6 py-4 hover:no-underline">
              <h2 className="text-lg font-semibold">Portfolio / Projects</h2>
            </AccordionTrigger>
            <AccordionContent>
              <Separator />
              <div className="px-6 py-4">
                <PortfolioSection />
              </div>
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </main>
    </div>
  );
}
