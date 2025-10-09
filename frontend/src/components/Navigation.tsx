import { AppBar, Box, Button, Toolbar, Typography } from '@mui/material';
import { Link, useLocation } from 'react-router-dom';
import { LoginedUserBox } from './LoginedUserBox';

export default function Navigation() {
  const location = useLocation();

  const isActive = (path: string) => {
    return location.pathname === path;
  };

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          <Button
            component={Link}
            to="/"
            color="inherit"
            sx={{
              textTransform: 'none',
              fontSize: '1.25rem',
              fontWeight: 600,
            }}
          >
            原神测验
          </Button>
        </Typography>

        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            component={Link}
            to="/"
            color="inherit"
            variant={isActive('/') ? 'outlined' : 'text'}
          >
            首页
          </Button>
          <Button
            component={Link}
            to="/quizzes"
            color="inherit"
            variant={isActive('/quizzes') ? 'outlined' : 'text'}
          >
            测验
          </Button>
          <Button
            component={Link}
            to="/vote"
            color="inherit"
            variant={isActive('/vote') ? 'outlined' : 'text'}
          >
            投票
          </Button>
          <Button
            component={Link}
            to="/users"
            color="inherit"
            variant={isActive('/rank') ? 'outlined' : 'text'}
          >
            排行榜
          </Button>
          <Button
            component={Link}
            to="/about"
            color="inherit"
            variant={isActive('/about') ? 'outlined' : 'text'}
          >
            关于
          </Button>
        </Box>
        <LoginedUserBox />
      </Toolbar>
    </AppBar>
  );
}
