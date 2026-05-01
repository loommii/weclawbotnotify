<template>
  <Dialog
    v-model:visible="dialogVisible"
    modal
    :closable="false"
    :style="{ width: '440px' }"
    :close-on-escape="false"
  >
    <template #header>
      <span class="text-base font-semibold text-color-emphasis">
        Application 创建成功
      </span>
    </template>

    <div class="flex flex-col items-center gap-6">
      <i class="pi pi-check-circle text-4xl" style="color: var(--p-primary-color)"></i>
      <div class="text-center">
        <p class="text-sm text-muted-color">
          Token 仅显示一次，请立即复制并妥善保存
        </p>
      </div>

      <div class="w-full">
        <div class="flex items-center gap-2 rounded-border bg-surface-700 p-3">
          <code class="flex-1 truncate font-mono text-sm text-color">
            {{ token }}
          </code>
          <Button
            icon="pi pi-copy"
            text
            rounded
            size="small"
            @click="handleCopy"
          />
        </div>
      </div>

      <Button
        label="我已保存，关闭"
        @click="handleClose"
      />
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useClipboard } from '@vueuse/core'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'

const props = defineProps<{
  visible: boolean
  token: string
  appName: string
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  closed: []
}>()

const dialogVisible = ref(props.visible)

watch(() => props.visible, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

const { copy } = useClipboard()

async function handleCopy() {
  await copy(props.token)
}

function handleClose() {
  dialogVisible.value = false
  emit('closed')
}
</script>
