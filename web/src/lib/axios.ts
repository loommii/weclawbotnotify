import axios from 'axios'
import type { AxiosError, InternalAxiosRequestConfig } from 'axios'
import { API_BASE_URL, TOKEN_KEY, REFRESH_TOKEN_KEY } from '@/config'
import { getItem, removeItem, setItem } from '@/lib/storage'
import { ApiError } from '@/types/api'
import type { ApiResponse } from '@/types/api'
import type { Router } from 'vue-router'

let router: Router | null = null
let clearAuthFn: (() => void) | null = null

export function injectRouter(r: Router) {
  router = r
}

export function injectClearAuth(fn: () => void) {
  clearAuthFn = fn
}

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getItem<string>(TOKEN_KEY)
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

let isRefreshing = false
let failedQueue: Array<{
  resolve: (token: string) => void
  reject: (err: unknown) => void
}> = []

function processQueue(error: unknown, token: string | null) {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) reject(error)
    else resolve(token!)
  })
  failedQueue = []
}

function redirectToLogin() {
  removeItem(TOKEN_KEY)
  removeItem(REFRESH_TOKEN_KEY)
  clearAuthFn?.()
  if (router) {
    router.replace('/login')
  }
}

api.interceptors.response.use(
  (response) => {
    const { code, msg, data } = response.data as ApiResponse
    if (code !== 0) {
      return Promise.reject(new ApiError(code, msg))
    }
    return data as typeof response
  },
  async (error: AxiosError<ApiResponse>) => {
    const config = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    const isAuthError = error.response?.data?.code === 100001 || error.response?.status === 401

    if (isAuthError && !config._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then((token) => {
          config.headers.Authorization = `Bearer ${token}`
          return api(config)
        })
      }

      config._retry = true
      isRefreshing = true

      const storedRefreshToken = getItem<string>(REFRESH_TOKEN_KEY)

      if (!storedRefreshToken) {
        redirectToLogin()
        return Promise.reject(error)
      }

      try {
        const res = await axios.post<ApiResponse<{ token: string; refresh_token: string }>>(
          `${API_BASE_URL}/auth/refresh`,
          { refresh_token: storedRefreshToken },
        )
        const { token, refresh_token: newRefreshToken } = res.data.data
        setItem(TOKEN_KEY, token)
        setItem(REFRESH_TOKEN_KEY, newRefreshToken)
        processQueue(null, token)
        config.headers.Authorization = `Bearer ${token}`
        return api(config)
      } catch (err) {
        processQueue(err, null)
        redirectToLogin()
        return Promise.reject(err)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  },
)

export default api
