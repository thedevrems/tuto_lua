import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  // wasmoon ships as CommonJS, so let esbuild pre-bundle it (default behaviour)
  // for proper named-export interop. The wasm itself is loaded via a `?url`
  // import in src/lib/lua.ts, independent of how the JS is bundled.
})
