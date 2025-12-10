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
import PlayerControlPanel from './PlayerControlPanel';
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
      } else if (event.type === EventType.EVENT_TYPE_DEAL_STARTED) {
        audioService.playDealStartSound();
      } else if (event.type === EventType.EVENT_TYPE_DEAL_ENDED && event.dealEnded) {
        const myTeam = playerSeat !== null ? playerSeat % 2 : -1;
        const isWinner = myTeam === event.dealEnded.winningTeam;
        audioService.playDealEndSound(isWinner);
      } else if (event.type === EventType.EVENT_TYPE_MATCH_ENDED && event.matchEnded) {
        const myTeam = playerSeat !== null ? playerSeat % 2 : -1;
        const isWinner = myTeam === event.matchEnded.winner;
        audioService.playMatchEndSound(isWinner);
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
      const actionData = message.data as GameActionData;

      if (actionData.timeout !== undefined) {
        setCountdown(actionData.timeout);
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
  }, [roomId, isConnected, playerSeat, setCurrentPhase, setCountdown, setPlayerView, setPlayerSeat, setMyTurn, setCurrentRoom]);

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
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-card rounded-lg p-8 text-center max-w-md mx-4">
          <h2 className="text-2xl font-bold text-foreground mb-4">游戏即将开始</h2>
          <div className="text-6xl font-bold text-primary mb-4">
            {countdown}
          </div>
          <p className="text-muted-foreground">请准备好开始游戏...</p>

          {/* Connection status indicator */}
          <div className="mt-4 flex items-center justify-center space-x-2">
            <div className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'
              }`} />
            <span className="text-sm text-muted-foreground">
              {isConnected ? '连接正常' : '连接断开'}
            </span>
          </div>
        </div>
      </div>
    );
  };

  const renderPlaying = () => {
    if (!playerViewData || !currentRoom || playerSeat === null) return null;

    const { teamLevels, dealLevel, plays, currentTurn, playStates, leader } = playerViewData;
    const players = currentRoom.players;

    const handleHint = (cards: Card[]) => {
      setSelectedCards(cards);
    };

    return (
      <div className="max-w-6xl mx-auto p-6">
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

        <div className="-mt-24 relative z-10">
          <PlayerControlPanel
            cards={hand}
            selectedCards={selectedCards}
            onCardSelect={setSelectedCards}
            currentLevel={dealLevel}
            canPlay={canPlay}
            isMyTurn={isMyTurn}
            turnDeadlineAtMs={turnDeadline?.deadlineAtMs || 0}
            onPlayCards={handlePlayCards}
            onPass={handlePass}
            onHint={handleHint}
            plays={plays || []}
            leader={leader}
            playerSeat={playerSeat}
            dealLevel={dealLevel || 2}
            disabled={false}
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
        <WaitingBoard
          players={currentRoom.players}
          currentPlayerSeat={playerSeat ?? 0}
          roomId={currentRoom.id}
          ownerId={currentRoom.owner}
          currentUserId={user?.id || ''}
          isConnected={isConnected}
          onStartGame={handleStartGame}
          onLeaveRoom={handleLeaveRoom}
          isStarting={isStarting}
          isLeaving={isLeaving}
        />
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

      case GamePageState.TRIBUTE_PHASE:
        return tributeData && currentRoom && playerSeat !== null ? (
          <TributeBoard
            tributeData={tributeData}
            players={currentRoom.players}
            currentPlayerSeat={playerSeat}
            teamLevels={playerViewData?.teamLevels || [2, 2]}
            currentLevel={playerViewData?.dealLevel || 2}
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
      <div className="flex items-center justify-center min-h-[400px]">
        <p className="text-muted-foreground">请先登录</p>
      </div>
    );
  }

  if (!currentRoom && currentPhase === GamePageState.WAITING_PLAYERS) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">加载房间信息...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {renderCurrentPhase()}
    </div>
  );
};

export default GamePage;

