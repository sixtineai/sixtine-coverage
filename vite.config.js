import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { viteSingleFile } from 'vite-plugin-singlefile'

const singleFile = process.env.SINGLE_FILE === '1'

export default defineConfig({
  plugins: [
    vue(),
    ...(singleFile ? [viteSingleFile()] : []),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  build: {
    ...(singleFile ? {
      assetsInlineLimit: Infinity,
      cssCodeSplit: false,
      rollupOptions: {
        output: { inlineDynamicImports: true },
      },
    } : {}),
  },
})
