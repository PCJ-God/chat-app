import { api } from './client'
import type { ApiResponse, Room } from '@/types'

export async function getMyRooms(): Promise<Room[]> {
  const { data } = await api.get<ApiResponse<Room[]>>('/rooms')
  return data.data
}

export interface CreateRoomPayload {
  name: string
  description?: string
  type?: string
  avatar?: string
}

export async function createRoom(payload: CreateRoomPayload): Promise<Room> {
  const { data } = await api.post<ApiResponse<Room>>('/rooms', payload)
  return data.data
}

export async function getRoom(id: number): Promise<Room> {
  const { data } = await api.get<ApiResponse<Room>>(`/rooms/${id}`)
  return data.data
}

export async function joinRoom(id: number): Promise<void> {
  await api.post(`/rooms/${id}/join`)
}

export async function leaveRoom(id: number): Promise<void> {
  await api.post(`/rooms/${id}/leave`)
}

export async function sendRoomMessage(id: number, content: string): Promise<void> {
  await api.post(`/rooms/${id}/messages`, { content })
}
