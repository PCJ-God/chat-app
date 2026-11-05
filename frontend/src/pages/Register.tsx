import { Button, Card, Form, Input, Typography, message } from 'antd'
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { register } from '@/api/auth'
import { extractErrorMessage } from '@/utils/error'

interface RegisterForm {
  username: string
  email: string
  password: string
  nickname?: string
}

export default function Register() {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const onFinish = async (values: RegisterForm) => {
    setLoading(true)
    try {
      await register(values)
      message.success('注册成功，请登录')
      navigate('/login')
    } catch (error) {
      message.error(extractErrorMessage(error, '注册失败，用户名或邮箱可能已存在'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-page">
      <Card title="注册" style={{ width: 380 }}>
        <Form<RegisterForm> layout="vertical" onFinish={onFinish}>
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="用户名" autoComplete="username" />
          </Form.Item>
          <Form.Item
            name="email"
            label="邮箱"
            rules={[{ required: true, type: 'email', message: '请输入有效邮箱' }]}
          >
            <Input placeholder="邮箱" autoComplete="email" />
          </Form.Item>
          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, min: 6, message: '密码至少 6 位' }]}
          >
            <Input.Password placeholder="密码" autoComplete="new-password" />
          </Form.Item>
          <Form.Item name="nickname" label="昵称">
            <Input placeholder="昵称（可选）" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              注册
            </Button>
          </Form.Item>
        </Form>
        <Typography.Text>
          已有账号？<Link to="/login">去登录</Link>
        </Typography.Text>
      </Card>
    </div>
  )
}
