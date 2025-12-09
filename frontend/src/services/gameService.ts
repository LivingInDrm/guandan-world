import { apiClient } from './api';
import { wsClient } from './websocket';
import { useAuthStore } from '../store/authStore';
import { useRoomStore } from '../store/roomStore';
import { useGameStore } from '../store/gameStore';
import { useTributeStore } from '../store/tributeStore';
import type { WSMessage } from '../types';

class GameService {
  private initialized = false;

  async initialize(): Promise<void> {
    if (this.initialized) return;

    // Set up WebSocket event handlers
    this.setupWebSocketHandlers();
    
    // Connect WebSocket if user is authenticated
    const { token } = useAuthStore.getState();
    if (token) {
      apiClient.setToken(token.access_token);
      wsClient.connect(token.access_token);
    }

    this.initialized = true;
  }

  private setupWebSocketHandlers(): void {
    // Connection status
    wsClient.onConnection((connected) => {
      useGameStore.getState().setConnected(connected);
    });

    // Room management messages
    wsClient.on('room_update', this.handleRoomUpdate.bind(this));
    wsClient.on('join_room', this.handleJoinRoom.bind(this));
    wsClient.on('leave_room', this.handleLeaveRoom.bind(this));

    // Game flow messages
    wsClient.on('game_prepare', this.handleGamePrepare.bind(this));
    wsClient.on('countdown', this.handleCountdown.bind(this));
    wsClient.on('game_begin', this.handleGameBegin.bind(this));

    // Game state messages
    wsClient.on('game_event', this.handleGameEvent.bind(this));
    // Note: player_view is handled directly in GamePage.tsx for better context and state management

    // Player management
    wsClient.on('player_timeout', this.handlePlayerTimeout.bind(this));
    wsClient.on('auto_play', this.handleAutoPlay.bind(this));
    wsClient.on('reconnect', this.handleReconnect.bind(this));

    // Error handling
    wsClient.onError((error) => {
      console.error('WebSocket error:', error);
      useGameStore.getState().setError('WebSocket connection error');
    });
  }

  // Authentication methods
  async login(username: string, password: string): Promise<boolean> {
    try {
      const { user, token } = await apiClient.login({ username, password });
      
      useAuthStore.getState().login(user, token);
      
      apiClient.setToken(token.access_token);
      wsClient.connect(token.access_token);
      
      return true;
    } catch (error) {
      console.error('Login failed:', error);
      useAuthStore.getState().setError(
        error instanceof Error ? error.message : 'Login failed'
      );
      return false;
    }
  }

  async register(username: string, password: string): Promise<boolean> {
    try {
      const { user, token } = await apiClient.register({ username, password });
      
      useAuthStore.getState().login(user, token);
      
      apiClient.setToken(token.access_token);
      wsClient.connect(token.access_token);
      
      return true;
    } catch (error) {
      console.error('Registration failed:', error);
      useAuthStore.getState().setError(
        error instanceof Error ? error.message : 'Registration failed'
      );
      return false;
    }
  }

  async logout(): Promise<void> {
    try {
      await apiClient.logout();
    } catch (error) {
      console.error('Logout API call failed:', error);
    } finally {
      // Always clean up local state
      useAuthStore.getState().logout();
      useRoomStore.getState().reset();
      useGameStore.getState().reset();
      useTributeStore.getState().reset();
      wsClient.disconnect();
    }
  }

  // Room management methods
  async loadRoomList(page: number = 1): Promise<void> {
    try {
      useRoomStore.getState().setLoading(true);
      const response = await apiClient.getRoomList(page);
      
      if (response.success && response.data) {
        useRoomStore.getState().setRoomList(response.data);
      }
    } catch (error) {
      console.error('Failed to load room list:', error);
      useRoomStore.getState().setError(
        error instanceof Error ? error.message : 'Failed to load rooms'
      );
    } finally {
      useRoomStore.getState().setLoading(false);
    }
  }

  async createRoom(): Promise<string | null> {
    try {
      const response = await apiClient.createRoom();
      if (response.success && response.data) {
        useRoomStore.getState().setCurrentRoom(response.data);
        return response.data.id;
      }
      return null;
    } catch (error) {
      console.error('Failed to create room:', error);
      useRoomStore.getState().setError(
        error instanceof Error ? error.message : 'Failed to create room'
      );
      return null;
    }
  }

  async joinRoom(roomId: string): Promise<boolean> {
    try {
      const response = await apiClient.joinRoom(roomId);
      if (response.success && response.data) {
        useRoomStore.getState().setCurrentRoom(response.data);
        return true;
      }
      return false;
    } catch (error) {
      console.error('Failed to join room:', error);
      useRoomStore.getState().setError(
        error instanceof Error ? error.message : 'Failed to join room'
      );
      return false;
    }
  }

  async leaveRoom(roomId: string): Promise<void> {
    try {
      await apiClient.leaveRoom(roomId);
      useRoomStore.getState().setCurrentRoom(null);
    } catch (error) {
      console.error('Failed to leave room:', error);
    }
  }

  async startGame(roomId: string): Promise<boolean> {
    try {
      const response = await apiClient.startGame(roomId);
      if (response.success) {
        return true;
      }
      return false;
    } catch (error) {
      console.error('Failed to start game:', error);
      useRoomStore.getState().setError(
        error instanceof Error ? error.message : 'Failed to start game'
      );
      return false;
    }
  }

  // WebSocket message handlers
  private handleRoomUpdate(message: WSMessage): void {
    const roomData = message.data;
    if (roomData) {
      useRoomStore.getState().updateRoomInList(roomData);
      
      // Update current room if it matches
      const currentRoom = useRoomStore.getState().currentRoom;
      if (currentRoom && currentRoom.id === roomData.id) {
        useRoomStore.getState().setCurrentRoom(roomData);
      }
    }
  }

  private handleJoinRoom(_message: WSMessage): void {
    // Room update will be handled by room_update message
  }

  private handleLeaveRoom(_message: WSMessage): void {
    // Room update will be handled by room_update message
  }

  private handleGamePrepare(message: WSMessage): void {
    useGameStore.getState().setInGame(true);
    useGameStore.getState().setLastMessage(message);
  }

  private handleCountdown(message: WSMessage): void {
    const countdown = message.data?.countdown;
    if (typeof countdown === 'number') {
      useGameStore.getState().setCountdown(countdown);
    }
  }

  private handleGameBegin(message: WSMessage): void {
    useGameStore.getState().setCountdown(null);
    useGameStore.getState().setLastMessage(message);
  }

  private handleGameEvent(message: WSMessage): void {
    // GameEvent is now sent directly in message.data (flattened structure)
    // Note: Do NOT set gameState here - it should only be updated from player_view messages
    // game_event is just a notification, not the full game state
    useGameStore.getState().setLastMessage(message);
  }

  // Note: handlePlayerView removed - it's handled directly in GamePage.tsx
  // GamePage has better context for handling player view updates and game phase transitions

  private handlePlayerTimeout(message: WSMessage): void {
    useGameStore.getState().setLastMessage(message);
  }

  private handleAutoPlay(message: WSMessage): void {
    useGameStore.getState().setLastMessage(message);
  }

  private handleReconnect(message: WSMessage): void {
    useGameStore.getState().setLastMessage(message);
  }

  // Utility methods
  get isConnected(): boolean {
    return wsClient.connected;
  }

  get isAuthenticated(): boolean {
    return useAuthStore.getState().isAuthenticated;
  }
}

// Create singleton instance
export const gameService = new GameService();

// Export the class for testing
export { GameService };