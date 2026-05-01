export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

export class ApiError extends Error {
  code: number

  constructor(code: number, msg: string) {
    super(msg)
    this.code = code
    this.name = 'ApiError'
  }
}
