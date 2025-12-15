import type {
  ApiResponse,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  Room,
  User
} from './types.js';

class ApiError extends Error {
  public status: number;
  public response?: unknown;

  constructor(
    message: string,
    status: number,
    response?: unknown
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
  private refreshTokenValue: string | null = null;
  private refreshPromise: Promise<AuthResponse> | null = null;
  private onUnauthorized: (() => void) | null = null;
  private onTokenRefreshed: ((auth: AuthResponse) => void) | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  setToken(token: string | null) {
    this.token = token;
  }

  setRefreshToken(token: string | null) {
    this.refreshTokenValue = token;
  }

  setOnUnauthorized(callback: () => void) {
    this.onUnauthorized = callback;
  }

  setOnTokenRefreshed(callback: (auth: AuthResponse) => void) {
    this.onTokenRefreshed = callback;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
    skipAuth: boolean = false,
    isRetry: boolean = false
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

      const data = await response.json() as Record<string, unknown>;

      if (!response.ok) {
        throw new ApiError(
          (data.message || data.error || 'Request failed') as string,
          response.status,
          data
        );
      }

      return {
        success: true,
        data: (data.data || data) as T,
        message: data.message as string | undefined
      };
    } catch (error) {
      if (error instanceof ApiError) {
        const isTokenExpired = error.status === 401 && 
          (error.response as Record<string, unknown>)?.error === 'token_expired';
        
        if (isTokenExpired && !isRetry && this.refreshTokenValue) {
          try {
            const auth = await this.refreshToken(this.refreshTokenValue);
            this.token = auth.token.access_token;
            this.refreshTokenValue = auth.token.refresh_token;
            if (this.onTokenRefreshed) {
              this.onTokenRefreshed(auth);
            }
            return this.request<T>(endpoint, options, skipAuth, true);
          } catch {
            if (this.onUnauthorized) {
              this.onUnauthorized();
            }
            throw error;
          }
        }
        
        if (isTokenExpired && this.onUnauthorized) {
          this.onUnauthorized();
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

  async joinRoomByCode(roomCode: string): Promise<Room> {
    const response = await this.request<{ room: Room }>('/api/rooms/join-by-code', {
      method: 'POST',
      body: JSON.stringify({ room_code: roomCode }),
    });
    return response.data!.room;
  }

  async leaveRoom(roomId: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/rooms/${roomId}/leave`, {
      method: 'POST',
    });
  }

  async startGame(roomId: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/rooms/${roomId}/start`, {
      method: 'POST',
    });
  }

  async playCards(roomId: string, playerSeat: number, deckIndexes: number[]): Promise<ApiResponse<void>> {
    return this.request<void>('/api/game/driver/play-decision', {
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
    return this.request<void>('/api/game/driver/play-decision', {
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
    return this.request<void>('/api/game/driver/tribute-select', {
      method: 'POST',
      body: JSON.stringify({
        room_id: roomId,
        player_seat: playerSeat,
        deck_index: deckIndex
      }),
    });
  }

  async returnTribute(roomId: string, playerSeat: number, deckIndex: number): Promise<ApiResponse<void>> {
    return this.request<void>('/api/game/driver/tribute-return', {
      method: 'POST',
      body: JSON.stringify({
        room_id: roomId,
        player_seat: playerSeat,
        deck_index: deckIndex
      }),
    });
  }

  async updateProfile(data: { nickname?: string }): Promise<ApiResponse<User>> {
    return this.request<User>('/api/user/profile', {
      method: 'PATCH',
      body: JSON.stringify(data),
    });
  }

  async getProfile(): Promise<User> {
    const response = await this.request<{ user: User }>('/api/user/profile');
    return response.data!.user;
  }
}

export { ApiClient, ApiError };
