import { getToken } from '@/utils/token'

export type RealtimeHandler = (event: unknown) => void

/**
 * WebSocket 客户端封装：自动带 token、断线指数退避重连、按事件分发。
 * 浏览器原生会自动回应服务端的 ping 控制帧，这里主要处理重连逻辑。
 */
export class RealtimeClient {
  private ws: WebSocket | null = null
  private readonly baseUrl: string
  private handlers = new Set<RealtimeHandler>()
  private reconnectAttempts = 0
  private shouldReconnect = true
  private reconnectTimer: number | null = null

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  connect(): void {
    this.shouldReconnect = true
    const token = getToken()
    if (!token) return

    const ws = new WebSocket(`${this.baseUrl}?token=${encodeURIComponent(token)}`)
    this.ws = ws

    ws.onopen = () => {
      this.reconnectAttempts = 0
    }

    ws.onmessage = (event) => {
      let data: unknown
      try {
        data = JSON.parse(event.data)
      } catch {
        return
      }
      this.handlers.forEach((handler) => handler(data))
    }

    ws.onclose = () => {
      this.ws = null
      if (this.shouldReconnect) {
        this.scheduleReconnect()
      }
    }

    ws.onerror = () => {
      ws.close()
    }
  }

  private scheduleReconnect(): void {
    const delay = Math.min(30000, 1000 * 2 ** this.reconnectAttempts)
    this.reconnectAttempts += 1
    this.reconnectTimer = window.setTimeout(() => this.connect(), delay)
  }

  send(data: unknown): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    }
  }

  on(handler: RealtimeHandler): () => void {
    this.handlers.add(handler)
    return () => this.handlers.delete(handler)
  }

  disconnect(): void {
    this.shouldReconnect = false
    if (this.reconnectTimer !== null) {
      window.clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    this.ws?.close()
    this.ws = null
  }
}

const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
export const realtime = new RealtimeClient(`${wsProtocol}://${window.location.host}/ws`)
