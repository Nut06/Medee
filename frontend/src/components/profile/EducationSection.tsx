import { Plus, Trash2, Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useProfile } from "@/hooks/useProfile";

export function EducationSection() {
  const { user } = useProfile();

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Add your educational background.
        </p>
        <Button size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add Education
        </Button>
      </div>

      <div className="space-y-3">
        {user?.educations && user.educations.length > 0 ? (
          user.educations.map((edu: any) => (
            <Card key={edu.id} className="p-4">
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="font-semibold">{edu.degree}</h3>
                  <p className="text-sm text-muted-foreground">{edu.school}</p>
                  {edu.fieldOfStudy && (
                    <p className="text-sm text-muted-foreground">
                      {edu.fieldOfStudy}
                    </p>
                  )}
                  <p className="text-xs text-muted-foreground mt-1">
                    {new Date(edu.startDate).toLocaleDateString()} -{" "}
                    {edu.endDate
                      ? new Date(edu.endDate).toLocaleDateString()
                      : "Present"}
                  </p>
                </div>
                <div className="flex gap-2">
                  <Button variant="ghost" size="icon">
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="icon">
                    <Trash2 className="h-4 w-4 text-destructive" />
                  </Button>
                </div>
              </div>
            </Card>
          ))
        ) : (
          <Card className="p-8 text-center">
            <p className="text-sm text-muted-foreground">
              No education added yet. Click "Add Education" to get started.
            </p>
          </Card>
        )}
      </div>
    </div>
  );
}
