import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Plus, Pencil, Trash2, Github, Globe } from "lucide-react";

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

const projectSchema = z.object({
  title: z.string().min(2, "Title is required"),
  description: z.string().min(2, "Description is required"),
  imageURL: z.string().url("Invalid URL").optional().or(z.literal("")),
  githubURL: z.string().url("Invalid URL").optional().or(z.literal("")),
  demoURL: z.string().url("Invalid URL").optional().or(z.literal("")),
});

type ProjectFormValues = z.infer<typeof projectSchema>;

export default function ProjectsSection() {
  const { user, onAddProject, onUpdateProject, onDeleteProject, isSaving } =
    useProfile();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const form = useForm<ProjectFormValues>({
    resolver: zodResolver(projectSchema),
    defaultValues: {
      title: "",
      description: "",
      imageURL: "",
      githubURL: "",
      demoURL: "",
    },
  });

  const onSubmit = async (data: ProjectFormValues) => {
    if (editingId) {
      await onUpdateProject(editingId, data);
    } else {
      await onAddProject(data);
    }
    setIsDialogOpen(false);
    form.reset();
    setEditingId(null);
  };

  const handleEdit = (project: any) => {
    setEditingId(project.id || null);
    form.reset({
      title: project.title || project.name, // Handle potential type mismatch
      description: project.description || project.detail,
      imageURL: project.imageURL || "",
      githubURL: project.githubURL || "",
      demoURL: project.demoURL || "",
    });
    setIsDialogOpen(true);
  };

  const handleAddNew = () => {
    setEditingId(null);
    form.reset({
      title: "",
      description: "",
      imageURL: "",
      githubURL: "",
      demoURL: "",
    });
    setIsDialogOpen(true);
  };

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Projects</CardTitle>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button variant="outline" size="sm" onClick={handleAddNew}>
              <Plus className="mr-2 h-4 w-4" />
              Add Project
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {editingId ? "Edit Project" : "Add Project"}
              </DialogTitle>
              <DialogDescription>Showcase your best work.</DialogDescription>
            </DialogHeader>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-4"
              >
                <FormField
                  control={form.control}
                  name="title"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Project Title</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. E-commerce Website"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Description</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Describe the project..."
                          className="resize-none"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="imageURL"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Image URL (Optional)</FormLabel>
                      <FormControl>
                        <Input placeholder="https://..." {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div className="grid grid-cols-2 gap-4">
                  <FormField
                    control={form.control}
                    name="githubURL"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>GitHub URL (Optional)</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="https://github.com/..."
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="demoURL"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Demo URL (Optional)</FormLabel>
                        <FormControl>
                          <Input placeholder="https://..." {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <DialogFooter>
                  <Button type="submit" disabled={isSaving}>
                    Save changes
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          </DialogContent>
        </Dialog>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* @ts-ignore */}
        {user?.projects?.length === 0 ? (
          <div className="text-center py-6 text-muted-foreground border-2 border-dashed rounded-lg">
            <p>No projects added yet.</p>
          </div>
        ) : (
          /* @ts-ignore */
          user?.projects?.map((project: any) => (
            <div
              key={project.id}
              className="flex flex-col md:flex-row gap-4 border rounded-lg p-4"
            >
              {project.imageURL && (
                <div className="w-full md:w-48 h-32 bg-muted rounded-md overflow-hidden flex-shrink-0">
                  <img
                    src={project.imageURL}
                    alt={project.title || project.name}
                    className="w-full h-full object-cover"
                  />
                </div>
              )}
              <div className="flex-1 space-y-2">
                <div className="flex justify-between items-start">
                  <h3 className="font-semibold text-lg">
                    {project.title || project.name}
                  </h3>
                  <div className="flex gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleEdit(project)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="text-destructive hover:text-destructive"
                      onClick={() => project.id && onDeleteProject(project.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
                <p className="text-sm text-muted-foreground">
                  {project.description || project.detail}
                </p>
                <div className="flex gap-3 pt-2">
                  {project.githubURL && (
                    <a
                      href={project.githubURL}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center text-xs text-primary hover:underline"
                    >
                      <Github className="mr-1 h-3 w-3" />
                      GitHub
                    </a>
                  )}
                  {project.demoURL && (
                    <a
                      href={project.demoURL}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center text-xs text-primary hover:underline"
                    >
                      <Globe className="mr-1 h-3 w-3" />
                      Live Demo
                    </a>
                  )}
                </div>
              </div>
            </div>
          ))
        )}
      </CardContent>
    </Card>
  );
}
