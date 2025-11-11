import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { useGameStore } from '../../store/gameStore';
import { wsClient } from '../../services/websocket';
import { apiClient } from '../../services/api';
import type { Room, Player, WSMessage, Card, GameActionData } from '../../types';
import { WS_MESSAGE_TYPES } from '../../types';
import GameBoard from './GameBoard';
import PlayerHand from './PlayerHand';
import GameControls from './GameControls';
import TributePhase from './TributePhase';
import DealResult from './DealResult';
import MatchResult from './MatchResult';

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
    gameState,
    currentPhase,
    selectedCards,
    dealResult,
    matchResult,
    tributeInfo,
    isMyTurn,
    setCountdown,
    setCurrentPhase,
    setSelectedCards,
    setDealResult,
    setMatchResult,
    setTributeInfo,
    setPlayerSeat,
    setGameState,
    setMyTurn
  } = useGameStore();

  const [room, setRoom] = useState<Room | null>(currentRoom);
  const [playerHand, setPlayerHand] = useState<Card[]>([]);
  const [canPlay, setCanPlay] = useState(false);
  const [isStarting, setIsStarting] = useState(false);
  const [isLeaving, setIsLeaving] = useState(false);
  const [turnTimeoutSeconds, setTurnTimeoutSeconds] = useState<number>(20);

  // 初始加载房间信息
  useEffect(() => {
    const loadRoomDetails = async () => {
      if (!roomId) return;

      setLoading(true);
      try {
        const response = await apiClient.getRoomDetails(roomId);
        if (response.success && response.data) {
          setRoom(response.data);
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
    } else if (currentRoom) {
      setRoom(currentRoom);
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
        setRoom(roomData);
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
      const { event_type, event_data, player_seat, ...restData } = message.data || {};
      // 兼容两种数据结构：优先使用event_data，回退到data本身
      const payload = event_data || restData;

      // 🎯 Detailed Game Event Logging for Debug
      console.group(`🎯 [GamePage] Processing Event: ${event_type}`);
      console.log('📌 Event Type:', event_type);
      console.log('👤 Player Seat:', player_seat);
      console.log('📦 Payload:', payload);
      console.log('🎮 Current Phase:', currentPhase);
      
      switch (event_type) {
        case 'tribute_started':
          console.log('➡️ Action: Set phase to TRIBUTE_PHASE');
          setCurrentPhase(GamePageState.TRIBUTE_PHASE);
          setTributeInfo(payload.tribute_info);
          break;
        case 'tribute_completed':
          console.log('➡️ Action: Set phase to PLAYING');
          setCurrentPhase(GamePageState.PLAYING);
          setTributeInfo(null);
          break;
        case 'deal_completed':
        case 'deal_ended':
          console.log('➡️ Action: Set phase to DEAL_RESULT');
          setCurrentPhase(GamePageState.DEAL_RESULT);
          setDealResult(payload.deal_result || payload);
          break;
        case 'match_completed':
        case 'match_ended':
          console.log('➡️ Action: Set phase to MATCH_RESULT');
          setCurrentPhase(GamePageState.MATCH_RESULT);
          setMatchResult(payload.match_result || payload);
          break;
        case 'trick_started':
          console.log('🎲 Trick Started:', payload.trick);
          // Note: Game state is updated from player_view messages, not from game_events
          // game_events are just notifications and don't contain full game state
          break;
        case 'player_played':
          console.log('🃏 Player Played:', {
            player: player_seat,
            cards: payload.cards
          });
          break;
        case 'player_passed':
          console.log('⏭️ Player Passed:', player_seat);
          break;
        case 'trick_ended':
        case 'trick_completed':
          console.log('✅ Trick Ended:', {
            winner: payload.winner,
            next_leader: payload.next_leader
          });
          break;
        case 'deal_started':
          console.log('🎴 Deal Started:', {
            level: payload.deal_level,
            team0_level: payload.team0_level,
            team1_level: payload.team1_level
          });
          break;
        case 'tribute_rules_set':
          console.log('📜 Tribute Rules:', payload);
          break;
        case 'tribute_immunity':
          console.log('🛡️ Tribute Immunity:', payload);
          break;
        case 'tribute_given':
          console.log('⬆️ Tribute Given:', {
            giver: payload.giver,
            receiver: payload.receiver,
            card: payload.card
          });
          break;
        case 'tribute_selected':
          console.log('✅ Tribute Selected:', payload);
          break;
        case 'return_tribute':
          console.log('⬇️ Return Tribute:', payload);
          break;
        default:
          console.log('ℹ️ Unhandled event type');
          break;
      }
      console.groupEnd();
    };

    // 玩家视角
    const handlePlayerView = (message: WSMessage) => {
      // Get the latest phase value from store instead of using closure
      const latestPhase = useGameStore.getState().currentPhase;
      
      console.log('👁️ [GamePage] Received PLAYER_VIEW, current phase:', latestPhase);
      
      console.log('🎮 [GamePage] Received player_view message', {
        timestamp: new Date().toISOString(),
        messageType: message.type,
        hasData: !!message.data,
        closurePhase: currentPhase,
        latestPhase: latestPhase,
        roomId
      });

      const data = message.data;
      if (!data) {
        console.warn('⚠️ [GamePage] player_view has no data, returning');
        return;
      }

      // GameDriver 发送的格式：data.player_view 包含实际数据
      const playerView = data.player_view;
      if (!playerView) {
        console.warn('⚠️ [GamePage] Received player_view message without player_view data:', data);
        return;
      }

      console.log('✅ [GamePage] player_view data structure:', {
        hasPlayerSeat: playerView.player_seat !== undefined,
        playerSeat: playerView.player_seat ?? data.player_seat,
        hasPlayerCards: !!playerView.player_cards,
        playerCardsCount: playerView.player_cards?.length || 0,
        hasGameState: !!playerView.game_state,
        eventType: data.event_type
      });

      // 设置玩家座位
      const seatNumber = playerView.player_seat ?? data.player_seat;
      console.log('👤 [GamePage] Setting player seat:', seatNumber);
      setPlayerSeat(seatNumber);
      
      // 设置手牌 - 从 player_cards 字段获取
      const cards = playerView.player_cards || [];
      console.log('🃏 [GamePage] Setting player hand:', {
        cardCount: cards.length,
        cards: cards.map((c: Card) => ({ id: c.id, suit: c.suit, rank: c.rank }))
      });
      setPlayerHand(cards);
      
      // 设置游戏状态
      if (playerView.game_state) {
        console.log('🎲 [GamePage] Game state received:', {
          hasCurrentMatch: !!playerView.game_state.current_match,
          matchStatus: playerView.game_state.current_match?.status,
          teamLevels: playerView.game_state.current_match?.team_levels,
          hasCurrentDeal: !!playerView.game_state.current_match?.current_deal,
          dealStatus: playerView.game_state.current_match?.current_deal?.status,
          hasTributePhase: playerView.game_state.current_match?.current_deal?.tribute_phase !== null
        });

        setGameState(playerView.game_state);
        
        // 根据当前阶段和游戏状态决定是否切换阶段
        // Use latest phase from store to avoid stale closure
        if (latestPhase === GamePageState.GAME_PREPARE) {
          console.log('🔄 [GamePage] Currently in GAME_PREPARE, checking phase transition...');
          
          // 判断是否有上贡阶段
          const currentDeal = playerView.game_state.current_match?.current_deal;
          const hasTribute = currentDeal?.tribute_phase !== null && currentDeal?.tribute_phase !== undefined;
          
          console.log('📊 [GamePage] Phase transition check:', {
            hasCurrentDeal: !!currentDeal,
            hasTribute,
            tributePhase: currentDeal?.tribute_phase,
            willTransitionTo: hasTribute ? 'TRIBUTE_PHASE' : 'PLAYING'
          });
          
          if (hasTribute) {
            console.log('🔀 [GamePage] ✅ Transitioning to TRIBUTE_PHASE');
            setCurrentPhase(GamePageState.TRIBUTE_PHASE);
          } else {
            console.log('🔀 [GamePage] ✅ Transitioning to PLAYING');
            setCurrentPhase(GamePageState.PLAYING);
          }
        } else {
          console.log('ℹ️ [GamePage] Not in GAME_PREPARE phase, skipping phase transition. Current phase:', latestPhase);
        }
      } else {
        console.warn('⚠️ [GamePage] player_view has no game_state, skipping game state update');
      }

      // 从 player_view 中计算 can_play 和 is_my_turn
      // 必须满足以下条件才能出牌：
      // 1. deal.status === 'playing' (出牌阶段，而非上贡阶段)
      // 2. 存在 currentTrick
      // 3. currentTrick.current_turn === playerSeat (轮到自己)
      // 4. 有手牌
      const currentDeal = playerView.game_state?.current_match?.current_deal;
      const currentTrick = currentDeal?.current_trick;
      const isPlayingPhase = currentDeal?.status === 'playing';
      const isMyTurnValue = isPlayingPhase && 
                           currentTrick !== null && 
                           currentTrick !== undefined &&
                           currentTrick.current_turn === playerView.player_seat;
      const canPlayValue = isMyTurnValue && (playerView.player_cards?.length || 0) > 0;
      
      console.log('🎯 [GamePage] Player action state:', {
        canPlay: canPlayValue,
        isMyTurn: isMyTurnValue,
        dealStatus: currentDeal?.status,
        isPlayingPhase: isPlayingPhase,
        currentTurn: currentTrick?.current_turn,
        playerSeat: playerView.player_seat,
        hasCards: (playerView.player_cards?.length || 0) > 0
      });
      
      setCanPlay(canPlayValue);
      setMyTurn(isMyTurnValue);

      console.log('✨ [GamePage] player_view processing completed');
    };

    const handleGameAction = (message: WSMessage) => {
      console.log('🎮 [GamePage] Received game_action:', message);
      const actionData = message.data as GameActionData;
      
      if (actionData.timeout !== undefined) {
        console.log(`⏱️ [GamePage] Setting timeout to ${actionData.timeout} seconds`);
        setTurnTimeoutSeconds(actionData.timeout);
      }
      
      if (actionData.player_seat === playerSeat) {
        console.log('🎯 [GamePage] This action is for current player');
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

    // 加入房间
    wsClient.send(WS_MESSAGE_TYPES.JOIN_ROOM, { room_id: roomId });

    // 清理函数
    return () => {
      wsClient.off(WS_MESSAGE_TYPES.ROOM_UPDATE, handleRoomUpdate);
      wsClient.off(WS_MESSAGE_TYPES.GAME_PREPARE, handleGamePrepare);
      wsClient.off(WS_MESSAGE_TYPES.COUNTDOWN, handleCountdown);
      wsClient.off(WS_MESSAGE_TYPES.GAME_BEGIN, handleGameBegin);
      wsClient.off(WS_MESSAGE_TYPES.GAME_EVENT, handleGameEvent);
      wsClient.off(WS_MESSAGE_TYPES.PLAYER_VIEW, handlePlayerView);
      wsClient.off(WS_MESSAGE_TYPES.GAME_ACTION, handleGameAction);
    };
  }, [roomId, isConnected, setCurrentPhase, setCountdown, setTributeInfo, setDealResult, setMatchResult, setGameState, setPlayerSeat, setMyTurn, setCurrentRoom]);

  // 开始游戏
  const handleStartGame = async () => {
    if (!room || !user || isStarting) return;

    // 检查是否是房主
    if (room.owner !== user.id) {
      setRoomError('只有房主可以开始游戏');
      return;
    }

    // 检查是否有4名玩家
    const playerCount = room.players.filter(p => p !== null).length;
    if (playerCount < 4) {
      setRoomError('需要4名玩家才能开始游戏');
      return;
    }

    setIsStarting(true);
    try {
      const response = await apiClient.startGame(room.id);
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
    if (!room || !user || isLeaving) return;

    setIsLeaving(true);
    try {
      await apiClient.leaveRoom(room.id);
      wsClient.send(WS_MESSAGE_TYPES.LEAVE_ROOM, { room_id: room.id });
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
    if (!room || !user || !canPlay || playerSeat === null) return;

    try {
      const cardIds = cards.map(c => c.id);
      const response = await apiClient.playCards(room.id, playerSeat, cardIds);
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
    if (!room || !user || !canPlay || playerSeat === null) return;

    try {
      const response = await apiClient.pass(room.id, playerSeat);
      if (!response.success) {
        setRoomError('操作失败');
      }
    } catch (error) {
      console.error('Failed to pass:', error);
      setRoomError('操作失败');
    }
  };

  // 选择贡牌
  const handleSelectTribute = async (cardId: string) => {
    if (!room || playerSeat === null) return;

    try {
      await apiClient.selectTribute(room.id, playerSeat, cardId);
    } catch (error) {
      console.error('Failed to select tribute:', error);
      setRoomError('选择贡牌失败');
    }
  };

  // 还贡
  const handleReturnTribute = async (cardId: string) => {
    if (!room || playerSeat === null) return;

    try {
      await apiClient.returnTribute(room.id, playerSeat, cardId);
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
    if (!room || !user) return;

    setMatchResult(null);
    setCurrentPhase(GamePageState.WAITING_PLAYERS);
    
    // 如果是房主，自动开始游戏
    if (room.owner === user.id) {
      await handleStartGame();
    }
  };

  // 返回大厅
  const handleReturnToLobby = async () => {
    await handleLeaveRoom();
  };

  // 渲染等待玩家界面
  const renderWaitingPlayers = () => {
    if (!room) return null;

    const getPlayerCount = () => {
      return room.players?.filter(p => p !== null).length || 0;
    };

    const isRoomOwner = () => {
      return user && room && room.owner === user.id;
    };

    const canStartGame = () => {
      return isRoomOwner() && getPlayerCount() === 4;
    };

    const renderPlayerSeat = (seatIndex: number) => {
      const player = room.players?.[seatIndex] || null;
      const isEmpty = !player;
      const isCurrentUser = player?.id === user?.id;
      const isOwner = player?.id === room.owner;

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
                  <div className={`w-2 h-2 rounded-full ${
                    player.online ? 'bg-green-500' : 'bg-gray-400'
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
              <p className="text-gray-600">房间ID: {room.id}</p>
            </div>
            <div className="text-right">
              <div className="text-sm text-gray-600">
                玩家数量: {getPlayerCount()}/4
              </div>
              <div className="flex items-center justify-end space-x-2 mt-1">
                <div className={`w-2 h-2 rounded-full ${
                  isConnected ? 'bg-green-500' : 'bg-red-500'
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
                className={`px-8 py-2 rounded-lg font-medium transition-colors ${
                  canStartGame()
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
            <div className={`w-3 h-3 rounded-full ${
              isConnected ? 'bg-green-500' : 'bg-red-500'
            }`} />
            <span className="text-sm text-gray-600">
              {isConnected ? '连接正常' : '连接断开'}
            </span>
          </div>
        </div>
      </div>
    );
  };

  // 渲染出牌界面
  const renderPlaying = () => {
    if (!gameState || !room || playerSeat === null) return null;

    // Ensure we have the minimum required game state data
    // Backend sends nested structure: gameState.current_match.current_deal
    const currentMatch = (gameState as any).current_match;
    const currentDeal = currentMatch?.current_deal;
    
    if (!currentMatch || !currentDeal) {
      return (
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-600">加载游戏数据...</p>
          </div>
        </div>
      );
    }

    const players = room.players;
    const trickInfo = currentDeal.current_trick || null;

    return (
      <div className="max-w-6xl mx-auto p-6 space-y-6">
        {/* Game Board */}
        <GameBoard
          gameState={gameState}
          players={players}
          currentPlayerSeat={playerSeat}
          trickInfo={trickInfo}
        />

        {/* Player Hand */}
        <PlayerHand
          cards={playerHand}
          selectedCards={selectedCards}
          onCardSelect={setSelectedCards}
          disabled={!canPlay}
        />

        {/* Game Controls */}
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
        return tributeInfo && room ? (
          <TributePhase
            tributePhase={tributeInfo}
            players={room.players}
            currentPlayerSeat={playerSeat || 0}
            playerHand={playerHand}
            onSelectTribute={handleSelectTribute}
            onReturnTribute={handleReturnTribute}
          />
        ) : null;
      
      case GamePageState.PLAYING:
        return renderPlaying();
      
      case GamePageState.DEAL_RESULT:
        return dealResult && room ? (
          <DealResult
            dealResult={dealResult}
            players={room.players.filter(p => p !== null) as Player[]}
            teamLevels={(gameState as any)?.current_match?.team_levels || [2, 2]}
            onContinue={handleContinue}
            onExit={handleReturnToLobby}
            isMatchFinished={false}
          />
        ) : null;
      
      case GamePageState.MATCH_RESULT:
        return matchResult ? (
          <MatchResult
            matchResult={matchResult}
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

  if (!room && currentPhase === GamePageState.WAITING_PLAYERS) {
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

