import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { TOKEN_KEY } from '@/config'
import { getItem } from '@/lib/storage'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
      },
      {
        path: 'clients',
        name: 'Clients',
        component: () => import('@/views/clients/ClientsView.vue'),
      },
      {
        path: 'applications',
        name: 'Applications',
        component: () => import('@/views/applications/ApplicationsView.vue'),
      },
      {
        path: 'messages',
        name: 'Messages',
        component: () => import('@/views/messages/MessagesView.vue'),
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/settings/SettingsView.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const token = getItem<string>(TOKEN_KEY)
  const requiresAuth = to.meta.requiresAuth

  if (requiresAuth && !token) {
    return { name: 'Login' }
  }

  if (requiresAuth && token) {
    const authStore = useAuthStore()
    if (!authStore.user) {
      await authStore.fetchUser()
      if (!authStore.user) {
        return { name: 'Login' }
      }
    }
  }

  if (!requiresAuth && token && (to.name === 'Login' || to.name === 'Register')) {
    return { name: 'Dashboard' }
  }
})

export default router
