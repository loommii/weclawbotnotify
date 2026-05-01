<template>
  <div class="mx-auto max-w-content p-6">
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold tracking-tight text-color-emphasis">
          Application
        </h1>
        <p class="mt-1 text-sm text-muted-color">
          管理推送应用凭证
        </p>
      </div>
      <Button
        v-if="applications.length > 0 || isFirstLoad"
        label="创建 Application"
        icon="pi pi-plus"
        @click="showCreateDialog = true"
      />
    </div>

    <ApplicationTable
      :applications="applications"
      :is-first-load="isFirstLoad"
      :is-refreshing="isRefreshing"
      :is-error="!!error"
      :first="first"
      :rows="pageSize"
      :total-records="paginationTotal"
      @create="showCreateDialog = true"
      @delete="handleDelete"
      @retry="handleRetry"
      @page-change="onPageChange"
    />

    <CreateAppDialog v-model:visible="showCreateDialog" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import { useQueryClient } from '@tanstack/vue-query'
import { useApplications, useDeleteApplication } from '@/composables/useApplications'
import { usePagination } from '@/composables/usePagination'
import ApplicationTable from '@/components/applications/ApplicationTable.vue'
import CreateAppDialog from '@/components/applications/CreateAppDialog.vue'
import type { ApplicationInfo } from '@/types/application'

const confirm = useConfirm()
const toast = useToast()
const queryClient = useQueryClient()

const { page, pageSize, first, total: paginationTotal, onPageChange, setTotal, reset } = usePagination()

const { data, isLoading, error } = useApplications(() => ({
  page: page.value,
  page_size: pageSize.value,
}))

watch(data, (newData) => {
  if (newData?.page_info) {
    setTotal(newData.page_info.total)
  }
}, { immediate: true })

const { mutateAsync } = useDeleteApplication()

const showCreateDialog = ref(false)

const isFirstLoad = computed(() => isLoading.value && !data.value)

const isRefreshing = computed(() => isLoading.value && !!data.value)

const applications = computed<ApplicationInfo[]>(() => {
  if (!data.value?.list) return []
  return data.value.list.map((item) => ({
    id: item.id,
    name: item.name,
    description: item.description,
    createdAt: item.created_at,
    lastUsedAt: item.last_used_at,
  }))
})

function handleRetry() {
  queryClient.invalidateQueries({
    queryKey: ['applications'],
  })
}

function handleDelete(app: ApplicationInfo) {
  confirm.require({
    message: `确定要删除 Application「${app.name}」吗？此操作不可撤销。`,
    header: '确认删除',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: '取消',
    acceptLabel: '删除',
    rejectClass: 'p-button-outlined',
    acceptClass: 'p-button-danger',
    async accept() {
      try {
        await mutateAsync(app.id)
        toast.add({
          severity: 'success',
          summary: '已删除',
          detail: `Application「${app.name}」已删除`,
          life: 3000,
        })
        queryClient.invalidateQueries({
          queryKey: ['applications'],
        })
      } catch {
        toast.add({
          severity: 'error',
          summary: '删除失败',
          detail: '请稍后重试',
          life: 3000,
        })
      }
    },
  })
}

watch(showCreateDialog, (visible) => {
  if (!visible) {
    reset()
    queryClient.invalidateQueries({
      queryKey: ['applications'],
    })
  }
})
</script>
