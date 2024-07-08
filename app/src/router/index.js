import { createRouter, createWebHistory } from 'vue-router'
import  {getCookie} from '../utils/helpers'
import { useStore } from 'vuex'



const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/Home', name:'home', component: () => import('../views/HomeView.vue'),meta: { requiresAuth: true } },
    { path: '/login', name:'login', component: () => import('../views/LoginView.vue'), meta:{hideNav: true } },
    { path: '/notes/new-note', name:'new note', component: () => import('../views/notes/NewNote.vue'), meta: { requiresAuth: true, hideNav: true  } },
  ]
})

// router.beforeEach((to, from, next) => {
// const store = useStore()
// if (to.matched.some(record => record.meta.requiresAuth)) {
//     if (!store.state.AccessToken) {
//       next({ name: 'login' });
//     } else {
//       next();
//     }
//   } else {
//     next();
//   }
// })

export default router