import type { Story } from "@ladle/react";
import CardDisplay from "./CardDisplay";
import type { Card } from "../../types/proto";

const makeCard = (suit: number, rank: number, deckIndex = 0): Card => ({
  suit,
  rank,
  deckIndex,
});

export const SuitsAndJokers: Story = () => (
  <div className="flex gap-4 p-8 bg-green-700">
    <CardDisplay card={makeCard(0, 14)} />
    <CardDisplay card={makeCard(1, 14)} />
    <CardDisplay card={makeCard(2, 14)} />
    <CardDisplay card={makeCard(3, 14)} />
    <CardDisplay card={makeCard(-1, 15)} />
    <CardDisplay card={makeCard(-1, 16)} />
  </div>
);

export const AllRanks: Story = () => (
  <div className="flex flex-wrap gap-2 p-8 bg-green-700">
    {[2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14].map((rank) => (
      <CardDisplay key={rank} card={makeCard(0, rank)} />
    ))}
  </div>
);

export const Sizes: Story = () => (
  <div className="flex items-end gap-4 p-8 bg-green-700">
    <CardDisplay card={makeCard(1, 12)} size="small" />
    <CardDisplay card={makeCard(1, 12)} size="normal" />
  </div>
);

export const Selected: Story = () => (
  <div className="flex gap-4 p-8 bg-green-700">
    <CardDisplay card={makeCard(0, 14)} />
    <CardDisplay card={makeCard(0, 14)} isSelected />
  </div>
);

const stackedCards = [
  makeCard(0, 12),
  makeCard(0, 13),
  makeCard(0, 14),
  makeCard(-1, 15),
  makeCard(-1, 16),
];

export const Stacked: Story = () => (
  <div className="flex p-8 bg-green-700">
    {stackedCards.map((card, i) => (
      <CardDisplay key={i} card={card} stackIndex={i} />
    ))}
  </div>
);

export const StackedSmall: Story = () => (
  <div className="flex p-8 bg-green-700">
    {stackedCards.map((card, i) => (
      <CardDisplay key={i} card={card} stackIndex={i} size="small" />
    ))}
  </div>
);
