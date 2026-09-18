import { createApp } from 'vue'
import naive from 'naive-ui'
import { createPinia } from 'pinia'
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'
import App from './App.vue'
import router from './router'

const pinia = createPinia()

createApp(App).use(router).use(naive).use(pinia).use(Toast).mount('#app')
