import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { Calculator } from "./Calculator";

vi.mock("../api/client", () => ({
  calculate: vi.fn(),
}));

import { calculate } from "../api/client";

describe("Calculator", () => {
  beforeEach(() => vi.mocked(calculate).mockReset());

  it("builds a numeric expression and sends it to the backend", async () => {
    vi.mocked(calculate).mockResolvedValue(5);
    render(<Calculator />);

    fireEvent.click(screen.getByRole("button", { name: "2" }));
    fireEvent.click(screen.getByRole("button", { name: "+" }));
    fireEvent.click(screen.getByRole("button", { name: "3" }));
    fireEvent.click(screen.getByRole("button", { name: "equals" }));

    await waitFor(() => expect(calculate).toHaveBeenCalledWith("add", [2, 3]));
    expect(screen.getByLabelText("Calculation result")).toHaveTextContent("5");
  });

  it("supports square root from the keypad", async () => {
    vi.mocked(calculate).mockResolvedValue(4);
    render(<Calculator />);

    fireEvent.click(screen.getByRole("button", { name: "1" }));
    fireEvent.click(screen.getByRole("button", { name: "6" }));
    fireEvent.click(screen.getByRole("button", { name: "√" }));

    await waitFor(() => expect(calculate).toHaveBeenCalledWith("sqrt", [16]));
  });

});
