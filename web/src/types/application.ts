export interface ApplicationInfo {
  id: number
  name: string
  description: string
  token?: string
  createdAt: string
  lastUsedAt: string | null
}

export interface ApplicationRaw {
  id: number
  name: string
  description: string
  created_at: string
  last_used_at: string | null
}

export interface CreateApplicationReq {
  name: string
  description?: string
}

export interface CreateApplicationResp {
  id: number
  name: string
  token: string
  created_at: string
}

export interface PageInfo {
  page: number
  pageSize: number
  total: number
  totalPages: number
}

export interface ListApplicationsReq {
  page: number
  page_size: number
}

export interface ListApplicationsResp {
  list: ApplicationRaw[]
  page_info: PageInfo
}
