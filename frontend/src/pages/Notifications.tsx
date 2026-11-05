import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Badge, Button, Empty, List, Spin, Typography, message } from 'antd'
import dayjs from 'dayjs'
import { getNotifications, markNotificationRead } from '@/api/notifications'
import { extractErrorMessage } from '@/utils/error'

export default function Notifications() {
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['notifications', 1],
    queryFn: () => getNotifications(1, 50),
  })

  const readMutation = useMutation({
    mutationFn: (id: number) => markNotificationRead(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['notifications', 1] }),
    onError: (error) => message.error(extractErrorMessage(error, '操作失败')),
  })

  if (isLoading) {
    return (
      <div style={{ padding: 24, textAlign: 'center' }}>
        <Spin />
      </div>
    )
  }

  const items = data?.items ?? []
  if (items.length === 0) {
    return <Empty style={{ marginTop: 120 }} description="暂无通知" />
  }

  return (
    <List
      header={
        <Typography.Title level={5} style={{ margin: 0 }}>
          通知
        </Typography.Title>
      }
      dataSource={items}
      renderItem={(item) => (
        <List.Item
          actions={
            item.is_read
              ? undefined
              : [
                  <Button
                    key="read"
                    size="small"
                    type="link"
                    loading={readMutation.isPending}
                    onClick={() => readMutation.mutate(item.id)}
                  >
                    标为已读
                  </Button>,
                ]
          }
        >
          <List.Item.Meta
            title={
              <Badge dot={!item.is_read} offset={[-2, 2]}>
                {item.title}
              </Badge>
            }
            description={
              <div>
                <div>{item.content}</div>
                <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                  {dayjs(item.created_at).format('YYYY-MM-DD HH:mm')}
                </Typography.Text>
              </div>
            }
          />
        </List.Item>
      )}
    />
  )
}
