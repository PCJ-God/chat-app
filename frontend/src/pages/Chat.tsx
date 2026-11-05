import { Empty } from 'antd'
import { useUIStore } from '@/stores/ui'
import MessageThread from '@/components/MessageThread'
import RoomThread from '@/components/RoomThread'

export default function Chat() {
  const activeUserId = useUIStore((state) => state.activeUserId)
  const activeRoomId = useUIStore((state) => state.activeRoomId)

  if (activeRoomId) {
    return <RoomThread roomId={activeRoomId} />
  }

  if (activeUserId) {
    return <MessageThread targetUserId={activeUserId} />
  }

  return <Empty style={{ marginTop: 120 }} description="选择左侧会话或群聊开始聊天" />
}
