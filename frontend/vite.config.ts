import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',  // 使服务器在所有网络接口上可用
    port: 3000,       // 可以设置为你需要的端口
  },
})
