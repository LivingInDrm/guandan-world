import React, { useEffect } from 'react';
import type { Story } from '@ladle/react';
import { MemoryRouter } from 'react-router-dom';
import type { Player } from '../../types';
import { useAuthStore } from '../../store/authStore';
import UserMenuFab from '../layout/UserMenuFab';
import WaitingBoard from './WaitingBoard';

const MockAuthProvider: React.FC<{ children: React.ReactNode; userId?: string }> = ({ 
  children, 
  userId = 'owner-1' 
}) => {
  const { login, logout } = useAuthStore();
  useEffect(() => {
    login(
      { id: userId, username: 'testuser', nickname: '测试用户', online: true },
      { access_token: 'mock-token', refresh_token: 'mock-refresh', expires_at: new Date(Date.now() + 3600000).toISOString(), user_id: userId }
    );
    return () => logout();
  }, [login, logout, userId]);
  return <>{children}</>;
};

const withProviders = (Story: React.ComponentType) => (
  <MemoryRouter>
    <MockAuthProvider>
      <Story />
    </MockAuthProvider>
  </MemoryRouter>
);

const withProvidersAsGuest = (Story: React.ComponentType) => (
  <MemoryRouter>
    <MockAuthProvider userId="p2">
      <Story />
    </MockAuthProvider>
  </MemoryRouter>
);

const makePlayer = (id: string, username: string, seat: number, online = true): Player => ({
  id,
  username,
  seat,
  online,
  auto_play: false,
});

const mockOwner = makePlayer('owner-1', '房主', 0);
const mockPlayers: Player[] = [
  mockOwner,
  makePlayer('p2', '玩家B', 1),
  makePlayer('p3', '玩家C', 2),
  makePlayer('p4', '玩家D', 3),
];

interface WaitingPhaseWrapperProps {
  players: (Player | null)[];
  currentPlayerSeat: number;
  roomCode?: string;
  ownerId?: string;
  currentUserId?: string;
  isConnected?: boolean;
  isStarting?: boolean;
  isLeaving?: boolean;
}

const WaitingPhaseWrapper: React.FC<WaitingPhaseWrapperProps> = ({
  players,
  currentPlayerSeat,
  roomCode = '8888',
  ownerId = 'owner-1',
  currentUserId = 'owner-1',
  isConnected = true,
  isStarting = false,
  isLeaving = false,
}) => {
  const noop = () => {};

  return (
    <div className="fixed inset-0 z-40 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
      <UserMenuFab />
      <div className="absolute inset-0 p-2">
        <WaitingBoard
          players={players}
          currentPlayerSeat={currentPlayerSeat}
          roomId="room-1"
          roomCode={roomCode}
          ownerId={ownerId}
          currentUserId={currentUserId}
          isConnected={isConnected}
          onStartGame={isStarting ? noop : () => console.log('开始游戏')}
          onLeaveRoom={isLeaving ? noop : () => console.log('离开房间')}
          isStarting={isStarting}
          isLeaving={isLeaving}
          className="h-full"
        />
      </div>
    </div>
  );
};

export const EmptyRoom: Story = () => (
  <WaitingPhaseWrapper
    players={[mockOwner, null, null, null]}
    currentPlayerSeat={0}
  />
);
EmptyRoom.meta = { description: '空房间 - 仅房主加入' };
EmptyRoom.decorators = [withProviders];

export const PartialPlayers: Story = () => (
  <WaitingPhaseWrapper
    players={[mockPlayers[0], mockPlayers[1], null, null]}
    currentPlayerSeat={0}
  />
);
PartialPlayers.meta = { description: '部分玩家 - 2人加入' };
PartialPlayers.decorators = [withProviders];

export const FullRoomOwner: Story = () => (
  <WaitingPhaseWrapper
    players={mockPlayers}
    currentPlayerSeat={0}
    currentUserId="owner-1"
  />
);
FullRoomOwner.meta = { description: '满员房间 - 房主视角，可开始游戏' };
FullRoomOwner.decorators = [withProviders];

export const FullRoomGuest: Story = () => (
  <WaitingPhaseWrapper
    players={mockPlayers}
    currentPlayerSeat={1}
    currentUserId="p2"
  />
);
FullRoomGuest.meta = { description: '满员房间 - 非房主视角，等待开始' };
FullRoomGuest.decorators = [withProvidersAsGuest];

export const StartingGame: Story = () => (
  <WaitingPhaseWrapper
    players={mockPlayers}
    currentPlayerSeat={0}
    isStarting={true}
  />
);
StartingGame.meta = { description: '开始游戏中状态' };
StartingGame.decorators = [withProviders];

export const LeavingRoom: Story = () => (
  <WaitingPhaseWrapper
    players={mockPlayers}
    currentPlayerSeat={0}
    isLeaving={true}
  />
);
LeavingRoom.meta = { description: '离开房间中状态' };
LeavingRoom.decorators = [withProviders];

export const WithOfflinePlayer: Story = () => (
  <WaitingPhaseWrapper
    players={[
      mockPlayers[0],
      makePlayer('p2', '玩家B', 1, false),
      mockPlayers[2],
      mockPlayers[3],
    ]}
    currentPlayerSeat={0}
  />
);
WithOfflinePlayer.meta = { description: '包含离线玩家' };
WithOfflinePlayer.decorators = [withProviders];
