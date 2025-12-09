import type { Story } from "@ladle/react";
import { useState } from "react";
import RoomList from "./RoomList";
import { RoomStatus, type RoomInfo, type Player } from "../../types";

const createPlayer = (id: string, seat: number, online = true): Player => ({
  id,
  username: `user_${id}`,
  nickname: `玩家${id}`,
  avatar_key: undefined,
  seat,
  online,
  auto_play: false,
});

const createRoom = (
  id: string,
  status: RoomStatus,
  playerCount: number,
  canJoin = true
): RoomInfo => {
  const players: Player[] = [];
  for (let i = 0; i < playerCount; i++) {
    players.push(createPlayer(`${id}-p${i}`, i, i !== 2));
  }
  return {
    id: `room-${id}`,
    status,
    player_count: playerCount,
    players,
    owner: players[0]?.id || "",
    can_join: canJoin,
  };
};

const defaultProps = {
  currentPage: 1,
  totalCount: 6,
  limit: 10,
  onPageChange: () => {},
  onJoinRoom: () => {},
  currentUserId: "current-user",
};

export const Loading: Story = () => (
  <div className="p-8 bg-background min-h-screen">
    <RoomList {...defaultProps} rooms={[]} isLoading={true} />
  </div>
);

export const Empty: Story = () => (
  <div className="p-8 bg-background min-h-screen">
    <RoomList {...defaultProps} rooms={[]} isLoading={false} totalCount={0} />
  </div>
);

export const Default: Story = () => {
  const rooms: RoomInfo[] = [
    createRoom("001", RoomStatus.WAITING, 2),
    createRoom("002", RoomStatus.WAITING, 1),
    createRoom("003", RoomStatus.PLAYING, 4, false),
    createRoom("004", RoomStatus.READY, 4, false),
    createRoom("005", RoomStatus.WAITING, 3),
    createRoom("006", RoomStatus.CLOSED, 0, false),
  ];

  return (
    <div className="p-8 bg-background min-h-screen">
      <RoomList {...defaultProps} rooms={rooms} isLoading={false} />
    </div>
  );
};

export const WaitingRooms: Story = () => {
  const rooms: RoomInfo[] = [
    createRoom("w01", RoomStatus.WAITING, 1),
    createRoom("w02", RoomStatus.WAITING, 2),
    createRoom("w03", RoomStatus.WAITING, 3),
  ];

  return (
    <div className="p-8 bg-background min-h-screen">
      <RoomList {...defaultProps} rooms={rooms} isLoading={false} totalCount={3} />
    </div>
  );
};

export const PlayingRooms: Story = () => {
  const rooms: RoomInfo[] = [
    createRoom("p01", RoomStatus.PLAYING, 4, false),
    createRoom("p02", RoomStatus.PLAYING, 4, false),
  ];

  return (
    <div className="p-8 bg-background min-h-screen">
      <RoomList {...defaultProps} rooms={rooms} isLoading={false} totalCount={2} />
    </div>
  );
};

export const MixedStatus: Story = () => {
  const rooms: RoomInfo[] = [
    createRoom("m01", RoomStatus.PLAYING, 4, false),
    createRoom("m02", RoomStatus.WAITING, 2),
    createRoom("m03", RoomStatus.CLOSED, 0, false),
    createRoom("m04", RoomStatus.READY, 4, false),
    createRoom("m05", RoomStatus.WAITING, 1),
  ];

  return (
    <div className="p-8 bg-background min-h-screen">
      <p className="text-muted-foreground mb-4 text-sm">
        房间按状态排序: 等待中 &gt; 准备中 &gt; 游戏中 &gt; 已关闭
      </p>
      <RoomList {...defaultProps} rooms={rooms} isLoading={false} totalCount={5} />
    </div>
  );
};

export const UserInRoom: Story = () => {
  const currentUserId = "user-in-room";
  const rooms: RoomInfo[] = [
    {
      id: "room-user",
      status: RoomStatus.WAITING,
      player_count: 2,
      players: [
        createPlayer("other-1", 0),
        { ...createPlayer(currentUserId, 1), nickname: "我" },
      ],
      owner: "other-1",
      can_join: false,
    },
    createRoom("other", RoomStatus.WAITING, 1),
  ];

  return (
    <div className="p-8 bg-background min-h-screen">
      <p className="text-muted-foreground mb-4 text-sm">
        第一个房间中当前用户已加入（显示"返回"按钮和高亮）
      </p>
      <RoomList
        {...defaultProps}
        rooms={rooms}
        isLoading={false}
        totalCount={2}
        currentUserId={currentUserId}
      />
    </div>
  );
};

export const WithPagination: Story = () => {
  const [page, setPage] = useState(1);
  const rooms: RoomInfo[] = [
    createRoom("pg1", RoomStatus.WAITING, 2),
    createRoom("pg2", RoomStatus.WAITING, 3),
  ];

  return (
    <div className="p-8 bg-background min-h-screen">
      <p className="text-muted-foreground mb-4 text-sm">当前页: {page}</p>
      <RoomList
        {...defaultProps}
        rooms={rooms}
        isLoading={false}
        currentPage={page}
        totalCount={50}
        limit={10}
        onPageChange={setPage}
      />
    </div>
  );
};

export const Interactive: Story = () => {
  const [page, setPage] = useState(1);
  const [lastAction, setLastAction] = useState<string>("等待操作");

  const rooms: RoomInfo[] = [
    createRoom("int1", RoomStatus.WAITING, 2),
    createRoom("int2", RoomStatus.WAITING, 1),
    createRoom("int3", RoomStatus.PLAYING, 4, false),
  ];

  return (
    <div className="p-8 bg-background min-h-screen">
      <p className="text-muted-foreground mb-4 text-sm">操作记录: {lastAction}</p>
      <RoomList
        {...defaultProps}
        rooms={rooms}
        isLoading={false}
        currentPage={page}
        totalCount={30}
        limit={10}
        onPageChange={(p) => {
          setPage(p);
          setLastAction(`切换到第 ${p} 页`);
        }}
        onJoinRoom={(roomId) => {
          setLastAction(`加入房间: ${roomId}`);
        }}
      />
    </div>
  );
};
