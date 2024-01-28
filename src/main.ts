import './global.css'
import App from "./routes/landing.svelte";

const app = new App({
  target: document.getElementById('app')!,
})

export default app