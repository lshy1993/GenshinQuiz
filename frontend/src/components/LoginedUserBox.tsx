import { Avatar, Box, Button, Typography } from '@mui/material';

export const LoginedUserBox = () => {
  // 假设有 user 状态
  const user = useUser(); // 你可以用 context、redux、swr 等管理
  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, ml: 2 }}>
      {user ? (
        <>
          <Avatar src={user.avatar} alt={user.name} sx={{ width: 32, height: 32 }} />
          <Typography>{user.name}</Typography>
          <Button onClick={logout}>退出</Button>
        </>
      ) : (
        <>
          <Button component={Link} to="/login" color="inherit">
            登录
          </Button>
          <Button component={Link} to="/register" color="inherit">
            注册
          </Button>
        </>
      )}
    </Box>
  );
};
