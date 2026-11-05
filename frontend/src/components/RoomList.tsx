import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Avatar, Button, Empty, Form, Input, List, Modal, Space, Spin, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { createRoom, getMyRooms, joinRoom } from '@/api/rooms'
import { useUIStore } from '@/stores/ui'
import { extractErrorMessage } from '@/utils/error'

interface CreateForm {
  name: string
  description?: string
}

interface JoinForm {
  roomId: number
}

export default function RoomList() {
  const activeRoomId = useUIStore((state) => state.activeRoomId)
  const setActiveRoom = useUIStore((state) => state.setActiveRoom)
  const queryClient = useQueryClient()

  const [createOpen, setCreateOpen] = useState(false)
  const [joinOpen, setJoinOpen] = useState(false)
  const [createForm] = Form.useForm<CreateForm>()
  const [joinForm] = Form.useForm<JoinForm>()

  const { data, isLoading } = useQuery({ queryKey: ['rooms'], queryFn: getMyRooms })

  const createMutation = useMutation({
    mutationFn: createRoom,
    onSuccess: (room) => {
      message.success('群聊创建成功')
      setCreateOpen(false)
      createForm.resetFields()
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
      setActiveRoom(room.id)
    },
    onError: (error) => message.error(extractErrorMessage(error, '创建失败')),
  })

  const joinMutation = useMutation({
    mutationFn: (id: number) => joinRoom(id),
    onSuccess: (_data, id) => {
      message.success('加入成功')
      setJoinOpen(false)
      joinForm.resetFields()
      queryClient.invalidateQueries({ queryKey: ['rooms'] })
      setActiveRoom(id)
    },
    onError: (error) => message.error(extractErrorMessage(error, '加入失败')),
  })

  const rooms = data ?? []

  return (
    <div>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: '8px 16px',
        }}
      >
        <span style={{ fontWeight: 600 }}>群聊</span>
        <Space size={0}>
          <Button size="small" type="text" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>
            创建
          </Button>
          <Button size="small" type="text" onClick={() => setJoinOpen(true)}>
            加入
          </Button>
        </Space>
      </div>

      {isLoading ? (
        <div style={{ padding: 16, textAlign: 'center' }}>
          <Spin size="small" />
        </div>
      ) : rooms.length === 0 ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无群聊" style={{ margin: '16px 0' }} />
      ) : (
        <List
          dataSource={rooms}
          renderItem={(room) => (
            <List.Item
              onClick={() => setActiveRoom(room.id)}
              style={{
                cursor: 'pointer',
                paddingInline: 16,
                background: activeRoomId === room.id ? '#e6f4ff' : undefined,
              }}
            >
              <List.Item.Meta
                avatar={<Avatar shape="square">{room.name?.[0] || '#'}</Avatar>}
                title={room.name}
                description={room.description || `房间 #${room.id}`}
              />
            </List.Item>
          )}
        />
      )}

      <Modal
        title="创建群聊"
        open={createOpen}
        onCancel={() => setCreateOpen(false)}
        onOk={() => createForm.submit()}
        confirmLoading={createMutation.isPending}
      >
        <Form form={createForm} layout="vertical" onFinish={(values) => createMutation.mutate(values)}>
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入群聊名称' }]}>
            <Input placeholder="群聊名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea placeholder="描述（可选）" rows={2} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="加入群聊"
        open={joinOpen}
        onCancel={() => setJoinOpen(false)}
        onOk={() => joinForm.submit()}
        confirmLoading={joinMutation.isPending}
      >
        <Form
          form={joinForm}
          layout="vertical"
          onFinish={(values) => joinMutation.mutate(Number(values.roomId))}
        >
          <Form.Item
            name="roomId"
            label="房间ID"
            rules={[{ required: true, message: '请输入房间ID' }]}
          >
            <Input type="number" placeholder="输入要加入的房间ID" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
