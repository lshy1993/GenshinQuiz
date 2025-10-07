import { useState } from 'react';
import './App.css';
import { Box, Button, Link } from '@mui/material';

function App() {
  const [count, setCount] = useState(0);

  return (
    <Box>
      <Link href="#">Home</Link>
    </Box>
  )
}

export default App
