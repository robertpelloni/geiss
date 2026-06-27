
// Omni-Workspace Standard Telemetry
const omniLogger = {
    info: (msg) => console.log(`[${new Date().toISOString()}] - [OMNI] - INFO - ${msg}`),
    error: (msg) => console.error(`[${new Date().toISOString()}] - [OMNI] - ERROR - ${msg}`)
};

import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3006,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/system': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
})
