import api from '@/lib/axios'
import type { CreateApplicationReq, CreateApplicationResp, ListApplicationsReq, ListApplicationsResp } from '@/types/application'

export const applicationService = {
  async create(data: CreateApplicationReq): Promise<CreateApplicationResp> {
    return api.post('/application/create', data) as unknown as CreateApplicationResp
  },

  async list(params: ListApplicationsReq): Promise<ListApplicationsResp> {
    return api.get('/application/list', { params }) as unknown as ListApplicationsResp
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/application/${id}`)
  },
}
