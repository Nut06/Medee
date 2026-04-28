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
import type { LoginRequest } from "@/utils/types/user.type";
import { useUserStore } from "@/stores/userStore";
import { useEffect } from "react";
import { initializeCsrf } from "@/services/auth.service";
const loginSchema = z.object({
  email: z
    .email({ message: "please enter a valid email" })
    .min(1, "Email is required"),
  password: z.string().min(6, "Password must be at least 6 characters."),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export const LoginForm = () => {
  const { email, password, setField, loading, error, loginUser } =
    useAuthForm();
  const { user } = useUserStore();

  const navigate = useNavigate();

  useEffect(() => {
      void initializeCsrf().catch(() => {
        
      });
    }, []);
    
  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    mode: "onSubmit",
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

    const success = await loginUser(payload);
    if (!success) {
      return;
    }

    if (user.companies && user.companies.length > 0) {
      navigate("/company");
    } else {
      navigate("/user");
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
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

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
