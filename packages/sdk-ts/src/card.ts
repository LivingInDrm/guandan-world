import { ProtoCard, Suit } from './types';

const NameMap: Record<number, string> = {
  11: 'Jack',
  12: 'Queen',
  13: 'King',
  14: 'Ace',
  15: 'Black Joker',
  16: 'Red Joker',
};

export class SdkCard {
  readonly rank: number;
  readonly rawRank: number;
  readonly suit: number;
  readonly level: number;
  readonly deckIndex: number;
  readonly name: string;

  constructor(rank: number, suit: number, level: number, deckIndex: number = 0) {
    this.suit = suit;
    this.level = level;
    this.deckIndex = deckIndex;

    this.rawRank = rank;
    if (rank === 14) {
      this.rawRank = 1;
    }
    if (rank === 1) {
      rank = 14;
    }
    this.rank = rank;

    if (rank >= 2 && rank <= 10) {
      this.name = String(rank);
    } else {
      this.name = NameMap[rank] ?? String(rank);
    }
  }

  isWildcard(): boolean {
    return this.rank === this.level && this.suit === Suit.Heart;
  }

  isBigJoker(): boolean {
    return this.rank === 16 && this.suit === Suit.Joker;
  }

  isSmallJoker(): boolean {
    return this.rank === 15 && this.suit === Suit.Joker;
  }

  isJoker(): boolean {
    return this.suit === Suit.Joker;
  }

  greaterThan(other: SdkCard): boolean {
    if (other.rank === this.level) {
      return this.rank >= 15;
    } else {
      if (this.rank === this.level) {
        return other.rank <= 14;
      } else {
        return this.rank > other.rank;
      }
    }
  }

  consecutiveGreaterThan(other: SdkCard): boolean {
    return this.rawRank > other.rawRank;
  }

  lessThan(other: SdkCard): boolean {
    if (other.greaterThan(this)) {
      return true;
    } else if (this.equals(other) && other.suit === Suit.Heart && this.suit !== Suit.Heart) {
      return true;
    }
    return false;
  }

  equals(other: SdkCard): boolean {
    return this.rank === other.rank;
  }

  toString(): string {
    if (this.suit !== Suit.Joker) {
      const suitName = this.getSuitName();
      return `${this.name} of ${suitName}`;
    }
    return this.name;
  }

  toShortString(): string {
    if (this.suit === Suit.Joker) {
      return this.rank === 15 ? 'SJ' : 'BJ';
    }

    let rankStr: string;
    switch (this.rank) {
      case 11:
        rankStr = 'J';
        break;
      case 12:
        rankStr = 'Q';
        break;
      case 13:
        rankStr = 'K';
        break;
      case 14:
        rankStr = 'A';
        break;
      default:
        rankStr = String(this.rank);
    }

    const suitStr = ['S', 'H', 'C', 'D'][this.suit] ?? '?';
    return rankStr + suitStr;
  }

  private getSuitName(): string {
    switch (this.suit) {
      case Suit.Spade:
        return 'Spade';
      case Suit.Heart:
        return 'Heart';
      case Suit.Club:
        return 'Club';
      case Suit.Diamond:
        return 'Diamond';
      default:
        return 'Joker';
    }
  }

  clone(): SdkCard {
    return new SdkCard(this.rank, this.suit, this.level, this.deckIndex);
  }
}

export function fromProtoCard(card: ProtoCard, level: number): SdkCard {
  return new SdkCard(card.rank, card.suit, level, card.deckIndex);
}

export function fromProtoCards(cards: ProtoCard[], level: number): SdkCard[] {
  return cards.map((c) => fromProtoCard(c, level));
}
