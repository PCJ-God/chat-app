import { useEffect, useMemo, useRef, useState } from 'react'
import { useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Button, Input, Spin, message } from 'antd'
import dayjs from 'dayjs'
import { getMessages } from '@/api/messages'
import { sendRoomMessage } from '@/api/rooms'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/utils/error'
import type { Message } from '@/types'

const { TextArea } = Input
const PAGE_SIZE = 30
const roomMsgKey = (roomId: number) => ['messages', 'room', roomId] as const

export default function RoomThread({ roomId }: { roomId: number }) {
  const selfId = useAuthStore((state) => state.user?.id)
  const queryClient = useQueryClient()
  const [draft, setDraft] = useState('')
  const bottomRef = useRef<HTMLDivElement>(null)

  const query = useInfiniteQuery({
    queryKey: roomMsgKey(roomId),
    queryFn: ({ pageParam }) =>
      getMessages({ room_id: roomId, page: pageParam, page_size: PAGE_SIZE }),
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      const loaded = allPages.reduce((sum, page) => sum + page.items.length, 0)
      return loaded < lastPage.total ? allPages.length + 1 : undefined
    },
  })

  const messages = useMemo(() => {
    const flat = query.data?.pages.flatMap((page) => page.items) ?? []
    return [...flat].reverse()
  }, [query.data])

  const sendMutation = useMutation({
    mutationFn: (content: string) => sendRoomMessage(roomId, content),
    onSuccess: () => {
      setDraft('')
      queryClient.invalidateQueries({ queryKey: roomMsgKey(roomId) })
    },
    onError: (error) => message.error(extractErrorMessage(error, '发送失败')),
  })

  const newestId = messages.length > 0 ? messages[messages.length - 1].id : null
  const prevNewestRef = useRef<number | null>(null)

  useEffect(() => {
    prevNewestRef.current = null
  }, [roomId])

  useEffect(() => {
    if (query.isFetchingNextPage) return
    if (newestId !== null && newestId !== prevNewestRef.current) {
      prevNewestRef.current = newestId
      bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
    }
  }, [newestId, query.isFetchingNextPage])

  const handleSend = () => {
    const content = draft.trim()
    if (!content || sendMutation.isPending) return
    sendMutation.mutate(content)
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: 'calc(100vh - 64px - 32px)' }}>
      <div style={{ flex: 1, overflowY: 'auto', paddingRight: 8 }}>
        {query.isLoading ? (
          <div style={{ padding: 24, textAlign: 'center' }}>
            <Spin />
          </div>
        ) : (
          <>
            {query.hasNextPage && (
              <div style={{ textAlign: 'center', marginBottom: 12 }}>
                <Button
                  size="small"
                  loading={query.isFetchingNextPage}
                  onClick={() => query.fetchNextPage()}
                >
                  加载更多历史
                </Button>
              </div>
            )}
            {messages.map((item) => (
              <RoomMessageItem key={item.id} message={item} selfId={selfId} />
            ))}
            <div ref={bottomRef} />
          </>
        )}
      </div>
      <div style={{ display: 'flex', gap: 8, paddingTop: 12 }}>
        <TextArea
          value={draft}
          onChange={(event) => setDraft(event.target.value)}
          placeholder="输入消息，Enter 发送，Shift+Enter 换行"
          autoSize={{ minRows: 1, maxRows: 4 }}
          onPressEnter={(event) => {
            if (!event.shiftKey) {
              event.preventDefault()
              handleSend()
            }
          }}
        />
        <Button type="primary" loading={sendMutation.isPending} onClick={handleSend}>
          发送
        </Button>
      </div>
    </div>
  )
}

function RoomMessageItem({ message: m, selfId }: { message: Message; selfId?: number }) {
  const mine = m.sender_id === selfId
  const name = m.sender?.nickname || m.sender?.username || `用户${m.sender_id}`
  return (
    <div style={{ display: 'flex', justifyContent: mine ? 'flex-end' : 'flex-start', marginBottom: 12 }}>
      <div style={{ maxWidth: '70%' }}>
        <div
          style={{
            fontSize: 12,
            color: '#999',
            textAlign: mine ? 'right' : 'left',
            marginBottom: 4,
          }}
        >
          {mine ? '' : `${name} · `}
          {dayjs(m.created_at).format('MM-DD HH:mm')}
        </div>
        <div
          style={{
            padding: '8px 12px',
            borderRadius: 8,
            background: mine ? '#1677ff' : '#fff',
            color: mine ? '#fff' : '#000',
            border: mine ? 'none' : '1px solid #f0f0f0',
            whiteSpace: 'pre-wrap',
            wordBreak: 'break-word',
          }}
        >
          {m.content}
        </div>
      </div>
    </div>
  )
}
