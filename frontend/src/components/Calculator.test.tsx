import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { Calculator } from "./Calculator";

vi.mock("../api/client", () => ({
  calculate: vi.fn(),
}));

import { calculate } from "../api/client";

describe("Calculator", () => {
  beforeEach(() => vi.mocked(calculate).mockReset());

  it("shows validation before contacting the backend", async () => {
    render(<Calculator />);
    fireEvent.click(screen.getByRole("button", { name: /add/i }));
    expect(screen.getByRole("alert")).toHaveTextContent("Enter valid numbers");
    expect(calculate).not.toHaveBeenCalled();
  });

  it("shows a successful result", async () => {
    vi.mocked(calculate).mockResolvedValue(5);
    render(<Calculator />);
    const inputs = screen.getAllByRole("textbox");
    fireEvent.change(inputs[0], { target: { value: "2" } });
    fireEvent.change(inputs[1], { target: { value: "3" } });
    fireEvent.click(screen.getByRole("button", { name: /add/i }));
    await waitFor(() => expect(screen.getByLabelText("Calculation result")).toHaveTextContent("5"));
  });

});
