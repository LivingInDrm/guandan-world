import { create } from 'zustand';
import type { Room, RoomInfo, RoomListResponse } from '../types';

interface RoomState {
  currentRoom: Room | null;
  roomList: RoomInfo[];
  totalCount: number;
  currentPage: number;
  limit: number;
  isLoading: boolean;
  error: string | null;
}

interface RoomActions {
  setCurrentRoom: (room: Room | null) => void;
  setRoomList: (response: RoomListResponse) => void;
  updateRoomInList: (roomInfo: RoomInfo) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  setPage: (page: number) => void;
  reset: () => void;
}

type RoomStore = RoomState & RoomActions;

const initialState: RoomState = {
  currentRoom: null,
  roomList: [],
  totalCount: 0,
  currentPage: 1,
  limit: 12,
  isLoading: false,
  error: null
};

export const useRoomStore = create<RoomStore>((set) => ({
  ...initialState,

  // Actions
  setCurrentRoom: (room: Room | null) => set({ currentRoom: room }),
  
  setRoomList: (response: RoomListResponse) => set({
    roomList: response.rooms,
    totalCount: response.total_count,
    currentPage: response.page,
    limit: response.limit
  }),
  
  updateRoomInList: (roomInfo: RoomInfo) => set((state) => ({
    roomList: state.roomList.map(room => 
      room.id === roomInfo.id ? roomInfo : room
    )
  })),
  
  setLoading: (loading: boolean) => set({ isLoading: loading }),
  
  setError: (error: string | null) => set({ error }),
  
  clearError: () => set({ error: null }),
  
  setPage: (page: number) => set({ currentPage: page }),
  
  reset: () => set(initialState)
}));