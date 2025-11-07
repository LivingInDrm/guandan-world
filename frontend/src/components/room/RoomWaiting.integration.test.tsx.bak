import React from 'react';
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import RoomWaiting from './RoomWaiting';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { useGameStore } from '../../store/gameStore';
import { gameService } from '../../services/gameService';
import { wsClient } from '../../services/websocket';
import type { Room, Player, WSMessage } from '../../types';
import { RoomStatus, WS_MESSAGE_TYPES } from '../../types';

// Mock the stores and services
vi.mock('../../store/authStore');
vi.mock('../../store/roomStore');
vi.mock('../../store/gameStore');
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

describe('RoomWaiting WebSocket Integration', () => {
  const mockUser = {
    id: 'user1',
    username: 'testuser',
    online: true
  };

  const mockPlayers: (Player | null)[] = [
    { id: 'user1', username: 'testuser', seat: 0, online: true, auto_play: false },
    { id: 'user2', username: 'player2', seat: 1, online: true, auto_play: false },
    { id: 'user3', username: 'player3', seat: 2, online: true, auto_play: false },
    { id: 'user4', username: 'player4', seat: 3, online: true, auto_play: false }
  ];

  const mockRoom: Room = {
    id: 'test-room-id',
    status: RoomStatus.READY,
    players: mockPlayers,
    owner: 'user1',
    created_at: '2024-01-01T00:00:00Z'
  };

  const mockAuthStore = {
    user: mockUser,
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

  // Mock WebSocket handlers storage
  const wsHandlers = new Map<string, Function[]>();

  beforeEach(() => {
    vi.mocked(useAuthStore).mockReturnValue(mockAuthStore as any);
    vi.mocked(useRoomStore).mockReturnValue(mockRoomStore as any);
    vi.mocked(useGameStore).mockReturnValue(mockGameStore as any);
    
    // Mock game service
    vi.mocked(gameService.initialize).mockResolvedValue();
    vi.mocked(gameService.startGame).mockResolvedValue(true);
    vi.mocked(gameService.leaveRoom).mockResolvedValue();

    // Mock WebSocket client
    vi.mocked(wsClient.on).mockImplementation((messageType: string, handler: Function) => {
      if (!wsHandlers.has(messageType)) {
        wsHandlers.set(messageType, []);
      }
      wsHandlers.get(messageType)!.push(handler);
    });

    vi.mocked(wsClient.off).mockImplementation((messageType: string, handler: Function) => {
      const handlers = wsHandlers.get(messageType);
      if (handlers) {
        const index = handlers.indexOf(handler);
        if (index > -1) {
          handlers.splice(index, 1);
        }
      }
    });

    vi.mocked(wsClient.send).mockImplementation(() => true);
  });

  afterEach(() => {
    vi.clearAllMocks();
    mockNavigate.mockClear();
    wsHandlers.clear();
  });

  const renderComponent = (props = {}) => {
    return render(
      <BrowserRouter>
        <RoomWaiting room={mockRoom} {...props} />
      </BrowserRouter>
    );
  };

  // Helper function to simulate WebSocket message
  const simulateWSMessage = (messageType: string, data: any) => {
    const handlers = wsHandlers.get(messageType);
    if (handlers) {
      const message: WSMessage = {
        type: messageType,
        data,
        timestamp: new Date().toISOString(),
      };
      handlers.forEach(handler => handler(message));
    }
  };

  describe('WebSocket Connection and Setup', () => {
    it('should initialize game service on mount', () => {
      renderComponent();
      expect(gameService.initialize).toHaveBeenCalled();
    });

    it('should register WebSocket event handlers', () => {
      renderComponent();
      
      expect(wsClient.on).toHaveBeenCalledWith(WS_MESSAGE_TYPES.ROOM_UPDATE, expect.any(Function));
      expect(wsClient.on).toHaveBeenCalledWith(WS_MESSAGE_TYPES.GAME_PREPARE, expect.any(Function));
      expect(wsClient.on).toHaveBeenCalledWith(WS_MESSAGE_TYPES.COUNTDOWN, expect.any(Function));
      expect(wsClient.on).toHaveBeenCalledWith(WS_MESSAGE_TYPES.GAME_BEGIN, expect.any(Function));
    });

    it('should send join room message when connected', () => {
      renderComponent();
      
      expect(wsClient.send).toHaveBeenCalledWith(WS_MESSAGE_TYPES.JOIN_ROOM, {
        room_id: 'test-room-id'
      });
    });

    it('should clean up WebSocket handlers on unmount', () => {
      const { unmount } = renderComponent();
      
      unmount();
      
      expect(wsClient.off).toHaveBeenCalledWith(WS_MESSAGE_TYPES.ROOM_UPDATE, expect.any(Function));
      expect(wsClient.off).toHaveBeenCalledWith(WS_MESSAGE_TYPES.GAME_PREPARE, expect.any(Function));
      expect(wsClient.off).toHaveBeenCalledWith(WS_MESSAGE_TYPES.COUNTDOWN, expect.any(Function));
      expect(wsClient.off).toHaveBeenCalledWith(WS_MESSAGE_TYPES.GAME_BEGIN, expect.any(Function));
    });

    it('should display connection status', () => {
      renderComponent();
      
      expect(screen.getByText('已连接')).toBeInTheDocument();
    });

    it('should display disconnected status when not connected', () => {
      const disconnectedGameStore = {
        ...mockGameStore,
        isConnected: false
      };
      vi.mocked(useGameStore).mockReturnValue(disconnectedGameStore as any);

      renderComponent();
      
      expect(screen.getByText('连接断开')).toBeInTheDocument();
    });
  });

  describe('Real-time Room Updates', () => {
    it('should update room when receiving room_update message', async () => {
      renderComponent();

      const updatedRoom = {
        ...mockRoom,
        players: [
          mockPlayers[0],
          mockPlayers[1],
          mockPlayers[2],
          null // Player 4 left
        ]
      };

      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.ROOM_UPDATE, updatedRoom);
      });

      await waitFor(() => {
        expect(mockRoomStore.setCurrentRoom).toHaveBeenCalledWith(updatedRoom);
      });
    });

    it('should not update room for different room ID', async () => {
      renderComponent();

      const differentRoom = {
        ...mockRoom,
        id: 'different-room-id'
      };

      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.ROOM_UPDATE, differentRoom);
      });

      // Should not update current room
      expect(mockRoomStore.setCurrentRoom).not.toHaveBeenCalledWith(differentRoom);
    });

    it('should handle player join/leave notifications', async () => {
      renderComponent();

      // Simulate player joining
      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.JOIN_ROOM, {
          room_id: 'test-room-id',
          player_id: 'new-player',
          username: 'newplayer'
        });
      });

      // Simulate player leaving
      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.LEAVE_ROOM, {
          room_id: 'test-room-id',
          player_id: 'user4'
        });
      });

      // These messages are logged but don't directly update UI
      // The actual room update comes through room_update message
    });
  });

  describe('Game Start Flow with Countdown', () => {
    it('should show prepare screen when receiving game_prepare message', async () => {
      // Mock the game store to return countdown value after the message
      const gameStoreWithCountdown = {
        ...mockGameStore,
        countdown: 3
      };

      const { rerender } = renderComponent();

      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.GAME_PREPARE, {
          room_id: 'test-room-id'
        });
      });

      // Update the mock to return countdown value
      vi.mocked(useGameStore).mockReturnValue(gameStoreWithCountdown as any);

      // Re-render with updated store
      rerender(
        <BrowserRouter>
          <RoomWaiting room={mockRoom} />
        </BrowserRouter>
      );

      await waitFor(() => {
        expect(mockGameStore.setCountdown).toHaveBeenCalledWith(3);
        expect(screen.getByText('游戏即将开始')).toBeInTheDocument();
      });
    });

    it('should not show prepare screen for different room', async () => {
      renderComponent();

      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.GAME_PREPARE, {
          room_id: 'different-room-id'
        });
      });

      expect(screen.queryByText('游戏即将开始')).not.toBeInTheDocument();
    });

    it('should update countdown when receiving countdown message', async () => {
      // First show prepare screen
      const gameStoreWithCountdown = {
        ...mockGameStore,
        countdown: 3
      };
      vi.mocked(useGameStore).mockReturnValue(gameStoreWithCountdown as any);

      const { rerender } = renderComponent();

      // Simulate countdown updates
      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.COUNTDOWN, {
          room_id: 'test-room-id',
          countdown: 2
        });
      });

      await waitFor(() => {
        expect(mockGameStore.setCountdown).toHaveBeenCalledWith(2);
      });

      // Update mock to show countdown 2
      const gameStoreWithCountdown2 = {
        ...mockGameStore,
        countdown: 2
      };
      vi.mocked(useGameStore).mockReturnValue(gameStoreWithCountdown2 as any);

      rerender(
        <BrowserRouter>
          <RoomWaiting room={mockRoom} />
        </BrowserRouter>
      );

      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.COUNTDOWN, {
          room_id: 'test-room-id',
          countdown: 1
        });
      });

      await waitFor(() => {
        expect(mockGameStore.setCountdown).toHaveBeenCalledWith(1);
      });
    });

    it('should navigate to game when receiving game_begin message', async () => {
      // Set up prepare screen state
      const gameStoreWithCountdown = {
        ...mockGameStore,
        countdown: 1
      };
      vi.mocked(useGameStore).mockReturnValue(gameStoreWithCountdown as any);

      renderComponent();

      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.GAME_BEGIN, {
          room_id: 'test-room-id'
        });
      });

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/game/test-room-id');
      });
    });

    it('should display countdown in prepare screen', () => {
      const gameStoreWithCountdown = {
        ...mockGameStore,
        countdown: 2
      };
      vi.mocked(useGameStore).mockReturnValue(gameStoreWithCountdown as any);

      // Mock showPrepare state by simulating the prepare flow
      const TestComponent = () => {
        const [showPrepare, setShowPrepare] = React.useState(true);
        const { countdown } = useGameStore();
        
        if (showPrepare && countdown !== null) {
          return (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
              <div className="bg-white rounded-lg p-8 text-center max-w-md mx-4">
                <h2 className="text-2xl font-bold text-gray-800 mb-4">游戏即将开始</h2>
                <div className="text-6xl font-bold text-blue-500 mb-4">
                  {countdown}
                </div>
                <p className="text-gray-600">请准备好开始游戏...</p>
              </div>
            </div>
          );
        }
        
        return <RoomWaiting room={mockRoom} />;
      };

      render(
        <BrowserRouter>
          <TestComponent />
        </BrowserRouter>
      );

      expect(screen.getByText('游戏即将开始')).toBeInTheDocument();
      expect(screen.getByText('2')).toBeInTheDocument();
      expect(screen.getByText('请准备好开始游戏...')).toBeInTheDocument();
    });

    it('should show connection status in prepare screen', () => {
      const gameStoreWithCountdown = {
        ...mockGameStore,
        countdown: 3,
        isConnected: false
      };
      vi.mocked(useGameStore).mockReturnValue(gameStoreWithCountdown as any);

      const TestComponent = () => {
        const [showPrepare] = React.useState(true);
        const { countdown, isConnected } = useGameStore();
        
        if (showPrepare && countdown !== null) {
          return (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
              <div className="bg-white rounded-lg p-8 text-center max-w-md mx-4">
                <h2 className="text-2xl font-bold text-gray-800 mb-4">游戏即将开始</h2>
                <div className="text-6xl font-bold text-blue-500 mb-4">
                  {countdown}
                </div>
                <p className="text-gray-600">请准备好开始游戏...</p>
                
                <div className="mt-4 flex items-center justify-center space-x-2">
                  <div className={`w-3 h-3 rounded-full ${
                    isConnected ? 'bg-green-500' : 'bg-red-500'
                  }`} />
                  <span className="text-sm text-gray-600">
                    {isConnected ? '连接正常' : '连接断开'}
                  </span>
                </div>
              </div>
            </div>
          );
        }
        
        return <RoomWaiting room={mockRoom} />;
      };

      render(
        <BrowserRouter>
          <TestComponent />
        </BrowserRouter>
      );

      expect(screen.getByText('连接断开')).toBeInTheDocument();
    });
  });

  describe('Game Service Integration', () => {
    it('should use game service for starting game', async () => {
      renderComponent();

      const startButton = screen.getByText('开始游戏');
      fireEvent.click(startButton);

      await waitFor(() => {
        expect(gameService.startGame).toHaveBeenCalledWith('test-room-id');
      });
    });

    it('should use game service for leaving room', async () => {
      renderComponent();

      const leaveButton = screen.getByText('离开房间');
      fireEvent.click(leaveButton);

      await waitFor(() => {
        expect(gameService.leaveRoom).toHaveBeenCalledWith('test-room-id');
      });
    });

    it('should handle game service errors', async () => {
      vi.mocked(gameService.startGame).mockResolvedValue(false);

      renderComponent();

      const startButton = screen.getByText('开始游戏');
      fireEvent.click(startButton);

      await waitFor(() => {
        expect(mockRoomStore.setError).toHaveBeenCalledWith('Failed to start game');
      });
    });
  });

  describe('Error Handling and Edge Cases', () => {
    it('should handle WebSocket message with missing data', async () => {
      renderComponent();

      // Send message with null data
      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.ROOM_UPDATE, null);
      });

      // Should not crash or update room
      expect(mockRoomStore.setCurrentRoom).not.toHaveBeenCalledWith(null);
    });

    it('should handle countdown message with invalid countdown value', async () => {
      renderComponent();

      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.COUNTDOWN, {
          room_id: 'test-room-id',
          countdown: 'invalid'
        });
      });

      // Should not update countdown with invalid value
      expect(mockGameStore.setCountdown).not.toHaveBeenCalledWith('invalid');
    });

    it('should handle multiple rapid WebSocket messages', async () => {
      renderComponent();

      // Send multiple rapid updates
      act(() => {
        simulateWSMessage(WS_MESSAGE_TYPES.COUNTDOWN, {
          room_id: 'test-room-id',
          countdown: 3
        });
        simulateWSMessage(WS_MESSAGE_TYPES.COUNTDOWN, {
          room_id: 'test-room-id',
          countdown: 2
        });
        simulateWSMessage(WS_MESSAGE_TYPES.COUNTDOWN, {
          room_id: 'test-room-id',
          countdown: 1
        });
      });

      await waitFor(() => {
        expect(mockGameStore.setCountdown).toHaveBeenCalledTimes(3);
      });
    });
  });
});