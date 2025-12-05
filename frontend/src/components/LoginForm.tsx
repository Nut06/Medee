import { zodResolver } from "@hookform/resolvers/zod";
import { BriefcaseBusiness, Lock, Mail } from "lucide-react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { z } from "zod";

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
import type { LoginRequest } from "@/utils/types/user.type";
import { loginLocal } from "@/services/authService";
import { useUserStore } from "@/stores/userStore";

const loginSchema = z.object({
  email: z.email("Please enter a valid email  address."),
  password: z.string().min(6, "Password must be at least 6 characters."),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export const LoginForm = () => {
  const { email, password, setField, loading, error } = useAuthForm();
  const { setUser } = useUserStore();
  
  const navigate = useNavigate();

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email,
      password,
    },
  });

  const handleLoginSubmit = async (values: LoginFormValues): Promise<void> => {
    const payload: LoginRequest = {
      email: values.email,
      password: values.password,
    };

    setField("email", values.email);
    setField("password", values.password);

    try {
      setField("loading", true);
      const response = await loginLocal(payload);

      // Update store
      setUser(response.user);
      setAuth(true);

      // Redirect based on company membership
      if (response.companies && response.companies.length > 0) {
        navigate("/company");
      } else {
        navigate("/user");
      }
    } catch (error) {
      // Error handling is done in useAuthForm, but we need to catch here to stop navigation
      console.error("Login failed:", error);
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
          <h1 className="text-2xl font-semibold">Welcome Back</h1>
          <p className="pt-1.5 text-sm text-muted-foreground">
            กรุณากรอก Email และ Password เพื่อเข้าสู่ระบบ
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

          <div className="flex items-center justify-end text-sm">
            <a href="#" className="font-medium text-primary hover:underline">
              ลืมรหัสผ่าน?
            </a>
          </div>

          {error && (
            <p className="text-sm font-medium text-destructive">
              {error || "Unable to log in. Please try again."}
            </p>
          )}

          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? "Logging In..." : "เข้าสู่ระบบ"}
          </Button>
        </form>
      </Form>

      <p className="flex gap-x-3 justify-self-center mt-8 text-center text-sm text-muted-foreground">
        ยังไม่มีบัญชี ?{" "}
        <Link
          to="/register"
          className="font-medium text-primary hover:underline"
        >
          ลงทะเบียน
        </Link>
      </p>
    </div>
  );
};
function setAuth(arg0: boolean) {
  throw new Error("Function not implemented.");
}

