import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { format } from "date-fns";
import { Briefcase, Plus, Pencil, Trash2 } from "lucide-react";

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
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useProfile } from "@/hooks/useProfile";
import type { WorkExperience } from "@/utils/types/user.type";

const experienceSchema = z.object({
  position: z.string().min(2, "Position is required"),
  companyName: z.string().min(2, "Company name is required"),
  startDate: z.string().min(1, "Start date is required"),
  endDate: z.string().optional(),
  description: z.string().optional(),
});

type ExperienceFormValues = z.infer<typeof experienceSchema>;

export default function ExperienceSection() {
  const { user, onAddExperience, onUpdateExperience, onDeleteExperience } =
    useProfile();
  const [isDialogOpen, setIsDialogOpen] = useState<boolean>(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const form = useForm<ExperienceFormValues>({
    resolver: zodResolver(experienceSchema),
    defaultValues: {
      position: "",
      companyName: "",
      startDate: "",
      endDate: "",
      description: "",
    },
  });

  const onSubmit = async (data: ExperienceFormValues) => {
    if (editingId) {
      await onUpdateExperience(editingId, data);
    } else {
      await onAddExperience(data);
    }
    setIsDialogOpen(false);
    form.reset();
    setEditingId(null);
  };

  const handleEdit = (exp: WorkExperience) => {
    setEditingId(exp.id || null);
    form.reset({
      position: exp.position,
      companyName: exp.companyName,
      startDate: exp.startDate
        ? new Date(exp.startDate).toISOString().split("T")[0]
        : "",
      endDate: exp.endDate
        ? new Date(exp.endDate).toISOString().split("T")[0]
        : "",
      description: exp.description || "",
    });
    setIsDialogOpen(true);
  };

  const handleAddNew = () => {
    setEditingId(null);
    form.reset({
      position: "",
      companyName: "",
      startDate: "",
      endDate: "",
      description: "",
    });
  };

  const handleDialogChange = (open: boolean) => {
    setIsDialogOpen(open);
    if (open && !editingId) {
      // Reset form when opening dialog for adding new experience
      handleAddNew();
    }
  };

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-row items-center justify-between">
        <div className="flex items-center gap-2">
          <Briefcase className="h-5 w-5 text-primary" />
          <CardTitle>Experience</CardTitle>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={handleDialogChange}>
          <DialogTrigger asChild>
            <Button variant="outline" size="sm">
              <Plus className="mr-2 h-4 w-4" />
              Add Experience
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {editingId ? "Edit Experience" : "Add Experience"}
              </DialogTitle>
              <DialogDescription>
                Add your work experience to your profile.
              </DialogDescription>
            </DialogHeader>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-4"
              >
                <FormField
                  control={form.control}
                  name="position"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Position</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. Senior Frontend Developer"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="companyName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Company Name</FormLabel>
                      <FormControl>
                        <Input placeholder="e.g. Google" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div className="grid grid-cols-2 gap-4">
                  <FormField
                    control={form.control}
                    name="startDate"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Start Date</FormLabel>
                        <FormControl>
                          <Input type="date" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="endDate"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>End Date</FormLabel>
                        <FormControl>
                          <Input type="date" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Description</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Describe your responsibilities and achievements..."
                          className="resize-none"
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
      </CardHeader>
      <CardContent className="space-y-6">
        {/* @ts-ignore */}
        {user.experiences && user.experiences.length === 0 ? (
          <div className="text-center py-12 text-muted-foreground border-2 border-dashed rounded-lg">
            <div className="flex justify-center mb-4">
              <div className="relative">
                <Briefcase className="h-12 w-12 opacity-50" />
                <div className="absolute -bottom-1 -right-1 bg-muted rounded-full p-1">
                  <Plus className="h-4 w-4" />
                </div>
              </div>
            </div>
            <h3 className="font-semibold text-lg text-foreground mb-2">
              No Experience Added
            </h3>
            <p className="text-sm mb-6 max-w-md mx-auto">
              Showcase your internships and work history to stand out.
            </p>
            <Button
              onClick={() => {
                handleAddNew();
                setIsDialogOpen(true);
              }}
              size="lg"
            >
              <Plus className="mr-2 h-4 w-4" />
              Add Experience
            </Button>
          </div>
        ) : (
          /* @ts-ignore */
          user.experiences.map((exp: WorkExperience) => (
            <div
              key={exp.id}
              className="relative pl-6 border-l-2 border-muted last:border-l-0 pb-6 last:pb-0"
            >
              <div className="absolute -left-[9px] top-0 h-4 w-4 rounded-full bg-primary" />
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{exp.position}</h3>
                  <p className="text-muted-foreground">{exp.companyName}</p>
                  <p className="text-sm text-muted-foreground mt-1">
                    {exp.startDate
                      ? format(new Date(exp.startDate), "MMM yyyy")
                      : ""}{" "}
                    -{" "}
                    {exp.endDate
                      ? format(new Date(exp.endDate), "MMM yyyy")
                      : "Present"}
                  </p>
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEdit(exp)}
                  >
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-destructive hover:text-destructive"
                    onClick={() => exp.id && onDeleteExperience(exp.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
              {exp.description && (
                <p className="mt-2 text-sm text-foreground/80 whitespace-pre-line">
                  {exp.description}
                </p>
              )}
            </div>
          ))
        )}
      </CardContent>
    </Card>
  );
}
