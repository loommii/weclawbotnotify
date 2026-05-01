<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-[368px]">
      <div class="mb-10 text-center">
        <i class="pi pi-shield mb-5 inline-block text-3xl" style="color: var(--p-primary-color)"></i>
        <h1 class="mb-1 text-xl font-semibold tracking-tight text-surface-900 dark:text-surface-0">
          登录
        </h1>
        <p class="text-sm text-muted-color">
          欢迎回来，请登录你的账号
        </p>
      </div>

      <form @submit.prevent="handleLogin" class="flex flex-col gap-5">
        <div class="flex flex-col gap-1.5">
          <label for="username" class="text-xs font-medium text-muted-color">
            用户名
          </label>
          <InputText
            id="username"
            v-model="form.username"
            type="text"
            placeholder="请输入用户名"
            autocomplete="username"
            fluid
            :invalid="!!errors.username"
          />
          <Message
            v-if="errors.username"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ errors.username }}
          </Message>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="password" class="text-xs font-medium text-muted-color">
            密码
          </label>
          <Password
            id="password"
            v-model="form.password"
            placeholder="请输入密码"
            autocomplete="current-password"
            :feedback="false"
            fluid
            :invalid="!!errors.password"
          />
          <Message
            v-if="errors.password"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ errors.password }}
          </Message>
        </div>

        <Message
          v-if="authError"
          severity="error"
          size="small"
          variant="simple"
        >
          {{ authError }}
        </Message>

        <Button
          type="submit"
          label="登录"
          :loading="loading"
          class="mt-1 w-full"
          size="large"
        />

        <div class="mt-3 text-center">
          <Button
            as="router-link"
            to="/register"
            label="还没有账号？注册"
            variant="link"
            size="small"
            class="text-sm"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Message from 'primevue/message'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/services/auth.service'
import { loginSchema } from '@/schemas/login.schema'
import type { LoginFormState } from '@/types/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive<LoginFormState>({
  username: '',
  password: '',
  remember: false,
})

const errors = reactive<{ username?: string; password?: string }>({})
const authError = ref('')
const loading = ref(false)

async function handleLogin() {
  errors.username = undefined
  errors.password = undefined
  authError.value = ''

  const result = loginSchema.safeParse(form)
  if (!result.success) {
    const fieldErrors = result.error.flatten().fieldErrors
    if (fieldErrors.username) errors.username = fieldErrors.username[0]
    if (fieldErrors.password) errors.password = fieldErrors.password[0]
    return
  }

  loading.value = true

  try {
    const resp = await authService.login({
      username: form.username,
      password: form.password,
    })
    authStore.setTokens(resp.token, resp.refreshToken)
    authStore.setUser(resp.user)
    router.replace('/')
  } catch (err: unknown) {
    if (err && typeof err === 'object' && 'code' in err) {
      const apiErr = err as { code: number; message: string }
      if (apiErr.code >= 110200 && apiErr.code < 110300) {
        authError.value = '用户名或密码错误'
      } else {
        authError.value = apiErr.message || '登录失败，请稍后重试'
      }
    } else {
      authError.value = '网络错误，请检查连接后重试'
    }
  } finally {
    loading.value = false
  }
}
</script>
