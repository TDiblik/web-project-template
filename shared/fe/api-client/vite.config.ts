import { defineConfig } from 'vite';
import path from 'node:path';

export default defineConfig({
  build: {
    minify: "oxc",
    lib: {
      entry: path.resolve(__dirname, 'src/index.ts'),
      name: 'ApiClient',
      fileName: 'index',
      formats: ['es']
    },
    rollupOptions: {
      external: (id) => !id.startsWith('.') && !path.isAbsolute(id)
    }
  }
});
