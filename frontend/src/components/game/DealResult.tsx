import React, { useState, useEffect, useRef } from 'react';
import type { Player } from '../../types';
import type { DealEndedPayload } from '../../types/generated/event';
import { VictoryType } from '../../types/proto';
import { Card, Badge, Button } from '@/components/ui-next';

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
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
      <Card 
        variant="elevated" 
        interactive={false}
        className="p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto shadow-elevation-3 border border-stroke-emphasis"
      >
        <div className="text-center mb-6">
          {isWinner ? (
            <div className="relative py-4">
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="w-32 h-32 bg-gradient-radial from-state-active/40 via-state-active/20 to-transparent rounded-full blur-xl animate-pulse"></div>
              </div>
              <div className="flex flex-col items-center">
                <h2 
                  className="relative text-5xl font-black tracking-[0.3em] bg-gradient-to-b from-state-active via-action-secondary to-action-secondary bg-clip-text text-transparent animate-bounce pl-[0.3em]"
                  style={{ 
                    textShadow: '0 0 30px hsla(var(--primitive-accent-500), 0.8), 0 0 60px hsla(var(--primitive-accent-500), 0.4)',
                    WebkitTextStroke: '1px hsla(var(--primitive-accent-700), 0.3)'
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
                className="text-5xl font-black tracking-[0.3em] bg-gradient-to-b from-fg-secondary via-action-neutral to-action-neutral bg-clip-text text-transparent"
                style={{ textShadow: '0 2px 10px hsla(var(--primitive-neutral-900), 0.2)' }}
              >
                失 败
              </h2>
              <div className="text-fg-secondary mt-2 text-sm">再接再厉！</div>
            </div>
          )}
        </div>

        <div className="grid grid-cols-2 gap-5 mb-6">
          <Card 
            variant="emphasis" 
            interactive={false}
            className="!bg-team-us p-4 shadow-elevation-3 border-2 border-state-active/60 relative overflow-hidden"
          >
            <div className="absolute top-0 right-0 w-16 h-16 bg-gradient-to-br from-state-active/20 to-transparent rounded-bl-3xl"></div>
            <h3 className="text-lg font-bold text-state-active mb-4 text-center tracking-wide flex items-center justify-center gap-2">
              <span className="text-2xl">🏆</span>
              队伍{winningTeam + 1} (胜方)
            </h3>
            <div className="space-y-2.5">
              {teamRankings[winningTeam].map(({ rank, player }) => (
                <div key={player.seat} className="flex justify-between items-center bg-white/5 rounded-sm px-3 py-1.5">
                  <span className="font-medium text-fg-inverse tracking-wide">{player.username}</span>
                  <Badge variant="landlord" size="sm">
                    第{rank}名
                  </Badge>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-3 border-t border-white/10">
              <div className="flex justify-between text-sm text-fg-inverse/70 mb-1">
                <span>当前等级</span>
                <span className="font-bold text-fg-inverse text-base">
                  {getLevelText(teamLevels[winningTeam])}
                </span>
              </div>
              <div className="flex justify-between text-sm text-fg-inverse/70">
                <span>升级</span>
                <span className="font-extrabold text-action-primary text-base">
                  +{dealResult.levelChange[winningTeam]}级
                </span>
              </div>
            </div>
          </Card>

          <Card 
            variant="base" 
            interactive={false}
            className="!bg-team-them p-4 shadow-elevation-3 border border-white/10 relative overflow-hidden"
          >
            <div className="absolute top-0 right-0 w-16 h-16 bg-gradient-to-br from-black/10 to-transparent rounded-bl-3xl"></div>
            <h3 className="text-lg font-bold text-fg-inverse mb-4 text-center tracking-wide">
              队伍{losingTeam + 1} (负方)
            </h3>
            <div className="space-y-2.5">
              {teamRankings[losingTeam].map(({ rank, player }) => (
                <div key={player.seat} className="flex justify-between items-center bg-black/10 rounded-sm px-3 py-1.5">
                  <span className="font-medium text-fg-inverse tracking-wide">{player.username}</span>
                  <Badge variant="neutral" size="sm">
                    第{rank}名
                  </Badge>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-3 border-t border-white/10">
              <div className="flex justify-between text-sm text-fg-inverse/70 mb-1">
                <span>当前等级</span>
                <span className="font-bold text-fg-inverse text-base">
                  {getLevelText(teamLevels[losingTeam])}
                </span>
              </div>
              <div className="flex justify-between text-sm text-fg-inverse/70">
                <span>升级</span>
                <span className="font-bold text-fg-inverse/50 text-base">
                  +{dealResult.levelChange[losingTeam]}级
                </span>
              </div>
            </div>
          </Card>
        </div>

        <Card variant="elevated" interactive={false} className="p-5 mb-5 border border-stroke shadow-elevation-1">
          <div className="grid grid-cols-3 gap-4">
            <div className="flex flex-col items-center bg-surface-emphasis rounded-md p-3 border border-stroke-emphasis">
              <span className="text-action-primary/70 text-xs font-medium mb-1">游戏时长</span>
              <span className="font-bold text-fg-primary text-lg">{formatDuration(dealResult.durationMs)}</span>
            </div>
            <div className="flex flex-col items-center bg-surface-emphasis rounded-md p-3 border border-stroke-emphasis">
              <span className="text-action-primary/70 text-xs font-medium mb-1">总轮次</span>
              <span className="font-bold text-fg-primary text-lg">{dealResult.trickCount}</span>
            </div>
            <div className="flex flex-col items-center bg-surface-emphasis rounded-md p-3 border border-stroke-emphasis">
              <span className="text-action-secondary/70 text-xs font-medium mb-1">胜利类型</span>
              <span className="font-bold text-action-secondary text-lg">{getVictoryTypeText(dealResult.victoryType)}</span>
            </div>
          </div>
        </Card>

        <Card variant="elevated" interactive={false} className="p-5 mb-6 border border-stroke shadow-elevation-1">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-surface-elevated text-fg-secondary border-b border-stroke">
                  <th className="text-left py-3 px-4 rounded-tl-sm font-medium">玩家</th>
                  <th className="text-center py-3 px-2 font-medium">排名</th>
                  <th className="text-center py-3 px-2 font-medium">出牌</th>
                  <th className="text-center py-3 px-2 font-medium">胜轮</th>
                  <th className="text-center py-3 px-4 rounded-tr-sm font-medium">过牌</th>
                </tr>
              </thead>
              <tbody>
                {dealResult.playerStats.map((stats, index) => {
                  const player = getPlayerBySeat(stats.playerSeat);
                  const isTopTwo = stats.finishRank <= 2;
                  return (
                    <tr 
                      key={stats.playerSeat} 
                      className={`border-b border-stroke/50 hover:bg-surface-emphasis/30 transition-colors duration-fast ${index % 2 === 0 ? 'bg-surface-base/40' : 'bg-transparent'}`}
                    >
                      <td className="py-3 px-4 font-medium text-fg-primary">
                        {player?.username || `玩家${stats.playerSeat + 1}`}
                      </td>
                      <td className="text-center py-3 px-2">
                        <Badge variant={isTopTwo ? "landlord" : "neutral"} size="sm">
                          第{stats.finishRank}名
                        </Badge>
                      </td>
                      <td className="text-center py-3 px-2 text-fg-primary font-medium">{stats.cardsPlayed}</td>
                      <td className="text-center py-3 px-2 text-fg-primary font-medium">{stats.tricksWon}</td>
                      <td className="text-center py-3 px-4 text-fg-secondary">{stats.passCount}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </Card>

        <div className="flex flex-col items-center space-y-5 pt-2">
          {countdown !== null && countdown > 0 && (
            <div className="bg-surface-base backdrop-blur-md rounded-full px-6 py-2 border border-stroke shadow-elevation-1 flex items-center gap-2">
              <div className="w-2 h-2 rounded-full bg-action-primary animate-pulse"></div>
              <span className="text-fg-secondary text-sm font-medium">下一局开始倒计时</span>
              <span className="text-action-primary font-bold text-xl tabular-nums">{countdown}</span>
              <span className="text-fg-secondary text-xs self-end mb-1">s</span>
            </div>
          )}
          <Button intent="neutral" size="lg" onClick={onExit}>
            返回大厅
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default DealResult;
