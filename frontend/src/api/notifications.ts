import { api } from './client'
import type { ApiResponse, Notification, Paginated } from '@/types'

export async function getNotifications(
  page = 1,
  pageSize = 20,
): Promise<Paginated<Notification>> {
  const { data } = await api.get<ApiResponse<Paginated<Notification>>>('/notifications', {
    params: { page, page_size: pageSize },
  })
  return data.data
}

export async function markNotificationRead(id: number): Promise<void> {
  await api.post(`/notifications/${id}/read`)
}
