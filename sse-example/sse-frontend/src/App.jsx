import React, { useEffect, useState } from "react";

function App() {
  const [items, setItems] = useState([]);
  const [validationResults, setValidationResults] = useState({});

  useEffect(() => {
    // Fetch initial data
    fetch("/data")
      .then((res) => res.json())
      .then(setItems);

    // Open SSE connection
    const evtSource = new EventSource("/events");
    evtSource.addEventListener("validationResult", (e) => {
      const data = e.data; // "id:1, valid:true"
      const parts = data.split(",");
      const idPart = parts[0].split(":")[1].trim();
      const validPart = parts[1].split(":")[1].trim() === "true";

      setValidationResults((old) => ({ ...old, [idPart]: validPart }));
    });

    return () => evtSource.close();
  }, []);

  return (
    <div>
      <h1>Data Items</h1>
      <ul>
        {items.map((it) => {
          const result = validationResults[it.id];
          return (
            <li
              key={it.id}
              style={{ color: result === false ? "red" : "black" }}
            >
              {it.value} {result !== undefined && (result ? "✓" : "✗")}
            </li>
          );
        })}
      </ul>
    </div>
  );
}

export default App;
