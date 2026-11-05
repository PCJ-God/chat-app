import { api } from './client'
import type { ApiResponse, AuthData, User } from '@/types'

export async function login(username: string, password: string): Promise<AuthData> {
  const { data } = await api.post<ApiResponse<AuthData>>('/auth/login', { username, password })
  return data.data
}

export interface RegisterPayload {
  username: string
  email: string
  password: string
  nickname?: string
}

export async function register(payload: RegisterPayload): Promise<User> {
  const { data } = await api.post<ApiResponse<User>>('/auth/register', payload)
  return data.data
}
