import { computed, unref, type Ref } from 'vue'
import { useQuery, useQueryClient, useMutation } from '@tanstack/vue-query'
import { applicationService } from '@/services/application.service'
import type { ListApplicationsReq } from '@/types/application'

export function useApplications(
  params: Ref<ListApplicationsReq> | (() => ListApplicationsReq) | ListApplicationsReq,
) {
  const paramsRef = computed(() => {
    if (typeof params === 'function') return params()
    return unref(params)
  })

  return useQuery({
    queryKey: ['applications', paramsRef],
    queryFn: () => applicationService.list(paramsRef.value),
    staleTime: 60_000,
  })
}

export function useDeleteApplication() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: applicationService.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['applications'] })
    },
  })
}
