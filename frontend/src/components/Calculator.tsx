import { useState } from "react";
import { calculate, type OperationName } from "../api/client";
import "./Calculator.css";

type BinaryOperation = Exclude<OperationName, "sqrt">;

const operationLabels: Record<OperationName, string> = {
  add: "+",
  subtract: "−",
  multiply: "×",
  divide: "÷",
  power: "xʸ",
  sqrt: "√",
  percentage: "%",
};

const binaryOperations: Array<{ operation: BinaryOperation; label: string }> = [
  { operation: "add", label: operationLabels.add },
  { operation: "subtract", label: operationLabels.subtract },
  { operation: "multiply", label: operationLabels.multiply },
  { operation: "divide", label: operationLabels.divide },
  { operation: "power", label: operationLabels.power },
  { operation: "percentage", label: operationLabels.percentage },
];

export function Calculator() {
  const [first, setFirst] = useState("");
  const [second, setSecond] = useState("");
  const [result, setResult] = useState("0");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function run(operation: OperationName) {
    const values = operation === "sqrt" ? [first] : [first, second];
    if (values.some((value) => value.trim() === "" || !Number.isFinite(Number(value)))) {
      setError("Enter valid numbers before calculating.");
      return;
    }
    setError("");
    setLoading(true);
    try {
      const value = await calculate(operation, values.map(Number));
      setResult(formatResult(value));
    } catch (calculationError) {
      setError(calculationError instanceof Error ? calculationError.message : "Calculation failed.");
    } finally {
      setLoading(false);
    }
  }

  function clear() {
    setFirst("");
    setSecond("");
    setResult("0");
    setError("");
  }

  return (
    <main className="calculator-page">
      <section className="calculator" aria-labelledby="calculator-title">
        <header className="calculator-header">
          <p className="eyebrow">SEZZLE</p>
          <h1 id="calculator-title">Calculator</h1>
        </header>
        <div className="display" aria-live="polite" aria-label="Calculation result">
          {loading ? "..." : result}
        </div>
        <div className="inputs">
          <label>
            <span>First number</span>
            <input inputMode="decimal" value={first} onChange={(event) => setFirst(event.target.value)} />
          </label>
          <label>
            <span>Second number</span>
            <input inputMode="decimal" value={second} onChange={(event) => setSecond(event.target.value)} />
          </label>
        </div>
        <div className="keypad" aria-label="Calculator operations">
          {binaryOperations.map(({ operation, label }) => (
            <button key={operation} type="button" onClick={() => void run(operation)} disabled={loading}>
              {label}
              <span>{operation}</span>
            </button>
          ))}
          <button type="button" onClick={() => void run("sqrt")} disabled={loading}>
            √<span>sqrt</span>
          </button>
          <button type="button" className="clear-button" onClick={clear} disabled={loading}>
            C<span>clear</span>
          </button>
        </div>
        {error && <p className="error-message" role="alert">{error}</p>}
      </section>
    </main>
  );
}

function formatResult(value: number): string {
  return Number.isInteger(value) ? String(value) : String(Number(value.toFixed(10)));
}
