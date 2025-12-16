import { randomUUID } from 'crypto';
import { fromProtoCards, fromCardList, compTypeToString } from '@guandan/sdk-ts';
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
import { createLogger, Logger } from './utils/Logger.js';

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
  private readonly logger: Logger;
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
    this.logger = createLogger();
    this.api = new ApiClient(config.serverUrl, this.logger);
    const wsUrl = config.serverUrl.replace(/^http/, 'ws') + '/ws';
    this.ws = new WebSocketClient({ url: wsUrl }, this.logger);
    this.strategy = new AIStrategy(config.level || 1);
  }

  private static readonly NICKNAMES = [
    '二牛', '狗子', '妞妞', '铁蛋', '猪猪', '胖胖',
    '阿童木', '熊大', '熊二', '我是小小鸟', '我是大大胖',
    '没头脑', '不开心'
  ];

  private getRandomNickname(): string {
    const index = Math.floor(Math.random() * AIPlayer.NICKNAMES.length);
    return AIPlayer.NICKNAMES[index];
  }

  async start(): Promise<void> {
    try {
      if (this.config.autoRegister) {
        this.logger.info('start', 'IDLE -> REGISTERING');
        this.state = AIState.REGISTERING;
        this.username = `ai_${randomUUID().slice(0, 8)}`;
        this.logger.setContext({ username: this.username });
        const password = this.config.password || 'Ai_Password_123';
        const auth = await this.api.register({ username: this.username, password });
        this.userId = auth.user.id;
        this.accessToken = auth.token.access_token;
        this.refreshToken = auth.token.refresh_token;
        this.api.setToken(this.accessToken);
        this.api.setRefreshToken(this.refreshToken);
        this.logger.info('start', `registered as ${this.username}`);

        if (!this.config.nickname) {
          this.config.nickname = this.getRandomNickname();
        }
      } else {
        this.logger.info('start', 'IDLE -> LOGGING_IN');
        this.state = AIState.LOGGING_IN;
        this.username = this.config.username!;
        this.logger.setContext({ username: this.username });
        const password = this.config.password || 'Ai_Password_123';
        const auth = await this.api.login({ username: this.username, password });
        this.userId = auth.user.id;
        this.accessToken = auth.token.access_token;
        this.refreshToken = auth.token.refresh_token;
        this.api.setToken(this.accessToken);
        this.api.setRefreshToken(this.refreshToken);
        this.logger.info('start', `logged in as ${this.username}`);
      }

      this.api.setOnTokenRefreshed((auth) => {
        this.logger.warn('onTokenRefreshed', 'token refreshed');
        this.accessToken = auth.token.access_token;
        this.refreshToken = auth.token.refresh_token;
        this.ws.reconnect(this.accessToken);
      });

      if (this.config.nickname) {
        await this.api.updateProfile({ nickname: this.config.nickname });
        this.logger.debug('start', `nickname set to ${this.config.nickname}`);
      }

      this.setupMessageHandlers();
      this.ws.connect(this.accessToken!);

      await this.waitForConnection();

      this.state = AIState.JOINING_ROOM;
      const room = await this.api.joinRoomByCode(this.config.roomCode);
      this.roomId = room.id;
      this.isOwner = room.owner === this.userId;
      this.state = AIState.WAITING_IN_ROOM;
      this.logger.setContext({ roomCode: this.config.roomCode });
      this.logger.info('start', `joined room=${this.config.roomCode} roomId=${this.roomId} isOwner=${this.isOwner}`);

      this.checkAndStartGame(room);
    } catch (error) {
      this.logger.error('start', `failed: ${error}`);
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
        this.logger.debug('onPlayerView', 'no view in message');
        return;
      }

      const seq = Number(view.seq);
      this.logger.debug('onPlayerView', `seq=${seq} dealStatus=${view.dealStatus} playerSeat=${view.playerSeat}`);

      if (seq <= this.lastSeq) {
        this.logger.debug('onPlayerView', `skipped (seq ${seq} <= lastSeq ${this.lastSeq})`);
        return;
      }
      this.lastSeq = seq;

      this.playerSeat = view.playerSeat;
      this.playerView = view;
      if (this.state !== AIState.PLAYING) {
        this.logger.info('onPlayerView', `${this.state} -> PLAYING`);
        this.state = AIState.PLAYING;
      }

      if (view.matchResult) {
        this.handleMatchEnd();
      }
    });

    this.ws.on('tribute_view', (msg: WSMessage) => {
      const data = msg.data as { tribute_view?: TributeView };
      const view = data.tribute_view;
      if (!view) return;

      this.logger.debug('onTributeView', `status=${view.status} receivers=${JSON.stringify(view.receivers)} givers=${JSON.stringify(view.givers)}`);
    });

    // ========== 行动指令消息 ==========

    this.ws.on('game_action', (msg: WSMessage) => {
      const wrapper = msg.data as { game_action?: GameActionData };
      const data = wrapper.game_action;
      if (!data || data.actionType === undefined) {
        this.logger.debug('onGameAction', 'no actionType in message');
        return;
      }

      this.logger.debug('onGameAction', `actionType=${data.actionType} playerSeat=${data.playerSeat} optionsCount=${data.options?.length || 0}`);

      // 如果还没有设置 playerSeat，从 game_action 中初始化
      if (this.playerSeat === null) {
        this.playerSeat = data.playerSeat;
        this.logger.debug('onGameAction', `initialized playerSeat=${data.playerSeat}`);
      }

      if (data.playerSeat !== this.playerSeat) {
        this.logger.debug('onGameAction', `ignored (playerSeat=${data.playerSeat} !== my seat=${this.playerSeat})`);
        return;
      }

      this.handleGameAction(data);
    });
  }

  private checkAndStartGame(room: Room): void {
    const playerCount = room.players.filter(p => p !== null).length;

    if (playerCount === 4 && this.isOwner && !this.startGameTimer) {
      this.logger.info('checkAndStartGame', 'room full, will start in 5s');
      this.startGameTimer = setTimeout(async () => {
        try {
          this.logger.info('checkAndStartGame', 'starting game...');
          await this.api.startGame(this.roomId!);
          this.logger.info('checkAndStartGame', 'game started');
        } catch (error) {
          this.logger.error('checkAndStartGame', `failed: ${error}`);
        }
      }, 5000);
    }
  }

  // ========== 行动处理方法 ==========

  private async handleGameAction(data: GameActionData): Promise<void> {
    if (this.actionInProgress) {
      this.logger.debug('handleGameAction', 'skipped (actionInProgress=true)');
      return;
    }

    this.logger.debug('handleGameAction', `executing actionType=${data.actionType}`);

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
      this.logger.debug('handlePlayDecision', 'no playerView');
      return;
    }

    this.actionInProgress = true;
    try {
      await this.delay(1000 + Math.random() * 2000);

      if (!this.playerView) {
        this.logger.debug('handlePlayDecision', 'no playerView after delay');
        return;
      }

      const dealLevel = this.playerView.dealLevel;
      const handCards = data.hand || this.playerView.playerCards;
      const hand = fromProtoCards(handCards, dealLevel);

      this.logger.debug('handlePlayDecision', `handSize=${hand.length}`);

      const isLeader = this.playerView.leader === this.playerSeat || this.playerView.plays.length === 0;

      let prevComp: CardComp | undefined;
      if (!isLeader && this.playerView.plays.length > 0) {
        for (let i = this.playerView.plays.length - 1; i >= 0; i--) {
          const lastPlay = this.playerView.plays[i];
          if (lastPlay.cards && lastPlay.cards.length > 0 && !lastPlay.isPass) {
            const lastCards = fromProtoCards(lastPlay.cards, dealLevel);
            prevComp = fromCardList(lastCards);
            const prevCardsStr = lastCards.map(c => c.toShortString()).join(', ');
            const prevTypeStr = compTypeToString(prevComp.getType());
            this.logger.info('handlePlayDecision', `prevComp cards=[${prevCardsStr}] type=${prevTypeStr} valid=${prevComp.isValid()}`);
            break;
          }
        }
      }

      const cardsToPlay = this.strategy.selectCardsToPlay(hand, isLeader, prevComp);

      if (cardsToPlay && cardsToPlay.length > 0) {
        const deckIndexes = cardsToPlay.map(c => c.deckIndex);
        const cardsStr = cardsToPlay.map(c => c.toShortString()).join(', ');
        this.logger.info('handlePlayDecision', `play cards=[${cardsStr}] deckIdx=[${deckIndexes.join(', ')}]`);
        await this.api.playCards(this.roomId!, this.playerSeat!, deckIndexes);
      } else {
        this.logger.info('handlePlayDecision', 'pass');
        await this.api.pass(this.roomId!, this.playerSeat!);
      }
    } catch (error) {
      this.logger.error('handlePlayDecision', `failed: ${error}`);
    } finally {
      this.actionInProgress = false;
    }
  }

  private async handleTributeSelection(data: GameActionData): Promise<void> {
    if (!data.options || data.options.length === 0) {
      this.logger.debug('handleTributeSelection', 'no options provided');
      return;
    }

    this.actionInProgress = true;
    try {
      const dealLevel = this.playerView?.dealLevel || 2;
      const poolCards = fromProtoCards(data.options, dealLevel);

      this.logger.debug('handleTributeSelection', `optionsCount=${poolCards.length}`);

      const tributeCard = this.strategy.selectTributeCard(poolCards);

      if (tributeCard) {
        await this.delay(1000 + Math.random() * 2000);
        this.logger.info('handleTributeSelection', `selectTribute card=${tributeCard.toShortString()} deckIdx=${tributeCard.deckIndex}`);
        await this.api.selectTribute(this.roomId!, this.playerSeat!, tributeCard.deckIndex);
      } else {
        this.logger.warn('handleTributeSelection', 'no card selected');
      }
    } catch (error) {
      this.logger.error('handleTributeSelection', `failed: ${error}`);
    } finally {
      this.actionInProgress = false;
    }
  }

  private async handleReturnTribute(data: GameActionData): Promise<void> {
    if (!data.options || data.options.length === 0) {
      this.logger.debug('handleReturnTribute', 'no options provided');
      return;
    }

    this.actionInProgress = true;
    try {
      const dealLevel = this.playerView?.dealLevel || 2;
      const hand = fromProtoCards(data.options, dealLevel);

      this.logger.debug('handleReturnTribute', `optionsCount=${hand.length}`);

      const returnCard = this.strategy.selectReturnTributeCard(hand);

      if (returnCard) {
        await this.delay(1000 + Math.random() * 2000);
        this.logger.info('handleReturnTribute', `returnTribute card=${returnCard.toShortString()} deckIdx=${returnCard.deckIndex}`);
        await this.api.returnTribute(this.roomId!, this.playerSeat!, returnCard.deckIndex);
      } else {
        this.logger.warn('handleReturnTribute', 'no card selected');
      }
    } catch (error) {
      this.logger.error('handleReturnTribute', `failed: ${error}`);
    } finally {
      this.actionInProgress = false;
    }
  }

  private async handleMatchEnd(): Promise<void> {
    this.logger.info('handleMatchEnd', `PLAYING -> LEAVING_ROOM (match ended)`);
    this.state = AIState.LEAVING_ROOM;
    try {
      await this.api.leaveRoom(this.roomId!);
      this.logger.info('handleMatchEnd', 'left room');
    } catch (error) {
      this.logger.error('handleMatchEnd', `failed to leave room: ${error}`);
    }
    this.reset();
    this.logger.info('handleMatchEnd', 'LEAVING_ROOM -> IDLE');
    this.state = AIState.IDLE;
  }

  private reset(): void {
    this.roomId = null;
    this.playerSeat = null;
    this.playerView = null;
    this.isOwner = false;
    this.lastSeq = -1;
    this.logger.setContext({ roomCode: null });
    if (this.startGameTimer) {
      clearTimeout(this.startGameTimer);
      this.startGameTimer = null;
    }
  }

  async stop(): Promise<void> {
    this.logger.info('stop', 'stopping...');
    this.state = AIState.STOPPED;
    if (this.roomId) {
      try {
        await this.api.leaveRoom(this.roomId);
      } catch {
      }
    }
    this.ws.disconnect();
    this.reset();
    this.logger.info('stop', 'stopped');
    this.logger.close();
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
