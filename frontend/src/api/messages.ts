import { api } from './client'
import type { ApiResponse, Message, Paginated } from '@/types'

export interface GetMessagesParams {
  target_id?: number
  room_id?: number
  page?: number
  page_size?: number
}

export async function getMessages(params: GetMessagesParams): Promise<Paginated<Message>> {
  const { data } = await api.get<ApiResponse<Paginated<Message>>>('/messages', { params })
  return data.data
}

export interface SendMessagePayload {
  receiver_id?: number
  room_id?: number
  content: string
  type?: string
}

export async function sendMessage(payload: SendMessagePayload): Promise<Message> {
  const { data } = await api.post<ApiResponse<Message>>('/messages', payload)
  return data.data
}

export async function markRead(id: number): Promise<void> {
  await api.post(`/messages/${id}/read`)
}

export async function getUnreadCount(): Promise<number> {
  const { data } = await api.get<ApiResponse<{ count: number }>>('/messages/unread/count')
  return data.data.count
}
