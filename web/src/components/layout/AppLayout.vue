<template>
  <div class="flex h-screen overflow-hidden">
    <!-- Sidebar -->
    <aside class="flex w-[220px] flex-col border-r border-surface-200 dark:border-surface-700 bg-surface-0 dark:bg-surface-900">
      <!-- Logo -->
      <div class="flex items-center gap-2.5 px-5 py-4">
        <i
          class="pi pi-shield text-xl"
          style="color: var(--p-primary-color)"
        ></i>
        <span class="text-base font-semibold tracking-tight text-surface-900 dark:text-surface-0">
          WeClawBotNotify
        </span>
      </div>

      <div class="border-b border-surface-200 dark:border-surface-700 mx-4"></div>

      <!-- Navigation -->
      <Menu :model="navItems" class="mt-2 border-none bg-transparent flex-1">
        <template #item="{ item, props }">
          <router-link
            v-if="item.route"
            v-slot="{ href, navigate, isExactActive }"
            :to="item.route"
            custom
            exact
          >
            <a
              :href="href"
              v-bind="props.action"
              @click="navigate"
              :class="[
                'flex items-center gap-3 px-4 py-2.5 mx-2 rounded-md cursor-pointer transition-all duration-150',
                isExactActive
                  ? 'bg-primary text-primary-contrast font-medium shadow-sm'
                  : 'text-surface-500 dark:text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800 hover:text-surface-700 dark:hover:text-surface-0',
              ]"
            >
              <span :class="[item.icon, 'text-lg']" />
              <span class="text-sm">{{ item.label }}</span>
            </a>
          </router-link>
        </template>
      </Menu>

      <!-- User area at bottom -->
      <div class="border-t border-surface-200 dark:border-surface-700 mx-4"></div>
      <div class="flex items-center justify-between px-5 py-3">
        <div class="flex items-center gap-2.5">
          <Avatar
            :label="userInitial"
            size="small"
            shape="circle"
            class="text-xs"
            style="background-color: var(--p-primary-100); color: var(--p-primary-700)"
          />
          <span class="text-sm text-surface-700 dark:text-surface-200 font-medium truncate max-w-[110px]">
            {{ authStore.user?.username || '用户' }}
          </span>
        </div>
        <Button
          icon="pi pi-sign-out"
          severity="secondary"
          variant="text"
          rounded
          size="small"
          @click="handleLogout"
          aria-label="退出登录"
        />
      </div>
    </aside>

    <!-- Main content -->
    <main class="flex flex-1 flex-col overflow-auto bg-surface-50 dark:bg-surface-950">
      <!-- Header bar -->
      <div class="flex items-center justify-between border-b border-surface-200 dark:border-surface-700 bg-surface-0 dark:bg-surface-900 px-6 py-3">
        <h2 class="text-sm font-medium text-surface-500 dark:text-surface-400">
          {{ pageTitle }}
        </h2>
        <div class="flex items-center gap-2">
          <Button
            :icon="isDark ? 'pi pi-sun' : 'pi pi-moon'"
            severity="secondary"
            variant="text"
            rounded
            size="small"
            aria-label="切换主题"
            @click="toggleTheme"
          />
        </div>
      </div>

      <!-- Page content -->
      <div class="flex-1 overflow-auto p-6">
        <RouterView />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter, useRoute, RouterView, RouterLink } from 'vue-router'
import Menu from 'primevue/menu'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const uiStore = useUiStore()

const isDark = computed(() => uiStore.theme === 'dark')

function toggleTheme() {
  const next = isDark.value ? 'light' : 'dark'
  uiStore.setTheme(next)
  document.documentElement.classList.toggle('p-dark', next === 'dark')
}

onMounted(() => {
  document.documentElement.classList.toggle('p-dark', isDark.value)
})

const navItems = computed(() => [
  {
    label: '仪表盘',
    icon: 'pi pi-chart-bar',
    route: '/',
  },
  {
    label: '客户端',
    icon: 'pi pi-mobile',
    route: '/clients',
  },
  {
    label: '应用',
    icon: 'pi pi-key',
    route: '/applications',
  },
  {
    label: '消息',
    icon: 'pi pi-inbox',
    route: '/messages',
  },
  {
    label: '设置',
    icon: 'pi pi-cog',
    route: '/settings',
  },
])

const userInitial = computed(() => {
  const name = authStore.user?.username
  return name ? name.charAt(0).toUpperCase() : 'U'
})

const pageTitle = computed(() => {
  const name = route.name as string
  const titles: Record<string, string> = {
    Dashboard: '仪表盘',
    Clients: '客户端管理',
    Applications: '应用管理',
    Messages: '消息管理',
    Settings: '设置',
  }
  return titles[name] || name
})

function handleLogout() {
  authStore.clearAuth()
  router.replace('/login')
}
</script>
