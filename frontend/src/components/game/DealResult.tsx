import React, { useState, useEffect, useRef } from 'react';
import type { Player } from '../../types';
import type { DealEndedPayload } from '../../types/generated/event';
import { VictoryType } from '../../types/proto';

interface DealResultProps {
  dealResult: DealEndedPayload;
  players: Player[];
  teamLevels: [number, number];
  onContinue: () => void;
  onExit: () => void;
  isMatchFinished: boolean;
  currentPlayerSeat: number;
}

const DealResult: React.FC<DealResultProps> = ({
  dealResult,
  players,
  teamLevels,
  onContinue,
  onExit,
  isMatchFinished,
  currentPlayerSeat
}) => {
  const [countdown, setCountdown] = useState<number | null>(null);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    const deadlineMs = dealResult.nextDealDeadlineMs;
    if (!deadlineMs || deadlineMs === 0) {
      setCountdown(null);
      return;
    }

    const updateCountdown = () => {
      const remaining = Math.max(0, Math.ceil((deadlineMs - Date.now()) / 1000));
      setCountdown(remaining);

      if (remaining === 0) {
        if (timerRef.current) {
          clearInterval(timerRef.current);
          timerRef.current = null;
        }
        onContinue();
      }
    };

    updateCountdown();
    timerRef.current = setInterval(updateCountdown, 1000);

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    };
  }, [dealResult.nextDealDeadlineMs, onContinue]);

  // Helper function to get team for player
  const getTeamForPlayer = (playerSeat: number): number => {
    return playerSeat % 2; // Team 0: seats 0,2; Team 1: seats 1,3
  };

  // Helper function to get player by seat
  const getPlayerBySeat = (seat: number): Player | undefined => {
    return players.find(p => p.seat === seat);
  };

  // Group rankings by team
  const getTeamRankings = () => {
    const teamRankings: { [team: number]: Array<{ rank: number; player: Player }> } = {
      0: [],
      1: []
    };

    dealResult.rankings.forEach((playerSeat, index) => {
      const team = getTeamForPlayer(playerSeat);
      const player = getPlayerBySeat(playerSeat);
      if (player) {
        teamRankings[team].push({
          rank: index + 1, // Convert to 1-based ranking
          player
        });
      }
    });

    return teamRankings;
  };

  // Get victory type display text
  const getVictoryTypeText = (victoryType: VictoryType): string => {
    switch (victoryType) {
      case VictoryType.VICTORY_TYPE_DOUBLE_DOWN:
        return '双下';
      case VictoryType.VICTORY_TYPE_SINGLE_LAST:
        return '单贡';
      case VictoryType.VICTORY_TYPE_PARTNER_LAST:
        return '对贡';
      default:
        return '胜利';
    }
  };

  // Get level display text
  const getLevelText = (level: number): string => {
    if (level <= 10) return level.toString();
    switch (level) {
      case 11: return 'J';
      case 12: return 'Q';
      case 13: return 'K';
      case 14: return 'A';
      default: return level.toString();
    }
  };

  // Format duration
  const formatDuration = (durationMs: number): string => {
    const minutes = Math.floor(durationMs / 60000);
    const seconds = Math.floor((durationMs % 60000) / 1000);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  const teamRankings = getTeamRankings();
  const winningTeam = dealResult.winningTeam;
  const losingTeam = 1 - winningTeam;
  const myTeam = currentPlayerSeat % 2;
  const isWinner = myTeam === winningTeam;

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-gradient-to-b from-[#EAF4EF] via-[#DDEEE5] to-[#D2E8DD] rounded-2xl p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto shadow-2xl border border-white/30">
        <div className="text-center mb-6">
          {isWinner ? (
            <div className="relative py-4">
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="w-32 h-32 bg-gradient-radial from-yellow-300/40 via-amber-400/20 to-transparent rounded-full blur-xl animate-pulse"></div>
              </div>
              <div className="flex flex-col items-center">
                <h2 
                  className="relative text-5xl font-black tracking-[0.3em] bg-gradient-to-b from-yellow-300 via-amber-400 to-orange-500 bg-clip-text text-transparent animate-bounce pl-[0.3em]"
                  style={{ 
                    textShadow: '0 0 30px rgba(251, 191, 36, 0.8), 0 0 60px rgba(251, 191, 36, 0.4)',
                    WebkitTextStroke: '1px rgba(180, 120, 0, 0.3)'
                  }}
                >
                  胜 利
                </h2>
                <div className="flex justify-center gap-2 mt-2">
                  <span className="text-2xl animate-bounce" style={{ animationDelay: '0.1s' }}>⭐</span>
                  <span className="text-3xl animate-bounce" style={{ animationDelay: '0.2s' }}>🏆</span>
                  <span className="text-2xl animate-bounce" style={{ animationDelay: '0.3s' }}>⭐</span>
                </div>
              </div>
            </div>
          ) : (
            <div className="relative py-4">
              <h2 
                className="text-5xl font-black tracking-[0.3em] bg-gradient-to-b from-slate-400 via-slate-500 to-slate-600 bg-clip-text text-transparent"
                style={{ textShadow: '0 2px 10px rgba(0,0,0,0.2)' }}
              >
                失 败
              </h2>
              <div className="text-gray-400 mt-2 text-sm">再接再厉！</div>
            </div>
          )}
        </div>

        {/* Team Rankings */}
        <div className="grid grid-cols-2 gap-5 mb-6">
          {/* Winning Team */}
          <div className="bg-gradient-to-b from-[#525E6B] to-[#3E4854] rounded-xl p-4 shadow-xl border-2 border-amber-400/60 relative overflow-hidden">
            <div className="absolute top-0 right-0 w-16 h-16 bg-gradient-to-br from-amber-400/20 to-transparent rounded-bl-3xl"></div>
            <h3 className="text-lg font-bold text-amber-400 mb-4 text-center tracking-wide flex items-center justify-center gap-2">
              <span className="text-2xl">🏆</span>
              队伍{winningTeam + 1} (胜方)
            </h3>
            <div className="space-y-2.5">
              {teamRankings[winningTeam].map(({ rank, player }) => (
                <div key={player.seat} className="flex justify-between items-center bg-white/5 rounded-lg px-3 py-1.5">
                  <span className="font-medium text-white tracking-wide">{player.username}</span>
                  <span className="text-xs bg-amber-400 text-gray-900 px-2.5 py-1 rounded-md font-bold shadow-sm">
                    第{rank}名
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-3 border-t border-white/10">
              <div className="flex justify-between text-sm text-gray-300 mb-1">
                <span>当前等级</span>
                <span className="font-bold text-white text-base">
                  {getLevelText(teamLevels[winningTeam])}
                </span>
              </div>
              <div className="flex justify-between text-sm text-gray-300">
                <span>升级</span>
                <span className="font-extrabold text-emerald-400 text-base">
                  +{dealResult.levelChange[winningTeam]}级
                </span>
              </div>
            </div>
          </div>

          {/* Losing Team */}
          <div className="bg-gradient-to-b from-[#9E3737] to-[#7B2A2A] rounded-xl p-4 shadow-xl border border-white/10 relative overflow-hidden">
            <div className="absolute top-0 right-0 w-16 h-16 bg-gradient-to-br from-black/10 to-transparent rounded-bl-3xl"></div>
            <h3 className="text-lg font-bold text-gray-100 mb-4 text-center tracking-wide">
              队伍{losingTeam + 1} (负方)
            </h3>
            <div className="space-y-2.5">
              {teamRankings[losingTeam].map(({ rank, player }) => (
                <div key={player.seat} className="flex justify-between items-center bg-black/10 rounded-lg px-3 py-1.5">
                  <span className="font-medium text-gray-100 tracking-wide">{player.username}</span>
                  <span className="text-xs bg-gray-200/20 text-gray-100 border border-gray-400/30 px-2.5 py-1 rounded-md font-medium">
                    第{rank}名
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-3 border-t border-white/10">
              <div className="flex justify-between text-sm text-gray-300 mb-1">
                <span>当前等级</span>
                <span className="font-bold text-gray-100 text-base">
                  {getLevelText(teamLevels[losingTeam])}
                </span>
              </div>
              <div className="flex justify-between text-sm text-gray-300">
                <span>升级</span>
                <span className="font-bold text-gray-400 text-base">
                  +{dealResult.levelChange[losingTeam]}级
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Deal Statistics */}
        <div className="bg-white/80 backdrop-blur-md rounded-xl p-5 mb-5 border border-white/50 shadow-sm">
          <div className="grid grid-cols-3 gap-4">
            <div className="flex flex-col items-center bg-emerald-50/50 rounded-xl p-3 border border-emerald-100/50">
              <span className="text-emerald-600/70 text-xs font-medium mb-1">游戏时长</span>
              <span className="font-bold text-emerald-900 text-lg">{formatDuration(dealResult.durationMs)}</span>
            </div>
            <div className="flex flex-col items-center bg-blue-50/50 rounded-xl p-3 border border-blue-100/50">
              <span className="text-blue-600/70 text-xs font-medium mb-1">总轮次</span>
              <span className="font-bold text-blue-900 text-lg">{dealResult.trickCount}</span>
            </div>
            <div className="flex flex-col items-center bg-amber-50/50 rounded-xl p-3 border border-amber-100/50">
              <span className="text-amber-600/70 text-xs font-medium mb-1">胜利类型</span>
              <span className="font-bold text-amber-700 text-lg">{getVictoryTypeText(dealResult.victoryType)}</span>
            </div>
          </div>
        </div>

        {/* Player Statistics */}
        <div className="bg-white/80 backdrop-blur-md rounded-xl p-5 mb-6 border border-white/50 shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-gray-50/80 text-gray-500 border-b border-gray-100">
                  <th className="text-left py-3 px-4 rounded-tl-lg font-medium">玩家</th>
                  <th className="text-center py-3 px-2 font-medium">排名</th>
                  <th className="text-center py-3 px-2 font-medium">出牌</th>
                  <th className="text-center py-3 px-2 font-medium">胜轮</th>
                  <th className="text-center py-3 px-4 rounded-tr-lg font-medium">过牌</th>
                </tr>
              </thead>
              <tbody>
                {dealResult.playerStats.map((stats, index) => {
                  const player = getPlayerBySeat(stats.playerSeat);
                  const isWinner = stats.finishRank <= 2;
                  return (
                    <tr 
                      key={stats.playerSeat} 
                      className={`border-b border-gray-50 hover:bg-emerald-50/30 transition-colors ${index % 2 === 0 ? 'bg-white/40' : 'bg-transparent'}`}
                    >
                      <td className="py-3 px-4 font-medium text-gray-700">
                        {player?.username || `玩家${stats.playerSeat + 1}`}
                      </td>
                      <td className="text-center py-3 px-2">
                        <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-bold shadow-sm ${
                          isWinner 
                            ? 'bg-amber-100 text-amber-700 ring-1 ring-amber-400/20' 
                            : 'bg-gray-100 text-gray-500'
                        }`}>
                          第{stats.finishRank}名
                        </span>
                      </td>
                      <td className="text-center py-3 px-2 text-gray-600 font-medium">{stats.cardsPlayed}</td>
                      <td className="text-center py-3 px-2 text-gray-600 font-medium">{stats.tricksWon}</td>
                      <td className="text-center py-3 px-4 text-gray-400">{stats.passCount}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-col items-center space-y-5 pt-2">
          {countdown !== null && countdown > 0 && (
            <div className="bg-white/90 backdrop-blur-md rounded-full px-6 py-2 border border-emerald-100 shadow-sm flex items-center gap-2">
              <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
              <span className="text-gray-500 text-sm font-medium">下一局开始倒计时</span>
              <span className="text-emerald-600 font-bold text-xl tabular-nums">{countdown}</span>
              <span className="text-gray-400 text-xs self-end mb-1">s</span>
            </div>
          )}
          <button
            onClick={onExit}
            className="px-10 py-3 bg-gradient-to-b from-slate-600 to-slate-700 text-white rounded-xl font-bold shadow-lg shadow-slate-300/50 hover:from-slate-500 hover:to-slate-600 hover:scale-105 active:scale-95 transition-all duration-200 border-t border-white/20"
          >
            返回大厅
          </button>
        </div>
      </div>
    </div>
  );
};

export default DealResult;