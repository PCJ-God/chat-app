import { create } from 'zustand'
import type { UserStatus } from '@/types'

interface PresenceState {
  // userId -> status
  statuses: Record<number, UserStatus>
  setStatus: (userId: number, status: UserStatus) => void
  reset: () => void
}

export const usePresenceStore = create<PresenceState>((set) => ({
  statuses: {},
  setStatus: (userId, status) =>
    set((state) => ({ statuses: { ...state.statuses, [userId]: status } })),
  reset: () => set({ statuses: {} }),
}))
