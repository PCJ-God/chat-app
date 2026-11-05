import type { UserStatus } from './index'

/** 客户端可发送 / 服务端可推送的消息类型 */
export type WSMessageType = 'message' | 'presence' | 'ping' | 'pong' | 'join' | 'leave'

/** 通用 WebSocket 报文（客户端发送用） */
export interface WSMessage {
  type: WSMessageType | string
  from?: number
  to?: number
  room_id?: number
  content?: string
  data?: unknown
  timestamp?: number
}

/** 服务端推送的聊天消息事件 */
export interface WSChatEvent {
  type: 'message'
  id: number
  sender_id: number
  receiver_id?: number
  room_id?: number
  content: string
  created_at: string
}

/** 服务端推送的在线状态变更事件 */
export interface WSPresenceEvent {
  type: 'presence'
  user_id: number
  status: UserStatus
}

export type WSEvent = WSChatEvent | WSPresenceEvent
