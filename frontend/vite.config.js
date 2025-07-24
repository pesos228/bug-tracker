
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

export default defineConfig(({ mode }) => {
  const appHost = process.env.VITE_APP_HOST || 'localhost';

  return {
    plugins: [react()],
    server: {
      host: true,
      port: 5173,
      hmr: {
        clientPort: 443,
        host: appHost,
        protocol: 'wss',
      },
      watch: {
        usePolling: true,
      },
      allowedHosts: [
        appHost
      ],
    }
  }
})