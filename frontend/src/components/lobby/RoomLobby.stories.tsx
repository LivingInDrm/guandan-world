import type { Story } from "@ladle/react";
import { MemoryRouter } from "react-router-dom";
import { useEffect } from "react";
import RoomLobby from "./RoomLobby";
import { useAuthStore } from "../../store/authStore";
import { useRoomStore } from "../../store/roomStore";
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

const mockUser = {
  id: "current-user",
  username: "testuser",
  nickname: "测试玩家",
  avatar_key: undefined,
  online: true,
};

const mockRooms: RoomInfo[] = [
  createRoom("001", RoomStatus.WAITING, 2),
  createRoom("002", RoomStatus.WAITING, 1),
  createRoom("003", RoomStatus.PLAYING, 4, false),
  createRoom("004", RoomStatus.READY, 4, false),
  createRoom("005", RoomStatus.WAITING, 3),
  createRoom("006", RoomStatus.CLOSED, 0, false),
];

const Wrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <MemoryRouter initialEntries={["/lobby"]}>
    <div className="bg-background min-h-screen">{children}</div>
  </MemoryRouter>
);

const setupStores = (options: {
  isAuthenticated?: boolean;
  roomList?: RoomInfo[];
  totalCount?: number;
  isLoading?: boolean;
  error?: string | null;
  currentPage?: number;
}) => {
  const {
    isAuthenticated = true,
    roomList = [],
    totalCount = 0,
    isLoading = false,
    error = null,
    currentPage = 1,
  } = options;

  if (isAuthenticated) {
    useAuthStore.setState({
      user: mockUser,
      isAuthenticated: true,
    });
  } else {
    useAuthStore.setState({
      user: null,
      isAuthenticated: false,
    });
  }

  useRoomStore.setState({
    roomList,
    totalCount,
    isLoading,
    error,
    currentPage,
    limit: 10,
  });
};

export const NotAuthenticated: Story = () => {
  useEffect(() => {
    setupStores({ isAuthenticated: false });
  }, []);

  return (
    <Wrapper>
      <RoomLobby />
    </Wrapper>
  );
};

export const Loading: Story = () => {
  useEffect(() => {
    setupStores({ isLoading: true, roomList: [] });
  }, []);

  return (
    <Wrapper>
      <RoomLobby />
    </Wrapper>
  );
};

export const Empty: Story = () => {
  useEffect(() => {
    setupStores({ roomList: [], totalCount: 0 });
  }, []);

  return (
    <Wrapper>
      <RoomLobby />
    </Wrapper>
  );
};

export const WithRooms: Story = () => {
  useEffect(() => {
    setupStores({ roomList: mockRooms, totalCount: mockRooms.length });
  }, []);

  return (
    <Wrapper>
      <RoomLobby />
    </Wrapper>
  );
};

export const WithError: Story = () => {
  useEffect(() => {
    setupStores({
      roomList: mockRooms.slice(0, 2),
      totalCount: 2,
      error: "网络连接失败，请稍后重试",
    });
  }, []);

  return (
    <Wrapper>
      <RoomLobby />
    </Wrapper>
  );
};

export const WithPagination: Story = () => {
  useEffect(() => {
    setupStores({
      roomList: mockRooms.slice(0, 4),
      totalCount: 50,
      currentPage: 2,
    });
  }, []);

  return (
    <Wrapper>
      <RoomLobby />
    </Wrapper>
  );
};

export const UserInRoom: Story = () => {
  const currentUserId = "current-user";
  const roomsWithUser: RoomInfo[] = [
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
    ...mockRooms.slice(0, 3),
  ];

  useEffect(() => {
    useAuthStore.setState({
      user: { ...mockUser, id: currentUserId },
      isAuthenticated: true,
    });
    useRoomStore.setState({
      roomList: roomsWithUser,
      totalCount: roomsWithUser.length,
      isLoading: false,
      error: null,
      currentPage: 1,
      limit: 10,
    });
  }, []);

  return (
    <Wrapper>
      <RoomLobby />
    </Wrapper>
  );
};

export const ManyRooms: Story = () => {
  const manyRooms: RoomInfo[] = [
    createRoom("m01", RoomStatus.WAITING, 1),
    createRoom("m02", RoomStatus.WAITING, 2),
    createRoom("m03", RoomStatus.WAITING, 3),
    createRoom("m04", RoomStatus.READY, 4, false),
    createRoom("m05", RoomStatus.PLAYING, 4, false),
    createRoom("m06", RoomStatus.WAITING, 1),
    createRoom("m07", RoomStatus.WAITING, 2),
    createRoom("m08", RoomStatus.PLAYING, 4, false),
  ];

  useEffect(() => {
    setupStores({ roomList: manyRooms, totalCount: 80 });
  }, []);

  return (
    <Wrapper>
      <RoomLobby />
    </Wrapper>
  );
};
