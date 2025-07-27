import type { 
  ApiResponse, 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse,
  RoomListResponse,
  CreateRoomRequest,
  Room
} from '../types';

// API configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

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

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  setToken(token: string | null) {
    this.token = token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };

    if (this.token) {
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
        throw error;
      }
      
      throw new ApiError(
        error instanceof Error ? error.message : 'Network error',
        0
      );
    }
  }

  // Authentication APIs
  async login(credentials: LoginRequest): Promise<ApiResponse<AuthResponse>> {
    return this.request<AuthResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async register(userData: RegisterRequest): Promise<ApiResponse<AuthResponse>> {
    return this.request<AuthResponse>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async logout(): Promise<ApiResponse<void>> {
    const response = await this.request<void>('/api/auth/logout', {
      method: 'POST',
    });
    this.token = null;
    return response;
  }

  // Room APIs
  async getRoomList(page: number = 1, limit: number = 12): Promise<ApiResponse<RoomListResponse>> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    
    return this.request<RoomListResponse>(`/api/rooms?${params}`);
  }

  async createRoom(roomData: CreateRoomRequest = {}): Promise<ApiResponse<Room>> {
    return this.request<Room>('/api/rooms', {
      method: 'POST',
      body: JSON.stringify(roomData),
    });
  }

  async joinRoom(roomId: string): Promise<ApiResponse<Room>> {
    return this.request<Room>(`/api/rooms/${roomId}/join`, {
      method: 'POST',
    });
  }

  async leaveRoom(roomId: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/rooms/${roomId}/leave`, {
      method: 'POST',
    });
  }

  async getRoomDetails(roomId: string): Promise<ApiResponse<Room>> {
    return this.request<Room>(`/api/rooms/${roomId}`);
  }

  async startGame(roomId: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/api/rooms/${roomId}/start`, {
      method: 'POST',
    });
  }

  // Health check
  async healthCheck(): Promise<ApiResponse<{ status: string }>> {
    return this.request<{ status: string }>('/healthz');
  }
}

// Create singleton instance
export const apiClient = new ApiClient();

// Export the class for testing
export { ApiClient, ApiError };