import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: 'dev.test'
  },
  plugins: [react()],
  optimizeDeps: {
    include: ["chat"],
  },
})
