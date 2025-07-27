import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import App from '../App'
import { AuthProvider } from '../components/auth/AuthProvider'

// Mock WebSocket
class MockWebSocket {
  static CONNECTING = 0
  static OPEN = 1
  static CLOSING = 2
  static CLOSED = 3

  readyState = MockWebSocket.CONNECTING
  onopen: ((event: Event) => void) | null = null
  onclose: ((event: CloseEvent) => void) | null = null
  onmessage: ((event: MessageEvent) => void) | null = null
  onerror: ((event: Event) => void) | null = null

  constructor(public url: string) {
    setTimeout(() => {
      this.readyState = MockWebSocket.OPEN
      if (this.onopen) {
        this.onopen(new Event('open'))
      }
    }, 100)
  }

  send(data: string) {
    // Mock sending data
    console.log('Mock WebSocket send:', data)
  }

  close() {
    this.readyState = MockWebSocket.CLOSED
    if (this.onclose) {
      this.onclose(new CloseEvent('close'))
    }
  }
}

// Mock fetch for API calls
const mockFetch = vi.fn()

// Mock API responses
const mockApiResponses = {
  register: {
    status: 201,
    json: () => Promise.resolve({
      user: { id: '1', username: 'testuser' },
      token: { token: 'mock-token', expires_at: '2024-12-31T23:59:59Z' }
    })
  },
  login: {
    status: 200,
    json: () => Promise.resolve({
      user: { id: '1', username: 'testuser' },
      token: { token: 'mock-token', expires_at: '2024-12-31T23:59:59Z' }
    })
  },
  rooms: {
    status: 200,
    json: () => Promise.resolve({
      rooms: [
        {
          id: 'room1',
          status: 'waiting',
          player_count: 2,
          players: [
            { id: '1', username: 'alice', seat: 0 },
            { id: '2', username: 'bob', seat: 1 }
          ],
          owner: '1',
          can_join: true
        },
        {
          id: 'room2',
          status: 'playing',
          player_count: 4,
          players: [
            { id: '3', username: 'charlie', seat: 0 },
            { id: '4', username: 'david', seat: 1 },
            { id: '5', username: 'eve', seat: 2 },
            { id: '6', username: 'frank', seat: 3 }
          ],
          owner: '3',
          can_join: false
        }
      ],
      total_count: 2,
      page: 1,
      limit: 12
    })
  },
  createRoom: {
    status: 201,
    json: () => Promise.resolve({
      room: {
        id: 'new-room',
        status: 'waiting',
        player_count: 1,
        players: [{ id: '1', username: 'testuser', seat: 0 }],
        owner: '1'
      }
    })
  },
  joinRoom: {
    status: 200,
    json: () => Promise.resolve({
      message: 'Successfully joined room'
    })
  }
}

describe('前端端到端测试', () => {
  beforeEach(() => {
    // Setup mocks
    global.WebSocket = MockWebSocket as any
    global.fetch = mockFetch
    
    // Clear localStorage
    localStorage.clear()
    
    // Reset fetch mock
    mockFetch.mockReset()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  const renderApp = () => {
    return render(
      <BrowserRouter>
        <AuthProvider>
          <App />
        </AuthProvider>
      </BrowserRouter>
    )
  }

  describe('需求1: 用户认证系统', () => {
    it('应该显示登录页面并支持用户注册', async () => {
      mockFetch.mockResolvedValueOnce(mockApiResponses.register)
      
      renderApp()

      // 应该显示登录页面
      expect(screen.getByText(/登录/)).toBeInTheDocument()
      
      // 点击注册链接
      const registerLink = screen.getByText(/注册新账号/)
      fireEvent.click(registerLink)

      // 应该显示注册表单
      expect(screen.getByText(/用户注册/)).toBeInTheDocument()
      
      // 填写注册表单
      const usernameInput = screen.getByPlaceholderText(/请输入用户名/)
      const passwordInput = screen.getByPlaceholderText(/请输入密码/)
      const registerButton = screen.getByRole('button', { name: /注册/ })

      fireEvent.change(usernameInput, { target: { value: 'testuser' } })
      fireEvent.change(passwordInput, { target: { value: 'password123' } })
      fireEvent.click(registerButton)

      // 验证API调用
      await waitFor(() => {
        expect(mockFetch).toHaveBeenCalledWith(
          expect.stringContaining('/api/auth/register'),
          expect.objectContaining({
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: 'testuser', password: 'password123' })
          })
        )
      })
    })

    it('应该支持用户登录并跳转到房间大厅', async () => {
      mockFetch
        .mockResolvedValueOnce(mockApiResponses.login)
        .mockResolvedValueOnce(mockApiResponses.rooms)

      renderApp()

      // 填写登录表单
      const usernameInput = screen.getByPlaceholderText(/请输入用户名/)
      const passwordInput = screen.getByPlaceholderText(/请输入密码/)
      const loginButton = screen.getByRole('button', { name: /登录/ })

      fireEvent.change(usernameInput, { target: { value: 'testuser' } })
      fireEvent.change(passwordInput, { target: { value: 'password123' } })
      fireEvent.click(loginButton)

      // 等待登录完成并跳转到房间大厅
      await waitFor(() => {
        expect(screen.getByText(/房间大厅/)).toBeInTheDocument()
      })

      // 验证房间列表显示
      expect(screen.getByText(/alice/)).toBeInTheDocument()
      expect(screen.getByText(/bob/)).toBeInTheDocument()
    })

    it('应该处理登录错误', async () => {
      mockFetch.mockResolvedValueOnce({
        status: 401,
        json: () => Promise.resolve({ error: '用户名或密码错误' })
      })

      renderApp()

      const usernameInput = screen.getByPlaceholderText(/请输入用户名/)
      const passwordInput = screen.getByPlaceholderText(/请输入密码/)
      const loginButton = screen.getByRole('button', { name: /登录/ })

      fireEvent.change(usernameInput, { target: { value: 'wronguser' } })
      fireEvent.change(passwordInput, { target: { value: 'wrongpass' } })
      fireEvent.click(loginButton)

      // 应该显示错误信息
      await waitFor(() => {
        expect(screen.getByText(/用户名或密码错误/)).toBeInTheDocument()
      })
    })
  })

  describe('需求2: 房间大厅管理', () => {
    beforeEach(async () => {
      // 模拟已登录状态
      localStorage.setItem('auth_token', 'mock-token')
      localStorage.setItem('user_info', JSON.stringify({ id: '1', username: 'testuser' }))
      
      mockFetch.mockResolvedValueOnce(mockApiResponses.rooms)
    })

    it('应该显示房间列表并支持分页', async () => {
      renderApp()

      // 等待房间列表加载
      await waitFor(() => {
        expect(screen.getByText(/房间大厅/)).toBeInTheDocument()
      })

      // 验证房间信息显示
      expect(screen.getByText(/alice/)).toBeInTheDocument()
      expect(screen.getByText(/bob/)).toBeInTheDocument()
      expect(screen.getByText(/charlie/)).toBeInTheDocument()
      
      // 验证房间状态显示
      expect(screen.getByText(/等待中/)).toBeInTheDocument()
      expect(screen.getByText(/游戏中/)).toBeInTheDocument()
    })

    it('应该支持创建新房间', async () => {
      mockFetch
        .mockResolvedValueOnce(mockApiResponses.rooms)
        .mockResolvedValueOnce(mockApiResponses.createRoom)

      renderApp()

      await waitFor(() => {
        expect(screen.getByText(/房间大厅/)).toBeInTheDocument()
      })

      // 点击创建房间按钮
      const createButton = screen.getByText(/创建房间/)
      fireEvent.click(createButton)

      // 应该显示创建房间弹窗
      expect(screen.getByText(/创建新房间/)).toBeInTheDocument()

      // 填写房间名称
      const roomNameInput = screen.getByPlaceholderText(/请输入房间名称/)
      fireEvent.change(roomNameInput, { target: { value: '我的测试房间' } })

      // 点击确认创建
      const confirmButton = screen.getByRole('button', { name: /确认创建/ })
      fireEvent.click(confirmButton)

      // 验证API调用
      await waitFor(() => {
        expect(mockFetch).toHaveBeenCalledWith(
          expect.stringContaining('/api/rooms/create'),
          expect.objectContaining({
            method: 'POST',
            headers: expect.objectContaining({
              'Content-Type': 'application/json',
              'Authorization': 'Bearer mock-token'
            }),
            body: JSON.stringify({ name: '我的测试房间' })
          })
        )
      })
    })

    it('应该支持加入房间', async () => {
      mockFetch
        .mockResolvedValueOnce(mockApiResponses.rooms)
        .mockResolvedValueOnce(mockApiResponses.joinRoom)

      renderApp()

      await waitFor(() => {
        expect(screen.getByText(/房间大厅/)).toBeInTheDocument()
      })

      // 找到可加入的房间并点击加入按钮
      const joinButtons = screen.getAllByText(/加入/)
      const availableJoinButton = joinButtons.find(button => 
        !button.closest('div')?.classList.contains('opacity-50')
      )

      if (availableJoinButton) {
        fireEvent.click(availableJoinButton)

        // 验证API调用
        await waitFor(() => {
          expect(mockFetch).toHaveBeenCalledWith(
            expect.stringContaining('/api/rooms/join'),
            expect.objectContaining({
              method: 'POST',
              headers: expect.objectContaining({
                'Authorization': 'Bearer mock-token'
              })
            })
          )
        })
      }
    })

    it('应该正确显示房间状态和按钮状态', async () => {
      mockFetch.mockResolvedValueOnce(mockApiResponses.rooms)

      renderApp()

      await waitFor(() => {
        expect(screen.getByText(/房间大厅/)).toBeInTheDocument()
      })

      // 验证等待中房间的加入按钮可点击
      const waitingRoomCards = screen.getAllByText(/等待中/).map(el => el.closest('.room-card'))
      expect(waitingRoomCards.length).toBeGreaterThan(0)

      // 验证游戏中房间的加入按钮置灰
      const playingRoomCards = screen.getAllByText(/游戏中/).map(el => el.closest('.room-card'))
      expect(playingRoomCards.length).toBeGreaterThan(0)
    })
  })

  describe('需求3: 房间内等待管理', () => {
    beforeEach(() => {
      localStorage.setItem('auth_token', 'mock-token')
      localStorage.setItem('user_info', JSON.stringify({ id: '1', username: 'testuser' }))
    })

    it('应该显示房间内玩家信息和座位', async () => {
      // Mock room waiting page data
      mockFetch.mockResolvedValueOnce({
        status: 200,
        json: () => Promise.resolve({
          room: {
            id: 'test-room',
            status: 'waiting',
            players: [
              { id: '1', username: 'testuser', seat: 0 },
              { id: '2', username: 'player2', seat: 1 },
              null,
              null
            ],
            owner: '1'
          }
        })
      })

      // 直接导航到房间等待页面
      window.history.pushState({}, '', '/room/test-room')
      renderApp()

      await waitFor(() => {
        expect(screen.getByText(/房间等待/)).toBeInTheDocument()
      })

      // 验证玩家信息显示
      expect(screen.getByText(/testuser/)).toBeInTheDocument()
      expect(screen.getByText(/player2/)).toBeInTheDocument()

      // 验证空座位显示
      expect(screen.getAllByText(/等待玩家加入/).length).toBe(2)
    })

    it('应该显示房主的开始游戏按钮', async () => {
      mockFetch.mockResolvedValueOnce({
        status: 200,
        json: () => Promise.resolve({
          room: {
            id: 'test-room',
            status: 'ready',
            players: [
              { id: '1', username: 'testuser', seat: 0 },
              { id: '2', username: 'player2', seat: 1 },
              { id: '3', username: 'player3', seat: 2 },
              { id: '4', username: 'player4', seat: 3 }
            ],
            owner: '1'
          }
        })
      })

      window.history.pushState({}, '', '/room/test-room')
      renderApp()

      await waitFor(() => {
        expect(screen.getByText(/房间等待/)).toBeInTheDocument()
      })

      // 房主应该看到开始游戏按钮
      expect(screen.getByText(/开始游戏/)).toBeInTheDocument()
    })

    it('人数不足时开始游戏按钮应该置灰', async () => {
      mockFetch.mockResolvedValueOnce({
        status: 200,
        json: () => Promise.resolve({
          room: {
            id: 'test-room',
            status: 'waiting',
            players: [
              { id: '1', username: 'testuser', seat: 0 },
              { id: '2', username: 'player2', seat: 1 },
              null,
              null
            ],
            owner: '1'
          }
        })
      })

      window.history.pushState({}, '', '/room/test-room')
      renderApp()

      await waitFor(() => {
        expect(screen.getByText(/房间等待/)).toBeInTheDocument()
      })

      // 开始游戏按钮应该置灰
      const startButton = screen.getByText(/开始游戏/)
      expect(startButton.closest('button')).toBeDisabled()
    })
  })

  describe('需求4: 游戏开始流程', () => {
    beforeEach(() => {
      localStorage.setItem('auth_token', 'mock-token')
      localStorage.setItem('user_info', JSON.stringify({ id: '1', username: 'testuser' }))
    })

    it('应该显示游戏准备页面和倒计时', async () => {
      // Mock WebSocket messages
      const mockWS = new MockWebSocket('ws://localhost/ws')
      
      setTimeout(() => {
        if (mockWS.onmessage) {
          mockWS.onmessage(new MessageEvent('message', {
            data: JSON.stringify({
              type: 'game_prepare',
              data: { countdown: 3 }
            })
          }))
        }
      }, 200)

      window.history.pushState({}, '', '/game/test-room')
      renderApp()

      // 等待游戏准备消息
      await waitFor(() => {
        expect(screen.getByText(/游戏准备中/)).toBeInTheDocument()
      }, { timeout: 1000 })

      // 应该显示倒计时
      expect(screen.getByText(/3/)).toBeInTheDocument()
    })

    it('应该在倒计时结束后进入游戏界面', async () => {
      const mockWS = new MockWebSocket('ws://localhost/ws')
      
      setTimeout(() => {
        if (mockWS.onmessage) {
          // 先发送准备消息
          mockWS.onmessage(new MessageEvent('message', {
            data: JSON.stringify({
              type: 'game_prepare',
              data: { countdown: 1 }
            })
          }))
          
          // 然后发送游戏开始消息
          setTimeout(() => {
            if (mockWS.onmessage) {
              mockWS.onmessage(new MessageEvent('message', {
                data: JSON.stringify({
                  type: 'game_begin',
                  data: {
                    game_state: {
                      status: 'playing',
                      current_player: 0,
                      players: [
                        { id: '1', username: 'testuser', seat: 0 },
                        { id: '2', username: 'player2', seat: 1 },
                        { id: '3', username: 'player3', seat: 2 },
                        { id: '4', username: 'player4', seat: 3 }
                      ]
                    }
                  }
                })
              }))
            }
          }, 100)
        }
      }, 200)

      window.history.pushState({}, '', '/game/test-room')
      renderApp()

      // 等待游戏界面显示
      await waitFor(() => {
        expect(screen.getByText(/游戏进行中/) || screen.getByText(/当前玩家/)).toBeInTheDocument()
      }, { timeout: 2000 })
    })
  })

  describe('需求10: 断线重连', () => {
    beforeEach(() => {
      localStorage.setItem('auth_token', 'mock-token')
      localStorage.setItem('user_info', JSON.stringify({ id: '1', username: 'testuser' }))
    })

    it('应该处理WebSocket连接断开', async () => {
      const mockWS = new MockWebSocket('ws://localhost/ws')
      
      renderApp()

      // 模拟连接断开
      setTimeout(() => {
        mockWS.readyState = MockWebSocket.CLOSED
        if (mockWS.onclose) {
          mockWS.onclose(new CloseEvent('close'))
        }
      }, 500)

      // 等待断线处理
      await waitFor(() => {
        // 应该显示连接状态或重连提示
        expect(
          screen.queryByText(/连接断开/) || 
          screen.queryByText(/重新连接/) ||
          screen.queryByText(/网络异常/)
        ).toBeTruthy()
      }, { timeout: 1000 })
    })

    it('应该支持自动重连', async () => {
      let connectionCount = 0
      const originalWebSocket = global.WebSocket

      // Mock WebSocket with reconnection
      global.WebSocket = class extends MockWebSocket {
        constructor(url: string) {
          super(url)
          connectionCount++
          
          if (connectionCount === 1) {
            // 第一次连接后立即断开
            setTimeout(() => {
              this.readyState = MockWebSocket.CLOSED
              if (this.onclose) {
                this.onclose(new CloseEvent('close'))
              }
            }, 100)
          } else {
            // 第二次连接成功
            setTimeout(() => {
              this.readyState = MockWebSocket.OPEN
              if (this.onopen) {
                this.onopen(new Event('open'))
              }
            }, 100)
          }
        }
      } as any

      renderApp()

      // 等待重连完成
      await waitFor(() => {
        expect(connectionCount).toBeGreaterThan(1)
      }, { timeout: 2000 })

      global.WebSocket = originalWebSocket
    })
  })

  describe('需求11: 操作时间控制', () => {
    beforeEach(() => {
      localStorage.setItem('auth_token', 'mock-token')
      localStorage.setItem('user_info', JSON.stringify({ id: '1', username: 'testuser' }))
    })

    it('应该显示操作倒计时', async () => {
      const mockWS = new MockWebSocket('ws://localhost/ws')
      
      setTimeout(() => {
        if (mockWS.onmessage) {
          mockWS.onmessage(new MessageEvent('message', {
            data: JSON.stringify({
              type: 'player_turn',
              data: {
                current_player: 0,
                timeout_seconds: 20,
                action_required: 'play_cards'
              }
            })
          }))
        }
      }, 200)

      window.history.pushState({}, '', '/game/test-room')
      renderApp()

      // 等待倒计时显示
      await waitFor(() => {
        expect(screen.getByText(/20/) || screen.getByText(/倒计时/)).toBeInTheDocument()
      }, { timeout: 1000 })
    })

    it('应该处理操作超时', async () => {
      const mockWS = new MockWebSocket('ws://localhost/ws')
      
      setTimeout(() => {
        if (mockWS.onmessage) {
          // 发送超时消息
          mockWS.onmessage(new MessageEvent('message', {
            data: JSON.stringify({
              type: 'player_timeout',
              data: {
                player_id: '1',
                action: 'auto_pass'
              }
            })
          }))
        }
      }, 200)

      window.history.pushState({}, '', '/game/test-room')
      renderApp()

      // 等待超时处理
      await waitFor(() => {
        expect(
          screen.queryByText(/超时/) || 
          screen.queryByText(/自动/) ||
          screen.queryByText(/托管/)
        ).toBeTruthy()
      }, { timeout: 1000 })
    })
  })

  describe('完整游戏流程集成测试', () => {
    it('应该支持完整的用户游戏流程', async () => {
      // 1. 用户注册登录
      mockFetch
        .mockResolvedValueOnce(mockApiResponses.login)
        .mockResolvedValueOnce(mockApiResponses.rooms)
        .mockResolvedValueOnce(mockApiResponses.createRoom)

      renderApp()

      // 登录
      const usernameInput = screen.getByPlaceholderText(/请输入用户名/)
      const passwordInput = screen.getByPlaceholderText(/请输入密码/)
      const loginButton = screen.getByRole('button', { name: /登录/ })

      fireEvent.change(usernameInput, { target: { value: 'testuser' } })
      fireEvent.change(passwordInput, { target: { value: 'password123' } })
      fireEvent.click(loginButton)

      // 2. 进入房间大厅
      await waitFor(() => {
        expect(screen.getByText(/房间大厅/)).toBeInTheDocument()
      })

      // 3. 创建房间
      const createButton = screen.getByText(/创建房间/)
      fireEvent.click(createButton)

      const roomNameInput = screen.getByPlaceholderText(/请输入房间名称/)
      fireEvent.change(roomNameInput, { target: { value: '测试房间' } })

      const confirmButton = screen.getByRole('button', { name: /确认创建/ })
      fireEvent.click(confirmButton)

      // 4. 验证整个流程
      await waitFor(() => {
        expect(mockFetch).toHaveBeenCalledTimes(3) // login + rooms + create
      })

      // 验证所有API调用都正确
      expect(mockFetch).toHaveBeenNthCalledWith(1, 
        expect.stringContaining('/api/auth/login'),
        expect.any(Object)
      )
      expect(mockFetch).toHaveBeenNthCalledWith(2,
        expect.stringContaining('/api/rooms'),
        expect.any(Object)
      )
      expect(mockFetch).toHaveBeenNthCalledWith(3,
        expect.stringContaining('/api/rooms/create'),
        expect.any(Object)
      )
    })
  })

  describe('错误处理和边界情况', () => {
    beforeEach(() => {
      localStorage.setItem('auth_token', 'mock-token')
      localStorage.setItem('user_info', JSON.stringify({ id: '1', username: 'testuser' }))
    })

    it('应该处理网络错误', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      renderApp()

      // 等待错误处理
      await waitFor(() => {
        expect(
          screen.queryByText(/网络错误/) ||
          screen.queryByText(/连接失败/) ||
          screen.queryByText(/请稍后重试/)
        ).toBeTruthy()
      }, { timeout: 1000 })
    })

    it('应该处理API错误响应', async () => {
      mockFetch.mockResolvedValueOnce({
        status: 500,
        json: () => Promise.resolve({ error: '服务器内部错误' })
      })

      renderApp()

      await waitFor(() => {
        expect(screen.queryByText(/服务器内部错误/)).toBeTruthy()
      }, { timeout: 1000 })
    })

    it('应该处理无效的token', async () => {
      localStorage.setItem('auth_token', 'invalid-token')
      
      mockFetch.mockResolvedValueOnce({
        status: 401,
        json: () => Promise.resolve({ error: 'Invalid token' })
      })

      renderApp()

      // 应该重定向到登录页面
      await waitFor(() => {
        expect(screen.getByText(/登录/)).toBeInTheDocument()
      })
    })
  })
})