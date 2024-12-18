import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsxPlugin from '@vitejs/plugin-vue-jsx'
import path from "path";
import { server } from 'typescript'
// https://vite.dev/config/
export default defineConfig({
  plugins: [
      vue(),
      vueJsxPlugin()
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },

  server:{
    proxy:{
      '/api':{
        target:'http://127.0.0.1:3000/',
        changeOrigin:true,
        rewrite:(path)=>path.replace(/^\/api/,'')
      }
    }
  }
})
