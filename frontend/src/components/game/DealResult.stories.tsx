import type { Story } from "@ladle/react";
import DealResult from "./DealResult";
import type { Player } from "../../types";
import type { DealEndedPayload, PlayerDealStats } from "../../types/generated/event";
import { VictoryType } from "../../types/proto";

const mockPlayers: Player[] = [
  { id: "user_1", username: "小明", seat: 0, online: true, auto_play: false },
  { id: "user_2", username: "小红", seat: 1, online: true, auto_play: false },
  { id: "user_3", username: "小刚", seat: 2, online: true, auto_play: false },
  { id: "user_4", username: "小丽", seat: 3, online: true, auto_play: false },
];

const mockPlayerStats: PlayerDealStats[] = [
  { playerSeat: 0, cardsPlayed: 27, tricksWon: 11, passCount: 0, finishRank: 1 },
  { playerSeat: 1, cardsPlayed: 27, tricksWon: 9, passCount: 11, finishRank: 2 },
  { playerSeat: 2, cardsPlayed: 26, tricksWon: 7, passCount: 23, finishRank: 3 },
  { playerSeat: 3, cardsPlayed: 17, tricksWon: 1, passCount: 30, finishRank: 4 },
];

const baseDealResult: DealEndedPayload = {
  dealLevel: 13,
  rankings: [0, 1, 2, 3],
  victoryType: VictoryType.VICTORY_TYPE_SINGLE_LAST,
  winningTeam: 0,
  levelChange: [2, 0],
  durationMs: 454000,
  trickCount: 28,
  playerStats: mockPlayerStats,
  nextDealDeadlineMs: Date.now() + 6000,
};

const teamLevels: [number, number] = [13, 2];

export const Victory: Story = () => (
  <DealResult
    dealResult={baseDealResult}
    players={mockPlayers}
    teamLevels={teamLevels}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={false}
    currentPlayerSeat={0}
  />
);
Victory.storyName = "胜利场景";

export const Defeat: Story = () => (
  <DealResult
    dealResult={baseDealResult}
    players={mockPlayers}
    teamLevels={teamLevels}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={false}
    currentPlayerSeat={1}
  />
);
Defeat.storyName = "失败场景";

export const MatchFinished: Story = () => (
  <DealResult
    dealResult={{
      ...baseDealResult,
      nextDealDeadlineMs: 0,
    }}
    players={mockPlayers}
    teamLevels={[14, 2]}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={true}
    currentPlayerSeat={0}
  />
);
MatchFinished.storyName = "比赛结束";

export const DoubleDown: Story = () => (
  <DealResult
    dealResult={{
      ...baseDealResult,
      rankings: [0, 2, 1, 3],
      victoryType: VictoryType.VICTORY_TYPE_DOUBLE_DOWN,
      levelChange: [3, 0],
      playerStats: [
        { playerSeat: 0, cardsPlayed: 27, tricksWon: 12, passCount: 0, finishRank: 1 },
        { playerSeat: 2, cardsPlayed: 25, tricksWon: 10, passCount: 5, finishRank: 2 },
        { playerSeat: 1, cardsPlayed: 20, tricksWon: 5, passCount: 15, finishRank: 3 },
        { playerSeat: 3, cardsPlayed: 15, tricksWon: 1, passCount: 25, finishRank: 4 },
      ],
    }}
    players={mockPlayers}
    teamLevels={[10, 5]}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={false}
    currentPlayerSeat={0}
  />
);
DoubleDown.storyName = "双下胜利";

export const NoCountdown: Story = () => (
  <DealResult
    dealResult={{
      ...baseDealResult,
      nextDealDeadlineMs: 0,
    }}
    players={mockPlayers}
    teamLevels={teamLevels}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={false}
    currentPlayerSeat={0}
  />
);
NoCountdown.storyName = "无倒计时";

export const PartnerLast: Story = () => (
  <DealResult
    dealResult={{
      ...baseDealResult,
      rankings: [0, 1, 3, 2],
      victoryType: VictoryType.VICTORY_TYPE_PARTNER_LAST,
      levelChange: [1, 0],
      playerStats: [
        { playerSeat: 0, cardsPlayed: 27, tricksWon: 10, passCount: 2, finishRank: 1 },
        { playerSeat: 1, cardsPlayed: 27, tricksWon: 9, passCount: 5, finishRank: 2 },
        { playerSeat: 3, cardsPlayed: 25, tricksWon: 6, passCount: 12, finishRank: 3 },
        { playerSeat: 2, cardsPlayed: 18, tricksWon: 3, passCount: 20, finishRank: 4 },
      ],
    }}
    players={mockPlayers}
    teamLevels={[8, 6]}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={false}
    currentPlayerSeat={0}
  />
);
PartnerLast.storyName = "对贡胜利";

export const HighLevels: Story = () => (
  <DealResult
    dealResult={{
      ...baseDealResult,
      dealLevel: 14,
      levelChange: [1, 0],
    }}
    players={mockPlayers}
    teamLevels={[14, 11]}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={false}
    currentPlayerSeat={0}
  />
);
HighLevels.storyName = "高等级显示 (A vs J)";

export const LongDuration: Story = () => (
  <DealResult
    dealResult={{
      ...baseDealResult,
      durationMs: 1800000,
      trickCount: 80,
      playerStats: [
        { playerSeat: 0, cardsPlayed: 27, tricksWon: 25, passCount: 5, finishRank: 1 },
        { playerSeat: 1, cardsPlayed: 27, tricksWon: 20, passCount: 15, finishRank: 2 },
        { playerSeat: 2, cardsPlayed: 27, tricksWon: 18, passCount: 25, finishRank: 3 },
        { playerSeat: 3, cardsPlayed: 27, tricksWon: 17, passCount: 30, finishRank: 4 },
      ],
    }}
    players={mockPlayers}
    teamLevels={teamLevels}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={false}
    currentPlayerSeat={0}
  />
);
LongDuration.storyName = "长时间游戏 (30分钟)";

export const CloseMatch: Story = () => (
  <DealResult
    dealResult={{
      ...baseDealResult,
      rankings: [0, 1, 2, 3],
      victoryType: VictoryType.VICTORY_TYPE_SINGLE_LAST,
      levelChange: [1, 0],
      playerStats: [
        { playerSeat: 0, cardsPlayed: 27, tricksWon: 8, passCount: 10, finishRank: 1 },
        { playerSeat: 1, cardsPlayed: 27, tricksWon: 8, passCount: 11, finishRank: 2 },
        { playerSeat: 2, cardsPlayed: 26, tricksWon: 7, passCount: 12, finishRank: 3 },
        { playerSeat: 3, cardsPlayed: 25, tricksWon: 5, passCount: 14, finishRank: 4 },
      ],
      durationMs: 720000,
      trickCount: 35,
    }}
    players={mockPlayers}
    teamLevels={[7, 7]}
    onContinue={() => console.log("Continue clicked")}
    onExit={() => console.log("Exit clicked")}
    isMatchFinished={false}
    currentPlayerSeat={0}
  />
);
CloseMatch.storyName = "势均力敌";
