import { Button, Card, Form, Input, Typography, message } from 'antd'
import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { login } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { extractErrorMessage } from '@/utils/error'

interface LoginForm {
  username: string
  password: string
}

export default function Login() {
  const [loading, setLoading] = useState(false)
  const setAuth = useAuthStore((state) => state.setAuth)
  const navigate = useNavigate()
  const location = useLocation()

  const onFinish = async (values: LoginForm) => {
    setLoading(true)
    try {
      const data = await login(values.username, values.password)
      setAuth(data.token, data.user)
      message.success('登录成功')
      const from = (location.state as { from?: string } | null)?.from || '/'
      navigate(from, { replace: true })
    } catch (error) {
      message.error(extractErrorMessage(error, '用户名或密码错误'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-page">
      <Card title="登录 Chat App" style={{ width: 380 }}>
        <Form<LoginForm> layout="vertical" onFinish={onFinish}>
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="用户名" autoComplete="username" />
          </Form.Item>
          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="密码" autoComplete="current-password" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              登录
            </Button>
          </Form.Item>
        </Form>
        <Typography.Text>
          还没有账号？<Link to="/register">去注册</Link>
        </Typography.Text>
      </Card>
    </div>
  )
}
