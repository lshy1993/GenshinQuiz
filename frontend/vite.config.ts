import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import svgr from 'vite-plugin-svgr'
import viteTsconfigPaths from 'vite-tsconfig-paths'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), viteTsconfigPaths(), svgr()],
  resolve: {
    preserveSymlinks: true,
    alias: {
      '@': '/src',
    }
  },
  server: {
    port: 3000,       // 可以设置为你需要的端口
    host: true,  // 使服务器在所有网络接口上可用
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
})
