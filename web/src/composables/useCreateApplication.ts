import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { applicationService } from '@/services/application.service'
import type { CreateApplicationReq } from '@/types/application'

export function useCreateApplication() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateApplicationReq) => applicationService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['applications'] })
    },
  })
}
