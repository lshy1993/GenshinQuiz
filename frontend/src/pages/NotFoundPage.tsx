import { Box, Typography, Button } from '@mui/material';
import { Link } from 'react-router-dom';

export default function NotFoundPage() {
  return (
    <Box sx={{ 
      p: 3, 
      textAlign: 'center', 
      minHeight: '60vh', 
      display: 'flex', 
      flexDirection: 'column', 
      justifyContent: 'center' 
    }}>
      <Typography variant="h1" color="primary" gutterBottom>
        404
      </Typography>
      <Typography variant="h4" gutterBottom>
        页面未找到
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        抱歉，您访问的页面不存在。
      </Typography>
      <Box>
        <Button
          component={Link}
          to="/"
          variant="contained"
          size="large"
          sx={{ mr: 2 }}
        >
          返回首页
        </Button>
        <Button
          component={Link}
          to="/quizzes"
          variant="outlined"
          size="large"
        >
          浏览测验
        </Button>
      </Box>
    </Box>
  );
}