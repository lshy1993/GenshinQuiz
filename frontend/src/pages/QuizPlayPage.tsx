import { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { 
  Box, 
  Typography, 
  Button, 
  Card, 
  CardContent, 
  Radio, 
  RadioGroup, 
  FormControlLabel, 
  FormControl, 
  LinearProgress,
  Alert
} from '@mui/material';
import { useGetQuiz } from '../api/genshinQuizAPI';

export default function QuizPlayPage() {
  const { id } = useParams<{ id: string }>();
  const quizId = Number(id);
  const { data: quiz, error } = useGetQuiz(quizId);
  
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [answers, setAnswers] = useState<Record<number, string>>({});
  const [isCompleted, setIsCompleted] = useState(false);
  const [score, setScore] = useState(0);

  if (error) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">加载测验失败: {error.message}</Typography>
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

  const currentQuestion = quiz.questions[currentQuestionIndex];
  const progress = ((currentQuestionIndex + 1) / quiz.questions.length) * 100;

  const handleAnswerChange = (questionId: number, answer: string) => {
    setAnswers(prev => ({
      ...prev,
      [questionId]: answer
    }));
  };

  const handleNext = () => {
    if (currentQuestionIndex < quiz.questions.length - 1) {
      setCurrentQuestionIndex(prev => prev + 1);
    } else {
      // 完成测验，计算分数
      const correctAnswers = quiz.questions.filter(question => 
        answers[question.id] === question.correct_answer
      ).length;
      const totalScore = (correctAnswers / quiz.questions.length) * 100;
      setScore(Math.round(totalScore));
      setIsCompleted(true);
    }
  };

  const handlePrevious = () => {
    if (currentQuestionIndex > 0) {
      setCurrentQuestionIndex(prev => prev - 1);
    }
  };

  if (isCompleted) {
    return (
      <Box sx={{ p: 3, maxWidth: 600, mx: 'auto' }}>
        <Card>
          <CardContent sx={{ textAlign: 'center' }}>
            <Typography variant="h4" gutterBottom>
              测验完成！
            </Typography>
            <Typography variant="h5" color="primary" gutterBottom>
              你的得分: {score}%
            </Typography>
            <Typography variant="body1" sx={{ mb: 3 }}>
              正确答案: {quiz.questions.filter(q => answers[q.id] === q.correct_answer).length} / {quiz.questions.length}
            </Typography>
            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center' }}>
              <Button
                component={Link}
                to={`/quizzes/${quiz.id}`}
                variant="contained"
              >
                查看详情
              </Button>
              <Button
                component={Link}
                to="/quizzes"
                variant="outlined"
              >
                更多测验
              </Button>
            </Box>
          </CardContent>
        </Card>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3, maxWidth: 800, mx: 'auto' }}>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" gutterBottom>
          {quiz.title}
        </Typography>
        <LinearProgress variant="determinate" value={progress} sx={{ mb: 2 }} />
        <Typography variant="body2" color="text.secondary">
          题目 {currentQuestionIndex + 1} / {quiz.questions.length}
        </Typography>
      </Box>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            {currentQuestion.question_text}
          </Typography>
          
          <FormControl component="fieldset" sx={{ width: '100%' }}>
            <RadioGroup
              value={answers[currentQuestion.id] || ''}
              onChange={(e) => handleAnswerChange(currentQuestion.id, e.target.value)}
            >
              {currentQuestion.options?.map((option, index) => (
                <FormControlLabel
                  key={index}
                  value={option}
                  control={<Radio />}
                  label={option}
                />
              ))}
            </RadioGroup>
          </FormControl>
        </CardContent>
      </Card>

      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
        <Button
          onClick={handlePrevious}
          disabled={currentQuestionIndex === 0}
          variant="outlined"
        >
          上一题
        </Button>
        
        <Button
          onClick={handleNext}
          variant="contained"
          disabled={!answers[currentQuestion.id]}
        >
          {currentQuestionIndex === quiz.questions.length - 1 ? '完成测验' : '下一题'}
        </Button>
      </Box>

      {quiz.time_limit && (
        <Alert severity="info" sx={{ mt: 2 }}>
          时间限制: {quiz.time_limit} 秒
        </Alert>
      )}
    </Box>
  );
}