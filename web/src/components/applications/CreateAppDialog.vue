<template>
  <Dialog
    v-model:visible="dialogVisible"
    header="创建 Application"
    modal
    :style="{ width: '440px' }"
    :close-on-escape="!isCreating"
    :closable="!isCreating"
  >
    <form @submit.prevent="handleSubmit" class="flex flex-col gap-5">
      <div class="flex flex-col gap-1.5">
        <label for="app-name" class="text-xs font-medium text-muted-color">
          应用名称
        </label>
        <InputText
          id="app-name"
          v-model="form.name"
          placeholder="请输入应用名称"
          fluid
          :invalid="!!errors.name"
          :disabled="isCreating"
        />
        <Message
          v-if="errors.name"
          severity="error"
          size="small"
          variant="simple"
        >
          {{ errors.name }}
        </Message>
      </div>

      <div class="flex flex-col gap-1.5">
        <label for="app-desc" class="text-xs font-medium text-muted-color">
          应用描述
          <span class="text-muted-color/60">（可选）</span>
        </label>
        <Textarea
          id="app-desc"
          v-model="form.description"
          placeholder="请输入应用描述"
          fluid
          rows="3"
          :invalid="!!errors.description"
          :disabled="isCreating"
        />
        <Message
          v-if="errors.description"
          severity="error"
          size="small"
          variant="simple"
        >
          {{ errors.description }}
        </Message>
      </div>

      <Message
        v-if="submitError"
        severity="error"
        size="small"
        variant="simple"
      >
        {{ submitError }}
      </Message>

      <div class="flex justify-end gap-2 pt-2">
        <Button
          label="取消"
          variant="outlined"
          :disabled="isCreating"
          @click="handleCancel"
        />
        <Button
          label="创建"
          type="submit"
          :loading="isCreating"
        />
      </div>
    </form>
  </Dialog>

  <TokenDisplay
    v-if="showToken"
    v-model:visible="showToken"
    :token="createdToken"
    :app-name="form.name"
    @closed="handleTokenClosed"
  />
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import Message from 'primevue/message'
import { createApplicationSchema } from '@/schemas/application.schema'
import { useCreateApplication } from '@/composables/useCreateApplication'
import TokenDisplay from './TokenDisplay.vue'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const dialogVisible = ref(props.visible)
const isCreating = ref(false)
const submitError = ref('')
const showToken = ref(false)
const createdToken = ref('')

const form = reactive({
  name: '',
  description: '',
})

const errors = reactive<{ name?: string; description?: string }>({})

const { mutateAsync } = useCreateApplication()

watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    resetForm()
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
  if (!val) {
    resetForm()
  }
})

function resetForm() {
  form.name = ''
  form.description = ''
  errors.name = undefined
  errors.description = undefined
  submitError.value = ''
  isCreating.value = false
}

async function handleSubmit() {
  errors.name = undefined
  errors.description = undefined
  submitError.value = ''

  const result = createApplicationSchema.safeParse(form)
  if (!result.success) {
    const fieldErrors = result.error.flatten().fieldErrors
    if (fieldErrors.name) errors.name = fieldErrors.name[0]
    if (fieldErrors.description) errors.description = fieldErrors.description[0]
    return
  }

  isCreating.value = true

  try {
    const resp = await mutateAsync({
      name: form.name,
      description: form.description || undefined,
    })
    createdToken.value = resp.token
    dialogVisible.value = false
    showToken.value = true
  } catch (err: unknown) {
    if (err && typeof err === 'object' && 'message' in err) {
      submitError.value = (err as { message: string }).message
    } else {
      submitError.value = '创建失败，请稍后重试'
    }
  } finally {
    isCreating.value = false
  }
}

function handleCancel() {
  dialogVisible.value = false
}

function handleTokenClosed() {
  showToken.value = false
  emit('update:visible', false)
}
</script>
