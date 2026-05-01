import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { applicationService } from '@/services/application.service'

export function useDeleteApplication() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: applicationService.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['applications'] })
    },
  })
}
