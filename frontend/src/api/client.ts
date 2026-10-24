export type OperationName =
  | "add"
  | "subtract"
  | "multiply"
  | "divide"
  | "power"
  | "sqrt"
  | "percentage";

type CalculateResponse = { result: number };
type ErrorResponse = { error?: { message?: string; code?: string } };

export async function calculate(
  operation: OperationName,
  operands: number[],
): Promise<number> {
  let response: Response;
  try {
    response = await fetch("/api/v1/calculate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ operation, operands }),
    });
  } catch {
    throw new Error("The calculator service is unavailable. Start the backend and try again.");
  }

  let payload: CalculateResponse & ErrorResponse;
  try {
    payload = (await response.json()) as CalculateResponse & ErrorResponse;
  } catch {
    throw new Error("The calculator service returned an invalid response. Try again.");
  }
  if (!response.ok) {
    throw new Error(payload.error?.message ?? "The calculation could not be completed.");
  }
  return payload.result;
}
