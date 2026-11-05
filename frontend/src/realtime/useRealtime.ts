import { useEffect } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { realtime } from './client'
import { useAuthStore } from '@/stores/auth'
import { usePresenceStore } from '@/stores/presence'
import type { WSChatEvent, WSPresenceEvent } from '@/types/ws'

/**
 * 建立并维持实时连接：
 * - presence 事件 -> 更新在线状态
 * - message 事件 -> 刷新对应会话（私聊/群聊）的消息缓存和未读计数
 */
export function useRealtime(): void {
  const token = useAuthStore((state) => state.token)
  const selfId = useAuthStore((state) => state.user?.id)
  const setStatus = usePresenceStore((state) => state.setStatus)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!token) return

    realtime.connect()

    const off = realtime.on((raw) => {
      const event = raw as { type?: string }

      if (event.type === 'presence') {
        const presence = raw as WSPresenceEvent
        setStatus(presence.user_id, presence.status)
        return
      }

      if (event.type === 'message') {
        const msg = raw as WSChatEvent
        if (msg.room_id != null) {
          // 群聊消息
          queryClient.invalidateQueries({ queryKey: ['messages', 'room', msg.room_id] })
        } else if (msg.receiver_id === selfId) {
          // 发给我的私聊消息：刷新与发送者的会话 + 未读数
          queryClient.invalidateQueries({ queryKey: ['messages', 'dm', msg.sender_id] })
          queryClient.invalidateQueries({ queryKey: ['unread-count'] })
        }
      }
    })

    return () => {
      off()
      realtime.disconnect()
    }
  }, [token, selfId, setStatus, queryClient])
}
