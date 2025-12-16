import React, { useState, useEffect, useRef } from 'react';
import type { Player } from '../../types';
import type { DealEndedPayload } from '../../types/generated/event';
import { Card, Button, Avatar } from '@/components/ui';
import { Trophy, Swords, Home } from 'lucide-react';
import { cn } from '@/lib/utils';
import { getAvatarUrl, getAvatarByUsername } from '../../utils/avatar';

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
  const [, setCountdown] = useState<number | null>(null);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Keep auto-continue functionality but hide countdown display
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

  const getRankLabel = (rank: number): string => {
    switch (rank) {
      case 1: return '头游';
      case 2: return '二游';
      case 3: return '三游';
      case 4: return '末游';
      default: return `第${rank}名`;
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

  const getPlayerAvatarUrl = (player: Player): string => {
    return getAvatarUrl(player.avatar_key) || getAvatarByUsername(player.username);
  };

  const teamRankings = getTeamRankings();
  const winningTeam = dealResult.winningTeam;
  const losingTeam = 1 - winningTeam;
  const myTeam = currentPlayerSeat % 2;
  const isWinner = myTeam === winningTeam;

  // PlayerRankCard component for displaying player avatar and rank
  const PlayerRankCard = ({ player, rank, isWinningTeam }: { player: Player; rank: number; isWinningTeam: boolean }) => (
    <div className="flex items-center gap-2 sm:gap-3">
      {/* 左侧：头像和名字 */}
      <div className="flex flex-col items-center gap-1">
        <Avatar
          src={getPlayerAvatarUrl(player)}
          alt={player.username}
          fallback={player.username.slice(0, 2)}
          size="lg"
          ringState={isWinningTeam ? "teamUs" : "teamThem"}
          className="w-12 h-12 sm:w-14 sm:h-14 md:w-16 md:h-16"
        />
        <span className={cn(
          "text-xs sm:text-sm px-2 py-0.5 rounded-full truncate max-w-[70px] sm:max-w-[80px]",
          "bg-black/40 border",
          isWinningTeam ? "border-state-active/30 text-white/90" : "border-white/10 text-white/60"
        )}>
          {player.username}
        </span>
      </div>
      {/* 右侧：排名标签 - 对齐头像中部 */}
      <span className={cn(
        "-mt-5 sm:-mt-6 text-lg sm:text-xl md:text-2xl font-black tracking-[0.15em] font-display",
        isWinningTeam
          ? "bg-gradient-to-b from-state-active via-[hsl(42,95%,58%)] to-[hsl(42,95%,42%)] bg-clip-text text-transparent"
          : "bg-gradient-to-b from-white/60 via-white/40 to-white/20 bg-clip-text text-transparent"
      )}>
        {getRankLabel(rank)}
      </span>
    </div>
  );

  // Level display component - 失败方仅显示原始等级，与胜方原始等级对齐
  const LevelDisplay = ({ level, change, isWinning }: { level: number; change: number; isWinning: boolean }) => (
    <div className="flex items-baseline gap-0.5 font-display">
      <span className={cn(
        "text-xl sm:text-2xl md:text-3xl font-bold",
        isWinning ? "text-white/80" : "text-white/50"
      )}>
        {getLevelText(level)}
      </span>
      {isWinning ? (
        <span className="text-3xl sm:text-4xl md:text-5xl font-black text-[hsl(158,55%,55%)]">
          +{change}
        </span>
      ) : (
        // 占位符保持对齐，但不显示内容
        <span className="text-3xl sm:text-4xl md:text-5xl font-black invisible">
          +0
        </span>
      )}
    </div>
  );

  return (
    <div className={cn(
      "fixed inset-0 flex items-center justify-center z-50",
      "bg-black/80 backdrop-blur-md",
      "px-3 py-6 sm:p-4 md:p-6",
    )}>
      <Card
        variant="elevated"
        interactive={false}
        className={cn(
          "p-4 sm:p-6 md:p-8 max-w-4xl w-full",
          "bg-gradient-to-br from-[hsl(158,20%,12%)] via-[hsl(158,25%,10%)] to-[hsl(158,30%,8%)]",
          "border border-[hsl(158,55%,25%)]/40",
          "shadow-[0_8px_60px_rgba(0,0,0,0.6),0_0_80px_hsla(158,55%,30%,0.1)]",
        )}
      >
        {/* 主内容区 - 水平布局 */}
        <div className="flex items-stretch gap-4 sm:gap-6 md:gap-8">

          {/* 左侧: 胜利/失败指示区 */}
          <div className={cn(
            "flex flex-col items-center justify-center px-3 sm:px-6 md:px-8 py-4 sm:py-6",
            "min-w-[100px] sm:min-w-[140px] md:min-w-[180px]",
            "border-r border-white/10"
          )}>
            {isWinner ? (
              <div className="relative flex flex-col items-center">
                {/* 背景光晕 */}
                <div className="absolute -inset-4 sm:-inset-8 flex items-center justify-center pointer-events-none">
                  <div className="w-32 sm:w-40 md:w-48 h-32 sm:h-40 md:h-48 bg-gradient-radial from-state-active/40 via-state-active/10 to-transparent rounded-full blur-3xl animate-pulse"></div>
                </div>
                <Trophy className="relative w-12 sm:w-16 md:w-20 h-12 sm:h-16 md:h-20 text-state-active mb-2 sm:mb-4 drop-shadow-[0_0_30px_hsla(42,95%,52%,0.7)]" />
                <h2
                  className={cn(
                    "relative text-3xl sm:text-4xl md:text-5xl font-black tracking-[0.2em] sm:tracking-[0.25em] font-display pl-[0.2em] sm:pl-[0.25em]",
                    "bg-gradient-to-b from-state-active via-[hsl(42,95%,58%)] to-[hsl(42,95%,42%)]",
                    "bg-clip-text text-transparent",
                  )}
                >
                  胜利
                </h2>
                <div className="relative flex items-center gap-1 sm:gap-2 mt-2 sm:mt-4 text-state-active/80 text-xs sm:text-sm">
                  <span className="text-base sm:text-lg">♠</span>
                  <span className="tracking-wider whitespace-nowrap">恭喜赢得本局</span>
                  <span className="text-base sm:text-lg">♥</span>
                </div>
              </div>
            ) : (
              <div className="flex flex-col items-center">
                <Swords className="w-10 sm:w-14 md:w-16 h-10 sm:h-14 md:h-16 text-white/30 mb-2 sm:mb-4" />
                <h2
                  className={cn(
                    "text-3xl sm:text-4xl md:text-5xl font-black tracking-[0.2em] sm:tracking-[0.25em] font-display pl-[0.2em] sm:pl-[0.25em]",
                    "bg-gradient-to-b from-white/60 via-white/40 to-white/20",
                    "bg-clip-text text-transparent",
                  )}
                >
                  失败
                </h2>
                <div className="text-white/40 mt-2 sm:mt-4 text-xs sm:text-sm tracking-wider">再接再厉</div>
              </div>
            )}
          </div>

          {/* 右侧: 队伍结果区 */}
          <div className="flex-1 flex flex-col gap-3 sm:gap-4">

            {/* 胜方队伍 */}
            <div className={cn(
              "p-3 sm:p-4 md:p-5 rounded-xl sm:rounded-2xl relative overflow-hidden",
              "bg-gradient-to-br from-state-active/15 via-state-active/5 to-transparent",
              "border border-state-active/30",
              "shadow-[inset_0_1px_0_hsla(42,95%,52%,0.2),0_4px_24px_hsla(42,95%,52%,0.15)]",
            )}>
              {/* 装饰光效 */}
              <div className="absolute top-0 right-0 w-20 sm:w-28 md:w-32 h-20 sm:h-28 md:h-32 bg-gradient-to-br from-state-active/20 to-transparent rounded-bl-full pointer-events-none"></div>

              <div className="relative flex items-center justify-between">
                {/* 玩家卡片 */}
                <div className="flex items-center gap-3 sm:gap-4 md:gap-6">
                  {teamRankings[winningTeam].map(({ rank, player }) => (
                    <PlayerRankCard
                      key={player.seat}
                      player={player}
                      rank={rank}
                      isWinningTeam={true}
                    />
                  ))}
                </div>

                {/* 等级显示 */}
                <div className="pr-1 sm:pr-2 md:pr-4">
                  <LevelDisplay
                    level={teamLevels[winningTeam]}
                    change={dealResult.levelChange[winningTeam]}
                    isWinning={true}
                  />
                </div>
              </div>
            </div>

            {/* 负方队伍 */}
            <div className={cn(
              "p-3 sm:p-4 md:p-5 rounded-xl sm:rounded-2xl relative overflow-hidden",
              "bg-black/30",
              "border border-white/10",
            )}>
              {/* 装饰 */}
              <div className="absolute top-0 right-0 w-16 sm:w-20 md:w-24 h-16 sm:h-20 md:h-24 bg-gradient-to-br from-white/5 to-transparent rounded-bl-full pointer-events-none"></div>

              <div className="relative flex items-center justify-between">
                {/* 玩家卡片 */}
                <div className="flex items-center gap-3 sm:gap-4 md:gap-6">
                  {teamRankings[losingTeam].map(({ rank, player }) => (
                    <PlayerRankCard
                      key={player.seat}
                      player={player}
                      rank={rank}
                      isWinningTeam={false}
                    />
                  ))}
                </div>

                {/* 等级显示 */}
                <div className="pr-1 sm:pr-2 md:pr-4">
                  <LevelDisplay
                    level={teamLevels[losingTeam]}
                    change={dealResult.levelChange[losingTeam]}
                    isWinning={false}
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 底部按钮 */}
        <div className="flex justify-center mt-4 sm:mt-6 md:mt-8">
          <Button
            intent="neutral"
            size="lg"
            onClick={onExit}
            className={cn(
              "px-6 sm:px-8 md:px-10 py-2 sm:py-3",
              "text-sm sm:text-base",
              "bg-gradient-to-b from-white/15 to-white/5",
              "hover:from-white/20 hover:to-white/10",
              "text-white/90 hover:text-white",
              "border border-white/20 hover:border-white/30",
              "shadow-[0_4px_12px_rgba(0,0,0,0.3),inset_0_1px_0_rgba(255,255,255,0.1)]",
              "transition-all duration-200",
            )}
          >
            <Home className="w-4 sm:w-5 h-4 sm:h-5 mr-1.5 sm:mr-2" />
            返回大厅
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default DealResult;
