import React, { useState } from 'react';

interface Props {
  initialCount?: number;
}

export const Counter: React.FC<Props> = ({ initialCount = 0 }) => {
  const [count, setCount] = useState(initialCount);

  return (
    <div className="counter">
      <h2>Count: {count}</h2>
      <button onClick={() => setCount(count + 1)}>
        Increment
      </button>
    </div>
  );
};
