import { api } from './client'
import type { ApiResponse, User } from '@/types'

export async function getUsers(): Promise<User[]> {
  const { data } = await api.get<ApiResponse<User[]>>('/users')
  return data.data
}

export async function getMe(): Promise<User> {
  const { data } = await api.get<ApiResponse<User>>('/users/me')
  return data.data
}
