import { afterEach, describe, expect, it, vi } from "vitest";
import { calculate } from "./client";

describe("calculate", () => {
  afterEach(() => vi.restoreAllMocks());

  it("returns a successful result", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(JSON.stringify({ result: 5 }), { status: 200 }),
    );
    await expect(calculate("add", [2, 3])).resolves.toBe(5);
  });

  it("returns a friendly message when the backend is unavailable", async () => {
    vi.spyOn(globalThis, "fetch").mockRejectedValue(new TypeError("offline"));
    await expect(calculate("add", [2, 3])).rejects.toThrow("service is unavailable");
  });

  it("uses the backend error message", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(JSON.stringify({ error: { message: "division by zero" } }), { status: 422 }),
    );
    await expect(calculate("divide", [2, 0])).rejects.toThrow("division by zero");
  });

  it("handles a non-JSON backend response", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response("<html>server error</html>", { status: 502 }),
    );
    await expect(calculate("add", [2, 3])).rejects.toThrow("invalid response");
  });
});
