import { Box, Typography, Card, CardContent, Grid, Button, Chip } from '@mui/material';
import { Link } from 'react-router-dom';
import { useGetQuizzes } from '../api/genshinQuizAPI';
import type { Quiz } from '../api/dto';

export default function QuizzesListPage() {
  const { data: quizzes, error } = useGetQuizzes();

  if (error) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">加载测验数据失败: {error.message}</Typography>
      </Box>
    );
  }

  const quizList: Quiz[] = quizzes?.quizzes || [];

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          所有测验
        </Typography>
        <Button
          component={Link}
          to="/quizzes/create"
          variant="contained"
        >
          创建新测验
        </Button>
      </Box>

      <Grid container spacing={3}>
        {quizList.map((quiz) => (
          <Grid item xs={12} md={6} lg={4} key={quiz.id}>
            <Card>
              <CardContent>
                <Typography variant="h6" component="h3" gutterBottom>
                  {quiz.title}
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  {quiz.description}
                </Typography>
                <Box sx={{ mb: 2 }}>
                  <Chip 
                    label={quiz.difficulty} 
                    color="primary" 
                    size="small" 
                    sx={{ mr: 1 }} 
                  />
                  <Chip 
                    label={quiz.category} 
                    color="secondary" 
                    size="small" 
                  />
                </Box>
                <Typography variant="caption" display="block" sx={{ mb: 2 }}>
                  题目数量: {quiz.questions.length}
                  {quiz.time_limit && ` | 时间限制: ${quiz.time_limit}秒`}
                </Typography>
                <Box sx={{ display: 'flex', gap: 1 }}>
                  <Button
                    component={Link}
                    to={`/quizzes/${quiz.id}`}
                    variant="contained"
                    size="small"
                  >
                    查看详情
                  </Button>
                  <Button
                    component={Link}
                    to={`/quizzes/${quiz.id}/play`}
                    variant="outlined"
                    size="small"
                  >
                    开始答题
                  </Button>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {quizList.length === 0 && (
        <Box sx={{ textAlign: 'center', mt: 4 }}>
          <Typography variant="h6" color="text.secondary">
            暂无测验数据
          </Typography>
          <Button
            component={Link}
            to="/quizzes/create"
            variant="contained"
            sx={{ mt: 2 }}
          >
            创建第一个测验
          </Button>
        </Box>
      )}
    </Box>
  );
}