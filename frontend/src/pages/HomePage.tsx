import { Box, Button, Card, CardContent, Grid, Typography } from '@mui/material';
import { Link } from 'react-router-dom';
import type { Quiz } from '../api/dto';
import { useGetQuizzes } from '../api/genshinQuizAPI';

export default function HomePage() {
  const { data: quizzes, error } = useGetQuizzes();

  if (error) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">加载测验数据失败: {error.message}</Typography>
      </Box>
    );
  }

  const featuredQuizzes: Quiz[] = quizzes?.quizzes?.slice(0, 3) || [];

  return (
    <>
      {/* 欢迎区域 */}
      <Box sx={{ textAlign: 'center', mb: 4 }}>
        <Typography variant="h3" component="h1" gutterBottom>
          原神知识测验
        </Typography>
        <Typography variant="h6" color="text.secondary" sx={{ mb: 3 }}>
          测试你对原神世界的了解程度
        </Typography>
        <Button component={Link} to="/quizzes" variant="contained" size="large" sx={{ mr: 2 }}>
          开始测验
        </Button>
        <Button component={Link} to="/about" variant="outlined" size="large">
          了解更多
        </Button>
      </Box>
      {/* 精选测验 */}
      <Typography variant="h4" component="h2" sx={{ mb: 3 }}>
        精选测验
      </Typography>
      <Grid container spacing={3}>
        {featuredQuizzes.map((quiz) => (
          <Grid item xs={12} md={4} key={quiz.id}>
            <Card>
              <CardContent>
                <Typography variant="h6" component="h3" gutterBottom>
                  {quiz.title}
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  {quiz.description}
                </Typography>
                <Typography variant="caption" display="block" sx={{ mb: 2 }}>
                  难度: {quiz.difficulty} | 分类: {quiz.category}
                </Typography>
                <Button
                  component={Link}
                  to={`/quizzes/${quiz.id}`}
                  variant="contained"
                  size="small"
                >
                  开始答题
                </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </>
  );
}
