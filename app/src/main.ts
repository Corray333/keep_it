import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import i18n from './i18n'
import ClickOutsideDirective from './clickOutside'

const app = createApp(App)
app.directive('click-outside', ClickOutsideDirective)
app.use(createPinia())
app.use(router)
app.use(i18n)

app.mount('#app')
