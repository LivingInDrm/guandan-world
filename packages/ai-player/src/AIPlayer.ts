import { randomUUID } from 'crypto';
import { fromProtoCards, fromCardList } from '@guandan/sdk-ts';
import type { CardComp } from '@guandan/sdk-ts';
import { ApiClient } from './services/ApiClient.js';
import { WebSocketClient } from './services/WebSocketClient.js';
import {
  Room,
  PlayerView,
  TributeView,
  WSMessage,
  GameActionData,
  GameActionType
} from './services/types.js';
import { AIStrategy } from './AIStrategy.js';
import type { AIPlayerConfig } from './config.js';

export enum AIState {
  IDLE = 'IDLE',
  REGISTERING = 'REGISTERING',
  LOGGING_IN = 'LOGGING_IN',
  JOINING_ROOM = 'JOINING_ROOM',
  WAITING_IN_ROOM = 'WAITING_IN_ROOM',
  PLAYING = 'PLAYING',
  LEAVING_ROOM = 'LEAVING_ROOM',
  STOPPED = 'STOPPED',
}

export interface AIPlayerInfo {
  id: string;
  config: AIPlayerConfig;
  state: AIState;
  username: string | null;
  roomId: string | null;
  playerSeat: number | null;
  createdAt: Date;
}

export class AIPlayer {
  private readonly id: string;
  private readonly createdAt: Date;
  private api: ApiClient;
  private ws: WebSocketClient;
  private strategy: AIStrategy;
  private config: AIPlayerConfig;

  private state: AIState = AIState.IDLE;
  private roomId: string | null = null;
  private playerSeat: number | null = null;
  private playerView: PlayerView | null = null;
  private isOwner: boolean = false;
  private userId: string | null = null;
  private username: string | null = null;
  private accessToken: string | null = null;
  private refreshToken: string | null = null;

  private startGameTimer: ReturnType<typeof setTimeout> | null = null;
  private actionInProgress: boolean = false;
  private lastSeq: number = -1;

  constructor(config: AIPlayerConfig) {
    this.id = randomUUID();
    this.createdAt = new Date();
    this.config = config;
    this.api = new ApiClient(config.serverUrl);
    const wsUrl = config.serverUrl.replace(/^http/, 'ws') + '/ws';
    this.ws = new WebSocketClient({ url: wsUrl });
    this.strategy = new AIStrategy(config.level || 1);
  }

  async start(): Promise<void> {
    try {
      if (this.config.autoRegister) {
        this.state = AIState.REGISTERING;
        this.username = `ai_${randomUUID().slice(0, 8)}`;
        const password = this.config.password || 'Ai_Password_123';
        console.log(`Registering as ${this.username}...`);
        const auth = await this.api.register({ username: this.username, password });
        this.userId = auth.user.id;
        this.accessToken = auth.token.access_token;
        this.refreshToken = auth.token.refresh_token;
        this.api.setToken(this.accessToken);
        this.api.setRefreshToken(this.refreshToken);
        console.log(`Registered successfully as ${this.username}`);
      } else {
        this.state = AIState.LOGGING_IN;
        this.username = this.config.username!;
        const password = this.config.password || 'Ai_Password_123';
        console.log(`Logging in as ${this.username}...`);
        const auth = await this.api.login({ username: this.username, password });
        this.userId = auth.user.id;
        this.accessToken = auth.token.access_token;
        this.refreshToken = auth.token.refresh_token;
        this.api.setToken(this.accessToken);
        this.api.setRefreshToken(this.refreshToken);
        console.log(`Logged in successfully as ${this.username}`);
      }

      this.api.setOnTokenRefreshed((auth) => {
        console.log('Token refreshed successfully');
        this.accessToken = auth.token.access_token;
        this.refreshToken = auth.token.refresh_token;
        this.ws.reconnect(this.accessToken);
      });

      if (this.config.nickname) {
        await this.api.updateProfile({ nickname: this.config.nickname });
        console.log(`Nickname set to ${this.config.nickname}`);
      }

      this.setupMessageHandlers();
      this.ws.connect(this.accessToken!);

      await this.waitForConnection();

      this.state = AIState.JOINING_ROOM;
      console.log(`Joining room with code ${this.config.roomCode}...`);
      const room = await this.api.joinRoomByCode(this.config.roomCode);
      this.roomId = room.id;
      this.isOwner = room.owner === this.userId;
      this.state = AIState.WAITING_IN_ROOM;
      console.log(`Joined room ${this.roomId}, isOwner: ${this.isOwner}`);

      this.checkAndStartGame(room);
    } catch (error) {
      console.error('Failed to start AI player:', error);
      throw error;
    }
  }

  private waitForConnection(): Promise<void> {
    return new Promise((resolve, reject) => {
      let resolved = false;

      const resolveOnce = () => {
        if (resolved) return;
        resolved = true;
        clearTimeout(timeout);
        resolve();
      };

      const timeout = setTimeout(() => {
        if (resolved) return;
        resolved = true;
        reject(new Error('WebSocket connection timeout'));
      }, 10000);

      const checkConnection = () => {
        if (this.ws.connected) {
          resolveOnce();
        } else if (!resolved) {
          setTimeout(checkConnection, 100);
        }
      };

      this.ws.onConnection((connected) => {
        if (connected) {
          resolveOnce();
        }
      });

      checkConnection();
    });
  }

  private setupMessageHandlers(): void {
    // ========== 状态同步消息 ==========

    this.ws.on('room_update', (msg: WSMessage) => {
      const data = msg.data as { room?: Room } | Room;
      const room = ('room' in data && data.room) ? data.room : data as Room;
      if (room && room.id) {
        this.checkAndStartGame(room);
      }
    });

    this.ws.on('player_view', (msg: WSMessage) => {
      const data = msg.data as { player_view?: PlayerView };
      const view = data.player_view;
      if (!view) {
        console.log('[DEBUG] player_view: no view in message');
        return;
      }

      const seq = Number(view.seq);
      console.log(`[DEBUG] player_view: seq=${seq}, dealStatus=${view.dealStatus}, playerSeat=${view.playerSeat}`);

      if (seq <= this.lastSeq) {
        console.log(`[DEBUG] player_view: skipped (seq ${seq} <= lastSeq ${this.lastSeq})`);
        return;
      }
      this.lastSeq = seq;

      this.playerSeat = view.playerSeat;
      this.playerView = view;
      this.state = AIState.PLAYING;

      if (view.matchResult) {
        this.handleMatchEnd();
      }
    });

    this.ws.on('tribute_view', (msg: WSMessage) => {
      const data = msg.data as { tribute_view?: TributeView };
      const view = data.tribute_view;
      if (!view) return;

      console.log(`[DEBUG] tribute_view: status=${view.status}, receivers=${JSON.stringify(view.receivers)}, givers=${JSON.stringify(view.givers)}`);
    });

    // ========== 行动指令消息 ==========

    this.ws.on('game_action', (msg: WSMessage) => {
      const wrapper = msg.data as { game_action?: GameActionData };
      const data = wrapper.game_action;
      if (!data || data.actionType === undefined) {
        console.log('[DEBUG] game_action: no actionType in message');
        return;
      }

      console.log(`[DEBUG] game_action: actionType=${data.actionType}, playerSeat=${data.playerSeat}, optionsCount=${data.options?.length || 0}`);

      if (data.playerSeat !== this.playerSeat) {
        console.log(`[DEBUG] game_action: ignored (playerSeat=${data.playerSeat} !== my seat=${this.playerSeat})`);
        return;
      }

      this.handleGameAction(data);
    });
  }

  private checkAndStartGame(room: Room): void {
    const playerCount = room.players.filter(p => p !== null).length;

    if (playerCount === 4 && this.isOwner && !this.startGameTimer) {
      console.log('Room is full, starting game in 5 seconds...');
      this.startGameTimer = setTimeout(async () => {
        try {
          console.log('Starting game...');
          await this.api.startGame(this.roomId!);
          console.log('Game started');
        } catch (error) {
          console.error('Failed to start game:', error);
        }
      }, 5000);
    }
  }

  // ========== 行动处理方法 ==========

  private async handleGameAction(data: GameActionData): Promise<void> {
    if (this.actionInProgress) {
      console.log(`[DEBUG] game_action: skipped (actionInProgress=true)`);
      return;
    }

    console.log(`[DEBUG] game_action: executing actionType=${data.actionType}`);

    switch (data.actionType) {
      case GameActionType.GAME_ACTION_TYPE_PLAY_DECISION:
        await this.handlePlayDecision(data);
        break;
      case GameActionType.GAME_ACTION_TYPE_TRIBUTE_SELECTION:
        await this.handleTributeSelection(data);
        break;
      case GameActionType.GAME_ACTION_TYPE_RETURN_TRIBUTE:
        await this.handleReturnTribute(data);
        break;
    }
  }

  private async handlePlayDecision(data: GameActionData): Promise<void> {
    if (!this.playerView) {
      console.log('[DEBUG] handlePlayDecision: no playerView');
      return;
    }

    this.actionInProgress = true;
    try {
      const dealLevel = this.playerView.dealLevel;
      const handCards = data.hand || this.playerView.playerCards;
      const hand = fromProtoCards(handCards, dealLevel);

      console.log(`[DEBUG] handlePlayDecision: handSize=${hand.length}`);

      const isLeader = this.playerView.leader === this.playerSeat || this.playerView.plays.length === 0;

      let prevComp: CardComp | undefined;
      if (!isLeader && this.playerView.plays.length > 0) {
        for (let i = this.playerView.plays.length - 1; i >= 0; i--) {
          const lastPlay = this.playerView.plays[i];
          if (lastPlay.cards && lastPlay.cards.length > 0 && !lastPlay.isPass) {
            const lastCards = fromProtoCards(lastPlay.cards, dealLevel);
            prevComp = fromCardList(lastCards);
            break;
          }
        }
      }

      const cardsToPlay = this.strategy.selectCardsToPlay(hand, isLeader, prevComp);

      await this.delay(500 + Math.random() * 1000);

      if (cardsToPlay && cardsToPlay.length > 0) {
        const deckIndexes = cardsToPlay.map(c => c.deckIndex);
        console.log(`[DEBUG] Playing cards: ${cardsToPlay.map(c => c.toShortString()).join(', ')}`);
        await this.api.playCards(this.roomId!, this.playerSeat!, deckIndexes);
        console.log('[DEBUG] playCards completed');
      } else {
        console.log('[DEBUG] Passing...');
        await this.api.pass(this.roomId!, this.playerSeat!);
        console.log('[DEBUG] pass completed');
      }
    } catch (error) {
      console.error('[DEBUG] Failed to play:', error);
    } finally {
      this.actionInProgress = false;
    }
  }

  private async handleTributeSelection(data: GameActionData): Promise<void> {
    if (!data.options || data.options.length === 0) {
      console.log('[DEBUG] handleTributeSelection: no options provided');
      return;
    }

    this.actionInProgress = true;
    try {
      const dealLevel = this.playerView?.dealLevel || 2;
      const poolCards = fromProtoCards(data.options, dealLevel);

      console.log(`[DEBUG] handleTributeSelection: optionsCount=${poolCards.length}`);

      const tributeCard = this.strategy.selectTributeCard(poolCards);

      if (tributeCard) {
        console.log(`[DEBUG] Selecting tribute card: ${tributeCard.toShortString()}, deckIndex=${tributeCard.deckIndex}`);
        await this.delay(500 + Math.random() * 1000);
        await this.api.selectTribute(this.roomId!, this.playerSeat!, tributeCard.deckIndex);
        console.log(`[DEBUG] selectTribute completed`);
      } else {
        console.log(`[DEBUG] handleTributeSelection: no card selected`);
      }
    } catch (error) {
      console.error('[DEBUG] Failed to select tribute:', error);
    } finally {
      this.actionInProgress = false;
    }
  }

  private async handleReturnTribute(data: GameActionData): Promise<void> {
    if (!data.options || data.options.length === 0) {
      console.log('[DEBUG] handleReturnTribute: no options provided');
      return;
    }

    this.actionInProgress = true;
    try {
      const dealLevel = this.playerView?.dealLevel || 2;
      const hand = fromProtoCards(data.options, dealLevel);

      console.log(`[DEBUG] handleReturnTribute: optionsCount=${hand.length}`);

      const returnCard = this.strategy.selectReturnTributeCard(hand);

      if (returnCard) {
        console.log(`[DEBUG] Returning tribute card: ${returnCard.toShortString()}, deckIndex=${returnCard.deckIndex}`);
        await this.delay(500 + Math.random() * 1000);
        await this.api.returnTribute(this.roomId!, this.playerSeat!, returnCard.deckIndex);
        console.log(`[DEBUG] returnTribute completed`);
      } else {
        console.log(`[DEBUG] handleReturnTribute: no card selected`);
      }
    } catch (error) {
      console.error('[DEBUG] Failed to return tribute:', error);
    } finally {
      this.actionInProgress = false;
    }
  }

  private async handleMatchEnd(): Promise<void> {
    console.log('Match ended, leaving room...');
    this.state = AIState.LEAVING_ROOM;
    try {
      await this.api.leaveRoom(this.roomId!);
      console.log('Left room successfully');
    } catch (error) {
      console.error('Failed to leave room:', error);
    }
    this.reset();
    this.state = AIState.IDLE;
  }

  private reset(): void {
    this.roomId = null;
    this.playerSeat = null;
    this.playerView = null;
    this.isOwner = false;
    this.lastSeq = -1;
    if (this.startGameTimer) {
      clearTimeout(this.startGameTimer);
      this.startGameTimer = null;
    }
  }

  async stop(): Promise<void> {
    console.log('Stopping AI player...');
    this.state = AIState.STOPPED;
    if (this.roomId) {
      try {
        await this.api.leaveRoom(this.roomId);
      } catch {
      }
    }
    this.ws.disconnect();
    this.reset();
    console.log('AI player stopped');
  }

  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  getState(): AIState {
    return this.state;
  }

  getId(): string {
    return this.id;
  }

  getInfo(): AIPlayerInfo {
    const { password: _, ...safeConfig } = this.config;
    return {
      id: this.id,
      config: safeConfig,
      state: this.state,
      username: this.username,
      roomId: this.roomId,
      playerSeat: this.playerSeat,
      createdAt: this.createdAt,
    };
  }
}
