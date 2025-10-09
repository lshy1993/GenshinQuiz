import { Box, Typography, Button } from '@mui/material';
import { Link } from 'react-router-dom';

export default function AboutPage() {
  return (
    <Box sx={{ p: 3, maxWidth: 800, mx: 'auto' }}>
      <Typography variant="h4" component="h1" gutterBottom>
        关于原神知识测验
      </Typography>
      
      <Typography variant="body1" paragraph>
        欢迎来到原神知识测验！这是一个专门为原神玩家设计的知识测试平台，
        让你能够测试自己对提瓦特大陆的了解程度。
      </Typography>

      <Typography variant="h6" gutterBottom sx={{ mt: 3 }}>
        功能特色
      </Typography>
      <Typography variant="body1" paragraph>
        • 丰富的题目类型：选择题、多选题、判断题等<br />
        • 多种难度等级：从入门到专家<br />
        • 分类测验：角色、武器、剧情、世界观等<br />
        • 实时评分系统<br />
        • 详细的答题分析
      </Typography>

      <Typography variant="h6" gutterBottom sx={{ mt: 3 }}>
        如何开始
      </Typography>
      <Typography variant="body1" paragraph>
        1. 浏览测验列表，选择你感兴趣的主题<br />
        2. 阅读测验详情和规则<br />
        3. 开始答题，挑战你的原神知识<br />
        4. 查看你的成绩和正确答案解析
      </Typography>

      <Box sx={{ mt: 4 }}>
        <Button
          component={Link}
          to="/quizzes"
          variant="contained"
          size="large"
          sx={{ mr: 2 }}
        >
          开始测验
        </Button>
        <Button
          component={Link}
          to="/"
          variant="outlined"
          size="large"
        >
          返回首页
        </Button>
      </Box>
    </Box>
  );
}