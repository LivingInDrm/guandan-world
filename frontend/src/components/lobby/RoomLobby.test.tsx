import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import RoomLobby from './RoomLobby';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { apiClient } from '../../services/api';
import type { User, RoomListResponse } from '../../types';

// Mock the stores
vi.mock('../../store/authStore');
vi.mock('../../store/roomStore');
vi.mock('../../services/api');

const mockUser: User = {
  id: 'user1',
  username: 'testuser',
  online: true
};

const mockRoomListResponse: RoomListResponse = {
  rooms: [
    {
      id: 'room1',
      status: 0, // WAITING
      player_count: 2,
      players: [
        { id: 'user1', username: 'player1', seat: 0, online: true, auto_play: false },
        { id: 'user2', username: 'player2', seat: 1, online: true, auto_play: false }
      ],
      owner: 'user1',
      can_join: true
    },
    {
      id: 'room2',
      status: 2, // PLAYING
      player_count: 4,
      players: [
        { id: 'user3', username: 'player3', seat: 0, online: true, auto_play: false },
        { id: 'user4', username: 'player4', seat: 1, online: true, auto_play: false },
        { id: 'user5', username: 'player5', seat: 2, online: true, auto_play: false },
        { id: 'user6', username: 'player6', seat: 3, online: true, auto_play: false }
      ],
      owner: 'user3',
      can_join: false
    }
  ],
  total_count: 2,
  page: 1,
  limit: 12
};

const renderRoomLobby = () => {
  return render(
    <BrowserRouter>
      <RoomLobby />
    </BrowserRouter>
  );
};

describe('RoomLobby', () => {
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
    vi.mocked(apiClient.getMyRoom).mockRejectedValue({ status: 404 });
    vi.mocked(apiClient.getRoomList).mockResolvedValue({
      success: true,
      data: mockRoomListResponse
    });
    vi.mocked(apiClient.createRoom).mockResolvedValue({
      success: true,
      data: { id: 'new-room', status: 0, players: [null, null, null, null], owner: 'user1', created_at: new Date().toISOString() }
    });
    vi.mocked(apiClient.joinRoom).mockResolvedValue({
      success: true,
      data: { id: 'room1', status: 0, players: [null, null, null, null], owner: 'user1', created_at: new Date().toISOString() }
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('renders room lobby with action buttons', async () => {
    renderRoomLobby();
    
    await waitFor(() => {
      expect(screen.getByText('快速开始')).toBeInTheDocument();
    });
    expect(screen.getByText('创建房间')).toBeInTheDocument();
    expect(screen.getByText('加入游戏')).toBeInTheDocument();
  });

  it('shows login prompt when user is not logged in', () => {
    vi.mocked(useAuthStore).mockReturnValue({ user: null } as any);
    
    renderRoomLobby();
    
    expect(screen.getByText('请先登录')).toBeInTheDocument();
  });

  it('loads room list on mount', async () => {
    renderRoomLobby();
    
    await waitFor(() => {
      expect(apiClient.getRoomList).toHaveBeenCalledWith(1, 12);
    });
  });

  it('displays error message when present', async () => {
    const errorMessage = '获取房间列表失败';
    vi.mocked(useRoomStore).mockReturnValue({
      ...mockRoomStore,
      error: errorMessage
    } as any);
    
    renderRoomLobby();
    
    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
  });

  it('clears error when close button is clicked', async () => {
    const errorMessage = '获取房间列表失败';
    vi.mocked(useRoomStore).mockReturnValue({
      ...mockRoomStore,
      error: errorMessage
    } as any);
    
    renderRoomLobby();
    
    await waitFor(() => {
      expect(screen.getByText('✕')).toBeInTheDocument();
    });
    
    const closeButton = screen.getByText('✕');
    fireEvent.click(closeButton);
    
    expect(mockRoomStore.clearError).toHaveBeenCalled();
  });

  it('opens create room modal when create button is clicked', async () => {
    renderRoomLobby();
    
    await waitFor(() => {
      expect(screen.getByText('创建房间')).toBeInTheDocument();
    });
    
    const createButton = screen.getByText('创建房间');
    fireEvent.click(createButton);
    
    expect(screen.getByText('创建新房间')).toBeInTheDocument();
  });

  it('creates room when modal is confirmed', async () => {
    renderRoomLobby();
    
    await waitFor(() => {
      expect(screen.getByText('创建房间')).toBeInTheDocument();
    });
    
    // Open modal
    const createButton = screen.getByText('创建房间');
    fireEvent.click(createButton);
    
    // Confirm creation
    const confirmButton = screen.getByText('确认创建');
    fireEvent.click(confirmButton);
    
    await waitFor(() => {
      expect(apiClient.createRoom).toHaveBeenCalled();
    });
  });

  it('handles API errors gracefully', async () => {
    const errorMessage = 'Network error';
    vi.mocked(apiClient.getRoomList).mockRejectedValue(new Error(errorMessage));
    
    renderRoomLobby();
    
    await waitFor(() => {
      expect(mockRoomStore.setError).toHaveBeenCalledWith(errorMessage);
    });
  });
});