import { createRouter, createWebHistory } from 'vue-router'
import  {getCookie} from '../utils/helpers'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/Home', name:'home', component: () => import('../views/HomeView.vue'),meta: { requiresAuth: true } },
    { path: '/login', name:'login', component: () => import('../views/LoginView.vue') },
  ]
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!getCookie("Authorization")) {
      next({ name: 'login' });
    } else {
      next();
    }
  } else {
    next();
  }
});

export default router