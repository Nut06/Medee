import { describe, it, expect } from "vitest";
import { sanitizeSearchQuery, validateUUID } from "./validation";

// ============================================================
// sanitizeSearchQuery Tests
// ============================================================

describe("sanitizeSearchQuery", () => {
  it("should trim whitespace", () => {
    expect(sanitizeSearchQuery("  hello  ")).toBe("hello");
  });

  it("should remove HTML injection characters", () => {
    expect(sanitizeSearchQuery('<script>alert("xss")</script>')).toBe(
      "scriptalert(xss)/script",
    );
  });

  it("should remove SQL injection characters", () => {
    expect(sanitizeSearchQuery("SELECT * FROM users WHERE name='admin'")).toBe(
      "SELECT * FROM users WHERE name=admin",
    );
  });

  it("should remove double quotes", () => {
    expect(sanitizeSearchQuery('hello "world"')).toBe("hello world");
  });

  it("should limit length to 100 characters", () => {
    const longString = "a".repeat(150);
    expect(sanitizeSearchQuery(longString).length).toBe(100);
  });

  it("should handle empty string", () => {
    expect(sanitizeSearchQuery("")).toBe("");
  });

  it("should handle normal search queries", () => {
    expect(sanitizeSearchQuery("golang developer")).toBe("golang developer");
  });

  it("should handle mixed dangerous characters", () => {
    expect(sanitizeSearchQuery("test<>\"'")).toBe("test");
  });
});

// ============================================================
// validateUUID Tests
// ============================================================

describe("validateUUID", () => {
  it("should validate a correct UUID v4", () => {
    expect(validateUUID("550e8400-e29b-41d4-a716-446655440000")).toBe(true);
  });

  it("should validate lowercase UUID", () => {
    expect(validateUUID("f47ac10b-58cc-4372-a567-0e02b2c3d479")).toBe(true);
  });

  it("should reject an invalid UUID", () => {
    expect(validateUUID("not-a-uuid")).toBe(false);
  });

  it("should reject empty string", () => {
    expect(validateUUID("")).toBe(false);
  });

  it("should reject UUID with wrong version (not v4)", () => {
    // UUID v1 should fail (third group starts with 1, not 4)
    expect(validateUUID("550e8400-e29b-11d4-a716-446655440000")).toBe(false);
  });

  it("should validate case-insensitive", () => {
    expect(validateUUID("550E8400-E29B-41D4-A716-446655440000")).toBe(true);
  });
});
