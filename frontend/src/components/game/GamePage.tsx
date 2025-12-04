import React, { useEffect, useState } from 'react';
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
import { EventType, eventTypeToJSON } from '../../types/generated/event';
import GameBoard from './GameBoard';
import PlayerHand from './PlayerHand';
import GameControls from './GameControls';
import TributeBoard from './tribute/TributeBoard';
import DealResult from './DealResult';
import MatchResult from './MatchResult';
import { useDealResultData, useMatchResultData } from '../../hooks/useResultData';
import { useTributeData } from '../../hooks/useTributeData';
import { usePlayerViewData, usePlayerHandData } from '../../hooks/useGameData';

// 游戏页面状态常量
const GamePageState = {
  WAITING_PLAYERS: 'waiting_players',
  GAME_PREPARE: 'game_prepare',
  TRIBUTE_PHASE: 'tribute_phase',
  PLAYING: 'playing',
  DEAL_RESULT: 'deal_result',
  MATCH_RESULT: 'match_result'
} as const;

// 事件类型 → GamePage Phase 映射
const EVENT_TO_PHASE_MAP: Partial<Record<EventType, string>> = {
  [EventType.EVENT_TYPE_DEAL_STARTED]: GamePageState.PLAYING,
  [EventType.EVENT_TYPE_TRIBUTE_STARTED]: GamePageState.TRIBUTE_PHASE,
  [EventType.EVENT_TYPE_TRIBUTE_COMPLETED]: GamePageState.PLAYING,
  [EventType.EVENT_TYPE_DEAL_ENDED]: GamePageState.DEAL_RESULT,
  [EventType.EVENT_TYPE_MATCH_ENDED]: GamePageState.MATCH_RESULT,
};

// 根据游戏事件更新 Phase
const updatePhaseFromEvent = (event: GameEvent, setCurrentPhase: (phase: string) => void) => {
  const newPhase = EVENT_TO_PHASE_MAP[event.type];
  if (newPhase) {
    console.log('[Phase Transition]', eventTypeToJSON(event.type), '→', newPhase);
    setCurrentPhase(newPhase);
  }
};

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
    setCountdown,
    setCurrentPhase,
    setDealResult,
    setMatchResult,
    setPlayerSeat,
    setPlayerView,
    setMyTurn
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
  const [turnTimeoutSeconds, setTurnTimeoutSeconds] = useState<number>(20);
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

    // 游戏事件
    const handleGameEvent = (message: WSMessage) => {
      const event: GameEvent = message.data as GameEvent;
      console.log('[game_event]', eventTypeToJSON(event.type), event);
      
      const roomIdFromMsg = message.data?.room_id;
      if (roomIdFromMsg && roomIdFromMsg !== roomId) return;

      // 统一处理 phase 转换
      updatePhaseFromEvent(event, setCurrentPhase);

      const tributeActions = useTributeStore.getState();

      // 处理事件数据
      switch (event.type) {
        case EventType.EVENT_TYPE_TRIBUTE_STARTED:
          if (event.tributeStarted) {
            tributeActions.handleTributeStarted(event.tributeStarted);
          }
          break;

        case EventType.EVENT_TYPE_TRIBUTE_EXEMPTED:
          if (event.tributeExempted) {
            tributeActions.handleTributeExempted(event.tributeExempted);
          }
          break;

        case EventType.EVENT_TYPE_TRIBUTE_CARD_SUBMITTED:
          if (event.tributeCardSubmitted?.submittedCard && event.actorSeat !== undefined) {
            tributeActions.handleCardSubmitted(
              event.actorSeat,
              event.tributeCardSubmitted.submittedCard as Card
            );
          }
          break;

        case EventType.EVENT_TYPE_TRIBUTE_CARD_SELECTED:
          if (event.tributeCardSelected?.selectedCard && event.actorSeat !== undefined) {
            tributeActions.handleCardSelected(
              event.actorSeat,
              event.tributeCardSelected.selectedCard as Card
            );
          }
          break;

        case EventType.EVENT_TYPE_TRIBUTE_CARD_RETURNED:
          if (event.tributeCardReturned && event.actorSeat !== undefined) {
            tributeActions.handleCardReturned(event.actorSeat, event.tributeCardReturned);
          }
          break;

        case EventType.EVENT_TYPE_TRIBUTE_COMPLETED:
          tributeActions.handleCompleted();
          break;
          
        case EventType.EVENT_TYPE_DEAL_ENDED:
          if (event.dealEnded) {
            setDealResult(event.dealEnded);
          }
          break;
          
        case EventType.EVENT_TYPE_MATCH_ENDED:
          if (event.matchEnded) {
            setMatchResult(event.matchEnded);
          }
          break;
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

      setPlayerView(protoPlayerView);
    };

    const handleGameAction = (message: WSMessage) => {
      console.log('[game_action]', message);
      const actionData = message.data as GameActionData;

      if (actionData.timeout !== undefined) {
        setTurnTimeoutSeconds(actionData.timeout);
      }
    };

    const handleTurnDeadline = (message: WSMessage) => {
      console.log('[turn_deadline]', message);
      const data = message.data as TurnDeadlineData;
      setTurnDeadline({
        playerSeat: data.player_seat,
        deadlineAtMs: data.deadline_at_ms,
      });
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
    };
  }, [roomId, isConnected, playerSeat, setCurrentPhase, setCountdown, setDealResult, setMatchResult, setPlayerView, setPlayerSeat, setMyTurn, setCurrentRoom]);

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
  const handleContinue = () => {
    setCurrentPhase(GamePageState.WAITING_PLAYERS);
    setDealResult(null);
  };

  // 再来一局（比赛结束后）
  const handlePlayAgain = async () => {
    if (!currentRoom || !user) return;

    setMatchResult(null);
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

  // 渲染等待玩家界面
  const renderWaitingPlayers = () => {
    if (!currentRoom) return null;

    const getPlayerCount = () => {
      return currentRoom.players?.filter(p => p !== null).length || 0;
    };

    const isRoomOwner = () => {
      return user && currentRoom && currentRoom.owner === user.id;
    };

    const canStartGame = () => {
      return isRoomOwner() && getPlayerCount() === 4;
    };

    const renderPlayerSeat = (seatIndex: number) => {
      const player = currentRoom.players?.[seatIndex] || null;
      const isEmpty = !player;
      const isCurrentUser = player?.id === user?.id;
      const isOwner = player?.id === currentRoom.owner;

      return (
        <div
          key={seatIndex}
          className={`
            relative p-4 rounded-lg border-2 min-h-[120px] flex flex-col items-center justify-center
            ${isEmpty
              ? 'border-dashed border-gray-300 bg-gray-50'
              : 'border-solid border-blue-300 bg-blue-50'
            }
            ${isCurrentUser ? 'ring-2 ring-blue-500' : ''}
          `}
        >
          {/* Seat number */}
          <div className="absolute top-2 left-2 text-xs text-gray-500 font-medium">
            座位 {seatIndex + 1}
          </div>

          {/* Owner badge */}
          {isOwner && (
            <div className="absolute top-2 right-2 bg-yellow-500 text-white text-xs px-2 py-1 rounded">
              房主
            </div>
          )}

          {isEmpty ? (
            <div className="text-center">
              <div className="w-12 h-12 bg-gray-200 rounded-full mb-2 flex items-center justify-center">
                <span className="text-gray-400 text-xl">+</span>
              </div>
              <span className="text-gray-500 text-sm">等待玩家</span>
            </div>
          ) : (
            <div className="text-center">
              <div className="w-12 h-12 bg-blue-500 rounded-full mb-2 flex items-center justify-center">
                <span className="text-white font-bold text-lg">
                  {player.username.charAt(0).toUpperCase()}
                </span>
              </div>
              <div className="space-y-1">
                <div className="font-medium text-gray-800">{player.username}</div>
                <div className="flex items-center justify-center space-x-2">
                  <div className={`w-2 h-2 rounded-full ${player.online ? 'bg-green-500' : 'bg-gray-400'
                    }`} />
                  <span className="text-xs text-gray-600">
                    {player.online ? '在线' : '离线'}
                  </span>
                </div>
                {player.auto_play && (
                  <div className="text-xs text-orange-600 bg-orange-100 px-2 py-1 rounded">
                    托管中
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      );
    };

    return (
      <div className="max-w-4xl mx-auto p-6">
        {/* Room header */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-2xl font-bold text-gray-800">房间等待</h1>
              <p className="text-gray-600">房间ID: {currentRoom.id}</p>
            </div>
            <div className="text-right">
              <div className="text-sm text-gray-600">
                玩家数量: {getPlayerCount()}/4
              </div>
              <div className="flex items-center justify-end space-x-2 mt-1">
                <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'
                  }`} />
                <span className="text-xs text-gray-500">
                  {isConnected ? '已连接' : '连接断开'}
                </span>
              </div>
            </div>
          </div>

          {/* Player seats grid */}
          <div className="grid grid-cols-2 gap-4 mb-6">
            {[0, 1, 2, 3].map(seatIndex => renderPlayerSeat(seatIndex))}
          </div>

          {/* Action buttons */}
          <div className="flex items-center justify-between">
            <button
              onClick={handleLeaveRoom}
              disabled={isLeaving}
              className="px-6 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isLeaving ? '离开中...' : '离开房间'}
            </button>

            {isRoomOwner() && (
              <button
                onClick={handleStartGame}
                disabled={!canStartGame() || isStarting}
                className={`px-8 py-2 rounded-lg font-medium transition-colors ${canStartGame()
                    ? 'bg-green-500 text-white hover:bg-green-600'
                    : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                  }`}
              >
                {isStarting ? '开始中...' : '开始游戏'}
              </button>
            )}
          </div>

          {/* Status messages */}
          {isRoomOwner() && getPlayerCount() < 4 && (
            <div className="mt-4 p-3 bg-yellow-100 border border-yellow-300 rounded-lg">
              <p className="text-yellow-800 text-sm">
                需要4名玩家才能开始游戏，当前有 {getPlayerCount()} 名玩家
              </p>
            </div>
          )}

          {!isRoomOwner() && (
            <div className="mt-4 p-3 bg-blue-100 border border-blue-300 rounded-lg">
              <p className="text-blue-800 text-sm">
                等待房主开始游戏...
              </p>
            </div>
          )}
        </div>
      </div>
    );
  };

  // 渲染游戏准备倒计时
  const renderGamePrepare = () => {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-white rounded-lg p-8 text-center max-w-md mx-4">
          <h2 className="text-2xl font-bold text-gray-800 mb-4">游戏即将开始</h2>
          <div className="text-6xl font-bold text-blue-500 mb-4">
            {countdown}
          </div>
          <p className="text-gray-600">请准备好开始游戏...</p>

          {/* Connection status indicator */}
          <div className="mt-4 flex items-center justify-center space-x-2">
            <div className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'
              }`} />
            <span className="text-sm text-gray-600">
              {isConnected ? '连接正常' : '连接断开'}
            </span>
          </div>
        </div>
      </div>
    );
  };

  const renderPlaying = () => {
    if (!playerViewData || !currentRoom || playerSeat === null) return null;

    const { teamLevels, dealLevel, plays, currentTurn, playStates } = playerViewData;
    const players = currentRoom.players;

    return (
      <div className="max-w-6xl mx-auto p-6 space-y-6">
        <GameBoard
          teamLevels={teamLevels || [2, 2]}
          currentLevel={dealLevel || 2}
          plays={plays || []}
          currentTurn={currentTurn ?? -1}
          players={players}
          currentPlayerSeat={playerSeat}
          playStates={playStates}
          turnDeadline={turnDeadline}
        />

        <PlayerHand
          cards={hand}
          selectedCards={selectedCards}
          onCardSelect={setSelectedCards}
        />

        <GameControls
          selectedCards={selectedCards}
          canPlay={canPlay}
          isMyTurn={isMyTurn}
          turnTimeoutSeconds={turnTimeoutSeconds}
          onPlayCards={handlePlayCards}
          onPass={handlePass}
          disabled={false}
        />
      </div>
    );
  };

  // 根据当前阶段渲染对应界面
  const renderCurrentPhase = () => {
    switch (currentPhase) {
      case GamePageState.WAITING_PLAYERS:
        return renderWaitingPlayers();

      case GamePageState.GAME_PREPARE:
        return (
          <>
            {renderWaitingPlayers()}
            {renderGamePrepare()}
          </>
        );

      case GamePageState.TRIBUTE_PHASE:
        return tributeData && currentRoom && playerSeat !== null ? (
          <TributeBoard
            tributeData={tributeData}
            players={currentRoom.players}
            currentPlayerSeat={playerSeat}
            playerHand={hand}
            selectedCards={selectedCards}
            onCardSelect={setSelectedCards}
            onSelectTribute={handleSelectTribute}
            onReturnTribute={handleReturnTribute}
          />
        ) : null;

      case GamePageState.PLAYING:
        return renderPlaying();

      case GamePageState.DEAL_RESULT:
        return dealResultData ? (
          <DealResult
            {...dealResultData}
            onContinue={handleContinue}
            onExit={handleReturnToLobby}
            isMatchFinished={false}
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
        return renderWaitingPlayers();
    }
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <p className="text-gray-600">请先登录</p>
      </div>
    );
  }

  if (!currentRoom && currentPhase === GamePageState.WAITING_PLAYERS) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-600">加载房间信息...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {renderCurrentPhase()}
    </div>
  );
};

export default GamePage;

