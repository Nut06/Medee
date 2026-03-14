import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect, vi } from "vitest";
import { http, HttpResponse } from "msw";
import { LoginForm } from "@/components/LoginForm";
import { server } from "@/test/setup";

// Helper to render LoginForm inside a Router (required for useNavigate and Link)
const renderLoginForm = () => {
  return render(
    <MemoryRouter initialEntries={["/login"]}>
      <LoginForm />
    </MemoryRouter>,
  );
};

// Helper: properly set input value and fire change event for browser mode
// In Playwright (real browser), fireEvent.input does NOT set element.value.
// We must set the value directly on the element first, then dispatch the event
// so that react-hook-form's onChange handler sees the correct value.
function typeInto(element: Element, value: string) {
  const nativeInputValueSetter = Object.getOwnPropertyDescriptor(
    window.HTMLInputElement.prototype,
    "value",
  )?.set;
  nativeInputValueSetter?.call(element, value);
  fireEvent.change(element, { target: { value } });
}

describe("LoginForm Integration Tests", () => {
  it("should render email and password fields", () => {
    renderLoginForm();
    expect(screen.getByPlaceholderText("you@example.com")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("********")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /เข้าสู่ระบบ/i }),
    ).toBeInTheDocument();
  });

  // it("should show validation error if email is invalid", async () => {
  //   renderLoginForm();
  //   typeInto(screen.getByPlaceholderText("you@example.com"), "not-an-email");
  //   fireEvent.click(screen.getByRole("button", { name: /เข้าสู่ระบบ/i }));
  //   await waitFor(() => {
  //     expect(
  //       screen.getByText(/please enter a valid email/i),
  //     ).toBeInTheDocument();
  //   });
  // });

  it("should show validation error if password is too short", async () => {
    renderLoginForm();
    typeInto(screen.getByPlaceholderText("you@example.com"), "user@test.com");
    typeInto(screen.getByPlaceholderText("********"), "123");
    fireEvent.click(screen.getByRole("button", { name: /เข้าสู่ระบบ/i }));
    await waitFor(() => {
      expect(screen.getByText(/at least 6 characters/i)).toBeInTheDocument();
    });
  });

  // it("should display error message when login API returns 401", async () => {
  //   // Override handler for this test only - simulate wrong credentials
  //   // Also override /auth/refresh to return 401 so the retry logic does not
  //   // re-attempt the login and swallow the error.
  //   server.use(
  //     http.post("*/auth/login", () =>
  //       HttpResponse.json({ message: "Invalid credentials" }, { status: 401 }),
  //     ),
  //     http.post("*/auth/refresh", () =>
  //       HttpResponse.json({ message: "Unauthorized" }, { status: 401 }),
  //     ),
  //   );

  //   renderLoginForm();
  //   typeInto(screen.getByPlaceholderText("you@example.com"), "wrong@test.com");
  //   typeInto(screen.getByPlaceholderText("********"), "wrongpass");
  //   fireEvent.click(screen.getByRole("button", { name: /เข้าสู่ระบบ/i }));

  //   await waitFor(() => {
  //     expect(screen.getByText(/Invalid credentials/i)).toBeInTheDocument();
  //   });
  // });

  it("should call login API with correct payload on valid submit", async () => {
    const loginSpy = vi.fn();
    server.use(
      http.post("*/auth/login", async ({ request }) => {
        const body = await request.json();
        loginSpy(body);
        return HttpResponse.json({
          user: { id: "1", firstName: "Test", email: "user@test.com" },
          companies: [],
          accessToken: "fake-token",
        });
      }),
    );

    renderLoginForm();
    typeInto(screen.getByPlaceholderText("you@example.com"), "user@test.com");
    typeInto(screen.getByPlaceholderText("********"), "password123");
    fireEvent.click(screen.getByRole("button", { name: /เข้าสู่ระบบ/i }));

    await waitFor(() => {
      expect(loginSpy).toHaveBeenCalledWith({
        email: "user@test.com",
        password: "password123",
      });
    });
  });
});
