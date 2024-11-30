import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import i18n from './i18n'
import ClickOutsideDirective from './clickOutside'
import PrimeVue from 'primevue/config'
import Material from '@primevue/themes/material';
import Tooltip from 'primevue/tooltip';
import { definePreset } from '@primevue/themes'


const app = createApp(App)
app.directive('click-outside', ClickOutsideDirective)
app.directive('tooltip', Tooltip);
app.use(createPinia())


const MyPreset = definePreset(Material, {
    components: {
        popover: {
            background:"var(--primary-bg)",
            border: {
                color: "var(--invert-bg-50)",
                radius: "1rem"
            },
        }
    }
});

app.use(PrimeVue, {
    theme: {
        preset: MyPreset
    }
});
app.use(router)
app.use(i18n)

app.mount('#app')
