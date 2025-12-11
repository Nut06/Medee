import { useState, useEffect } from "react";
import { Pencil, Plus, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { useProfile } from "@/hooks/useProfile";
import type { Skill } from "@/utils/types/user.type";

// Standalone Skills Section for Edit Profile Page
export function SkillsSection() {
  const { user, onUpdateSkills, isSaving } = useProfile();
  const [tempSkills, setTempSkills] = useState<Skill[]>([]);
  const [newSkill, setNewSkill] = useState("");

  useEffect(() => {
    setTempSkills(user?.skills || []);
  }, [user?.skills]);

  const handleAddSkill = () => {
    if (newSkill.trim()) {
      if (!tempSkills.some((s) => s.name === newSkill.trim())) {
        setTempSkills([...tempSkills, { name: newSkill.trim() }]);
      }
      setNewSkill("");
    }
  };

  const handleRemoveSkill = (skillName: string) => {
    setTempSkills(tempSkills.filter((s) => s.name !== skillName));
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleAddSkill();
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <Input
          placeholder="Type a skill and press Enter"
          value={newSkill}
          onChange={(e) => setNewSkill(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <Button onClick={handleAddSkill} size="icon">
          <Plus className="h-4 w-4" />
        </Button>
      </div>
      <div className="flex flex-wrap gap-2 min-h-[100px] p-4 border rounded-md bg-muted/50">
        {tempSkills.length === 0 ? (
          <p className="text-sm text-muted-foreground w-full text-center self-center">
            No skills added yet.
          </p>
        ) : (
          tempSkills.map((skill, index) => (
            <Badge key={index} variant="secondary" className="gap-1">
              {skill.name}
              <button
                onClick={() => handleRemoveSkill(skill.name || "")}
                className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
              >
                <X className="h-3 w-3 text-muted-foreground hover:text-foreground" />
              </button>
            </Badge>
          ))
        )}
      </div>
      <p className="text-sm text-muted-foreground">
        Add up to 15 skills. Press Enter or click + to add.
      </p>

      {/* Save Button */}
      <div className="flex justify-end pt-4">
        <Button onClick={() => onUpdateSkills(tempSkills)} disabled={isSaving}>
          {isSaving ? "Saving..." : "Save"}
        </Button>
      </div>
    </div>
  );
}

// Card version for Profile Page
export default function SkillsSectionCard() {
  const { user, onUpdateSkills, isSaving } = useProfile();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [tempSkills, setTempSkills] = useState<Skill[]>([]);
  const [newSkill, setNewSkill] = useState("");

  const handleOpen = () => {
    setTempSkills(user?.skills || []);
    setIsDialogOpen(true);
  };

  const handleAddSkill = () => {
    if (newSkill.trim()) {
      if (!tempSkills.some((s) => s.name === newSkill.trim())) {
        setTempSkills([...tempSkills, { name: newSkill.trim() }]);
      }
      setNewSkill("");
    }
  };

  const handleRemoveSkill = (skillName: string) => {
    setTempSkills(tempSkills.filter((s) => s.name !== skillName));
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleAddSkill();
    }
  };

  const handleSave = async () => {
    await onUpdateSkills(tempSkills);
    setIsDialogOpen(false);
  };

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Skills</CardTitle>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button variant="outline" size="sm" onClick={handleOpen}>
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>Edit Skills</DialogTitle>
              <DialogDescription>
                Add or remove skills from your profile.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="flex gap-2">
                <Input
                  placeholder="Add a skill (e.g. React)"
                  value={newSkill}
                  onChange={(e) => setNewSkill(e.target.value)}
                  onKeyDown={handleKeyDown}
                />
                <Button onClick={handleAddSkill} size="icon">
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
              <div className="flex flex-wrap gap-2 min-h-[100px] p-4 border rounded-md bg-muted/50">
                {tempSkills.length === 0 ? (
                  <p className="text-sm text-muted-foreground w-full text-center self-center">
                    No skills added yet.
                  </p>
                ) : (
                  tempSkills.map((skill, index) => (
                    <Badge key={index} variant="secondary" className="gap-1">
                      {skill.name}
                      <button
                        onClick={() => handleRemoveSkill(skill.name || "")}
                        className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                      >
                        <X className="h-3 w-3 text-muted-foreground hover:text-foreground" />
                      </button>
                    </Badge>
                  ))
                )}
              </div>
            </div>
            <DialogFooter>
              <Button onClick={handleSave} disabled={isSaving}>
                Save changes
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </CardHeader>
      <CardContent>
        <div className="flex flex-wrap gap-2">
          {user?.skills && user.skills.length > 0 ? (
            user.skills.map((skill, index) => (
              <Badge key={index} variant="secondary">
                {skill.name}
              </Badge>
            ))
          ) : (
            <p className="text-muted-foreground italic">No skills added yet.</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
