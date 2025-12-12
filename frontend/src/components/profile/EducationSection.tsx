import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Plus, Trash2, Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useProfile } from "@/hooks/useProfile";
import type { Education } from "@/utils/types/user.type";

const educationSchema = z.object({
  instituteName: z.string().min(2, "Institute name is required"),
  degree: z.string().min(2, "Degree is required"),
  fieldOfStudy: z.string().min(2, "Field of study is required"),
  graduationYear: z
    .string()
    .refine(
      (val) => {
        if (!val) return true; // optional
        const num = parseInt(val);
        return !isNaN(num) && num >= 1900 && num <= 2100;
      },
      { message: "Year must be between 1900 and 2100" }
    )
    .optional(),
});

type EducationFormValues = z.infer<typeof educationSchema>;

export function EducationSection() {
  const { user, onAddEducation, onUpdateEducation, onDeleteEducation } =
    useProfile();
  const [isDialogOpen, setIsDialogOpen] = useState<boolean>(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const form = useForm<EducationFormValues>({
    resolver: zodResolver(educationSchema),
    defaultValues: {
      instituteName: "",
      degree: "",
      fieldOfStudy: "",
      graduationYear: "",
    },
  });

  const onSubmit = async (data: EducationFormValues) => {
    const payload = {
      instituteName: data.instituteName,
      degree: data.degree,
      fieldOfStudy: data.fieldOfStudy,
      graduationYear: data.graduationYear
        ? parseInt(data.graduationYear)
        : undefined,
    };

    if (editingId) {
      await onUpdateEducation(editingId, payload);
    } else {
      await onAddEducation(payload);
    }
    setIsDialogOpen(false);
    form.reset();
    setEditingId(null);
  };

  const handleEdit = (edu: Education) => {
    setEditingId(edu.id || null);
    form.reset({
      instituteName: edu.instituteName,
      degree: edu.degree,
      fieldOfStudy: edu.fieldOfStudy,
      graduationYear: edu.graduationYear ? String(edu.graduationYear) : "",
    });
    setIsDialogOpen(true);
  };

  const handleAddNew = () => {
    setEditingId(null);
    form.reset({
      instituteName: "",
      degree: "",
      fieldOfStudy: "",
      graduationYear: "",
    });
  };

  const handleDialogChange = (open: boolean) => {
    setIsDialogOpen(open);
    if (open && !editingId) {
      handleAddNew();
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Add your educational background.
        </p>
        <Dialog open={isDialogOpen} onOpenChange={handleDialogChange}>
          <DialogTrigger asChild>
            <Button size="sm">
              <Plus className="mr-2 h-4 w-4" />
              Add Education
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {editingId ? "Edit Education" : "Add Education"}
              </DialogTitle>
              <DialogDescription>
                Add your educational background to your profile.
              </DialogDescription>
            </DialogHeader>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-4"
              >
                <FormField
                  control={form.control}
                  name="instituteName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Institute Name</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. Stanford University"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="degree"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Degree</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. Bachelor of Science"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="fieldOfStudy"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Field of Study</FormLabel>
                      <FormControl>
                        <Input placeholder="e.g. Computer Science" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="graduationYear"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Graduation Year (optional)</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          placeholder="e.g. 2020"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <DialogFooter>
                  <Button type="submit">Save changes</Button>
                </DialogFooter>
              </form>
            </Form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="space-y-3">
        {user.educations && user.educations.length > 0 ? (
          user.educations.map((edu: any) => (
            <Card key={edu.id} className="p-4">
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="font-semibold">{edu.degree}</h3>
                  <p className="text-sm text-muted-foreground">
                    {edu.instituteName}
                  </p>
                  {edu.fieldOfStudy && (
                    <p className="text-sm text-muted-foreground">
                      {edu.fieldOfStudy}
                    </p>
                  )}
                  {edu.graduationYear && (
                    <p className="text-xs text-muted-foreground mt-1">
                      Graduated: {edu.graduationYear}
                    </p>
                  )}
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEdit(edu)}
                  >
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => edu.id && onDeleteEducation(edu.id)}
                  >
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
