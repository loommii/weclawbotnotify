import api from '@/lib/axios'
import type { LoginReq, LoginResp, RegisterReq, RegisterResp, RefreshReq, RefreshResp } from '@/types/auth'
import { TOKEN_KEY, REFRESH_TOKEN_KEY } from '@/config'
import { setItem } from '@/lib/storage'

function persistAuth(data: { token: string; refreshToken: string }) {
  setItem(TOKEN_KEY, data.token)
  setItem(REFRESH_TOKEN_KEY, data.refreshToken)
}

function mapLoginResponse(raw: {
  token: string
  refresh_token: string
  user: { id: number; username: string; created_at: string }
}): LoginResp {
  return {
    token: raw.token,
    refreshToken: raw.refresh_token,
    user: {
      id: raw.user.id,
      username: raw.user.username,
      createdAt: raw.user.created_at,
    },
  }
}

function mapRegisterResponse(raw: {
  token: string
  refresh_token: string
  user: { id: number; username: string; created_at: string }
}): RegisterResp {
  return {
    token: raw.token,
    refreshToken: raw.refresh_token,
    user: {
      id: raw.user.id,
      username: raw.user.username,
      createdAt: raw.user.created_at,
    },
  }
}

function mapRefreshResponse(raw: { token: string; refresh_token: string }): RefreshResp {
  return { token: raw.token, refreshToken: raw.refresh_token }
}

export const authService = {
  async login(data: LoginReq): Promise<LoginResp> {
    const raw = await api.post('/auth/login', {
      username: data.username,
      password: data.password,
    }) as unknown as {
      token: string
      refresh_token: string
      user: { id: number; username: string; created_at: string }
    }
    const resp = mapLoginResponse(raw)
    persistAuth(resp)
    return resp
  },

  async register(data: RegisterReq): Promise<RegisterResp> {
    const raw = await api.post('/auth/register', {
      username: data.username,
      password: data.password,
    }) as unknown as {
      token: string
      refresh_token: string
      user: { id: number; username: string; created_at: string }
    }
    const resp = mapRegisterResponse(raw)
    persistAuth(resp)
    return resp
  },

  async refresh(data: RefreshReq): Promise<RefreshResp> {
    const raw = await api.post('/auth/refresh', {
      refresh_token: data.refreshToken,
    }) as unknown as {
      token: string
      refresh_token: string
    }
    const resp = mapRefreshResponse(raw)
    persistAuth(resp)
    return resp
  },
}
