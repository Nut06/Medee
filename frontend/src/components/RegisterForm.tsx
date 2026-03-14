import { zodResolver } from "@hookform/resolvers/zod";
import { BriefcaseBusiness, Lock, Mail } from "lucide-react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { z } from "zod";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useAuthForm } from "@/hooks/userAuthForm";
import { register } from "@/services/auth.service";
import type { RegisterRequest } from "@/utils/types/user.type";
import { useUserStore } from "@/stores/userStore";

const registerSchema = z.object({
  firstName: z.string({ error: "กรุณากรอกชื่อ" }).min(1, "กรุณากรอกชื่อ"),
  lastName: z.string({ error: "กรุณากรอกนามสกุล" }).min(1, "กรุณากรอกนามสกุล"),
  email: z.email("กรุณากรอก Email"),
  password: z.string({ error: "กรุณากรอกรหัสผ่าน" }).min(6, "กรุณากรอกรหัสผ่าน"),
});

type RegisterFormValues = z.infer<typeof registerSchema>;

export const RegisterForm = () => {
  const { email, password, setField, loading, error } = useAuthForm();
  const navigate = useNavigate();

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    mode: "onSubmit",
    defaultValues: {
      firstName: "",
      lastName: "",
      email,
      password,
    },
  });

  const handleLoginSubmit = async (
    values: RegisterFormValues
  ): Promise<void> => {
    const payload: RegisterRequest = {
      email: values.email,
      password: values.password,
      firstName: values.firstName,
      lastName: values.lastName,
    };

    setField("email", values.email);
    setField("firstName", values.firstName);
    setField("lastName", values.lastName);
    setField("password", values.password);

    try {
      setField("loading", true);
      const response = await register(payload);

      // Update store
      useUserStore.getState().setUser(response.user);

      // Redirect based on company membership
      if (response.companies && response.companies.length > 0) {
        navigate("/company");
      } else {
        navigate("/user");
      }
    } catch (error) {
      console.error("Registration failed:", error);
    } finally {
      setField("loading", false);
    }
  };

  return (
    <div className="mx-auto w-full max-w-sm rounded-3xl border bg-background px-8 py-10 shadow-lg">
      <div className="mb-8 space-y-3 text-center">
        <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-primary/10 text-primary">
          <BriefcaseBusiness className="size-6" />
        </div>
        <div>
          <h1 className="text-2xl font-semibold">ยินดีต้อนรับสู่ Medee</h1>
          <p className="text-sm text-muted-foreground">
            กรุณากรอกข้อมูลเพื่อสมัครเป็นสมาชิก
          </p>
        </div>
      </div>

      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(handleLoginSubmit)}
          className="space-y-6"
        >
          <FormField
            control={form.control}
            name="firstName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>First Name</FormLabel>
                <FormControl>
                  <div className="relative">
                    <Mail className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      type="name"
                      placeholder="สมชาย"
                      className="pl-10"
                      autoComplete="name"
                      {...field}
                      onChange={(event) => {
                        field.onChange(event);
                        setField("firstName", event.target.value);
                      }}
                    />
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="lastName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Last Name</FormLabel>
                <FormControl>
                  <div className="relative">
                    <Mail className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      type="name"
                      placeholder="ใจดี"
                      className="pl-10"
                      autoComplete="name"
                      {...field}
                      onChange={(event) => {
                        field.onChange(event);
                        setField("lastName", event.target.value);
                      }}
                    />
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email Address</FormLabel>
                <FormControl>
                  <div className="relative">
                    <Mail className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      type="email"
                      placeholder="you@example.com"
                      className="pl-10"
                      autoComplete="email"
                      {...field}
                      onChange={(event) => {
                        field.onChange(event);
                        setField("email", event.target.value);
                      }}
                    />
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <div className="relative">
                    <Lock className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      type="password"
                      placeholder="********"
                      className="pl-10"
                      autoComplete="current-password"
                      {...field}
                      onChange={(event) => {
                        field.onChange(event);
                        setField("password", event.target.value);
                      }}
                    />
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          {error && (
            <Alert variant="destructive">
              <AlertDescription>
                {error || "Unable to register. Please try again."}
              </AlertDescription>
            </Alert>
          )}

          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? "Logging In..." : "ลงทะเบียน"}
          </Button>
        </form>
      </Form>

      <p className="flex gap-x-3 justify-self-center mt-8 text-center text-sm text-muted-foreground">
        มีบัญชีแล้ว ?{" "}
        <Link to="/login" className="font-medium text-primary hover:underline">
          เข้าสู่ระบบ
        </Link>
      </p>
    </div>
  );
};
