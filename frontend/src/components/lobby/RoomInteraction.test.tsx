import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import RoomLobby from './RoomLobby';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { apiClient } from '../../services/api';
import type { User, RoomListResponse, RoomInfo } from '../../types';

// Mock the stores
vi.mock('../../store/authStore');
vi.mock('../../store/roomStore');
vi.mock('../../services/api');

const mockUser: User = {
  id: 'user1',
  username: 'testuser',
  online: true
};

// Create mock data for pagination testing (25 rooms to test pagination)
const createMockRooms = (count: number): RoomInfo[] => {
  return Array.from({ length: count }, (_, index) => ({
    id: `room${index + 1}`,
    status: index % 3, // Mix of different statuses
    player_count: Math.floor(Math.random() * 4) + 1,
    players: [
      { id: `user${index + 1}`, username: `player${index + 1}`, seat: 0, online: true, auto_play: false }
    ],
    owner: `user${index + 1}`,
    can_join: index % 2 === 0 // Half can join, half cannot
  }));
};

const mockRoomsPage1 = createMockRooms(12);
const mockRoomsPage2 = createMockRooms(12);
const mockRoomsPage3 = createMockRooms(1); // Last page with 1 room

const renderRoomLobby = () => {
  return render(
    <BrowserRouter>
      <RoomLobby />
    </BrowserRouter>
  );
};

describe('Room Interaction Integration Tests', () => {
  const mockAuthStore = {
    user: mockUser,
  };

  const mockRoomStore = {
    roomList: [],
    totalCount: 0,
    currentPage: 1,
    limit: 12,
    isLoading: false,
    error: null,
    setRoomList: vi.fn(),
    setLoading: vi.fn(),
    setError: vi.fn(),
    clearError: vi.fn(),
    setPage: vi.fn(),
  };

  beforeEach(() => {
    vi.mocked(useAuthStore).mockReturnValue(mockAuthStore as any);
    vi.mocked(useRoomStore).mockReturnValue(mockRoomStore as any);
    
    // Mock API responses for different pages
    vi.mocked(apiClient.getRoomList).mockImplementation((page = 1, limit = 12) => {
      let rooms: RoomInfo[];
      let totalCount = 25;
      
      switch (page) {
        case 1:
          rooms = mockRoomsPage1;
          break;
        case 2:
          rooms = mockRoomsPage2;
          break;
        case 3:
          rooms = mockRoomsPage3;
          break;
        default:
          rooms = [];
      }
      
      return Promise.resolve({
        success: true,
        data: {
          rooms,
          total_count: totalCount,
          page,
          limit
        }
      });
    });
    
    vi.mocked(apiClient.joinRoom).mockResolvedValue({
      success: true,
      data: { id: 'room1', status: 0, players: [null, null, null, null], owner: 'user1', created_at: new Date().toISOString() }
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Pagination Functionality', () => {
    it('loads 12 rooms per page by default', async () => {
      renderRoomLobby();
      
      await waitFor(() => {
        expect(apiClient.getRoomList).toHaveBeenCalledWith(1, 12);
      });
    });

    it('handles page changes correctly', async () => {
      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: mockRoomsPage1,
        totalCount: 25,
        currentPage: 1
      } as any);
      
      renderRoomLobby();
      
      // Should show pagination for 25 total rooms (3 pages)
      expect(screen.getByText('下一页')).toBeInTheDocument();
      
      // Click next page
      const nextButton = screen.getByText('下一页');
      fireEvent.click(nextButton);
      
      expect(mockRoomStore.setPage).toHaveBeenCalledWith(2);
    });

    it('disables pagination buttons appropriately', () => {
      // Test first page
      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: mockRoomsPage1,
        totalCount: 25,
        currentPage: 1
      } as any);
      
      renderRoomLobby();
      
      const prevButton = screen.getByText('上一页');
      const nextButton = screen.getByText('下一页');
      
      expect(prevButton).toBeDisabled();
      expect(nextButton).not.toBeDisabled();
    });

    it('shows correct page numbers', () => {
      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: mockRoomsPage1,
        totalCount: 25,
        currentPage: 1
      } as any);
      
      renderRoomLobby();
      
      expect(screen.getByText('1')).toBeInTheDocument();
      expect(screen.getByText('2')).toBeInTheDocument();
      expect(screen.getByText('3')).toBeInTheDocument();
    });
  });

  describe('Room Join Button State Control', () => {
    it('shows join button as clickable for joinable rooms', () => {
      const joinableRoom: RoomInfo = {
        id: 'room1',
        status: 0, // WAITING
        player_count: 2,
        players: [
          { id: 'other1', username: 'other1', seat: 0, online: true, auto_play: false },
          { id: 'other2', username: 'other2', seat: 1, online: true, auto_play: false }
        ],
        owner: 'other1',
        can_join: true
      };

      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: [joinableRoom],
        totalCount: 1
      } as any);
      
      renderRoomLobby();
      
      const joinButton = screen.getByText('加入房间');
      expect(joinButton).not.toBeDisabled();
      expect(joinButton).toHaveClass('bg-blue-600');
    });

    it('shows join button as disabled for full rooms', () => {
      const fullRoom: RoomInfo = {
        id: 'room1',
        status: 1, // READY (full)
        player_count: 4,
        players: [
          { id: 'other1', username: 'other1', seat: 0, online: true, auto_play: false },
          { id: 'other2', username: 'other2', seat: 1, online: true, auto_play: false },
          { id: 'other3', username: 'other3', seat: 2, online: true, auto_play: false },
          { id: 'other4', username: 'other4', seat: 3, online: true, auto_play: false }
        ],
        owner: 'other1',
        can_join: false
      };

      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: [fullRoom],
        totalCount: 1
      } as any);
      
      renderRoomLobby();
      
      const disabledButton = screen.getByText('房间已满');
      expect(disabledButton).toBeDisabled();
      expect(disabledButton).toHaveClass('bg-gray-100');
    });

    it('shows different button state for rooms in game', () => {
      const playingRoom: RoomInfo = {
        id: 'room1',
        status: 2, // PLAYING
        player_count: 4,
        players: [
          { id: 'other1', username: 'other1', seat: 0, online: true, auto_play: false },
          { id: 'other2', username: 'other2', seat: 1, online: true, auto_play: false },
          { id: 'other3', username: 'other3', seat: 2, online: true, auto_play: false },
          { id: 'other4', username: 'other4', seat: 3, online: true, auto_play: false }
        ],
        owner: 'other1',
        can_join: false
      };

      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: [playingRoom],
        totalCount: 1
      } as any);
      
      renderRoomLobby();
      
      // Find the button specifically (not the status badge or statistics)
      const gameButtons = screen.getAllByText('游戏中');
      const gameButton = gameButtons.find(button => button.tagName === 'BUTTON');
      
      expect(gameButton).toBeDisabled();
      expect(gameButton).toHaveClass('bg-gray-100');
    });

    it('shows "already in room" for user\'s own room', () => {
      const userRoom: RoomInfo = {
        id: 'room1',
        status: 0, // WAITING
        player_count: 1,
        players: [
          { id: 'user1', username: 'testuser', seat: 0, online: true, auto_play: false } // Current user
        ],
        owner: 'user1',
        can_join: true
      };

      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: [userRoom],
        totalCount: 1
      } as any);
      
      renderRoomLobby();
      
      const inRoomButton = screen.getByText('已在房间内');
      expect(inRoomButton).toBeDisabled();
      expect(inRoomButton).toHaveClass('bg-gray-100');
    });
  });

  describe('Room Join Functionality', () => {
    it('calls join room API when join button is clicked', async () => {
      const joinableRoom: RoomInfo = {
        id: 'room1',
        status: 0,
        player_count: 2,
        players: [
          { id: 'other1', username: 'other1', seat: 0, online: true, auto_play: false },
          { id: 'other2', username: 'other2', seat: 1, online: true, auto_play: false }
        ],
        owner: 'other1',
        can_join: true
      };

      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: [joinableRoom],
        totalCount: 1
      } as any);
      
      renderRoomLobby();
      
      const joinButton = screen.getByText('加入房间');
      fireEvent.click(joinButton);
      
      await waitFor(() => {
        expect(apiClient.joinRoom).toHaveBeenCalledWith('room1');
      });
    });

    it('handles join room API errors', async () => {
      const joinableRoom: RoomInfo = {
        id: 'room1',
        status: 0,
        player_count: 2,
        players: [
          { id: 'other1', username: 'other1', seat: 0, online: true, auto_play: false },
          { id: 'other2', username: 'other2', seat: 1, online: true, auto_play: false }
        ],
        owner: 'other1',
        can_join: true
      };

      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: [joinableRoom],
        totalCount: 1
      } as any);

      vi.mocked(apiClient.joinRoom).mockRejectedValue(new Error('Room is full'));
      
      renderRoomLobby();
      
      const joinButton = screen.getByText('加入房间');
      fireEvent.click(joinButton);
      
      await waitFor(() => {
        expect(mockRoomStore.setError).toHaveBeenCalledWith('Room is full');
      });
    });
  });

  describe('Real-time Updates', () => {
    it('sets up auto-refresh interval on mount', () => {
      vi.useFakeTimers();
      
      renderRoomLobby();
      
      // Initial load
      expect(apiClient.getRoomList).toHaveBeenCalledTimes(1);
      
      // Fast-forward 5 seconds
      vi.advanceTimersByTime(5000);
      
      // Should have called again for auto-refresh
      expect(apiClient.getRoomList).toHaveBeenCalledTimes(2);
      
      vi.useRealTimers();
    });

    it('shows loading indicator during refresh', () => {
      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: mockRoomsPage1,
        totalCount: 12,
        isLoading: true
      } as any);
      
      renderRoomLobby();
      
      expect(screen.getByText('更新中...')).toBeInTheDocument();
    });

    it('manual refresh button works correctly', async () => {
      vi.mocked(useRoomStore).mockReturnValue({
        ...mockRoomStore,
        roomList: mockRoomsPage1,
        totalCount: 12
      } as any);
      
      renderRoomLobby();
      
      const refreshButton = screen.getByText('手动刷新');
      fireEvent.click(refreshButton);
      
      await waitFor(() => {
        expect(apiClient.getRoomList).toHaveBeenCalledTimes(2); // Once on mount, once on manual refresh
      });
    });
  });
});