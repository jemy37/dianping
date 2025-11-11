import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Vite dev server configuration with API proxy.
// The backend from this repo serves APIs under `/api` on port 8080.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      // Proxy API calls to backend during development to avoid CORS.
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
        // No path rewrite is needed as backend already serves under /api
      }
    }
  }
})

