<template>
  <div class="mx-auto max-w-content">
    <p class="mb-8 text-sm text-muted-color">
      管理你的账号安全
    </p>

    <button
      class="w-full rounded-border border border-surface-200 dark:border-surface-700 bg-surface-0 dark:bg-surface-800 p-4 flex items-center justify-between cursor-pointer transition-colors duration-150 hover:bg-surface-50 dark:hover:bg-surface-700 text-left"
      @click="showPasswordDialog = true"
    >
      <div class="flex items-center gap-3">
        <i class="pi pi-lock text-lg text-surface-500 dark:text-surface-400"></i>
        <div>
          <h3 class="text-sm font-medium text-surface-900 dark:text-surface-0">
            修改密码
          </h3>
          <p class="text-xs text-muted-color">
            定期修改密码可以提高账号安全性
          </p>
        </div>
      </div>
      <i class="pi pi-chevron-right text-surface-400 dark:text-surface-500"></i>
    </button>

    <Dialog
      v-model:visible="showPasswordDialog"
      header="修改密码"
      :modal="true"
      :style="{ width: '440px' }"
      :draggable="false"
    >
      <form @submit.prevent="handleChangePassword" class="flex flex-col gap-5">
        <div class="flex flex-col gap-1.5">
          <label for="oldPassword" class="text-xs font-medium text-muted-color">
            旧密码
          </label>
          <Password
            id="oldPassword"
            v-model="form.oldPassword"
            placeholder="请输入旧密码"
            autocomplete="current-password"
            :feedback="false"
            fluid
            :invalid="!!errors.oldPassword"
            :toggle-mask="true"
          />
          <Message
            v-if="errors.oldPassword"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ errors.oldPassword }}
          </Message>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="newPassword" class="text-xs font-medium text-muted-color">
            新密码
          </label>
          <Password
            id="newPassword"
            v-model="form.newPassword"
            placeholder="请输入新密码"
            autocomplete="new-password"
            :feedback="false"
            fluid
            :invalid="!!errors.newPassword"
            :toggle-mask="true"
          />
          <Message
            v-if="errors.newPassword"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ errors.newPassword }}
          </Message>
          <p class="text-xs text-muted-color">
            至少8位，包含大小写字母和数字
          </p>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="confirmNewPassword" class="text-xs font-medium text-muted-color">
            确认新密码
          </label>
          <Password
            id="confirmNewPassword"
            v-model="form.confirmNewPassword"
            placeholder="请再次输入新密码"
            autocomplete="new-password"
            :feedback="false"
            fluid
            :invalid="!!errors.confirmNewPassword"
            :toggle-mask="true"
          />
          <Message
            v-if="errors.confirmNewPassword"
            severity="error"
            size="small"
            variant="simple"
          >
            {{ errors.confirmNewPassword }}
          </Message>
        </div>

        <Message
          v-if="formError"
          severity="error"
          size="small"
          variant="simple"
        >
          {{ formError }}
        </Message>
      </form>

      <template #footer>
        <Button
          label="取消"
          severity="secondary"
          text
          @click="showPasswordDialog = false"
        />
        <Button
          label="修改密码"
          :loading="loading"
          @click="handleChangePassword"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Message from 'primevue/message'
import Dialog from 'primevue/dialog'
import { useAuthStore } from '@/stores/auth'
import { userService } from '@/services/user.service'
import { passwordSchema } from '@/schemas/password.schema'

interface ChangePasswordFormState {
  oldPassword: string
  newPassword: string
  confirmNewPassword: string
}

const router = useRouter()
const queryClient = useQueryClient()
const authStore = useAuthStore()

const showPasswordDialog = ref(false)

const form = reactive<ChangePasswordFormState>({
  oldPassword: '',
  newPassword: '',
  confirmNewPassword: '',
})

const errors = reactive<{ oldPassword?: string; newPassword?: string; confirmNewPassword?: string }>({})
const formError = ref('')
const loading = ref(false)

async function handleChangePassword() {
  errors.oldPassword = undefined
  errors.newPassword = undefined
  errors.confirmNewPassword = undefined
  formError.value = ''

  const result = passwordSchema.safeParse(form)
  if (!result.success) {
    const fieldErrors = result.error.flatten().fieldErrors
    if (fieldErrors.oldPassword) errors.oldPassword = fieldErrors.oldPassword[0]
    if (fieldErrors.newPassword) errors.newPassword = fieldErrors.newPassword.join('，')
    if (fieldErrors.confirmNewPassword) errors.confirmNewPassword = fieldErrors.confirmNewPassword[0]
    return
  }

  loading.value = true

  try {
    await userService.changePassword({
      oldPassword: form.oldPassword,
      newPassword: form.newPassword,
    })

    showPasswordDialog.value = false
    authStore.clearAuth()
    queryClient.clear()
    router.replace('/login')
  } catch (err: unknown) {
    if (err && typeof err === 'object' && 'code' in err) {
      const apiErr = err as { code: number; message: string }
      switch (apiErr.code) {
        case 110202:
          formError.value = '旧密码错误'
          break
        case 110400:
          formError.value = '密码强度不足，至少8位且包含大小写字母和数字'
          break
        default:
          formError.value = apiErr.message || '修改密码失败，请稍后重试'
      }
    } else {
      formError.value = '网络错误，请检查连接后重试'
    }
  } finally {
    loading.value = false
  }
}
</script>
