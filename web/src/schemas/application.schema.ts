import { z } from 'zod'

export const createApplicationSchema = z.object({
  name: z
    .string()
    .min(1, '请输入应用名称')
    .max(50, '应用名称不能超过50个字符'),
  description: z
    .string()
    .max(200, '应用描述不能超过200个字符')
    .optional()
    .or(z.literal('')),
})

export type CreateApplicationSchema = z.infer<typeof createApplicationSchema>
