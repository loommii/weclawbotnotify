import axios from 'axios'
import { API_BASE_URL } from '@/config'

export interface HealthResp {
  serviceName: string
  message: string
  timestamp: number
}

const raw = axios.create({
  baseURL: API_BASE_URL,
  timeout: 5000,
})

export const healthService = {
  async check(): Promise<HealthResp> {
    const resp = await raw.get('/health')
    const data = resp.data
    return {
      serviceName: data.service_name,
      message: data.message,
      timestamp: data.timestamp,
    }
  },
}
