import type { Story } from '@ladle/react';
import WaitingBoard from './WaitingBoard';
import type { Player } from '../../types';

const makePlayer = (id: string, username: string, seat: number, online = true): Player => ({
  id,
  username,
  seat,
  online,
  auto_play: false,
});

const mockPlayers: Player[] = [
  makePlayer('owner-1', '房主', 0),
  makePlayer('p2', '玩家B', 1),
  makePlayer('p3', '玩家C', 2),
  makePlayer('p4', '玩家D', 3),
];

const noop = () => {};

export const EmptyRoom: Story = () => (
  <WaitingBoard
    players={[null, null, null, null]}
    currentPlayerSeat={0}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="owner-1"
    isConnected={true}
    onStartGame={noop}
    onLeaveRoom={noop}
    isStarting={false}
    isLeaving={false}
  />
);
EmptyRoom.meta = {
  description: '空房间 - 无玩家加入',
};

export const OnePlayer: Story = () => (
  <WaitingBoard
    players={[mockPlayers[0], null, null, null]}
    currentPlayerSeat={0}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="owner-1"
    isConnected={true}
    onStartGame={noop}
    onLeaveRoom={noop}
    isStarting={false}
    isLeaving={false}
  />
);
OnePlayer.meta = {
  description: '单人房间 - 仅房主加入',
};

export const PartialPlayers: Story = () => (
  <WaitingBoard
    players={[mockPlayers[0], mockPlayers[1], null, null]}
    currentPlayerSeat={0}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="owner-1"
    isConnected={true}
    onStartGame={noop}
    onLeaveRoom={noop}
    isStarting={false}
    isLeaving={false}
  />
);
PartialPlayers.meta = {
  description: '部分玩家 - 2人加入',
};

export const FullRoom: Story = () => (
  <WaitingBoard
    players={mockPlayers}
    currentPlayerSeat={0}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="owner-1"
    isConnected={true}
    onStartGame={noop}
    onLeaveRoom={noop}
    isStarting={false}
    isLeaving={false}
  />
);
FullRoom.meta = {
  description: '满员房间 - 4人全部加入',
};

export const OwnerView: Story = () => (
  <WaitingBoard
    players={mockPlayers}
    currentPlayerSeat={0}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="owner-1"
    isConnected={true}
    onStartGame={() => console.log('开始游戏')}
    onLeaveRoom={() => console.log('离开房间')}
    isStarting={false}
    isLeaving={false}
  />
);
OwnerView.meta = {
  description: '房主视角 - 可以开始游戏',
};

export const NonOwnerView: Story = () => (
  <WaitingBoard
    players={mockPlayers}
    currentPlayerSeat={1}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="p2"
    isConnected={true}
    onStartGame={noop}
    onLeaveRoom={() => console.log('离开房间')}
    isStarting={false}
    isLeaving={false}
  />
);
NonOwnerView.meta = {
  description: '非房主视角 - 等待开始',
};

export const StartingGame: Story = () => (
  <WaitingBoard
    players={mockPlayers}
    currentPlayerSeat={0}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="owner-1"
    isConnected={true}
    onStartGame={noop}
    onLeaveRoom={noop}
    isStarting={true}
    isLeaving={false}
  />
);
StartingGame.meta = {
  description: '开始游戏中状态',
};

export const LeavingRoom: Story = () => (
  <WaitingBoard
    players={mockPlayers}
    currentPlayerSeat={0}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="owner-1"
    isConnected={true}
    onStartGame={noop}
    onLeaveRoom={noop}
    isStarting={false}
    isLeaving={true}
  />
);
LeavingRoom.meta = {
  description: '离开房间中状态',
};

export const OfflinePlayer: Story = () => (
  <WaitingBoard
    players={[
      mockPlayers[0],
      makePlayer('p2', '玩家B', 1, false),
      mockPlayers[2],
      mockPlayers[3],
    ]}
    currentPlayerSeat={0}
    roomId="room-1"
    roomCode="8888"
    ownerId="owner-1"
    currentUserId="owner-1"
    isConnected={true}
    onStartGame={noop}
    onLeaveRoom={noop}
    isStarting={false}
    isLeaving={false}
  />
);
OfflinePlayer.meta = {
  description: '包含离线玩家',
};
