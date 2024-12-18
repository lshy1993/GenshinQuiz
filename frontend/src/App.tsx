import { useState } from 'react';
import './App.css';
import { Box, Button, Link } from '@mui/material';

function App() {
  const [count, setCount] = useState(0);

  return (
    <Box>
      <Link href="#">Home</Link>
      <Link href="localhost:3082/api-docs">API</Link>
      <Button onClick={() => setCount((count) => count + 1)}>
        count is {count}
      </Button>
    </Box>
  )
}

export default App
