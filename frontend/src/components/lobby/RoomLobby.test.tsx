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

  it('renders room lobby header correctly', () => {
    renderRoomLobby();
    
    expect(screen.getByText('房间大厅')).toBeInTheDocument();
    expect(screen.getByText('欢迎，testuser！选择一个房间开始游戏')).toBeInTheDocument();
    expect(screen.getByText('创建房间')).toBeInTheDocument();
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

  it('displays error message when present', () => {
    const errorMessage = '获取房间列表失败';
    vi.mocked(useRoomStore).mockReturnValue({
      ...mockRoomStore,
      error: errorMessage
    } as any);
    
    renderRoomLobby();
    
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });

  it('clears error when close button is clicked', () => {
    const errorMessage = '获取房间列表失败';
    vi.mocked(useRoomStore).mockReturnValue({
      ...mockRoomStore,
      error: errorMessage
    } as any);
    
    renderRoomLobby();
    
    const closeButton = screen.getByText('✕');
    fireEvent.click(closeButton);
    
    expect(mockRoomStore.clearError).toHaveBeenCalled();
  });

  it('opens create room modal when create button is clicked', () => {
    renderRoomLobby();
    
    const createButton = screen.getByText('创建房间');
    fireEvent.click(createButton);
    
    expect(screen.getByText('创建新房间')).toBeInTheDocument();
  });

  it('creates room when modal is confirmed', async () => {
    renderRoomLobby();
    
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

  it('refreshes room list manually when refresh button is clicked', async () => {
    renderRoomLobby();
    
    const refreshButton = screen.getByText('手动刷新');
    fireEvent.click(refreshButton);
    
    await waitFor(() => {
      expect(apiClient.getRoomList).toHaveBeenCalledTimes(2); // Once on mount, once on manual refresh
    });
  });

  it('displays room statistics correctly', () => {
    vi.mocked(useRoomStore).mockReturnValue({
      ...mockRoomStore,
      roomList: mockRoomListResponse.rooms,
      totalCount: mockRoomListResponse.total_count
    } as any);
    
    renderRoomLobby();
    
    expect(screen.getByText('2')).toBeInTheDocument(); // Total count
    expect(screen.getByText('总房间数')).toBeInTheDocument();
    
    // Check statistics section exists
    const statisticsSection = screen.getByText('总房间数').closest('.bg-gray-50');
    expect(statisticsSection).toBeInTheDocument();
    
    // Check that statistics labels exist (using getAllByText since they appear in both stats and room cards)
    const waitingTexts = screen.getAllByText('等待中');
    const playingTexts = screen.getAllByText('游戏中');
    
    expect(waitingTexts.length).toBeGreaterThan(0);
    expect(playingTexts.length).toBeGreaterThan(0);
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