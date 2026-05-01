import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { TOKEN_KEY, REFRESH_TOKEN_KEY } from '@/config'
import { getItem, removeItem, setItem } from '@/lib/storage'
import { userService } from '@/services/user.service'
import type { UserInfo } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(getItem(TOKEN_KEY))
  const refreshToken = ref<string | null>(getItem(REFRESH_TOKEN_KEY))
  const user = ref<UserInfo | null>(null)

  const isAuthenticated = computed(() => !!accessToken.value)

  function setTokens(token: string, refresh: string) {
    accessToken.value = token
    refreshToken.value = refresh
    setItem(TOKEN_KEY, token)
    setItem(REFRESH_TOKEN_KEY, refresh)
  }

  function setUser(u: UserInfo) {
    user.value = u
  }

  async function fetchUser() {
    if (!accessToken.value || user.value) return
    try {
      const u = await userService.getProfile()
      user.value = u
    } catch {
      clearAuth()
    }
  }

  function clearAuth() {
    accessToken.value = null
    refreshToken.value = null
    user.value = null
    removeItem(TOKEN_KEY)
    removeItem(REFRESH_TOKEN_KEY)
  }

  return {
    accessToken,
    refreshToken,
    user,
    isAuthenticated,
    setTokens,
    setUser,
    fetchUser,
    clearAuth,
  }
})
