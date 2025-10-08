import { useState } from 'react';
import './App.css';
import { Box, Button, Link } from '@mui/material';
import { useGetQuizzes, useGetUsers } from './api/genshinQuizAPI';

function App() {
  const { data: quizzes, error } = useGetQuizzes();
  const {data: users, error: userErr} = useGetUsers();
  if (error||userErr) {
    return <div>Error: {error?.message || userErr?.message}</div>;
  }
  return (
    <Box>
      <Link href="#">Home</Link>
      <pre>{JSON.stringify(quizzes, null, 2)}</pre>
      <pre>{JSON.stringify(users, null, 2)}</pre>
    </Box>
  );
}

export default App;
