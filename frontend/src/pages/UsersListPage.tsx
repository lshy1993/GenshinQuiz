import { Box, Typography, Card, CardContent, Grid, Button, Chip } from '@mui/material';
import { Link } from 'react-router-dom';
import { useGetUsers } from '../api/genshinQuizAPI';
import type { User } from '../api/dto';

export default function UsersListPage() {
  const { data: users, error } = useGetUsers();

  if (error) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography color="error">加载用户数据失败: {error.message}</Typography>
      </Box>
    );
  }

  const userList: User[] = users?.users || [];

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          用户管理
        </Typography>
        <Button
          component={Link}
          to="/users/create"
          variant="contained"
        >
          创建新用户
        </Button>
      </Box>

      <Grid container spacing={3}>
        {userList.map((user) => (
          <Grid item xs={12} md={6} lg={4} key={user.id}>
            <Card>
              <CardContent>
                <Typography variant="h6" component="h3" gutterBottom>
                  {user.username}
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  {user.email}
                </Typography>
                <Box sx={{ mb: 2 }}>
                  <Chip 
                    label={`积分: ${user.total_score || 0}`} 
                    color="primary" 
                    size="small" 
                    sx={{ mr: 1 }}
                  />
                  <Chip 
                    label={`完成: ${user.quizzes_completed || 0}次`} 
                    color="secondary" 
                    size="small" 
                  />
                </Box>
                <Typography variant="caption" display="block" sx={{ mb: 2 }}>
                  注册时间: {new Date(user.created_at).toLocaleDateString()}
                </Typography>
                <Box sx={{ display: 'flex', gap: 1 }}>
                  <Button
                    component={Link}
                    to={`/users/${user.id}`}
                    variant="contained"
                    size="small"
                  >
                    查看详情
                  </Button>
                  <Button
                    component={Link}
                    to={`/users/${user.id}/edit`}
                    variant="outlined"
                    size="small"
                  >
                    编辑
                  </Button>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {userList.length === 0 && (
        <Box sx={{ textAlign: 'center', mt: 4 }}>
          <Typography variant="h6" color="text.secondary">
            暂无用户数据
          </Typography>
        </Box>
      )}
    </Box>
  );
}