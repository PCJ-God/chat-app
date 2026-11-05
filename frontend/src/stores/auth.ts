import { create } from 'zustand'
import type { User } from '@/types'
import {
  getToken,
  setToken,
  clearToken,
  getStoredUser,
  setStoredUser,
  clearStoredUser,
} from '@/utils/token'

interface AuthState {
  token: string | null
  user: User | null
  setAuth: (token: string, user: User) => void
  setUser: (user: User) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: getToken(),
  user: getStoredUser(),
  setAuth: (token, user) => {
    setToken(token)
    setStoredUser(user)
    set({ token, user })
  },
  setUser: (user) => {
    setStoredUser(user)
    set({ user })
  },
  logout: () => {
    clearToken()
    clearStoredUser()
    set({ token: null, user: null })
  },
}))
