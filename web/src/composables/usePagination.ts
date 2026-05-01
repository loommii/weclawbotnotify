import { ref, computed, watch } from 'vue'

const DEFAULT_PAGE_SIZE = 10

export function usePagination(initialPageSize = DEFAULT_PAGE_SIZE) {
  const page = ref(1)
  const pageSize = ref(initialPageSize)
  const total = ref(0)

  const first = computed(() => (page.value - 1) * pageSize.value)

  const totalPages = computed(() => {
    if (total.value === 0) return 0
    return Math.ceil(total.value / pageSize.value)
  })

  function onPageChange(event: { page: number; rows: number }) {
    page.value = event.page + 1
    pageSize.value = event.rows
  }

  function setTotal(value: number) {
    total.value = value
  }

  function reset() {
    page.value = 1
    total.value = 0
  }

  watch(pageSize, () => {
    page.value = 1
  })

  return {
    page,
    pageSize,
    first,
    total,
    totalPages,
    onPageChange,
    setTotal,
    reset,
  }
}
