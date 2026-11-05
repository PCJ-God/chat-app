import { useQuery } from '@tanstack/react-query'
import { Avatar, Badge, Empty, List, Spin } from 'antd'
import { getUsers } from '@/api/users'
import { useAuthStore } from '@/stores/auth'
import { usePresenceStore } from '@/stores/presence'
import { useUIStore } from '@/stores/ui'
import type { UserStatus } from '@/types'

const statusText: Record<UserStatus, string> = {
  online: '在线',
  away: '离开',
  offline: '离线',
}

const statusColor: Record<UserStatus, string> = {
  online: '#52c41a',
  away: '#faad14',
  offline: '#bfbfbf',
}

export default function ConversationList() {
  const selfId = useAuthStore((state) => state.user?.id)
  const statuses = usePresenceStore((state) => state.statuses)
  const activeUserId = useUIStore((state) => state.activeUserId)
  const setActiveUser = useUIStore((state) => state.setActiveUser)

  const { data, isLoading } = useQuery({ queryKey: ['users'], queryFn: getUsers })

  const users = (data ?? []).filter((user) => user.id !== selfId)

  if (isLoading) {
    return (
      <div style={{ padding: 24, textAlign: 'center' }}>
        <Spin />
      </div>
    )
  }

  if (users.length === 0) {
    return <Empty description="暂无其他用户" style={{ marginTop: 40 }} />
  }

  return (
    <List
      dataSource={users}
      renderItem={(user) => {
        const status: UserStatus = statuses[user.id] ?? user.status ?? 'offline'
        return (
          <List.Item
            onClick={() => setActiveUser(user.id)}
            style={{
              cursor: 'pointer',
              paddingInline: 16,
              background: activeUserId === user.id ? '#e6f4ff' : undefined,
            }}
          >
            <List.Item.Meta
              avatar={
                <Badge dot color={statusColor[status]} offset={[-2, 32]}>
                  <Avatar>{(user.nickname || user.username || 'U')[0]}</Avatar>
                </Badge>
              }
              title={user.nickname || user.username}
              description={
                <span style={{ color: status === 'online' ? '#52c41a' : '#999' }}>
                  {statusText[status]}
                </span>
              }
            />
          </List.Item>
        )
      }}
    />
  )
}
