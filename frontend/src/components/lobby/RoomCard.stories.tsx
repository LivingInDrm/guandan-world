import type { Story } from "@ladle/react";
import { useState } from "react";
import RoomCard from "./RoomCard";
import { RoomStatus, type Player, type RoomInfo } from "../../types";

const mockPlayer = (
  id: string,
  username: string,
  seat: number = 0,
  online: boolean = true
): Player => ({
  id,
  username,
  seat,
  online,
  auto_play: false,
});

const mockRoom = (overrides: Partial<RoomInfo> = {}): RoomInfo => ({
  id: "room-abc123456",
  status: RoomStatus.WAITING,
  player_count: 2,
  players: [
    mockPlayer("user-1", "玩家一", 0),
    mockPlayer("user-2", "玩家二", 1),
  ],
  owner: "user-1",
  can_join: true,
  ...overrides,
});

export const Default: Story = () => (
  <div className="p-8 bg-gray-800">
    <div className="w-80">
      <RoomCard
        room={mockRoom()}
        onJoinRoom={(id) => console.log("Join room:", id)}
        currentUserId="user-999"
      />
    </div>
  </div>
);

export const AllRoomStatuses: Story = () => (
  <div className="p-8 bg-gray-800 grid grid-cols-2 gap-4">
    <div>
      <p className="text-white text-sm mb-2">WAITING</p>
      <RoomCard
        room={mockRoom({ status: RoomStatus.WAITING })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
    <div>
      <p className="text-white text-sm mb-2">READY</p>
      <RoomCard
        room={mockRoom({ status: RoomStatus.READY, player_count: 4, can_join: false })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
    <div>
      <p className="text-white text-sm mb-2">PLAYING</p>
      <RoomCard
        room={mockRoom({ status: RoomStatus.PLAYING, player_count: 4, can_join: false })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
    <div>
      <p className="text-white text-sm mb-2">CLOSED</p>
      <RoomCard
        room={mockRoom({ status: RoomStatus.CLOSED, can_join: false })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
  </div>
);

export const FullRoom: Story = () => (
  <div className="p-8 bg-gray-800">
    <div className="w-80">
      <RoomCard
        room={mockRoom({
          player_count: 4,
          players: [
            mockPlayer("user-1", "玩家一", 0),
            mockPlayer("user-2", "玩家二", 1),
            mockPlayer("user-3", "玩家三", 2),
            mockPlayer("user-4", "玩家四", 3),
          ],
          can_join: false,
        })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
  </div>
);

export const EmptyRoom: Story = () => (
  <div className="p-8 bg-gray-800">
    <div className="w-80">
      <RoomCard
        room={mockRoom({
          player_count: 0,
          players: [],
        })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
  </div>
);

export const AsOwner: Story = () => (
  <div className="p-8 bg-gray-800">
    <div className="w-80">
      <RoomCard
        room={mockRoom()}
        onJoinRoom={() => {}}
        currentUserId="user-1"
      />
    </div>
  </div>
);

export const AsPlayer: Story = () => (
  <div className="p-8 bg-gray-800">
    <div className="w-80">
      <RoomCard
        room={mockRoom()}
        onJoinRoom={() => {}}
        currentUserId="user-2"
      />
    </div>
  </div>
);

export const CannotJoin: Story = () => (
  <div className="p-8 bg-gray-800 flex gap-4">
    <div className="w-80">
      <p className="text-white text-sm mb-2">房间已满</p>
      <RoomCard
        room={mockRoom({
          player_count: 4,
          players: [
            mockPlayer("user-1", "玩家一", 0),
            mockPlayer("user-2", "玩家二", 1),
            mockPlayer("user-3", "玩家三", 2),
            mockPlayer("user-4", "玩家四", 3),
          ],
          can_join: false,
        })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
    <div className="w-80">
      <p className="text-white text-sm mb-2">游戏中</p>
      <RoomCard
        room={mockRoom({
          status: RoomStatus.PLAYING,
          player_count: 4,
          can_join: false,
        })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
  </div>
);

export const OnlineOfflinePlayers: Story = () => (
  <div className="p-8 bg-gray-800">
    <div className="w-80">
      <RoomCard
        room={mockRoom({
          player_count: 4,
          players: [
            mockPlayer("user-1", "在线玩家", 0, true),
            mockPlayer("user-2", "离线玩家", 1, false),
            mockPlayer("user-3", "在线玩家", 2, true),
            mockPlayer("user-4", "离线玩家", 3, false),
          ],
          can_join: false,
        })}
        onJoinRoom={() => {}}
        currentUserId="user-999"
      />
    </div>
  </div>
);

export const Interactive: Story = () => {
  const [lastAction, setLastAction] = useState<string>("");

  return (
    <div className="p-8 bg-gray-800">
      <div className="flex gap-4">
        <div className="w-80">
          <p className="text-white text-sm mb-2">可加入</p>
          <RoomCard
            room={mockRoom()}
            onJoinRoom={(id) => setLastAction(`加入房间: ${id}`)}
            currentUserId="user-999"
          />
        </div>
        <div className="w-80">
          <p className="text-white text-sm mb-2">已在房间</p>
          <RoomCard
            room={mockRoom()}
            onJoinRoom={(id) => setLastAction(`返回房间: ${id}`)}
            currentUserId="user-2"
          />
        </div>
      </div>
      {lastAction && (
        <div className="mt-4 p-3 bg-green-800 text-white rounded">
          {lastAction}
        </div>
      )}
    </div>
  );
};
