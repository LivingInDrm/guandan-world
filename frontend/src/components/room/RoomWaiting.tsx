import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { useRoomStore } from '../../store/roomStore';
import { useGameStore } from '../../store/gameStore';
import { gameService } from '../../services/gameService';
import { wsClient } from '../../services/websocket';
import { apiClient } from '../../services/api';
import type { Room, Player, RoomStatus, WSMessage } from '../../types';
import { RoomStatus as RoomStatusEnum, WS_MESSAGE_TYPES } from '../../types';

interface RoomWaitingProps {
  room?: Room;
  onStartGame?: () => void;
  onLeaveRoom?: () => void;
}

const RoomWaiting: React.FC<RoomWaitingProps> = ({ 
  room: propRoom, 
  onStartGame, 
  onLeaveRoom 
}) => {
  const navigate = useNavigate();
  const { roomId } = useParams<{ roomId: string }>();
  const { user } = useAuthStore();
  const { currentRoom, setCurrentRoom, setError, setLoading } = useRoomStore();
  const { countdown, isConnected, setCountdown } = useGameStore();
  const [isStarting, setIsStarting] = useState(false);
  const [isLeaving, setIsLeaving] = useState(false);
  const [showPrepare, setShowPrepare] = useState(false);

  // Use prop room or store room
  const room = propRoom || currentRoom;

  useEffect(() => {
    // Initialize game service
    gameService.initialize();

    // If no room in props and we have a roomId, fetch room details
    if (!propRoom && roomId && !currentRoom) {
      loadRoomDetails();
    }

    // Set up WebSocket handlers for room-specific events
    const handleRoomUpdate = (message: WSMessage) => {
      if (message.data && message.data.id === (room?.id || roomId)) {
        setCurrentRoom(message.data);
      }
    };

    const handleGamePrepare = (message: WSMessage) => {
      if (message.data?.room_id === (room?.id || roomId)) {
        setShowPrepare(true);
        setCountdown(3); // Start 3-second countdown
      }
    };

    const handleCountdown = (message: WSMessage) => {
      if (message.data?.room_id === (room?.id || roomId)) {
        const countdownValue = message.data.countdown;
        if (typeof countdownValue === 'number') {
          setCountdown(countdownValue);
        }
      }
    };

    const handleGameBegin = (message: WSMessage) => {
      if (message.data?.room_id === (room?.id || roomId)) {
        setShowPrepare(false);
        setCountdown(null);
        // Navigate to game page
        navigate(`/game/${room?.id || roomId}`);
      }
    };

    // Register WebSocket handlers
    wsClient.on(WS_MESSAGE_TYPES.ROOM_UPDATE, handleRoomUpdate);
    wsClient.on(WS_MESSAGE_TYPES.GAME_PREPARE, handleGamePrepare);
    wsClient.on(WS_MESSAGE_TYPES.COUNTDOWN, handleCountdown);
    wsClient.on(WS_MESSAGE_TYPES.GAME_BEGIN, handleGameBegin);

    // Join room via WebSocket if we have a room
    if ((room?.id || roomId) && isConnected) {
      wsClient.send(WS_MESSAGE_TYPES.JOIN_ROOM, { room_id: room?.id || roomId });
    }

    // Cleanup function
    return () => {
      wsClient.off(WS_MESSAGE_TYPES.ROOM_UPDATE, handleRoomUpdate);
      wsClient.off(WS_MESSAGE_TYPES.GAME_PREPARE, handleGamePrepare);
      wsClient.off(WS_MESSAGE_TYPES.COUNTDOWN, handleCountdown);
      wsClient.off(WS_MESSAGE_TYPES.GAME_BEGIN, handleGameBegin);
    };
  }, [roomId, propRoom, currentRoom, room?.id, isConnected, navigate, setCurrentRoom, setCountdown]);

  const loadRoomDetails = async () => {
    if (!roomId) return;
    
    setLoading(true);
    try {
      const response = await apiClient.getRoomDetails(roomId);
      if (response.success && response.data) {
        setCurrentRoom(response.data);
      } else {
        setError('Failed to load room details');
        navigate('/lobby');
      }
    } catch (error) {
      console.error('Failed to load room:', error);
      setError('Failed to load room details');
      navigate('/lobby');
    } finally {
      setLoading(false);
    }
  };

  const handleStartGame = async () => {
    if (!room || !user || isStarting) return;
    
    // Check if user is room owner
    if (room.owner !== user.id) {
      setError('Only room owner can start the game');
      return;
    }

    // Check if room has 4 players
    const playerCount = room.players.filter(p => p !== null).length;
    if (playerCount < 4) {
      setError('Need 4 players to start the game');
      return;
    }

    setIsStarting(true);
    try {
      if (onStartGame) {
        onStartGame();
      } else {
        // Use game service which handles both API and WebSocket
        const success = await gameService.startGame(room.id);
        if (!success) {
          setError('Failed to start game');
        }
        // Game start flow will be handled by WebSocket events
        // Component will receive game_prepare, countdown, and game_begin messages
      }
    } catch (error) {
      console.error('Failed to start game:', error);
      setError('Failed to start game');
    } finally {
      setIsStarting(false);
    }
  };

  const handleLeaveRoom = async () => {
    if (!room || !user || isLeaving) return;
    
    setIsLeaving(true);
    try {
      if (onLeaveRoom) {
        onLeaveRoom();
      } else {
        // Use game service which handles both API and WebSocket
        await gameService.leaveRoom(room.id);
        navigate('/lobby');
      }
    } catch (error) {
      console.error('Failed to leave room:', error);
      setError('Failed to leave room');
    } finally {
      setIsLeaving(false);
    }
  };

  const renderPlayerSeat = (seatIndex: number) => {
    const player = room?.players[seatIndex] || null;
    const isEmpty = !player;
    const isCurrentUser = player?.id === user?.id;
    const isOwner = player?.id === room?.owner;

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

  const getPlayerCount = () => {
    return room?.players.filter(p => p !== null).length || 0;
  };

  const isRoomOwner = () => {
    return user && room && room.owner === user.id;
  };

  const canStartGame = () => {
    return isRoomOwner() && getPlayerCount() === 4 && room?.status === RoomStatusEnum.READY;
  };

  if (!room) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-600">加载房间信息...</p>
        </div>
      </div>
    );
  }

  // Show prepare screen with countdown
  if (showPrepare && countdown !== null) {
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
  }

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
            <div className={`text-sm font-medium ${
              room.status === RoomStatusEnum.WAITING ? 'text-yellow-600' :
              room.status === RoomStatusEnum.READY ? 'text-green-600' :
              room.status === RoomStatusEnum.PLAYING ? 'text-blue-600' :
              'text-gray-600'
            }`}>
              状态: {
                room.status === RoomStatusEnum.WAITING ? '等待中' :
                room.status === RoomStatusEnum.READY ? '准备就绪' :
                room.status === RoomStatusEnum.PLAYING ? '游戏中' :
                '已关闭'
              }
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

export default RoomWaiting;