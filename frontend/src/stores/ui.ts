import { create } from 'zustand'

interface UIState {
  // 当前选中的私聊对象 userId
  activeUserId: number | null
  // 当前选中的群聊房间 roomId
  activeRoomId: number | null
  setActiveUser: (userId: number | null) => void
  setActiveRoom: (roomId: number | null) => void
}

export const useUIStore = create<UIState>((set) => ({
  activeUserId: null,
  activeRoomId: null,
  setActiveUser: (userId) => set({ activeUserId: userId, activeRoomId: null }),
  setActiveRoom: (roomId) => set({ activeRoomId: roomId, activeUserId: null }),
}))
