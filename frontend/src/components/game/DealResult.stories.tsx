import type { Story } from "@ladle/react";
import DealResult from "./DealResult";
import type { Player } from "../../types";
import type { DealEndedPayload, PlayerDealStats } from "../../types/generated/event";
import { VictoryType } from "../../types/proto";

const mockPlayers: Player[] = [
  { id: "user_1", username: "lxc", seat: 0, online: true, auto_play: false },
  { id: "user_2", username: "ai_player_1", seat: 1, online: true, auto_play: true },
  { id: "user_3", username: "ai_player_2", seat: 2, online: true, auto_play: true },
  { id: "user_4", username: "ai_player_3", seat: 3, online: true, auto_play: true },
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
