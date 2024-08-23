import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3001, // 将端口号修改为你想要的端口号
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    }
  }
});

