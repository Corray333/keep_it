import './assets/main.css'
import 'primeicons/primeicons.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import i18n from './i18n'
import ClickOutsideDirective from './clickOutside'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura';
import Tooltip from 'primevue/tooltip';
import { definePreset } from '@primevue/themes'


const app = createApp(App)
app.directive('click-outside', ClickOutsideDirective)
app.directive('tooltip', Tooltip);
app.use(createPinia())


const MyPreset = definePreset(Aura, {
    components: {
        popover: {
            background:"var(--primary-bg)",
            border: {
                color: "var(--invert-bg-50)",
                radius: "1rem"
            },
        },
        contextmenu:{
            background: "var(--primary-bg)",
            item:{
                color: "var(--invert-bg)",
                focus:{
                    color: "var(--invert-bg)",
                    background: "var(--secondary-bg)",
                },
                icon: {
                    color: "var(--invert-bg)",
                    focus:{
                        color: "var(--invert-bg)",
                    },
                }
            }
        },
        skeleton:{
                background: "var(--secondary-bg)",
                animation:{
                background: "var(--secondary-bg)",
            }
        },
        toolbar:{
            background: "var(--primary-bg)",
            border: {
                radius: "1rem",
                width:"4px"
            }
        },
        tree:{
            background: "var(--primary-bg)",
            color: "var(--invert-bg)",
            node:{
                border:{
                    radius: "999px"
                },
                color: "var(--invert-bg)",
                toggle:{
                    button:{
                        color: "var(--invert-bg)",
                        hover:{
                            color: "var(--invert-bg)",
                            background: "var(--secondary-bg)"
                        }
                    }
                },
                hover:{
                    background: "var(--secondary-bg)",
                    color: "var(--invert-bg)",
                },
                selected:{
                    background: "var(--secondary-bg)",
                    color: "var(--invert-bg)",
                }
            },
        }
    },
    treeselect: {
        background: "var(--primary-bg)",
        color: "var(--primary-bg)"
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
