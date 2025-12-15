import WebSocket from 'ws';
import type { WSMessage, WSMessageType } from './types.js';

export type WSEventHandler = (message: WSMessage) => void;
export type WSConnectionHandler = (connected: boolean) => void;
export type WSErrorHandler = (error: Error) => void;

interface WSClientOptions {
  url: string;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  heartbeatInterval?: number;
}

class WebSocketClient {
  private ws: WebSocket | null = null;
  private url: string;
  private token: string | undefined;
  private reconnectInterval: number;
  private maxReconnectAttempts: number;
  private heartbeatInterval: number;
  private reconnectAttempts: number = 0;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private isConnected: boolean = false;
  private isReconnecting: boolean = false;
  private shouldReconnect: boolean = true;

  private messageHandlers: Map<WSMessageType | '*', WSEventHandler[]> = new Map();
  private connectionHandlers: WSConnectionHandler[] = [];
  private errorHandlers: WSErrorHandler[] = [];

  constructor(options: WSClientOptions) {
    this.url = options.url;
    this.reconnectInterval = options.reconnectInterval || 3000;
    this.maxReconnectAttempts = options.maxReconnectAttempts || 10;
    this.heartbeatInterval = options.heartbeatInterval || 30000;
  }

  connect(token?: string): void {
    if (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
      return;
    }

    this.token = token;
    this.shouldReconnect = true;

    try {
      const wsUrl = token ? `${this.url}?token=${token}` : this.url;
      this.ws = new WebSocket(wsUrl);

      this.ws.on('open', this.handleOpen.bind(this));
      this.ws.on('message', this.handleMessage.bind(this));
      this.ws.on('close', this.handleClose.bind(this));
      this.ws.on('error', this.handleError.bind(this));
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

  reconnect(token?: string): void {
    console.log('Reconnecting WebSocket with new token...');
    this.disconnect();
    this.connect(token);
  }

  send(type: WSMessageType, data: unknown): boolean {
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

  on(messageType: WSMessageType | '*', handler: WSEventHandler): void {
    if (!this.messageHandlers.has(messageType)) {
      this.messageHandlers.set(messageType, []);
    }
    this.messageHandlers.get(messageType)!.push(handler);
  }

  off(messageType: WSMessageType | '*', handler: WSEventHandler): void {
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

  private handleOpen(): void {
    console.log('WebSocket connected');
    this.isConnected = true;
    this.isReconnecting = false;
    this.reconnectAttempts = 0;
    this.clearTimers();
    this.startHeartbeat();

    this.connectionHandlers.forEach(handler => handler(true));
  }

  private handleMessage(data: WebSocket.RawData): void {
    try {
      const message: WSMessage = JSON.parse(data.toString());

      if (message.type === 'pong') {
        return;
      }

      console.log(`[WS DEBUG] Received message type: ${message.type}`);

      const handlers = this.messageHandlers.get(message.type as WSMessageType);
      if (handlers) {
        handlers.forEach(handler => handler(message));
      } else {
        console.log(`[WS DEBUG] No handler registered for message type: ${message.type}`);
      }

      const genericHandlers = this.messageHandlers.get('*');
      if (genericHandlers) {
        genericHandlers.forEach(handler => handler(message));
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error);
    }
  }

  private handleClose(): void {
    console.log('WebSocket disconnected');
    this.isConnected = false;
    this.clearTimers();

    this.connectionHandlers.forEach(handler => handler(false));

    if (this.shouldReconnect && !this.isReconnecting) {
      this.scheduleReconnect();
    }
  }

  private handleError(error: Error): void {
    console.error('WebSocket error:', error);
    this.errorHandlers.forEach(handler => handler(error));
  }

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
    console.log(`Reconnecting... attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);

    this.reconnectTimer = setTimeout(() => {
      this.connect(this.token);
    }, this.reconnectInterval);
  }

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

export { WebSocketClient };
