import { Plus, Trash2, Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useProfile } from "@/hooks/useProfile";

export function WorkExperienceSection() {
  const { user } = useProfile();

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Add your work experience to showcase your professional journey.
        </p>
        <Button size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add Experience
        </Button>
      </div>

      <div className="space-y-3">
        {user.experiences && user.experiences.length > 0 ? (
          user.experiences.map((exp: any) => (
            <Card key={exp.id} className="p-4">
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="font-semibold">{exp.position}</h3>
                  <p className="text-sm text-muted-foreground">{exp.company}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {new Date(exp.startDate).toLocaleDateString()} -{" "}
                    {exp.endDate
                      ? new Date(exp.endDate).toLocaleDateString()
                      : "Present"}
                  </p>
                  {exp.description && (
                    <p className="text-sm mt-2">{exp.description}</p>
                  )}
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
              No work experience added yet. Click "Add Experience" to get
              started.
            </p>
          </Card>
        )}
      </div>
    </div>
  );
}
