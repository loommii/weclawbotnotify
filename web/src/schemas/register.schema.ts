import { z } from 'zod'

export const registerSchema = z.object({
  username: z
    .string()
    .min(1, '请输入用户名'),
  password: z.string().superRefine((val, ctx) => {
    if (val.length < 8) {
      ctx.addIssue({ code: z.ZodIssueCode.custom, message: '密码至少8位' })
    }
    if (!/[A-Z]/.test(val)) {
      ctx.addIssue({ code: z.ZodIssueCode.custom, message: '密码需包含大写字母' })
    }
    if (!/[a-z]/.test(val)) {
      ctx.addIssue({ code: z.ZodIssueCode.custom, message: '密码需包含小写字母' })
    }
    if (!/[0-9]/.test(val)) {
      ctx.addIssue({ code: z.ZodIssueCode.custom, message: '密码需包含数字' })
    }
  }),
  confirmPassword: z
    .string()
    .min(1, '请再次输入密码'),
}).refine(
  (data) => data.password === data.confirmPassword,
  {
    message: '两次输入的密码不一致',
    path: ['confirmPassword'],
  },
)

export type RegisterSchema = z.infer<typeof registerSchema>
