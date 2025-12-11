import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': 'http://localhost:8080',
      '/auth': 'http://localhost:8080',
      '/preflight': 'http://localhost:8080',
      '/eventupload': 'http://localhost:8080',
      '/ruledownload': 'http://localhost:8080',
      '/postflight': 'http://localhost:8080',
    },
  },
})
