import { useState } from 'react';
import './App.css';
import { Box, Button, Link } from '@mui/material';
import { useGetQuizzes } from './api/genshinQuizAPI';

function App() {
  const { data: quizzes, error } = useGetQuizzes();
  if(error) {
    return <div>Error: {error.message}</div>
  }
  return (
    <Box>
      <Link href="#">Home</Link>
      <pre>{JSON.stringify(quizzes, null, 2)}</pre>
    </Box>
  )
}

export default App
