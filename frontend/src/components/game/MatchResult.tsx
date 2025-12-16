import React from 'react';
import type { Player } from '../../types';
import type { MatchEndedPayload } from '../../types/generated/event';
import { Trophy, Clock, Swords, Medal, RefreshCw, Home } from 'lucide-react';
import { Button } from '@/components/ui';
import { cn } from '@/lib/utils';

interface MatchResultProps {
  matchResult: MatchEndedPayload;
  players: Player[];
  onReturnToLobby: () => void;
  onPlayAgain: () => void;
}

const MatchResult: React.FC<MatchResultProps> = ({
  matchResult,
  players,
  onReturnToLobby,
  onPlayAgain
}) => {
  // Helper function to get team for player
  const getTeamForPlayer = (playerSeat: number): number => {
    return playerSeat % 2; // Team 0: seats 0,2; Team 1: seats 1,3
  };

  // Group players by team
  const getTeamPlayers = () => {
    const teamPlayers: { [team: number]: Player[] } = {
      0: [],
      1: []
    };

    players.forEach(player => {
      const team = getTeamForPlayer(player.seat);
      teamPlayers[team].push(player);
    });

    return teamPlayers;
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
    const hours = Math.floor(durationMs / 3600000);
    const minutes = Math.floor((durationMs % 3600000) / 60000);
    const seconds = Math.floor((durationMs % 60000) / 1000);
    
    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  const teamPlayers = getTeamPlayers();
  const winningTeam = matchResult.winner;
  const losingTeam = 1 - winningTeam;

  return (
    <div className={cn(
      "fixed inset-0 flex items-center justify-center z-50",
      "bg-black/70 backdrop-blur-md",
    )}>
      <div className={cn(
        "max-w-3xl w-full mx-4 max-h-[90vh] overflow-y-auto",
        "p-8 rounded-2xl",
        "bg-gradient-to-b from-[hsl(158,25%,15%)] to-[hsl(158,30%,10%)]",
        "border border-[hsl(158,55%,30%)]/30",
        "shadow-[0_8px_32px_rgba(0,0,0,0.5),0_0_80px_hsla(42,95%,52%,0.15)]",
      )}>
        {/* 标题区域 */}
        <div className="text-center mb-8">
          <div className="relative py-6">
            {/* 背景光晕 */}
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="w-48 h-48 bg-gradient-radial from-state-active/50 via-state-active/20 to-transparent rounded-full blur-3xl animate-pulse"></div>
            </div>
            <div className="relative flex flex-col items-center">
              <Trophy className="w-20 h-20 text-state-active mb-4 drop-shadow-[0_0_30px_hsla(42,95%,52%,0.7)]" />
              <h1
                className={cn(
                  "text-4xl font-black tracking-[0.2em] font-display",
                  "bg-gradient-to-b from-state-active via-[hsl(42,95%,60%)] to-[hsl(42,95%,45%)]",
                  "bg-clip-text text-transparent",
                )}
                style={{
                  textShadow: '0 0 40px hsla(42, 95%, 52%, 0.8)',
                }}
              >
                比赛结束
              </h1>
              <div className="flex items-center gap-3 mt-4 text-white/80 text-lg">
                <span className="text-state-active">♠</span>
                <span>队伍{winningTeam + 1}获得最终胜利！</span>
                <span className="text-state-active">♥</span>
              </div>
              <div className="text-white/50 mt-2 text-sm">恭喜率先达到A级！</div>
            </div>
          </div>
        </div>

        {/* 队伍结果卡片 */}
        <div className="grid grid-cols-2 gap-6 mb-8">
          {/* 冠军队伍 */}
          <div className={cn(
            "p-5 rounded-xl relative overflow-hidden",
            "bg-gradient-to-br from-state-active/25 via-state-active/15 to-transparent",
            "border-2 border-state-active/50",
            "shadow-[0_4px_24px_hsla(42,95%,52%,0.25)]",
          )}>
            <div className="absolute top-0 right-0 w-24 h-24 bg-gradient-to-br from-state-active/30 to-transparent rounded-bl-full"></div>
            <div className="text-center mb-4">
              <h2 className="text-2xl font-bold text-state-active mb-1 flex items-center justify-center gap-2">
                <Trophy className="w-6 h-6" />
                队伍{winningTeam + 1}
              </h2>
              <div className="text-sm font-medium text-state-active/80">冠军队伍</div>
            </div>
            <div className="space-y-3">
              {teamPlayers[winningTeam].map(player => (
                <div key={player.seat} className={cn(
                  "flex justify-between items-center px-4 py-3 rounded-lg",
                  "bg-black/30 border border-white/10",
                )}>
                  <span className="font-medium text-white">{player.username}</span>
                  <span className={cn(
                    "text-xs px-2 py-1 rounded-full",
                    "bg-state-active/20 text-state-active border border-state-active/30",
                  )}>
                    座位{player.seat + 1}
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-4 border-t border-white/10">
              <div className="flex justify-between items-center">
                <span className="text-white/70">最终等级</span>
                <span className="text-3xl font-black text-state-active font-display">
                  {getLevelText(matchResult.finalLevels[winningTeam])}
                </span>
              </div>
            </div>
          </div>

          {/* 亚军队伍 */}
          <div className={cn(
            "p-5 rounded-xl relative overflow-hidden",
            "bg-black/30",
            "border border-white/10",
          )}>
            <div className="absolute top-0 right-0 w-20 h-20 bg-gradient-to-br from-white/5 to-transparent rounded-bl-full"></div>
            <div className="text-center mb-4">
              <h2 className="text-2xl font-bold text-white/70 mb-1 flex items-center justify-center gap-2">
                <Medal className="w-5 h-5" />
                队伍{losingTeam + 1}
              </h2>
              <div className="text-sm font-medium text-white/50">亚军队伍</div>
            </div>
            <div className="space-y-3">
              {teamPlayers[losingTeam].map(player => (
                <div key={player.seat} className={cn(
                  "flex justify-between items-center px-4 py-3 rounded-lg",
                  "bg-black/30 border border-white/5",
                )}>
                  <span className="font-medium text-white/80">{player.username}</span>
                  <span className={cn(
                    "text-xs px-2 py-1 rounded-full",
                    "bg-white/5 text-white/50 border border-white/10",
                  )}>
                    座位{player.seat + 1}
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-4 border-t border-white/10">
              <div className="flex justify-between items-center">
                <span className="text-white/50">最终等级</span>
                <span className="text-3xl font-black text-white/70 font-display">
                  {getLevelText(matchResult.finalLevels[losingTeam])}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* 比赛统计 */}
        <div className={cn(
          "p-5 mb-8 rounded-xl",
          "bg-black/30 border border-white/10",
        )}>
          <h3 className="text-lg font-bold text-white/80 mb-4 text-center">比赛统计</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className={cn(
              "flex flex-col items-center p-4 rounded-lg",
              "bg-black/20 border border-white/5",
            )}>
              <Swords className="w-5 h-5 text-[hsl(158,55%,50%)] mb-2" />
              <div className="text-2xl font-bold text-white font-display mb-1">
                {matchResult.totalDeals}
              </div>
              <div className="text-xs text-white/50">总局数</div>
            </div>
            <div className={cn(
              "flex flex-col items-center p-4 rounded-lg",
              "bg-black/20 border border-white/5",
            )}>
              <Clock className="w-5 h-5 text-[hsl(158,55%,50%)] mb-2" />
              <div className="text-2xl font-bold text-white font-display mb-1">
                {formatDuration(matchResult.durationMs)}
              </div>
              <div className="text-xs text-white/50">总时长</div>
            </div>
            <div className={cn(
              "flex flex-col items-center p-4 rounded-lg",
              "bg-black/20 border border-white/5",
            )}>
              <Medal className="w-5 h-5 text-state-active mb-2" />
              <div className="text-2xl font-bold text-white font-display mb-1">
                {getLevelText(matchResult.finalLevels[0])} vs {getLevelText(matchResult.finalLevels[1])}
              </div>
              <div className="text-xs text-white/50">最终等级</div>
            </div>
            <div className={cn(
              "flex flex-col items-center p-4 rounded-lg",
              "bg-black/20 border border-white/5",
            )}>
              <Trophy className="w-5 h-5 text-state-active mb-2" />
              <div className="text-2xl font-bold text-state-active font-display mb-1">
                队伍{matchResult.winner + 1}
              </div>
              <div className="text-xs text-white/50">获胜队伍</div>
            </div>
          </div>
        </div>

        {/* 感谢信息 */}
        <div className={cn(
          "p-5 mb-8 rounded-xl text-center",
          "bg-gradient-to-r from-[hsl(158,55%,20%)]/30 via-black/30 to-[hsl(42,95%,30%)]/20",
          "border border-[hsl(158,55%,30%)]/20",
        )}>
          <h3 className="text-lg font-bold text-white/90 mb-2">
            感谢参与本次掼蛋对战！
          </h3>
          <p className="text-white/60 mb-4">
            经过 {matchResult.totalDeals} 局激烈的对战，
            队伍{winningTeam + 1} 成功率先达到A级，获得最终胜利！
          </p>
          <div className="flex justify-center items-center gap-6 text-sm text-white/40">
            <span className="flex items-center gap-1">
              <span className="text-[hsl(158,55%,50%)]">♣</span>
              精彩对局
            </span>
            <span className="flex items-center gap-1">
              <span className="text-suit-red">♥</span>
              友谊第一
            </span>
            <span className="flex items-center gap-1">
              <span className="text-state-active">♦</span>
              比赛第二
            </span>
          </div>
        </div>

        {/* 操作按钮 */}
        <div className="flex justify-center gap-4">
          <Button
            onClick={onPlayAgain}
            size="lg"
            className={cn(
              "px-8",
              "bg-gradient-to-r from-[hsl(158,55%,35%)] to-[hsl(158,55%,42%)]",
              "hover:from-[hsl(158,55%,40%)] hover:to-[hsl(158,55%,47%)]",
              "text-white font-bold",
              "border-none",
              "shadow-[0_4px_16px_hsla(158,55%,42%,0.4)]",
              "hover:shadow-[0_6px_20px_hsla(158,55%,42%,0.5)]",
            )}
          >
            <RefreshCw className="w-5 h-5 mr-2" />
            再来一局
          </Button>
          <Button
            onClick={onReturnToLobby}
            size="lg"
            className={cn(
              "px-8",
              "bg-white/10 hover:bg-white/20",
              "text-white/80 hover:text-white",
              "border border-white/20 hover:border-white/30",
            )}
          >
            <Home className="w-5 h-5 mr-2" />
            返回大厅
          </Button>
        </div>
      </div>
    </div>
  );
};

export default MatchResult;