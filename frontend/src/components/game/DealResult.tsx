import React, { useState, useEffect, useRef } from 'react';
import type { Player } from '../../types';
import type { DealEndedPayload } from '../../types/generated/event';
import { VictoryType } from '../../types/proto';
import { Card, Badge, Button } from '@/components/ui';
import { Trophy, Clock, Repeat, Swords, ArrowUp, Home } from 'lucide-react';
import { cn } from '@/lib/utils';

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
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  isMatchFinished: _isMatchFinished,
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

  const getTeamForPlayer = (playerSeat: number): number => {
    return playerSeat % 2;
  };

  const getPlayerBySeat = (seat: number): Player | undefined => {
    return players.find(p => p.seat === seat);
  };

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
          rank: index + 1,
          player
        });
      }
    });

    return teamRankings;
  };

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
    <div className={cn(
      "fixed inset-0 flex items-center justify-center z-50",
      "bg-black/70 backdrop-blur-md",
    )}>
      <Card
        variant="elevated"
        interactive={false}
        className={cn(
          "p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto",
          "bg-gradient-to-b from-[hsl(158,25%,15%)] to-[hsl(158,30%,10%)]",
          "border border-[hsl(158,55%,30%)]/30",
          "shadow-[0_8px_32px_rgba(0,0,0,0.5),0_0_60px_hsla(158,55%,30%,0.15)]",
        )}
      >
        {/* 标题区域 */}
        <div className="text-center mb-6">
          {isWinner ? (
            <div className="relative py-6">
              {/* 背景光晕 */}
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="w-40 h-40 bg-gradient-radial from-state-active/50 via-state-active/20 to-transparent rounded-full blur-2xl animate-pulse"></div>
              </div>
              <div className="relative flex flex-col items-center">
                <Trophy className="w-16 h-16 text-state-active mb-3 drop-shadow-[0_0_20px_hsla(42,95%,52%,0.6)]" />
                <h2
                  className={cn(
                    "text-5xl font-black tracking-[0.3em] font-display pl-[0.3em]",
                    "bg-gradient-to-b from-state-active via-[hsl(42,95%,60%)] to-[hsl(42,95%,45%)]",
                    "bg-clip-text text-transparent",
                  )}
                  style={{
                    textShadow: '0 0 40px hsla(42, 95%, 52%, 0.8), 0 0 80px hsla(42, 95%, 52%, 0.4)',
                  }}
                >
                  胜 利
                </h2>
                <div className="flex items-center gap-2 mt-3 text-state-active/80 text-sm">
                  <span>♠</span>
                  <span>恭喜赢得本局</span>
                  <span>♥</span>
                </div>
              </div>
            </div>
          ) : (
            <div className="relative py-6">
              <div className="relative flex flex-col items-center">
                <Swords className="w-12 h-12 text-white/40 mb-3" />
                <h2
                  className={cn(
                    "text-5xl font-black tracking-[0.3em] font-display pl-[0.3em]",
                    "bg-gradient-to-b from-white/70 via-white/50 to-white/30",
                    "bg-clip-text text-transparent",
                  )}
                  style={{ textShadow: '0 2px 10px rgba(0, 0, 0, 0.3)' }}
                >
                  失 败
                </h2>
                <div className="text-white/50 mt-2 text-sm">再接再厉！</div>
              </div>
            </div>
          )}
        </div>

        {/* 队伍结果卡片 */}
        <div className="grid grid-cols-2 gap-5 mb-6">
          {/* 胜方队伍 */}
          <div className={cn(
            "p-4 rounded-xl relative overflow-hidden",
            "bg-gradient-to-br from-state-active/20 via-state-active/10 to-transparent",
            "border-2 border-state-active/40",
            "shadow-[0_4px_20px_hsla(42,95%,52%,0.2)]",
          )}>
            <div className="absolute top-0 right-0 w-20 h-20 bg-gradient-to-br from-state-active/20 to-transparent rounded-bl-full"></div>
            <h3 className="text-lg font-bold text-state-active mb-4 text-center tracking-wide flex items-center justify-center gap-2">
              <Trophy className="w-5 h-5" />
              队伍{winningTeam + 1} (胜方)
            </h3>
            <div className="space-y-2.5">
              {teamRankings[winningTeam].map(({ rank, player }) => (
                <div key={player.seat} className={cn(
                  "flex justify-between items-center px-3 py-2 rounded-lg",
                  "bg-black/20 border border-white/10",
                )}>
                  <span className="font-medium text-white tracking-wide">{player.username}</span>
                  <Badge variant="landlord" size="sm">
                    第{rank}名
                  </Badge>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-3 border-t border-white/10">
              <div className="flex justify-between text-sm text-white/70 mb-1">
                <span>当前等级</span>
                <span className="font-bold text-white text-base font-display">
                  {getLevelText(teamLevels[winningTeam])}
                </span>
              </div>
              <div className="flex justify-between text-sm text-white/70">
                <span className="flex items-center gap-1">
                  <ArrowUp className="w-3 h-3 text-[hsl(158,55%,50%)]" />
                  升级
                </span>
                <span className="font-extrabold text-[hsl(158,55%,50%)] text-base">
                  +{dealResult.levelChange[winningTeam]}级
                </span>
              </div>
            </div>
          </div>

          {/* 负方队伍 */}
          <div className={cn(
            "p-4 rounded-xl relative overflow-hidden",
            "bg-black/30",
            "border border-white/10",
          )}>
            <div className="absolute top-0 right-0 w-16 h-16 bg-gradient-to-br from-white/5 to-transparent rounded-bl-3xl"></div>
            <h3 className="text-lg font-bold text-white/70 mb-4 text-center tracking-wide">
              队伍{losingTeam + 1} (负方)
            </h3>
            <div className="space-y-2.5">
              {teamRankings[losingTeam].map(({ rank, player }) => (
                <div key={player.seat} className={cn(
                  "flex justify-between items-center px-3 py-2 rounded-lg",
                  "bg-black/20 border border-white/5",
                )}>
                  <span className="font-medium text-white/80 tracking-wide">{player.username}</span>
                  <Badge variant="neutral" size="sm">
                    第{rank}名
                  </Badge>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-3 border-t border-white/10">
              <div className="flex justify-between text-sm text-white/50 mb-1">
                <span>当前等级</span>
                <span className="font-bold text-white/70 text-base font-display">
                  {getLevelText(teamLevels[losingTeam])}
                </span>
              </div>
              <div className="flex justify-between text-sm text-white/50">
                <span>升级</span>
                <span className="font-bold text-white/50 text-base">
                  +{dealResult.levelChange[losingTeam]}级
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* 统计信息 */}
        <div className={cn(
          "p-4 mb-5 rounded-xl",
          "bg-black/30 border border-white/10",
        )}>
          <div className="grid grid-cols-3 gap-4">
            <div className={cn(
              "flex flex-col items-center p-3 rounded-lg",
              "bg-black/20 border border-white/5",
            )}>
              <Clock className="w-5 h-5 text-[hsl(158,55%,50%)] mb-1" />
              <span className="text-white/50 text-xs font-medium mb-1">游戏时长</span>
              <span className="font-bold text-white text-lg font-display">{formatDuration(dealResult.durationMs)}</span>
            </div>
            <div className={cn(
              "flex flex-col items-center p-3 rounded-lg",
              "bg-black/20 border border-white/5",
            )}>
              <Repeat className="w-5 h-5 text-[hsl(158,55%,50%)] mb-1" />
              <span className="text-white/50 text-xs font-medium mb-1">总轮次</span>
              <span className="font-bold text-white text-lg font-display">{dealResult.trickCount}</span>
            </div>
            <div className={cn(
              "flex flex-col items-center p-3 rounded-lg",
              "bg-black/20 border border-white/5",
            )}>
              <Swords className="w-5 h-5 text-state-active mb-1" />
              <span className="text-white/50 text-xs font-medium mb-1">胜利类型</span>
              <span className="font-bold text-state-active text-lg font-display">{getVictoryTypeText(dealResult.victoryType)}</span>
            </div>
          </div>
        </div>

        {/* 玩家统计表格 */}
        <div className={cn(
          "p-4 mb-6 rounded-xl overflow-hidden",
          "bg-black/30 border border-white/10",
        )}>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className={cn(
                  "bg-gradient-to-r from-[hsl(158,55%,25%)] to-[hsl(158,55%,20%)]",
                  "border-b border-[hsl(158,55%,30%)]/30",
                )}>
                  <th className="text-left py-3 px-4 font-medium text-white/80">玩家</th>
                  <th className="text-center py-3 px-2 font-medium text-white/80">排名</th>
                  <th className="text-center py-3 px-2 font-medium text-white/80">出牌</th>
                  <th className="text-center py-3 px-2 font-medium text-white/80">胜轮</th>
                  <th className="text-center py-3 px-4 font-medium text-white/80">过牌</th>
                </tr>
              </thead>
              <tbody>
                {dealResult.playerStats.map((stats, index) => {
                  const player = getPlayerBySeat(stats.playerSeat);
                  const isTopTwo = stats.finishRank <= 2;
                  return (
                    <tr
                      key={stats.playerSeat}
                      className={cn(
                        "border-b border-white/5",
                        "hover:bg-white/5 transition-colors",
                        index % 2 === 0 ? 'bg-black/20' : 'bg-transparent',
                      )}
                    >
                      <td className="py-3 px-4 font-medium text-white">
                        {player?.username || `玩家${stats.playerSeat + 1}`}
                      </td>
                      <td className="text-center py-3 px-2">
                        <Badge variant={isTopTwo ? "landlord" : "neutral"} size="sm">
                          第{stats.finishRank}名
                        </Badge>
                      </td>
                      <td className="text-center py-3 px-2 text-white font-medium">{stats.cardsPlayed}</td>
                      <td className="text-center py-3 px-2 text-white font-medium">{stats.tricksWon}</td>
                      <td className="text-center py-3 px-4 text-white/60">{stats.passCount}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>

        {/* 底部按钮区域 */}
        <div className="flex flex-col items-center space-y-5 pt-2">
          {countdown !== null && countdown > 0 && (
            <div className={cn(
              "flex items-center gap-3 px-6 py-3 rounded-full",
              "bg-black/40 backdrop-blur-sm",
              "border border-state-active/30",
              "shadow-[0_0_20px_hsla(42,95%,52%,0.2)]",
            )}>
              <div className="w-2 h-2 rounded-full bg-state-active animate-pulse"></div>
              <span className="text-white/70 text-sm font-medium">下一局开始倒计时</span>
              <span className="text-state-active font-bold text-2xl tabular-nums font-display">{countdown}</span>
              <span className="text-white/50 text-xs self-end mb-1">s</span>
            </div>
          )}
          <Button
            intent="neutral"
            size="lg"
            onClick={onExit}
            className={cn(
              "bg-white/10 hover:bg-white/20",
              "text-white/80 hover:text-white",
              "border border-white/20 hover:border-white/30",
            )}
          >
            <Home className="w-5 h-5 mr-2" />
            返回大厅
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default DealResult;
