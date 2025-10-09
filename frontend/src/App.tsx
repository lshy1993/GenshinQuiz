import { Box } from '@mui/material';
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';

import { AdminRouteComponent } from './admin/AdminRouteComponent';
import Navigation from './components/Navigation';

import AboutPage from './pages/AboutPage';
import HomePage from './pages/HomePage';
import NotFoundPage from './pages/NotFoundPage';
import QuizDetailPage from './pages/QuizDetailPage';
import QuizPlayPage from './pages/QuizPlayPage';
import QuizzesListPage from './pages/QuizzesListPage';
import UsersListPage from './pages/UsersListPage';

function App() {
  return (
    <Router>
      <Navigation />
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          overflowY: 'auto',
          bgcolor: 'background.default',
        }}
      >
        <Routes>
          <Route path="/admin" element={<AdminRouteComponent />} />
          <Route path="/" element={<HomePage />} />
          <Route path="/quizzes" element={<QuizzesListPage />} />
          <Route path="/quizzes/:id" element={<QuizDetailPage />} />
          <Route path="/quizzes/:id/play" element={<QuizPlayPage />} />
          <Route path="/users" element={<UsersListPage />} />
          <Route path="/about" element={<AboutPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </Box>
    </Router>
  );
}

export default App;
