import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import RoomWaiting from './RoomWaiting';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { useGameStore } from '../../store/gameStore';
import { apiClient } from '../../services/api';
import { gameService } from '../../services/gameService';
import { wsClient } from '../../services/websocket';
import type { Room, Player } from '../../types';
import { RoomStatus } from '../../types';

// Mock the stores and services
vi.mock('../../store/authStore');
vi.mock('../../store/roomStore');
vi.mock('../../store/gameStore');
vi.mock('../../services/api');
vi.mock('../../services/gameService');
vi.mock('../../services/websocket');

// Mock react-router-dom
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useParams: () => ({ roomId: 'test-room-id' }),
  };
});

describe('RoomWaiting', () => {
  const mockUser = {
    id: 'user1',
    username: 'testuser',
    online: true
  };

  const mockPlayers: (Player | null)[] = [
    { id: 'user1', username: 'testuser', seat: 0, online: true, auto_play: false },
    { id: 'user2', username: 'player2', seat: 1, online: true, auto_play: false },
    null,
    null
  ];

  const mockRoom: Room = {
    id: 'test-room-id',
    status: RoomStatus.WAITING,
    players: mockPlayers,
    owner: 'user1',
    created_at: '2024-01-01T00:00:00Z'
  };

  const mockAuthStore = {
    user: mockUser,
    token: {
      token: 'mock-token',
      expires_at: '2024-12-31T23:59:59Z',
      user_id: 'user1'
    },
  };

  const mockRoomStore = {
    currentRoom: mockRoom,
    setCurrentRoom: vi.fn(),
    setError: vi.fn(),
    setLoading: vi.fn(),
  };

  const mockGameStore = {
    countdown: null,
    isConnected: true,
    setCountdown: vi.fn(),
  };

  beforeEach(() => {
    vi.mocked(useAuthStore).mockReturnValue(mockAuthStore as any);
    vi.mocked(useRoomStore).mockReturnValue(mockRoomStore as any);
    vi.mocked(useGameStore).mockReturnValue(mockGameStore as any);
    vi.mocked(apiClient.getRoomDetails).mockResolvedValue({
      success: true,
      data: mockRoom
    });
    vi.mocked(apiClient.startGame).mockResolvedValue({ success: true });
    vi.mocked(apiClient.leaveRoom).mockResolvedValue({ success: true });
    
    // Mock game service
    vi.mocked(gameService.initialize).mockResolvedValue();
    vi.mocked(gameService.startGame).mockResolvedValue(true);
    vi.mocked(gameService.leaveRoom).mockResolvedValue();
    
    // Mock WebSocket client
    vi.mocked(wsClient.on).mockImplementation(() => {});
    vi.mocked(wsClient.off).mockImplementation(() => {});
    vi.mocked(wsClient.send).mockImplementation(() => true);
  });

  afterEach(() => {
    vi.clearAllMocks();
    mockNavigate.mockClear();
  });

  const renderComponent = (props = {}) => {
    return render(
      <BrowserRouter>
        <RoomWaiting {...props} />
      </BrowserRouter>
    );
  };

  describe('Room Display', () => {
    it('should display room information correctly', () => {
      renderComponent({ room: mockRoom });

      expect(screen.getByText('房间等待')).toBeInTheDocument();
      expect(screen.getByText('房间ID: test-room-id')).toBeInTheDocument();
      expect(screen.getByText('玩家数量: 2/4')).toBeInTheDocument();
      expect(screen.getByText('状态: 等待中')).toBeInTheDocument();
    });

    it('should display all 4 seats', () => {
      renderComponent({ room: mockRoom });

      expect(screen.getByText('座位 1')).toBeInTheDocument();
      expect(screen.getByText('座位 2')).toBeInTheDocument();
      expect(screen.getByText('座位 3')).toBeInTheDocument();
      expect(screen.getByText('座位 4')).toBeInTheDocument();
    });

    it('should show player information in occupied seats', () => {
      renderComponent({ room: mockRoom });

      expect(screen.getByText('testuser')).toBeInTheDocument();
      expect(screen.getByText('player2')).toBeInTheDocument();
      expect(screen.getAllByText('等待玩家')).toHaveLength(2);
    });

    it('should show owner badge for room owner', () => {
      renderComponent({ room: mockRoom });

      expect(screen.getByText('房主')).toBeInTheDocument();
    });

    it('should show online status for players', () => {
      renderComponent({ room: mockRoom });

      expect(screen.getAllByText('在线')).toHaveLength(2);
    });

    it('should show auto-play status when player is in auto-play mode', () => {
      const roomWithAutoPlay = {
        ...mockRoom,
        players: [
          { ...mockPlayers[0]!, auto_play: true },
          mockPlayers[1],
          null,
          null
        ]
      };

      renderComponent({ room: roomWithAutoPlay });

      expect(screen.getByText('托管中')).toBeInTheDocument();
    });
  });

  describe('Room Owner Controls', () => {
    it('should show start game button for room owner', () => {
      renderComponent({ room: mockRoom });

      expect(screen.getByText('开始游戏')).toBeInTheDocument();
    });

    it('should not show start game button for non-owner', () => {
      const nonOwnerStore = {
        ...mockAuthStore,
        user: { ...mockUser, id: 'user2' }
      };
      vi.mocked(useAuthStore).mockReturnValue(nonOwnerStore as any);

      renderComponent({ room: mockRoom });

      expect(screen.queryByText('开始游戏')).not.toBeInTheDocument();
      expect(screen.getByText('等待房主开始游戏...')).toBeInTheDocument();
    });

    it('should disable start game button when less than 4 players', () => {
      renderComponent({ room: mockRoom });

      const startButton = screen.getByText('开始游戏');
      expect(startButton).toBeDisabled();
      expect(screen.getByText('需要4名玩家才能开始游戏，当前有 2 名玩家')).toBeInTheDocument();
    });

    it('should enable start game button when room is ready', () => {
      const readyRoom = {
        ...mockRoom,
        status: RoomStatus.READY,
        players: [
          mockPlayers[0],
          mockPlayers[1],
          { id: 'user3', username: 'player3', seat: 2, online: true, auto_play: false },
          { id: 'user4', username: 'player4', seat: 3, online: true, auto_play: false }
        ]
      };

      renderComponent({ room: readyRoom });

      const startButton = screen.getByText('开始游戏');
      expect(startButton).not.toBeDisabled();
    });

    it('should call start game API when start button is clicked', async () => {
      const readyRoom = {
        ...mockRoom,
        status: RoomStatus.READY,
        players: [
          mockPlayers[0],
          mockPlayers[1],
          { id: 'user3', username: 'player3', seat: 2, online: true, auto_play: false },
          { id: 'user4', username: 'player4', seat: 3, online: true, auto_play: false }
        ]
      };

      renderComponent({ room: readyRoom });

      const startButton = screen.getByText('开始游戏');
      fireEvent.click(startButton);

      await waitFor(() => {
        expect(gameService.startGame).toHaveBeenCalledWith('test-room-id');
      });
    });

    it('should call onStartGame prop when provided', async () => {
      const onStartGame = vi.fn();
      const readyRoom = {
        ...mockRoom,
        status: RoomStatus.READY,
        players: [
          mockPlayers[0],
          mockPlayers[1],
          { id: 'user3', username: 'player3', seat: 2, online: true, auto_play: false },
          { id: 'user4', username: 'player4', seat: 3, online: true, auto_play: false }
        ]
      };

      renderComponent({ room: readyRoom, onStartGame });

      const startButton = screen.getByText('开始游戏');
      fireEvent.click(startButton);

      await waitFor(() => {
        expect(onStartGame).toHaveBeenCalled();
      });
      expect(apiClient.startGame).not.toHaveBeenCalled();
    });
  });

  describe('Leave Room Functionality', () => {
    it('should show leave room button', () => {
      renderComponent({ room: mockRoom });

      expect(screen.getByText('离开房间')).toBeInTheDocument();
    });

    it('should call leave room API when leave button is clicked', async () => {
      renderComponent({ room: mockRoom });

      const leaveButton = screen.getByText('离开房间');
      fireEvent.click(leaveButton);

      await waitFor(() => {
        expect(gameService.leaveRoom).toHaveBeenCalledWith('test-room-id');
      });
    });

    it('should navigate to lobby after leaving room', async () => {
      renderComponent({ room: mockRoom });

      const leaveButton = screen.getByText('离开房间');
      fireEvent.click(leaveButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/lobby');
      });
    });

    it('should call onLeaveRoom prop when provided', async () => {
      const onLeaveRoom = vi.fn();

      renderComponent({ room: mockRoom, onLeaveRoom });

      const leaveButton = screen.getByText('离开房间');
      fireEvent.click(leaveButton);

      await waitFor(() => {
        expect(onLeaveRoom).toHaveBeenCalled();
      });
      expect(apiClient.leaveRoom).not.toHaveBeenCalled();
    });

    it('should disable leave button while leaving', async () => {
      renderComponent({ room: mockRoom });

      const leaveButton = screen.getByText('离开房间');
      fireEvent.click(leaveButton);

      expect(screen.getByText('离开中...')).toBeInTheDocument();
      expect(screen.getByText('离开中...')).toBeDisabled();
    });
  });

  describe('Error Handling', () => {
    it('should handle start game API error', async () => {
      vi.mocked(gameService.startGame).mockResolvedValue(false);

      const readyRoom = {
        ...mockRoom,
        status: RoomStatus.READY,
        players: [
          mockPlayers[0],
          mockPlayers[1],
          { id: 'user3', username: 'player3', seat: 2, online: true, auto_play: false },
          { id: 'user4', username: 'player4', seat: 3, online: true, auto_play: false }
        ]
      };

      renderComponent({ room: readyRoom });

      const startButton = screen.getByText('开始游戏');
      fireEvent.click(startButton);

      await waitFor(() => {
        expect(mockRoomStore.setError).toHaveBeenCalledWith('Failed to start game');
      });
    });

    it('should handle leave room API error', async () => {
      vi.mocked(gameService.leaveRoom).mockRejectedValue(new Error('API Error'));

      renderComponent({ room: mockRoom });

      const leaveButton = screen.getByText('离开房间');
      fireEvent.click(leaveButton);

      await waitFor(() => {
        expect(mockRoomStore.setError).toHaveBeenCalledWith('Failed to leave room');
      });
    });

    it('should prevent non-owner from starting game', async () => {
      const nonOwnerStore = {
        ...mockAuthStore,
        user: { ...mockUser, id: 'user2' }
      };
      vi.mocked(useAuthStore).mockReturnValue(nonOwnerStore as any);

      const readyRoom = {
        ...mockRoom,
        status: RoomStatus.READY,
        players: [
          mockPlayers[0],
          mockPlayers[1],
          { id: 'user3', username: 'player3', seat: 2, online: true, auto_play: false },
          { id: 'user4', username: 'player4', seat: 3, online: true, auto_play: false }
        ]
      };

      // Mock the component to have a start button for testing
      const TestComponent = () => {
        const handleStart = () => {
          // This should trigger the owner check
        };
        return (
          <div>
            <RoomWaiting room={readyRoom} />
            <button onClick={handleStart}>Test Start</button>
          </div>
        );
      };

      render(
        <BrowserRouter>
          <TestComponent />
        </BrowserRouter>
      );

      // Non-owner should not see start button
      expect(screen.queryByText('开始游戏')).not.toBeInTheDocument();
    });
  });

  describe('Loading States', () => {
    it('should show loading state when room is not available', () => {
      const emptyRoomStore = {
        ...mockRoomStore,
        currentRoom: null
      };
      vi.mocked(useRoomStore).mockReturnValue(emptyRoomStore as any);

      renderComponent();

      expect(screen.getByText('加载房间信息...')).toBeInTheDocument();
    });

    it('should show starting state when starting game', async () => {
      const readyRoom = {
        ...mockRoom,
        status: RoomStatus.READY,
        players: [
          mockPlayers[0],
          mockPlayers[1],
          { id: 'user3', username: 'player3', seat: 2, online: true, auto_play: false },
          { id: 'user4', username: 'player4', seat: 3, online: true, auto_play: false }
        ]
      };

      // Mock API to delay response
      vi.mocked(apiClient.startGame).mockImplementation(() => 
        new Promise(resolve => setTimeout(() => resolve({ success: true }), 100))
      );

      renderComponent({ room: readyRoom });

      const startButton = screen.getByText('开始游戏');
      fireEvent.click(startButton);

      expect(screen.getByText('开始中...')).toBeInTheDocument();
      expect(screen.getByText('开始中...')).toBeDisabled();
    });
  });

  describe('Room Status Display', () => {
    it('should display correct status for different room states', () => {
      const testCases = [
        { status: RoomStatus.WAITING, expected: '等待中' },
        { status: RoomStatus.READY, expected: '准备就绪' },
        { status: RoomStatus.PLAYING, expected: '游戏中' },
        { status: RoomStatus.CLOSED, expected: '已关闭' }
      ];

      testCases.forEach(({ status, expected }) => {
        const testRoom = { ...mockRoom, status };
        const { unmount } = renderComponent({ room: testRoom });
        
        expect(screen.getByText(`状态: ${expected}`)).toBeInTheDocument();
        
        unmount();
      });
    });
  });
});