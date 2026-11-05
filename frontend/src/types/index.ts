export type UserStatus = 'online' | 'offline' | 'away'

export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  status: UserStatus
  last_seen?: string
  created_at: string
  updated_at: string
}

export interface Message {
  id: number
  content: string
  type: string
  status: string
  created_at: string
  sender_id: number
  receiver_id?: number
  room_id?: number
  sender?: User
  receiver?: User
}

export interface RoomMember {
  id: number
  room_id: number
  user_id: number
  role: 'owner' | 'admin' | 'member'
  joined_at: string
  user?: User
}

export interface Room {
  id: number
  name: string
  description: string
  type: string
  avatar: string
  created_by: number
  created_at: string
  updated_at: string
  members?: RoomMember[]
  creator?: User
}

export interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  is_read: boolean
  created_at: string
}

export interface AuthData {
  token: string
  user: User
}

/** 后端统一响应结构 */
export interface ApiResponse<T> {
  success: boolean
  message?: string
  data: T
}

/** 分页响应 */
export interface Paginated<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}
