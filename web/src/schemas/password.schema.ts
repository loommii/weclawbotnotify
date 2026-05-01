import { z } from 'zod'

export const passwordSchema = z.object({
  oldPassword: z
    .string()
    .min(1, '请输入旧密码'),
  newPassword: z.string().superRefine((val, ctx) => {
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
  confirmNewPassword: z
    .string()
    .min(1, '请再次输入新密码'),
}).refine(
  (data) => data.newPassword === data.confirmNewPassword,
  {
    message: '两次输入的新密码不一致',
    path: ['confirmNewPassword'],
  },
)

export type PasswordSchema = z.infer<typeof passwordSchema>
