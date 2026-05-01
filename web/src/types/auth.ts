export interface UserInfo {
  id: number
  username: string
  createdAt: string
}

export interface LoginReq {
  username: string
  password: string
}

export interface LoginResp {
  user: UserInfo
  token: string
  refreshToken: string
}

export interface RegisterReq {
  username: string
  password: string
}

export interface RegisterResp {
  user: UserInfo
  token: string
  refreshToken: string
}

export interface RefreshReq {
  refreshToken: string
}

export interface RefreshResp {
  token: string
  refreshToken: string
}

export interface LoginFormState {
  username: string
  password: string
  remember: boolean
}

export interface ChangePasswordReq {
  oldPassword: string
  newPassword: string
}
