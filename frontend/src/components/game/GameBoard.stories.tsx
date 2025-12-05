import type { Story } from "@ladle/react";
import GameBoard from "./GameBoard";
import type { Player } from "../../types";
import type { Card, PlayAction } from "../../types/proto";
import { CompType } from "../../types/proto";

const makeCard = (suit: number, rank: number, deckIndex = 0): Card => ({
  suit,
  rank,
  deckIndex,
});

const makePlayer = (id: string, username: string, seat: number): Player => ({
  id,
  username,
  seat,
  online: true,
  auto_play: false,
});

const makePlayAction = (
  playerSeat: number,
  cards: Card[],
  isPass = false
): PlayAction => ({
  playerSeat,
  cards,
  compType: CompType.COMP_TYPE_UNSPECIFIED,
  timestampMs: Date.now(),
  isPass,
});

const mockPlayers: Player[] = [
  makePlayer("p1", "玩家A", 0),
  makePlayer("p2", "玩家B", 1),
  makePlayer("p3", "玩家C", 2),
  makePlayer("p4", "玩家D", 3),
];

const emptyPlayers: (Player | null)[] = [null, null, null, null];

const pairOfKings: Card[] = [
  makeCard(0, 13, 0),
  makeCard(1, 13, 1),
];

const straight: Card[] = [
  makeCard(0, 5, 0),
  makeCard(1, 6, 1),
  makeCard(2, 7, 2),
  makeCard(3, 8, 3),
  makeCard(0, 9, 4),
];

const bomb: Card[] = [
  makeCard(0, 14, 0),
  makeCard(1, 14, 1),
  makeCard(2, 14, 2),
  makeCard(3, 14, 3),
];

const withOnePlayed: PlayAction[] = [
  makePlayAction(1, pairOfKings),
];

const withPlayAndPass: PlayAction[] = [
  makePlayAction(1, pairOfKings),
  makePlayAction(2, [], true),
];

const mixedPlays: PlayAction[] = [
  makePlayAction(1, straight),
  makePlayAction(2, [], true),
  makePlayAction(3, bomb),
];

export const EmptyBoard: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameBoard
      teamLevels={[2, 2]}
      currentLevel={2}
      plays={[]}
      currentTurn={0}
      players={emptyPlayers}
      currentPlayerSeat={0}
    />
  </div>
);

export const AllPlayersJoined: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameBoard
      teamLevels={[2, 2]}
      currentLevel={2}
      plays={[]}
      currentTurn={0}
      players={mockPlayers}
      currentPlayerSeat={0}
    />
  </div>
);

export const CurrentPlayerTurn: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameBoard
      teamLevels={[5, 3]}
      currentLevel={5}
      plays={[]}
      currentTurn={0}
      players={mockPlayers}
      currentPlayerSeat={0}
      playStates={[0, 0, 0, 0]}
    />
  </div>
);

export const WithPlayedCards: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameBoard
      teamLevels={[7, 5]}
      currentLevel={7}
      plays={withOnePlayed}
      currentTurn={2}
      players={mockPlayers}
      currentPlayerSeat={0}
      playStates={[0, 1, 0, 0]}
    />
  </div>
);

export const WithPassedPlayers: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameBoard
      teamLevels={[10, 8]}
      currentLevel={10}
      plays={withPlayAndPass}
      currentTurn={3}
      players={mockPlayers}
      currentPlayerSeat={0}
      playStates={[0, 1, 2, 0]}
    />
  </div>
);

export const MixedStates: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameBoard
      teamLevels={[12, 11]}
      currentLevel={12}
      plays={mixedPlays}
      currentTurn={0}
      players={mockPlayers}
      currentPlayerSeat={0}
      playStates={[0, 1, 2, 1]}
    />
  </div>
);

export const HighLevels: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameBoard
      teamLevels={[14, 13]}
      currentLevel={14}
      plays={[]}
      currentTurn={0}
      players={mockPlayers}
      currentPlayerSeat={0}
    />
  </div>
);
