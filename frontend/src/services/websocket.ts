import type { WSMessage, WSMessageType } from '../types';

export type WSEventHandler = (message: WSMessage) => void;
export type WSConnectionHandler = (connected: boolean) => void;
export type WSErrorHandler = (error: Event) => void;

interface WSClientOptions {
  url?: string;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  heartbeatInterval?: number;
}

class WebSocketClient {
  private ws: WebSocket | null = null;
  private url: string;
  private reconnectInterval: number;
  private maxReconnectAttempts: number;
  private heartbeatInterval: number;
  private reconnectAttempts: number = 0;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private heartbeatTimer: NodeJS.Timeout | null = null;
  private isConnected: boolean = false;
  private isReconnecting: boolean = false;
  private shouldReconnect: boolean = true;

  // Event handlers
  private messageHandlers: Map<WSMessageType, WSEventHandler[]> = new Map();
  private connectionHandlers: WSConnectionHandler[] = [];
  private errorHandlers: WSErrorHandler[] = [];

  constructor(options: WSClientOptions = {}) {
    this.url = options.url || this.getWebSocketURL();
    this.reconnectInterval = options.reconnectInterval || 3000;
    this.maxReconnectAttempts = options.maxReconnectAttempts || 10;
    this.heartbeatInterval = options.heartbeatInterval || 30000;
  }

  private getWebSocketURL(): string {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = import.meta.env.VITE_WS_HOST || window.location.host.replace(':5173', ':8080');
    return `${protocol}//${host}/ws`;
  }

  connect(token?: string): void {
    if (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
      return;
    }

    this.shouldReconnect = true;
    
    try {
      const wsUrl = token ? `${this.url}?token=${token}` : this.url;
      this.ws = new WebSocket(wsUrl);
      
      this.ws.onopen = this.handleOpen.bind(this);
      this.ws.onmessage = this.handleMessage.bind(this);
      this.ws.onclose = this.handleClose.bind(this);
      this.ws.onerror = this.handleError.bind(this);
    } catch (error) {
      console.error('WebSocket connection failed:', error);
      this.scheduleReconnect();
    }
  }

  disconnect(): void {
    this.shouldReconnect = false;
    this.clearTimers();
    
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    
    this.isConnected = false;
    this.reconnectAttempts = 0;
  }

  send(type: WSMessageType, data: any): boolean {
    if (!this.isConnected || !this.ws) {
      console.warn('WebSocket not connected, message not sent:', { type, data });
      return false;
    }

    const message: WSMessage = {
      type,
      data,
      timestamp: new Date().toISOString(),
    };

    try {
      this.ws.send(JSON.stringify(message));
      return true;
    } catch (error) {
      console.error('Failed to send WebSocket message:', error);
      return false;
    }
  }

  // Event handler registration
  on(messageType: WSMessageType, handler: WSEventHandler): void {
    if (!this.messageHandlers.has(messageType)) {
      this.messageHandlers.set(messageType, []);
    }
    this.messageHandlers.get(messageType)!.push(handler);
  }

  off(messageType: WSMessageType, handler: WSEventHandler): void {
    const handlers = this.messageHandlers.get(messageType);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  onConnection(handler: WSConnectionHandler): void {
    this.connectionHandlers.push(handler);
  }

  onError(handler: WSErrorHandler): void {
    this.errorHandlers.push(handler);
  }

  // WebSocket event handlers
  private handleOpen(): void {
    console.log('WebSocket connected');
    this.isConnected = true;
    this.isReconnecting = false;
    this.reconnectAttempts = 0;
    this.clearTimers();
    this.startHeartbeat();
    
    this.connectionHandlers.forEach(handler => handler(true));
  }

  private handleMessage(event: MessageEvent): void {
    try {
      const message: WSMessage = JSON.parse(event.data);
      
      // Handle heartbeat response
      if (message.type === 'pong') {
        return;
      }
      
      // Dispatch to registered handlers
      const handlers = this.messageHandlers.get(message.type as WSMessageType);
      if (handlers) {
        handlers.forEach(handler => handler(message));
      }
      
      // Also dispatch to generic message handlers
      const genericHandlers = this.messageHandlers.get('*' as WSMessageType);
      if (genericHandlers) {
        genericHandlers.forEach(handler => handler(message));
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error);
    }
  }

  private handleClose(event: CloseEvent): void {
    console.log('WebSocket disconnected:', event.code, event.reason);
    this.isConnected = false;
    this.clearTimers();
    
    this.connectionHandlers.forEach(handler => handler(false));
    
    if (this.shouldReconnect && !this.isReconnecting) {
      this.scheduleReconnect();
    }
  }

  private handleError(event: Event): void {
    console.error('WebSocket error:', event);
    this.errorHandlers.forEach(handler => handler(event));
  }

  // Reconnection logic
  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    if (this.isReconnecting) {
      return;
    }

    this.isReconnecting = true;
    this.reconnectAttempts++;
    
    console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
    
    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, this.reconnectInterval);
  }

  // Heartbeat mechanism
  private startHeartbeat(): void {
    this.heartbeatTimer = setInterval(() => {
      if (this.isConnected) {
        this.send('ping' as WSMessageType, {});
      }
    }, this.heartbeatInterval);
  }

  private clearTimers(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  // Getters
  get connected(): boolean {
    return this.isConnected;
  }

  get reconnecting(): boolean {
    return this.isReconnecting;
  }

  get readyState(): number {
    return this.ws?.readyState ?? WebSocket.CLOSED;
  }
}

// Create singleton instance
export const wsClient = new WebSocketClient();

// Export the class for testing
export { WebSocketClient };