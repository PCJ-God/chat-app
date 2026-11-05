import { createBrowserRouter } from 'react-router-dom'
import Login from '@/pages/Login'
import Register from '@/pages/Register'
import Chat from '@/pages/Chat'
import Notifications from '@/pages/Notifications'
import NotFound from '@/pages/NotFound'
import ChatLayout from '@/layouts/ChatLayout'
import ProtectedRoute from './ProtectedRoute'

export const router = createBrowserRouter([
  { path: '/login', element: <Login /> },
  { path: '/register', element: <Register /> },
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <ChatLayout />
      </ProtectedRoute>
    ),
    children: [
      { index: true, element: <Chat /> },
      { path: 'notifications', element: <Notifications /> },
    ],
  },
  { path: '*', element: <NotFound /> },
])
