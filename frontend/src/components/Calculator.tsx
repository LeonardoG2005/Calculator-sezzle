import { useState } from "react";
import { calculate, type OperationName } from "../api/client";
import "./Calculator.css";

type ButtonItem =
  | { kind: "action"; label: string; operation: OperationName }
  | { kind: "digit"; label: string }
  | { kind: "decimal"; label: string }
  | { kind: "clear"; label: string }
  | { kind: "backspace"; label: string }
  | { kind: "equals"; label: string };

const keypad: ButtonItem[][] = [
  [
    { kind: "clear", label: "C" },
    { kind: "action", label: "√", operation: "sqrt" },
    { kind: "action", label: "xʸ", operation: "power" },
    { kind: "action", label: "%", operation: "percentage" },
  ],
  [
    { kind: "digit", label: "7" },
    { kind: "digit", label: "8" },
    { kind: "digit", label: "9" },
    { kind: "action", label: "÷", operation: "divide" },
  ],
  [
    { kind: "digit", label: "4" },
    { kind: "digit", label: "5" },
    { kind: "digit", label: "6" },
    { kind: "action", label: "×", operation: "multiply" },
  ],
  [
    { kind: "digit", label: "1" },
    { kind: "digit", label: "2" },
    { kind: "digit", label: "3" },
    { kind: "action", label: "−", operation: "subtract" },
  ],
  [
    { kind: "backspace", label: "◀" },
    { kind: "digit", label: "0" },
    { kind: "decimal", label: "." },
    { kind: "action", label: "+", operation: "add" },
  ],
];

export function Calculator() {
  const [display, setDisplay] = useState("0");
  const [storedValue, setStoredValue] = useState<number | null>(null);
  const [pendingOperation, setPendingOperation] = useState<OperationName | null>(null);
  const [expression, setExpression] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  function appendDigit(digit: string) {
    setError("");
    setDisplay((current) => (current === "0" ? digit : current + digit));
  }

  function appendDecimal() {
    setError("");
    setDisplay((current) => (current.includes(".") ? current : `${current}.`));
  }

  function clear() {
    setDisplay("0");
    setStoredValue(null);
    setPendingOperation(null);
    setExpression("");
    setError("");
  }

  function backspace() {
    setError("");
    setDisplay((current) => (current.length > 1 ? current.slice(0, -1) : "0"));
  }

  function chooseOperation(operation: OperationName) {
    if (operation === "sqrt") {
      void runCalculation(operation, [Number(display)]);
      return;
    }
    if (operation === "percentage") {
      if (storedValue === null) {
        setStoredValue(Number(display));
        setPendingOperation(operation);
        setExpression(`${display} %`);
        setDisplay("0");
        return;
      }
      void runCalculation(operation, [storedValue, Number(display)]);
      return;
    }
    if (operation === "power" || operation === "add" || operation === "subtract" ||
      operation === "multiply" || operation === "divide") {
      setStoredValue(Number(display));
      setPendingOperation(operation);
      setExpression(`${display} ${symbolFor(operation)}`);
      setDisplay("0");
    }
  }

  function equals() {
    if (storedValue === null || pendingOperation === null) {
      return;
    }
    void runCalculation(pendingOperation, [storedValue, Number(display)]);
  }

  async function runCalculation(operation: OperationName, operands: number[]) {
    if (operands.some((value) => !Number.isFinite(value))) {
      setError("Enter valid numbers before calculating.");
      return;
    }
    setError("");
    setLoading(true);
    try {
      const value = await calculate(operation, operands);
      setDisplay(formatResult(value));
      setExpression(`${formatResult(operands[0])} ${symbolFor(operation)} ${operands[1] ?? ""}`);
      setStoredValue(null);
      setPendingOperation(null);
    } catch (calculationError) {
      setError(calculationError instanceof Error ? calculationError.message : "Calculation failed.");
    } finally {
      setLoading(false);
    }
  }

  function handleButton(button: ButtonItem) {
    if (loading) {
      return;
    }
    if (button.kind === "digit") appendDigit(button.label);
    if (button.kind === "decimal") appendDecimal();
    if (button.kind === "clear") clear();
    if (button.kind === "backspace") backspace();
    if (button.kind === "equals") equals();
    if (button.kind === "action") chooseOperation(button.operation);
  }

  return (
    <main className="calculator-page">
      <section className="calculator" aria-labelledby="calculator-title">
        <div className="display" aria-live="polite" aria-label="Calculation result">
          {error ? (
            <span className="display-error" role="alert">{error}</span>
          ) : (
            <>
              <span className="expression">{expression}</span>
              <strong>{loading ? "..." : display}</strong>
            </>
          )}
        </div>
        <div className="keypad" aria-label="Calculator keypad">
          {keypad.flat().map((button, index) => (
            <button
              key={`${button.kind}-${button.label}-${index}`}
              className={button.kind === "action"
                ? `operation-button${button.operation === "sqrt" ? " sqrt-button" : ""}`
                : ""}
              type="button"
              aria-label={button.kind === "equals" ? "equals" : button.label}
              onClick={() => handleButton(button)}
              disabled={loading}
            >
              {button.label}
            </button>
          ))}
          <button className="equals-button" type="button" aria-label="equals" onClick={equals} disabled={loading}>
            =
          </button>
        </div>
      </section>
    </main>
  );
}

function symbolFor(operation: OperationName): string {
  return { add: "+", subtract: "−", multiply: "×", divide: "÷", power: "xʸ", percentage: "%", sqrt: "√" }[operation];
}

function formatResult(value: number): string {
  return Number.isInteger(value) ? String(value) : String(Number(value.toFixed(10)));
}
