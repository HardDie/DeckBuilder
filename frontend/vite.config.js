import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import eslint from 'vite-plugin-eslint'
const fs = require('fs')
const path = require('path')

// Serves docs/wiki/images/foursouls for the screenshot script (?screenshot=1).
function screenshotAssets() {
  const root = path.resolve(__dirname, '../docs/wiki/images/foursouls')
  return {
    name: 'screenshot-assets',
    configureServer(server) {
      server.middlewares.use('/screenshot-assets', (req, res, next) => {
        const rel = decodeURIComponent((req.url || '/').split('?')[0]).replace(/^\/+/, '')
        const file = path.resolve(root, rel)
        if (!file.startsWith(root + path.sep) && file !== root) {
          next()
          return
        }
        if (!fs.existsSync(file) || !fs.statSync(file).isFile()) {
          next()
          return
        }
        res.setHeader('Content-Type', 'image/png')
        fs.createReadStream(file).pipe(res)
      })
    },
  }
}

export default defineConfig({
  plugins: [vue(), eslint(), screenshotAssets()],
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:5000/',
        changeOrigin: true,
        secure: false,
      },
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    outDir: 'dist',
  },
})
