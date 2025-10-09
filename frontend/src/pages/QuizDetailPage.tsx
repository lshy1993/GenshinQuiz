import { useParams, Link } from 'react-router-dom';
import { Box, Typography, Button, Card, CardContent, Chip, List, ListItem, ListItemText } from '@mui/material';
import { useGetQuiz } from '../api/genshinQuizAPI';

export default function QuizDetailPage() {
  const { id } = useParams<{ id: string }>();
  const quizId = Number(id);
  const { data: quiz, error } = useGetQuiz(quizId);

  if (error) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">加载测验详情失败: {error.message}</Typography>
        <Button component={Link} to="/quizzes" sx={{ mt: 2 }}>
          返回测验列表
        </Button>
      </Box>
    );
  }

  if (!quiz) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography>加载中...</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Button component={Link} to="/quizzes" sx={{ mb: 2 }}>
        ← 返回测验列表
      </Button>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h4" component="h1" gutterBottom>
            {quiz.title}
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
            {quiz.description}
          </Typography>
          
          <Box sx={{ mb: 2 }}>
            <Chip 
              label={quiz.difficulty} 
              color="primary" 
              sx={{ mr: 1 }} 
            />
            <Chip 
              label={quiz.category} 
              color="secondary" 
              sx={{ mr: 1 }} 
            />
          </Box>

          <Typography variant="body2" sx={{ mb: 2 }}>
            题目数量: {quiz.questions.length}
            {quiz.time_limit && ` | 时间限制: ${quiz.time_limit}秒`}
          </Typography>

          <Typography variant="caption" color="text.secondary">
            创建时间: {new Date(quiz.created_at).toLocaleString()}
          </Typography>
        </CardContent>
      </Card>

      <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
        <Button
          component={Link}
          to={`/quizzes/${quiz.id}/play`}
          variant="contained"
          size="large"
        >
          开始答题
        </Button>
        <Button
          component={Link}
          to={`/quizzes/${quiz.id}/edit`}
          variant="outlined"
        >
          编辑测验
        </Button>
      </Box>

      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            题目预览
          </Typography>
          <List>
            {quiz.questions.map((question, index) => (
              <ListItem key={question.id} divider>
                <ListItemText
                  primary={`${index + 1}. ${question.question_text}`}
                  secondary={`类型: ${question.question_type} | 分数: ${question.points || 10}`}
                />
              </ListItem>
            ))}
          </List>
        </CardContent>
      </Card>
    </Box>
  );
}