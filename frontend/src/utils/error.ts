import axios from 'axios'

/**
 * 从 axios 错误中提取后端返回的 message，取不到则回退到默认文案。
 */
export function extractErrorMessage(error: unknown, fallback = '操作失败'): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { message?: string } | undefined
    if (data?.message) return data.message
    return error.message || fallback
  }
  if (error instanceof Error) {
    return error.message || fallback
  }
  return fallback
}
