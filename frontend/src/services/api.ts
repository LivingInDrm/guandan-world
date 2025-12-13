import type { 
  ApiResponse, 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse,
  RoomListResponse,
  CreateRoomRequest,
  Room,
  User
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

class ApiError extends Error {
  public status: number;
  public response?: any;

  constructor(
    message: string,
    status: number,
    response?: any
  ) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.response = response;
  }
}

class ApiClient {
  private baseURL: string;
  private token: string | null = null;
  private refreshPromise: Promise<AuthResponse> | null = null;
  private onUnauthorized: (() => void) | null = null;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  setToken(token: string | null) {
    this.token = token;
  }

  setOnUnauthorized(callback: () => void) {
    this.onUnauthorized = callback;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
    skipAuth: boolean = false
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };

    if (this.token && !skipAuth) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new ApiError(
          data.message || data.error || 'Request failed',
          response.status,
          data
        );
      }

      return {
        success: true,
        data: data.data || data,
        message: data.message
      };
    } catch (error) {
      if (error instanceof ApiError) {
        if (error.status === 401 && error.response?.error === 'token_expired') {
          if (this.onUnauthorized) {
            this.onUnauthorized();
          }
        }
        throw error;
      }
      
      throw new ApiError(
        error instanceof Error ? error.message : 'Network error',
        0
      );
    }
  }

  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    }, true);
    return response.data!;
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    }, true);
    return response.data!;
  }

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    this.refreshPromise = (async () => {
      try {
        const response = await this.request<AuthResponse>('/api/auth/refresh', {
          method: 'POST',
          body: JSON.stringify({ refresh_token: refreshToken }),
        }, true);
        return response.data!;
      } finally {
        this.refreshPromise = null;
      }
    })();

    return this.refreshPromise;
  }

  async logout(refreshToken?: string): Promise<ApiResponse<void>> {
    const response = await this.request<void>('/api/auth/logout', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken || '' }),
    });
    this.token = null;
    return response;
  }

  async getRoomList(page: number = 1, limit: number = 12): Promise<ApiResponse<RoomListResponse>> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    
    return this.request<RoomListResponse>(`/api/rooms?${params}`);
  }

  async createRoom(roomData: CreateRoomRequest = {}): Promise<ApiResponse<Room>> {
    const response = await this.request<{ room: Room }>('/api/rooms/create', {
      method: 'POST',
      body: JSON.stringify(roomData),
    });
    return {
      ...response,
      data: response.data?.room as Room
    };
  }

  async joinRoom(roomId: string): Promise<ApiResponse<Room>> {
    const response = await this.request<{ room: Room }>(`/api/rooms/${roomId}/join`, {
      method: 'POST',
    });
    return {
      ...response,
      data: response.data?.room as Room
    };
  }

  async joinRoomByCode(roomCode: string): Promise<ApiResponse<Room>> {
    const response = await this.request<{ room: Room }>('/api/rooms/join-by-code', {
      method: 'POST',
      body: JSON.stringify({ room_code: roomCode }),
    });
    return {
      ...response,
      data: response.data?.room as Room
    };
  }

  async quickJoin(): Promise<ApiResponse<Room>> {
    const response = await this.request<{ room: Room }>('/api/rooms/quick-join', {
      method: 'POST',
    });
    return {
      ...response,
      data: response.data?.room as Room
    };
  }

  async leaveRoom(roomId: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/rooms/${roomId}/leave`, {
      method: 'POST',
    });
  }

  async getRoomDetails(roomId: string): Promise<ApiResponse<Room>> {
    const response = await this.request<{ room: Room }>(`/api/rooms/${roomId}`);
    return {
      ...response,
      data: response.data?.room as Room
    };
  }

  async getMyRoom(): Promise<ApiResponse<Room>> {
    const response = await this.request<{ room: Room }>('/api/rooms/my');
    return {
      ...response,
      data: response.data?.room as Room
    };
  }

  async startGame(roomId: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/rooms/${roomId}/start`, {
      method: 'POST',
    });
  }

  async playCards(roomId: string, playerSeat: number, deckIndexes: number[]): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/game/driver/play-decision`, {
      method: 'POST',
      body: JSON.stringify({
        room_id: roomId,
        player_seat: playerSeat,
        action: 'play',
        deck_indexes: deckIndexes
      }),
    });
  }

  async pass(roomId: string, playerSeat: number): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/game/driver/play-decision`, {
      method: 'POST',
      body: JSON.stringify({
        room_id: roomId,
        player_seat: playerSeat,
        action: 'pass',
        deck_indexes: []
      }),
    });
  }

  async selectTribute(roomId: string, playerSeat: number, deckIndex: number): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/game/driver/tribute-select`, {
      method: 'POST',
      body: JSON.stringify({
        room_id: roomId,
        player_seat: playerSeat,
        deck_index: deckIndex
      }),
    });
  }

  async returnTribute(roomId: string, playerSeat: number, deckIndex: number): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/game/driver/tribute-return`, {
      method: 'POST',
      body: JSON.stringify({
        room_id: roomId,
        player_seat: playerSeat,
        deck_index: deckIndex
      }),
    });
  }

  async updateProfile(data: { nickname?: string; avatar?: File }): Promise<ApiResponse<User>> {
    const url = `${this.baseURL}/api/user/profile`;
    const formData = new FormData();
    
    if (data.nickname !== undefined) {
      formData.append('nickname', data.nickname);
    }
    if (data.avatar) {
      formData.append('avatar', data.avatar);
    }

    const headers: Record<string, string> = {};
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    try {
      const response = await fetch(url, {
        method: 'PATCH',
        headers,
        body: formData,
      });

      const responseData = await response.json();

      if (!response.ok) {
        throw new ApiError(
          responseData.message || responseData.error || 'Update failed',
          response.status,
          responseData
        );
      }

      return {
        success: true,
        data: responseData.user || responseData.data?.user,
        message: responseData.message
      };
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      throw new ApiError(
        error instanceof Error ? error.message : 'Network error',
        0
      );
    }
  }

  async getProfile(): Promise<ApiResponse<User>> {
    const response = await this.request<{ user: User }>('/api/user/profile');
    return {
      ...response,
      data: response.data?.user as User
    };
  }

  async healthCheck(): Promise<ApiResponse<{ status: string }>> {
    return this.request<{ status: string }>('/healthz');
  }
}

export const apiClient = new ApiClient();

export { ApiClient, ApiError };
