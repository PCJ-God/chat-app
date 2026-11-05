import { Layout, Avatar, Dropdown, Typography, Badge, Button } from 'antd'
import { BellOutlined, LogoutOutlined } from '@ant-design/icons'
import { Outlet, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth'
import { usePresenceStore } from '@/stores/presence'
import { useRealtime } from '@/realtime/useRealtime'
import { getUnreadCount } from '@/api/messages'
import ConversationList from '@/components/ConversationList'
import RoomList from '@/components/RoomList'

const { Header, Sider, Content } = Layout

export default function ChatLayout() {
  const navigate = useNavigate()
  const user = useAuthStore((state) => state.user)
  const logout = useAuthStore((state) => state.logout)
  const resetPresence = usePresenceStore((state) => state.reset)

  useRealtime()

  const { data: unread = 0 } = useQuery({
    queryKey: ['unread-count'],
    queryFn: getUnreadCount,
    refetchInterval: 30000,
  })

  const handleLogout = () => {
    resetPresence()
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <Layout style={{ height: '100vh' }}>
      <Header
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          background: '#fff',
          borderBottom: '1px solid #f0f0f0',
          paddingInline: 20,
        }}
      >
        <Typography.Title level={5} style={{ margin: 0 }}>
          Chat App
        </Typography.Title>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
          <Badge count={unread} size="small">
            <Button type="text" icon={<BellOutlined />} onClick={() => navigate('/notifications')} />
          </Badge>
          <Dropdown
            menu={{
              items: [{ key: 'logout', icon: <LogoutOutlined />, label: '退出登录' }],
              onClick: ({ key }) => {
                if (key === 'logout') handleLogout()
              },
            }}
          >
            <span style={{ cursor: 'pointer' }}>
              <Avatar>{user?.nickname || user?.username || 'U'}</Avatar>
            </span>
          </Dropdown>
        </div>
      </Header>
      <Layout>
        <Sider width={280} theme="light" style={{ borderRight: '1px solid #f0f0f0', overflow: 'auto' }}>
          <div style={{ padding: '12px 16px', fontWeight: 600 }}>私聊</div>
          <ConversationList />
          <div style={{ borderTop: '1px solid #f0f0f0', marginTop: 8 }}>
            <RoomList />
          </div>
        </Sider>
        <Content style={{ padding: 16, background: '#fafafa', overflow: 'hidden' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
