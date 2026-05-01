<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-[368px]">
      <div class="mb-10 text-center">
        <i class="pi pi-user-plus mb-5 inline-block text-3xl" style="color: var(--p-primary-color)"></i>
        <h1 class="mb-1 text-xl font-semibold tracking-tight text-surface-900 dark:text-surface-0">
          注册
        </h1>
        <p class="text-sm text-muted-color">
          创建你的管理员账号
        </p>
      </div>

      <form @submit.prevent="handleRegister" class="flex flex-col gap-5">
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
            autocomplete="new-password"
            :feedback="false"
            fluid
            :invalid="!!errors.password"
            :toggle-mask="true"
          />
          <Message
            v-if="errors.password"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ errors.password }}
          </Message>
          <p class="text-xs text-muted-color">
            至少8位，包含大小写字母和数字
          </p>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="confirmPassword" class="text-xs font-medium text-muted-color">
            确认密码
          </label>
          <Password
            id="confirmPassword"
            v-model="form.confirmPassword"
            placeholder="请再次输入密码"
            autocomplete="new-password"
            :feedback="false"
            fluid
            :invalid="!!errors.confirmPassword"
            :toggle-mask="true"
          />
          <Message
            v-if="errors.confirmPassword"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ errors.confirmPassword }}
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
          label="注册"
          :loading="loading"
          class="mt-1 w-full"
          size="large"
        />

        <div class="mt-3 text-center">
          <Button
            as="router-link"
            to="/login"
            label="已有账号？登录"
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
import { registerSchema } from '@/schemas/register.schema'

interface RegisterFormState {
  username: string
  password: string
  confirmPassword: string
}

const router = useRouter()
const authStore = useAuthStore()

const form = reactive<RegisterFormState>({
  username: '',
  password: '',
  confirmPassword: '',
})

const errors = reactive<{ username?: string; password?: string; confirmPassword?: string }>({})
const authError = ref('')
const loading = ref(false)

async function handleRegister() {
  errors.username = undefined
  errors.password = undefined
  errors.confirmPassword = undefined
  authError.value = ''

  const result = registerSchema.safeParse(form)
  if (!result.success) {
    const fieldErrors = result.error.flatten().fieldErrors
    if (fieldErrors.username) errors.username = fieldErrors.username[0]
    if (fieldErrors.password) errors.password = fieldErrors.password.join('，')
    if (fieldErrors.confirmPassword) errors.confirmPassword = fieldErrors.confirmPassword[0]
    return
  }

  loading.value = true

  try {
    const resp = await authService.register({
      username: form.username,
      password: form.password,
    })
    authStore.setTokens(resp.token, resp.refreshToken)
    authStore.setUser(resp.user)
    router.replace('/')
  } catch (err: unknown) {
    if (err && typeof err === 'object' && 'code' in err) {
      const apiErr = err as { code: number; message: string }
      switch (apiErr.code) {
        case 110100:
          authError.value = '用户名或密码不能为空'
          break
        case 110101:
          authError.value = '该用户名已被注册'
          break
        case 110107:
          authError.value = '注册已关闭，系统仅支持单用户'
          break
        case 110400:
          authError.value = '密码强度不足，至少8位且包含大小写字母和数字'
          break
        default:
          authError.value = apiErr.message || '注册失败，请稍后重试'
      }
    } else {
      authError.value = '网络错误，请检查连接后重试'
    }
  } finally {
    loading.value = false
  }
}
</script>
