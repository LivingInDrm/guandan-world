import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import RoomList from './RoomList';
import { RoomStatus, type RoomInfo } from '../../types';

const mockRooms: RoomInfo[] = [
  {
    id: 'room1',
    status: RoomStatus.WAITING,
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
    status: RoomStatus.PLAYING,
    player_count: 4,
    players: [
      { id: 'user3', username: 'player3', seat: 0, online: true, auto_play: false },
      { id: 'user4', username: 'player4', seat: 1, online: true, auto_play: false },
      { id: 'user5', username: 'player5', seat: 2, online: true, auto_play: false },
      { id: 'user6', username: 'player6', seat: 3, online: true, auto_play: false }
    ],
    owner: 'user3',
    can_join: false
  },
  {
    id: 'room3',
    status: RoomStatus.WAITING,
    player_count: 1,
    players: [
      { id: 'user7', username: 'player7', seat: 0, online: true, auto_play: false }
    ],
    owner: 'user7',
    can_join: true
  }
];

const defaultProps = {
  rooms: mockRooms,
  isLoading: false,
  currentPage: 1,
  totalCount: 3,
  limit: 12,
  onPageChange: vi.fn(),
  onJoinRoom: vi.fn(),
  currentUserId: 'current-user'
};

describe('RoomList', () => {
  it('renders room rows correctly', () => {
    render(<RoomList {...defaultProps} />);
    
    expect(screen.getByText('#room1')).toBeInTheDocument();
    expect(screen.getByText('#room2')).toBeInTheDocument();
    expect(screen.getByText('#room3')).toBeInTheDocument();
  });

  it('sorts rooms correctly - waiting rooms first, then by player count', () => {
    render(<RoomList {...defaultProps} />);
    
    const roomIds = screen.getAllByText(/^#room/);
    
    expect(roomIds[0]).toHaveTextContent('#room1');
    expect(roomIds[1]).toHaveTextContent('#room3');
    expect(roomIds[2]).toHaveTextContent('#room2');
  });

  it('shows loading skeleton when loading and no rooms', () => {
    render(<RoomList {...defaultProps} rooms={[]} isLoading={true} />);
    
    const skeletonRows = document.querySelectorAll('tr.animate-pulse');
    expect(skeletonRows.length).toBeGreaterThan(0);
  });

  it('shows empty state when no rooms and not loading', () => {
    render(<RoomList {...defaultProps} rooms={[]} isLoading={false} totalCount={0} />);
    
    expect(screen.getByText('暂无房间')).toBeInTheDocument();
    expect(screen.getByText('成为第一个创建房间的玩家吧！')).toBeInTheDocument();
  });

  it('shows loading indicator when refreshing with existing rooms', () => {
    render(<RoomList {...defaultProps} isLoading={true} />);
    
    expect(screen.getByText('更新中...')).toBeInTheDocument();
  });

  it('shows pagination when total pages > 1', () => {
    render(<RoomList {...defaultProps} totalCount={25} />);
    
    expect(screen.getByText('上一页')).toBeInTheDocument();
    expect(screen.getByText('下一页')).toBeInTheDocument();
    expect(screen.getByText('1')).toBeInTheDocument();
  });

  it('does not show pagination when total pages <= 1', () => {
    render(<RoomList {...defaultProps} totalCount={5} />);
    
    expect(screen.queryByText('上一页')).not.toBeInTheDocument();
    expect(screen.queryByText('下一页')).not.toBeInTheDocument();
  });

  it('displays player count correctly', () => {
    render(<RoomList {...defaultProps} />);
    
    expect(screen.getByText('2/4')).toBeInTheDocument();
    expect(screen.getByText('4/4')).toBeInTheDocument();
    expect(screen.getByText('1/4')).toBeInTheDocument();
  });

  it('displays player names in mini cards', () => {
    render(<RoomList {...defaultProps} />);
    
    expect(screen.getByText('player1')).toBeInTheDocument();
    expect(screen.getByText('player2')).toBeInTheDocument();
  });

  it('shows empty slots for rooms not full', () => {
    render(<RoomList {...defaultProps} />);
    
    const emptySlots = screen.getAllByText('空位');
    expect(emptySlots.length).toBeGreaterThan(0);
  });
});
