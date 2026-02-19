import { describe, it, expect } from "vitest";
import { calculateProfileCompleteness, getMissingFields } from "./profileUtils";
import type { User } from "./types/user.type";

// ============================================================
// Helper to create test user
// ============================================================

const createFullUser = (overrides: Partial<User> = {}): User =>
  ({
    id: "123",
    firstName: "John",
    lastName: "Doe",
    email: "john@example.com",
    bio: "A developer",
    tagline: "Full Stack Dev",
    avatarURL: "https://example.com/avatar.jpg",
    phoneNumber: "0812345678",
    linkedInURL: "https://linkedin.com/in/john",
    githubURL: "https://github.com/john",
    websiteURL: "https://john.dev",
    skills: [{ id: "1", name: "Go" }],
    experiences: [
      {
        id: "1",
        position: "Engineer",
        companyName: "Corp",
        startDate: "2024-01-01",
      },
    ],
    projects: [{ id: "1", title: "Project" }],
    ...overrides,
  }) as User;

// ============================================================
// calculateProfileCompleteness Tests
// ============================================================

describe("calculateProfileCompleteness", () => {
  it("should return 0 for null user", () => {
    expect(calculateProfileCompleteness(null)).toBe(0);
  });

  it("should return 100 for fully completed profile", () => {
    const user = createFullUser();
    expect(calculateProfileCompleteness(user)).toBe(100);
  });

  it("should return partial percentage for incomplete profile", () => {
    const user = createFullUser({
      bio: undefined,
      tagline: undefined,
      avatarURL: undefined,
      skills: [],
      experiences: [],
      projects: [],
    });

    const result = calculateProfileCompleteness(user);
    expect(result).toBeGreaterThan(0);
    expect(result).toBeLessThan(100);
  });

  it("should count array fields (skills, experiences, projects)", () => {
    const withArrays = createFullUser();
    const withoutArrays = createFullUser({
      skills: [],
      experiences: [],
      projects: [],
    });

    expect(calculateProfileCompleteness(withArrays)).toBeGreaterThan(
      calculateProfileCompleteness(withoutArrays),
    );
  });

  it("should return 0 for user with no filled fields", () => {
    const emptyUser = {
      id: "123",
      firstName: "",
      lastName: "",
      email: "",
    } as User;

    expect(calculateProfileCompleteness(emptyUser)).toBe(0);
  });
});

// ============================================================
// getMissingFields Tests
// ============================================================

describe("getMissingFields", () => {
  it("should return empty array for null user", () => {
    expect(getMissingFields(null)).toEqual([]);
  });

  it("should return empty array for fully complete user", () => {
    const user = createFullUser();
    expect(getMissingFields(user)).toEqual([]);
  });

  it("should list missing basic fields", () => {
    const user = createFullUser({
      bio: undefined,
      avatarURL: undefined,
    });

    const missing = getMissingFields(user);
    expect(missing).toContain("Bio");
    expect(missing).toContain("Profile Picture");
  });

  it("should list missing array fields", () => {
    const user = createFullUser({
      skills: [],
      experiences: [],
      projects: [],
    });

    const missing = getMissingFields(user);
    expect(missing).toContain("Skills");
    expect(missing).toContain("Work Experience");
    expect(missing).toContain("Projects");
  });

  it("should treat empty strings as missing", () => {
    const user = createFullUser({
      firstName: "",
      email: "",
    });

    const missing = getMissingFields(user);
    expect(missing).toContain("First Name");
    expect(missing).toContain("Email");
  });
});
