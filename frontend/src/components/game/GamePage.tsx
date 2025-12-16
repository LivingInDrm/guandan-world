import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { useGameStore } from '../../store/gameStore';
import { useTributeStore } from '../../store/tributeStore';
import { wsClient } from '../../services/websocket';
import { apiClient } from '../../services/api';
import type { WSMessage, GameActionData, TurnDeadlineData } from '../../types';
import { WS_MESSAGE_TYPES, DealStatus } from '../../types';
import type { GameEvent, Card, PlayerView as ProtoPlayerView } from '../../types/proto';
import { eventTypeToJSON, EventType } from '../../types/generated/event';
import { audioService } from '../../services/audioService';
import GameBoard from './GameBoard';
import WaitingBoard from './WaitingBoard';
import PlayerHand from './PlayerHand';
import GameControls from './GameControls';
import Header from '../layout/Header';
import TributeBoard from './tribute/TributeBoard';
import TributeControls from './tribute/TributeControls';
import { TributeStatus } from '../../types/generated/view';
import DealResult from './DealResult';
import MatchResult from './MatchResult';
import { useDealResultData, useMatchResultData } from '../../hooks/useResultData';
import { useTributeData } from '../../hooks/useTributeData';
import { usePlayerViewData, usePlayerHandData } from '../../hooks/useGameData';
import OrientationGuard from '../ui/OrientationGuard';

// 游戏页面状态常量
const GamePageState = {
  WAITING_PLAYERS: 'waiting_players',
  GAME_PREPARE: 'game_prepare',
  TRIBUTE_PHASE: 'tribute_phase',
  PLAYING: 'playing',
  DEAL_RESULT: 'deal_result',
  MATCH_RESULT: 'match_result'
} as const;

const GamePage: React.FC = () => {
  const navigate = useNavigate();
  const { roomId } = useParams<{ roomId: string }>();
  const { user } = useAuthStore();
  const { currentRoom, setCurrentRoom, setError: setRoomError, setLoading } = useRoomStore();
  const {
    isConnected,
    countdown,
    playerSeat,
    currentPhase,
    isMyTurn,
    myFinishRank,
    setCountdown,
    setCurrentPhase,
    setPlayerSeat,
    setPlayerView,
    setMyTurn,
    setMyFinishRank
  } = useGameStore();

  const dealResultData = useDealResultData();
  const matchResultData = useMatchResultData();
  const tributeData = useTributeData();
  const playerViewData = usePlayerViewData();
  const { hand } = usePlayerHandData();

  const [selectedCards, setSelectedCards] = useState<Card[]>([]);
  const [canPlay, setCanPlay] = useState(false);
  const [isStarting, setIsStarting] = useState(false);
  const [isLeaving, setIsLeaving] = useState(false);
  const [turnDeadline, setTurnDeadline] = useState<{
    playerSeat: number;
    deadlineAtMs: number;
  } | null>(null);

  // 初始加载房间信息
  useEffect(() => {
    const loadRoomDetails = async () => {
      if (!roomId) return;

      setLoading(true);
      try {
        const response = await apiClient.getRoomDetails(roomId);
        if (response.success && response.data) {
          setCurrentRoom(response.data);
          // Only set to WAITING_PLAYERS if we're not already in a game phase
          // This prevents resetting the phase when room updates during gameplay
          if (currentPhase === GamePageState.WAITING_PLAYERS || !currentPhase) {
            setCurrentPhase(GamePageState.WAITING_PLAYERS);
          }
        } else {
          setRoomError('加载房间信息失败');
          navigate('/lobby');
        }
      } catch (error) {
        console.error('Failed to load room:', error);
        setRoomError('加载房间信息失败');
        navigate('/lobby');
      } finally {
        setLoading(false);
      }
    };

    if (!currentRoom && roomId) {
      loadRoomDetails();
    }
  }, [roomId, currentRoom, setCurrentRoom, setRoomError, setLoading, navigate, setCurrentPhase, currentPhase]);

  // WebSocket 事件监听
  useEffect(() => {
    if (!roomId || !isConnected) return;

    // 房间更新
    const handleRoomUpdate = (message: WSMessage) => {
      // 后端发送的数据结构: { action: "...", room: {...}, player_id: "..." }
      const roomData = message.data?.room || message.data;
      // 只有当房间数据有id且匹配当前房间时，才更新
      if (roomData?.id && roomData.id === roomId) {
        setCurrentRoom(roomData);
      }
    };

    // 游戏准备
    const handleGamePrepare = (message: WSMessage) => {
      // 若无room_id字段或匹配当前房间，则接受
      if (!message.data?.room_id || message.data.room_id === roomId) {
        setCurrentPhase(GamePageState.GAME_PREPARE);
        setCountdown(3);
      }
    };

    // 倒计时
    const handleCountdown = (message: WSMessage) => {
      // 若无room_id字段或匹配当前房间，则接受
      if (!message.data?.room_id || message.data.room_id === roomId) {
        const countdownValue = message.data.countdown;
        if (typeof countdownValue === 'number') {
          setCountdown(countdownValue);
        }
      }
    };

    // 游戏开始
    const handleGameBegin = (message: WSMessage) => {
      // 若无room_id字段或匹配当前房间，则接受
      if (!message.data?.room_id || message.data.room_id === roomId) {
        setCountdown(null);
        // 不立即切换阶段，等待 player_view 消息
        // player_view 会携带完整的游戏状态，收到后再切换阶段
      }
    };

    // 游戏事件（仅日志，phase 由 PlayerView 驱动）
    const handleGameEvent = (message: WSMessage) => {
      const event: GameEvent = message.data as GameEvent;
      console.log('[game_event]', eventTypeToJSON(event.type), event);

      if (event.type === EventType.EVENT_TYPE_PLAYER_PLAYED && event.playerPlayed) {
        audioService.playCardSound(event.playerPlayed.compType);
      } else if (event.type === EventType.EVENT_TYPE_PLAYER_PASSED) {
        audioService.playPassSound();
      } else if (event.type === EventType.EVENT_TYPE_PLAYER_FINISHED && event.playerFinished) {
        audioService.playRankSound(event.playerFinished.finishRank);
        if (event.actorSeat === playerSeat) {
          setMyFinishRank(event.playerFinished.finishRank);
        }
      } else if (event.type === EventType.EVENT_TYPE_DEAL_STARTED) {
        audioService.playDealStartSound();
        setMyFinishRank(null);
      } else if (event.type === EventType.EVENT_TYPE_DEAL_ENDED && event.dealEnded) {
        const myTeam = playerSeat !== null ? playerSeat % 2 : -1;
        const isWinner = myTeam === event.dealEnded.winningTeam;
        audioService.playDealEndSound(isWinner);
      } else if (event.type === EventType.EVENT_TYPE_MATCH_ENDED && event.matchEnded) {
        const myTeam = playerSeat !== null ? playerSeat % 2 : -1;
        const isWinner = myTeam === event.matchEnded.winner;
        audioService.playMatchEndSound(isWinner);
        setMyFinishRank(null);
      }
    };

    // 玩家视角
    const handlePlayerView = (message: WSMessage) => {
      console.log('[player_view]', message.data);

      const data = message.data;
      if (!data || !data.player_view) {
        console.warn('[player_view] no data, returning');
        return;
      }

      const protoPlayerView: ProtoPlayerView = data.player_view;
      
      setPlayerSeat(protoPlayerView.playerSeat);

      const isPlayingPhase = protoPlayerView.dealStatus === DealStatus.PLAYING;
      const isMyTurnValue = isPlayingPhase &&
        protoPlayerView.currentTurn === protoPlayerView.playerSeat;
      const handLen = protoPlayerView.playerCards?.length ?? 0;
      const canPlayValue = isMyTurnValue && handLen > 0;

      setCanPlay(canPlayValue);
      setMyTurn(isMyTurnValue);

      // 根据 dealStatus 推导 currentPhase（仅在变化时更新）
      let newPhase: string | null = null;
      switch (protoPlayerView.dealStatus) {
        case DealStatus.TRIBUTE:
          newPhase = GamePageState.TRIBUTE_PHASE;
          break;
        case DealStatus.PLAYING:
          newPhase = GamePageState.PLAYING;
          break;
        case DealStatus.FINISHED:
          newPhase = protoPlayerView.matchResult
            ? GamePageState.MATCH_RESULT
            : GamePageState.DEAL_RESULT;
          break;
      }
      if (newPhase && newPhase !== useGameStore.getState().currentPhase) {
        console.log('[Phase Transition]', protoPlayerView.dealStatus, '→', newPhase);
        setCurrentPhase(newPhase);
      }

      setPlayerView(protoPlayerView);
    };

    const handleGameAction = (message: WSMessage) => {
      console.log('[game_action]', message);
      const wrapper = message.data as { game_action?: GameActionData };
      const actionData = wrapper.game_action;
      if (!actionData) return;

      if (actionData.timeout !== undefined) {
        setCountdown(actionData.timeout);
      }
    };

    const handleTurnDeadline = (message: WSMessage) => {
      console.log('[turn_deadline]', message);
      const wrapper = message.data as { turn_deadline?: TurnDeadlineData };
      const data = wrapper.turn_deadline;
      if (!data) return;
      
      setTurnDeadline({
        playerSeat: data.playerSeat,
        deadlineAtMs: data.deadlineAtMs,
      });
    };

    const handleTributeView = (message: WSMessage) => {
      console.log('[tribute_view]', message.data);
      const tributeView = message.data?.tribute_view;
      if (tributeView) {
        useTributeStore.getState().setTributeView(tributeView);
      }
    };

    // 注册所有监听器
    wsClient.on(WS_MESSAGE_TYPES.ROOM_UPDATE, handleRoomUpdate);
    wsClient.on(WS_MESSAGE_TYPES.GAME_PREPARE, handleGamePrepare);
    wsClient.on(WS_MESSAGE_TYPES.COUNTDOWN, handleCountdown);
    wsClient.on(WS_MESSAGE_TYPES.GAME_BEGIN, handleGameBegin);
    wsClient.on(WS_MESSAGE_TYPES.GAME_EVENT, handleGameEvent);
    wsClient.on(WS_MESSAGE_TYPES.PLAYER_VIEW, handlePlayerView);
    wsClient.on(WS_MESSAGE_TYPES.GAME_ACTION, handleGameAction);
    wsClient.on(WS_MESSAGE_TYPES.TURN_DEADLINE, handleTurnDeadline);
    wsClient.on(WS_MESSAGE_TYPES.TRIBUTE_VIEW, handleTributeView);

    // 清理函数
    return () => {
      wsClient.off(WS_MESSAGE_TYPES.ROOM_UPDATE, handleRoomUpdate);
      wsClient.off(WS_MESSAGE_TYPES.GAME_PREPARE, handleGamePrepare);
      wsClient.off(WS_MESSAGE_TYPES.COUNTDOWN, handleCountdown);
      wsClient.off(WS_MESSAGE_TYPES.GAME_BEGIN, handleGameBegin);
      wsClient.off(WS_MESSAGE_TYPES.GAME_EVENT, handleGameEvent);
      wsClient.off(WS_MESSAGE_TYPES.PLAYER_VIEW, handlePlayerView);
      wsClient.off(WS_MESSAGE_TYPES.GAME_ACTION, handleGameAction);
      wsClient.off(WS_MESSAGE_TYPES.TURN_DEADLINE, handleTurnDeadline);
      wsClient.off(WS_MESSAGE_TYPES.TRIBUTE_VIEW, handleTributeView);
    };
  }, [roomId, isConnected, playerSeat, setCurrentPhase, setCountdown, setPlayerView, setPlayerSeat, setMyTurn, setCurrentRoom, setMyFinishRank]);

  // 连接成功后请求同步游戏状态（防抖 300ms，避免网络抖动时重复请求）
  useEffect(() => {
    if (!isConnected || !roomId) return;

    const timer = setTimeout(() => {
      wsClient.send(WS_MESSAGE_TYPES.SYNC_GAME_STATE, { room_id: roomId });
    }, 300);

    return () => clearTimeout(timer);
  }, [isConnected, roomId]);

  // 开始游戏
  const handleStartGame = async () => {
    if (!currentRoom || !user || isStarting) return;

    // 检查是否是房主
    if (currentRoom.owner !== user.id) {
      setRoomError('只有房主可以开始游戏');
      return;
    }

    // 检查是否有4名玩家
    const playerCount = currentRoom.players.filter(p => p !== null).length;
    if (playerCount < 4) {
      setRoomError('需要4名玩家才能开始游戏');
      return;
    }

    setIsStarting(true);
    try {
      const response = await apiClient.startGame(currentRoom.id);
      if (!response.success) {
        setRoomError('开始游戏失败');
      }
    } catch (error) {
      console.error('Failed to start game:', error);
      setRoomError('开始游戏失败');
    } finally {
      setIsStarting(false);
    }
  };

  // 离开房间
  const handleLeaveRoom = async () => {
    if (!currentRoom || !user || isLeaving) return;

    setIsLeaving(true);
    try {
      await apiClient.leaveRoom(currentRoom.id);
      navigate('/lobby', { state: { shouldRefresh: true } });
    } catch (error) {
      console.error('Failed to leave room:', error);
      setRoomError('离开房间失败');
    } finally {
      setIsLeaving(false);
    }
  };

  // 出牌
  const handlePlayCards = async (cards: Card[]) => {
    if (!currentRoom || !user || !canPlay || playerSeat === null) return;

    try {
      const deckIndexes = cards.map(c => c.deckIndex);
      const response = await apiClient.playCards(currentRoom.id, playerSeat, deckIndexes);
      if (response.success) {
        setSelectedCards([]);
      } else {
        setRoomError('出牌失败');
      }
    } catch (error) {
      console.error('Failed to play cards:', error);
      setRoomError('出牌失败');
    }
  };

  // 不出
  const handlePass = async () => {
    if (!currentRoom || !user || !canPlay || playerSeat === null) return;

    setSelectedCards([]);

    try {
      const response = await apiClient.pass(currentRoom.id, playerSeat);
      if (!response.success) {
        setRoomError('操作失败');
      }
    } catch (error) {
      console.error('Failed to pass:', error);
      setRoomError('操作失败');
    }
  };

  // 选择贡牌
  const handleSelectTribute = async (deckIndex: number) => {
    if (!currentRoom || playerSeat === null) return;

    try {
      const response = await apiClient.selectTribute(currentRoom.id, playerSeat, deckIndex);
      if (!response.success) {
        setRoomError('选择贡牌失败');
      }
    } catch (error) {
      console.error('Failed to select tribute:', error);
      setRoomError('选择贡牌失败');
    }
  };

  // 还贡
  const handleReturnTribute = async (deckIndex: number) => {
    if (!currentRoom || playerSeat === null) return;

    try {
      const response = await apiClient.returnTribute(currentRoom.id, playerSeat, deckIndex);
      if (response.success) {
        setSelectedCards([]);
      } else {
        setRoomError('还贡失败');
      }
    } catch (error) {
      console.error('Failed to return tribute:', error);
      setRoomError('还贡失败');
    }
  };

  // 继续游戏（局结算后）
  const handleContinue = useCallback(() => {
    setCurrentPhase(GamePageState.WAITING_PLAYERS);
  }, [setCurrentPhase]);

  // 再来一局（比赛结束后）
  const handlePlayAgain = async () => {
    if (!currentRoom || !user) return;

    setCurrentPhase(GamePageState.WAITING_PLAYERS);

    // 如果是房主，自动开始游戏
    if (currentRoom.owner === user.id) {
      await handleStartGame();
    }
  };

  // 返回大厅
  const handleReturnToLobby = async () => {
    await handleLeaveRoom();
  };

  // 渲染游戏准备倒计时
  const renderGamePrepare = () => {
    return (
      <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
        <div className="relative bg-gradient-to-b from-surface-elevated to-surface-base rounded-2xl p-10 text-center max-w-md mx-4 border border-stroke/50 shadow-card-elevated">
          {/* 顶部装饰线 */}
          <div className="absolute inset-x-0 top-0 h-1 rounded-t-2xl bg-gradient-to-r from-transparent via-state-active to-transparent" />

          {/* 装饰性扑克牌花色 */}
          <div className="flex items-center justify-center gap-3 mb-4 opacity-30">
            <span className="text-2xl">♠</span>
            <span className="text-2xl text-suit-red">♥</span>
            <span className="text-2xl">♣</span>
            <span className="text-2xl text-suit-red">♦</span>
          </div>

          <h2 className="text-2xl font-display font-bold bg-gradient-primary bg-clip-text text-transparent mb-6">
            游戏即将开始
          </h2>

          {/* 倒计时数字 */}
          <div className="relative mb-6">
            <div className="text-7xl font-display font-bold bg-gradient-active bg-clip-text text-transparent animate-pulse-glow">
              {countdown}
            </div>
            {/* 光晕效果 */}
            <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
              <div className="w-24 h-24 rounded-full bg-state-active/20 blur-xl animate-pulse" />
            </div>
          </div>

          <p className="text-fg-secondary mb-6">请准备好开始游戏...</p>

          {/* 连接状态指示器 */}
          <div className="flex items-center justify-center gap-2 px-4 py-2 rounded-full bg-surface-elevated/50 border border-stroke/30 inline-flex mx-auto">
            <div className={`w-2.5 h-2.5 rounded-full ${isConnected
              ? 'bg-action-primary shadow-jade'
              : 'bg-action-danger shadow-ruby'
              }`}
            />
            <span className="text-sm text-fg-secondary">
              {isConnected ? '连接正常' : '连接断开'}
            </span>
          </div>
        </div>
      </div>
    );
  };

  const renderPlaying = () => {
    if (!playerViewData || !currentRoom || playerSeat === null) return null;

    const { teamLevels, dealLevel, plays, currentTurn, leader } = playerViewData;
    const players = currentRoom.players;

    return (
      <div className="fixed inset-0 z-40 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
        {/* 可收起的 Header */}
        <Header collapsible />

        {/* 牌桌 - 占满屏幕 */}
        <div className="absolute inset-0 p-2">
          <GameBoard
            teamLevels={teamLevels || [2, 2]}
            currentLevel={dealLevel || 2}
            plays={plays || []}
            currentTurn={currentTurn ?? -1}
            players={players}
            currentPlayerSeat={playerSeat}
            turnDeadline={turnDeadline}
            className="h-full"
          />
        </div>

        {/* 控制按钮 - 紧贴出牌区下边缘 */}
        <div
          className="absolute left-1/2 z-20"
          style={{
            transform: 'translateX(-50%)',
            // 定位: 50% + 出牌区偏移 + 出牌区高度/2
            top: 'calc(50% + var(--play-area-offset-y, 0) + var(--table-center-height) / 2)',
          }}
        >
          <GameControls
            selectedCards={selectedCards}
            canPlay={canPlay}
            isMyTurn={isMyTurn}
            turnDeadlineAtMs={turnDeadline?.deadlineAtMs || 0}
            onPlayCards={handlePlayCards}
            onPass={handlePass}
            onHint={(cards) => setSelectedCards(cards)}
            handCards={hand}
            plays={plays || []}
            leader={leader}
            playerSeat={playerSeat}
            dealLevel={dealLevel || 2}
            disabled={false}
          />
        </div>

        {/* 手牌区域 - 固定在底部，浮在牌桌上方 */}
        <div className="absolute bottom-0 left-0 right-0 z-10">
          <PlayerHand
            cards={hand}
            selectedCards={selectedCards}
            onCardSelect={setSelectedCards}
            currentLevel={dealLevel}
            finishRank={myFinishRank}
          />
        </div>
      </div>
    );
  };

  // 根据当前阶段渲染对应界面
  const renderCurrentPhase = () => {
    const renderWaitingBoard = () => {
      if (!currentRoom) return null;
      return (
        <div className="fixed inset-0 z-40 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
          <Header collapsible />
          <div className="absolute inset-0 p-2">
            <WaitingBoard
              players={currentRoom.players}
              currentPlayerSeat={playerSeat ?? 0}
              roomId={currentRoom.id}
              roomCode={currentRoom.room_code || ''}
              ownerId={currentRoom.owner}
              currentUserId={user?.id || ''}
              isConnected={isConnected}
              onStartGame={handleStartGame}
              onLeaveRoom={handleLeaveRoom}
              isStarting={isStarting}
              isLeaving={isLeaving}
              className="h-full"
            />
          </div>
        </div>
      );
    };

    switch (currentPhase) {
      case GamePageState.WAITING_PLAYERS:
        return renderWaitingBoard();

      case GamePageState.GAME_PREPARE:
        return (
          <>
            {renderWaitingBoard()}
            {renderGamePrepare()}
          </>
        );

      case GamePageState.TRIBUTE_PHASE: {
        if (!tributeData || !currentRoom || playerSeat === null) return null;
        const canReturnTribute = 
          tributeData.status === TributeStatus.TRIBUTE_STATUS_RETURNING &&
          tributeData.receivers.includes(playerSeat);
        return (
          <div className="fixed inset-0 z-40 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
            <Header collapsible />
            <div className="absolute inset-0 p-2">
              <TributeBoard
                tributeData={tributeData}
                players={currentRoom.players}
                currentPlayerSeat={playerSeat}
                teamLevels={playerViewData?.teamLevels || [2, 2]}
                currentLevel={playerViewData?.dealLevel || 2}
                onSelectTribute={handleSelectTribute}
                className="h-full"
              />
            </div>
            {canReturnTribute && (
              <div
                className="absolute left-1/2 z-20"
                style={{
                  transform: 'translateX(-50%)',
                  top: 'calc(50% + var(--play-area-offset-y, 0) + var(--table-center-height) / 2)',
                }}
              >
                <TributeControls
                  selectedCards={selectedCards}
                  canReturnTribute={true}
                  turnDeadlineAtMs={turnDeadline?.deadlineAtMs || 0}
                  onReturnTribute={() => {
                    if (selectedCards.length === 1) {
                      handleReturnTribute(selectedCards[0].deckIndex);
                    }
                  }}
                  onHint={setSelectedCards}
                  handCards={hand}
                  dealLevel={playerViewData?.dealLevel || 2}
                />
              </div>
            )}
            {hand.length > 0 && (
              <div className="absolute bottom-0 left-0 right-0 z-10">
                <PlayerHand
                  cards={hand}
                  selectedCards={selectedCards}
                  onCardSelect={setSelectedCards}
                  currentLevel={playerViewData?.dealLevel || 2}
                />
              </div>
            )}
          </div>
        );
      }

      case GamePageState.PLAYING:
        return renderPlaying();

      case GamePageState.DEAL_RESULT:
        return dealResultData ? (
          <DealResult
            {...dealResultData}
            onContinue={handleContinue}
            onExit={handleReturnToLobby}
            isMatchFinished={false}
            currentPlayerSeat={playerSeat ?? 0}
          />
        ) : null;

      case GamePageState.MATCH_RESULT:
        return matchResultData ? (
          <MatchResult
            {...matchResultData}
            onReturnToLobby={handleReturnToLobby}
            onPlayAgain={handlePlayAgain}
          />
        ) : null;

      default:
        return renderWaitingBoard();
    }
  };

  if (!user) {
    return (
      <OrientationGuard>
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-center">
            <div className="flex items-center justify-center gap-2 mb-4 opacity-40">
              <span className="text-3xl">♠</span>
              <span className="text-3xl text-suit-red">♥</span>
              <span className="text-3xl">♣</span>
              <span className="text-3xl text-suit-red">♦</span>
            </div>
            <p className="text-fg-secondary">请先登录</p>
          </div>
        </div>
      </OrientationGuard>
    );
  }

  if (!currentRoom && currentPhase === GamePageState.WAITING_PLAYERS) {
    return (
      <OrientationGuard>
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-center">
            {/* 加载动画 */}
            <div className="relative w-16 h-16 mx-auto mb-6">
              <div className="absolute inset-0 rounded-full border-4 border-state-active/20" />
              <div className="absolute inset-0 rounded-full border-4 border-transparent border-t-state-active animate-spin" />
              <div className="absolute inset-2 rounded-full bg-surface-elevated/50 flex items-center justify-center">
                <span className="text-state-active text-lg">♠</span>
              </div>
            </div>
            <p className="text-fg-secondary font-medium">加载房间信息...</p>
          </div>
        </div>
      </OrientationGuard>
    );
  }

  return (
    <OrientationGuard>
      <div className="min-h-screen">
        {renderCurrentPhase()}
      </div>
    </OrientationGuard>
  );
};

export default GamePage;

