import api from '@/lib/axios'
import type { UserInfo, ChangePasswordReq } from '@/types/auth'

export const userService = {
  async getProfile(): Promise<UserInfo> {
    const raw = await api.get('/user/profile') as unknown as {
      user: { id: number; username: string; created_at: string }
    }
    return {
      id: raw.user.id,
      username: raw.user.username,
      createdAt: raw.user.created_at,
    }
  },

  async changePassword(data: ChangePasswordReq): Promise<void> {
    await api.post('/user/password', {
      old_password: data.oldPassword,
      new_password: data.newPassword,
    })
  },
}
