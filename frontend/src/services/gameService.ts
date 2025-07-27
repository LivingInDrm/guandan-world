import { apiClient } from './api';
import { wsClient } from './websocket';
import { useAuthStore } from '../store/authStore';
import { useRoomStore } from '../store/roomStore';
import { useGameStore } from '../store/gameStore';
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
      apiClient.setToken(token.token);
      wsClient.connect(token.token);
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
    wsClient.on('player_view', this.handlePlayerView.bind(this));

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
      const response = await apiClient.login({ username, password });
      if (response.success && response.data) {
        const { user, token } = response.data;
        
        // Update auth store
        useAuthStore.getState().login(user, token);
        
        // Set API token and connect WebSocket
        apiClient.setToken(token.token);
        wsClient.connect(token.token);
        
        return true;
      }
      return false;
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
      const response = await apiClient.register({ username, password });
      if (response.success && response.data) {
        const { user, token } = response.data;
        
        // Update auth store
        useAuthStore.getState().login(user, token);
        
        // Set API token and connect WebSocket
        apiClient.setToken(token.token);
        wsClient.connect(token.token);
        
        return true;
      }
      return false;
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
        
        // Send WebSocket join message
        wsClient.send('join_room', { room_id: roomId });
        
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
      
      // Send WebSocket leave message
      wsClient.send('leave_room', { room_id: roomId });
      
      useRoomStore.getState().setCurrentRoom(null);
    } catch (error) {
      console.error('Failed to leave room:', error);
    }
  }

  async startGame(roomId: string): Promise<boolean> {
    try {
      const response = await apiClient.startGame(roomId);
      if (response.success) {
        // Send WebSocket start game message
        wsClient.send('start_game', { room_id: roomId });
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

  private handleJoinRoom(message: WSMessage): void {
    console.log('Player joined room:', message.data);
    // Room update will be handled by room_update message
  }

  private handleLeaveRoom(message: WSMessage): void {
    console.log('Player left room:', message.data);
    // Room update will be handled by room_update message
  }

  private handleGamePrepare(message: WSMessage): void {
    console.log('Game prepare:', message.data);
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
    console.log('Game begin:', message.data);
    useGameStore.getState().setCountdown(null);
    useGameStore.getState().setLastMessage(message);
  }

  private handleGameEvent(message: WSMessage): void {
    console.log('Game event:', message.data);
    useGameStore.getState().setGameState(message.data);
    useGameStore.getState().setLastMessage(message);
  }

  private handlePlayerView(message: WSMessage): void {
    const playerData = message.data;
    if (playerData) {
      useGameStore.getState().setPlayerSeat(playerData.seat);
      useGameStore.getState().setMyTurn(playerData.is_my_turn || false);
      useGameStore.getState().setGameState(playerData.game_state);
    }
  }

  private handlePlayerTimeout(message: WSMessage): void {
    console.log('Player timeout:', message.data);
    useGameStore.getState().setLastMessage(message);
  }

  private handleAutoPlay(message: WSMessage): void {
    console.log('Auto play:', message.data);
    useGameStore.getState().setLastMessage(message);
  }

  private handleReconnect(message: WSMessage): void {
    console.log('Reconnect:', message.data);
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