import React, { useRef, useCallback, useState, useEffect } from 'react';
import type { Player, Card, TributePair } from '../../../types';
import type { PlayerPosition } from '../GameTable';
import { WS_MESSAGE_TYPES } from '../../../types';
import { TributeType } from '../../../types/generated/event';
import { TributeStatus } from '../../../types/generated/view';
import { EventType } from '../../../types/generated/event';
import { wsClient } from '../../../services/websocket';
import GameTable from '../GameTable';
import PlayerCard from '../PlayerCard';
import TeamLevelDisplay from '../TeamLevelDisplay';
import TributePool from './TributePool';
import TributeRoleBadge from './TributeRoleBadge';
import CardFlyAnimation from './CardFlyAnimation';
import PlayerHand from '../PlayerHand';

interface FlyingCard {
  id: string;
  card: Card;
  fromSeat: number | 'pool';
  toSeat: number | 'pool';
  fromPoolSlot?: number;
  toPoolSlot?: number;
}

interface TributeData {
  status: TributeStatus;
  tributeType: TributeType;
  givers: number[];
  receivers: number[];
  tributePairs: TributePair[];
  poolCards: Card[];
  isImmune: boolean;
}

interface TributeBoardProps {
  tributeData: TributeData;
  players: (Player | null)[];
  currentPlayerSeat: number;
  teamLevels: [number, number];
  currentLevel: number;
  playerHand: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onSelectTribute: (deckIndex: number) => void;
  onReturnTribute: (deckIndex: number) => void;
}

const TributeBoard: React.FC<TributeBoardProps> = ({
  tributeData,
  players,
  currentPlayerSeat,
  teamLevels,
  currentLevel,
  playerHand,
  selectedCards,
  onCardSelect,
  onSelectTribute,
  onReturnTribute,
}) => {
  const boardRef = useRef<HTMLDivElement>(null);
  const poolRef = useRef<HTMLDivElement>(null);
  const poolSlotRefs = useRef<{ [slot: number]: HTMLDivElement | null }>({});
  
  const [flyingCards, setFlyingCards] = useState<FlyingCard[]>([]);

  const { 
    status, 
    tributeType,
    givers,
    receivers,
    tributePairs,
    poolCards,
  } = tributeData;

  const isGiver = useCallback((seat: number) => givers.includes(seat), [givers]);
  const isReceiver = useCallback((seat: number) => receivers.includes(seat), [receivers]);
  
  const hasSubmitted = useCallback((seat: number) => {
    return tributePairs.some(pair => pair.giver === seat && pair.tributeCard);
  }, [tributePairs]);
  
  const hasReceived = useCallback((seat: number) => {
    return tributePairs.some(pair => pair.receiver === seat && pair.tributeCard);
  }, [tributePairs]);

  const getCurrentSelectingSeat = useCallback((): number | null => {
    if (status !== TributeStatus.TRIBUTE_STATUS_SELECTING) return null;
    for (const receiver of receivers) {
      const received = tributePairs.some(pair => pair.receiver === receiver && pair.tributeCard);
      if (!received) return receiver;
    }
    return null;
  }, [status, receivers, tributePairs]);

  const currentSelectingSeat = getCurrentSelectingSeat();

  const canSelectFromPool = 
    status === TributeStatus.TRIBUTE_STATUS_SELECTING && 
    currentSelectingSeat === currentPlayerSeat;
  
  const canReturnTribute = 
    status === TributeStatus.TRIBUTE_STATUS_RETURNING && 
    isReceiver(currentPlayerSeat);

  const getPositionForSeat = useCallback((
    seat: number | 'pool',
    poolSlot?: number
  ): { x: number; y: number } => {
    if (!boardRef.current) return { x: 0, y: 0 };
    const boardRect = boardRef.current.getBoundingClientRect();

    if (seat === 'pool') {
      const slotElement = poolSlotRefs.current[poolSlot ?? -1];
      if (slotElement) {
        const slotRect = slotElement.getBoundingClientRect();
        return {
          x: slotRect.left + slotRect.width / 2,
          y: slotRect.top + slotRect.height / 2,
        };
      }
      if (poolRef.current) {
        const poolRect = poolRef.current.getBoundingClientRect();
        return {
          x: poolRect.left + poolRect.width / 2,
          y: poolRect.top + poolRect.height / 2,
        };
      }
      return {
        x: boardRect.left + boardRect.width / 2,
        y: boardRect.top + boardRect.height / 2,
      };
    }

    const relativeSeat = (seat - currentPlayerSeat + 4) % 4;
    switch (relativeSeat) {
      case 0:
        return { x: boardRect.left + boardRect.width / 2, y: boardRect.bottom - 40 };
      case 1:
        return { x: boardRect.left + 40, y: boardRect.top + boardRect.height / 2 };
      case 2:
        return { x: boardRect.left + boardRect.width / 2, y: boardRect.top + 40 };
      case 3:
        return { x: boardRect.right - 40, y: boardRect.top + boardRect.height / 2 };
      default:
        return { x: 0, y: 0 };
    }
  }, [currentPlayerSeat]);

  const removeFlyingCard = useCallback((id: string) => {
    setFlyingCards(prev => prev.filter(fc => fc.id !== id));
  }, []);

  useEffect(() => {
    const handleGameEvent = (message: { data: unknown }) => {
      const event = message.data as {
        type?: EventType;
        actorSeat?: number;
        tributeCardSubmitted?: { submittedCard?: Card };
        tributeCardSelected?: { selectedCard?: Card };
        tributeCardReturned?: { returnedCard?: Card; targetPlayer?: number };
      };
      if (!event || !event.type) return;

      switch (event.type) {
        case EventType.EVENT_TYPE_TRIBUTE_CARD_SUBMITTED: {
          if (event.tributeCardSubmitted?.submittedCard && event.actorSeat !== undefined) {
            const card = event.tributeCardSubmitted.submittedCard;
            const actorSeat = event.actorSeat;
            
            const occupiedSlots = new Set(
              flyingCards
                .filter(fc => fc.toSeat === 'pool' && fc.toPoolSlot !== undefined)
                .map(fc => fc.toPoolSlot)
            );
            const maxSlots = tributeType === TributeType.TRIBUTE_TYPE_DOUBLE_DOWN ? 2 : 1;
            let toPoolSlot = 0;
            for (let i = 0; i < maxSlots; i++) {
              if (!occupiedSlots.has(i) && !poolCards[i]) {
                toPoolSlot = i;
                break;
              }
            }

            const flyingCard: FlyingCard = {
              id: `submit-${actorSeat}-${Date.now()}`,
              card,
              fromSeat: actorSeat,
              toSeat: 'pool',
              toPoolSlot,
            };
            setFlyingCards(prev => [...prev, flyingCard]);
          }
          break;
        }

        case EventType.EVENT_TYPE_TRIBUTE_CARD_SELECTED: {
          if (event.tributeCardSelected?.selectedCard && event.actorSeat !== undefined) {
            const card = event.tributeCardSelected.selectedCard;
            const actorSeat = event.actorSeat;
            
            const fromPoolSlot = poolCards.findIndex(
              c => c && c.deckIndex === card.deckIndex
            );

            const flyingCard: FlyingCard = {
              id: `select-${actorSeat}-${Date.now()}`,
              card,
              fromSeat: 'pool',
              toSeat: actorSeat,
              fromPoolSlot: fromPoolSlot >= 0 ? fromPoolSlot : undefined,
            };
            setFlyingCards(prev => [...prev, flyingCard]);
          }
          break;
        }

        case EventType.EVENT_TYPE_TRIBUTE_CARD_RETURNED: {
          if (event.tributeCardReturned?.returnedCard && event.actorSeat !== undefined) {
            const card = event.tributeCardReturned.returnedCard;
            const actorSeat = event.actorSeat;
            const targetPlayer = event.tributeCardReturned.targetPlayer;
            
            if (targetPlayer === undefined) break;

            const flyingCard: FlyingCard = {
              id: `return-${actorSeat}-${Date.now()}`,
              card,
              fromSeat: actorSeat,
              toSeat: targetPlayer,
            };
            setFlyingCards(prev => [...prev, flyingCard]);
          }
          break;
        }

        case EventType.EVENT_TYPE_TRIBUTE_STARTED:
        case EventType.EVENT_TYPE_TRIBUTE_COMPLETED: {
          setFlyingCards([]);
          break;
        }
      }
    };

    wsClient.on(WS_MESSAGE_TYPES.GAME_EVENT, handleGameEvent);
    return () => {
      wsClient.off(WS_MESSAGE_TYPES.GAME_EVENT, handleGameEvent);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tributeType, poolCards]);

  const handlePoolCardSelect = (card: Card) => {
    if (canSelectFromPool) {
      onSelectTribute(card.deckIndex);
    }
  };

  const handleReturnCard = () => {
    if (canReturnTribute && selectedCards.length === 1) {
      onReturnTribute(selectedCards[0].deckIndex);
    }
  };

  const handleSlotRefReady = useCallback((slotIndex: number, element: HTMLDivElement | null) => {
    poolSlotRefs.current[slotIndex] = element;
  }, []);

  const renderPlayer = (player: Player | null, position: PlayerPosition, seat: number) => {
    if (position === 'bottom') {
      return null;
    }

    const isGiverRole = isGiver(seat);
    const isReceiverRole = isReceiver(seat);
    const role = isGiverRole ? 'giver' : isReceiverRole ? 'receiver' : null;
    const isCurrentSelector = currentSelectingSeat === seat;

    return (
      <PlayerCard
        player={player}
        position={position}
        isHighlighted={isCurrentSelector}
        statusSlot={
          role ? (
            <TributeRoleBadge
              role={role}
              isSubmitted={hasSubmitted(seat)}
              isReceived={hasReceived(seat)}
              isCurrentSelector={isCurrentSelector}
            />
          ) : undefined
        }
      />
    );
  };

  const renderCenter = () => {
    const flyingToPoolDeckIndexes = flyingCards
      .filter(fc => fc.toSeat === 'pool')
      .map(fc => fc.card.deckIndex);

    return (
      <div ref={poolRef} className="absolute inset-0 flex items-center justify-center">
        <TributePool
          poolCards={poolCards}
          canSelect={canSelectFromPool}
          onSelectCard={handlePoolCardSelect}
          tributeType={tributeType}
          messages={[]}
          onSlotRefReady={handleSlotRefReady}
          flyingToPoolDeckIndexes={flyingToPoolDeckIndexes}
        />
      </div>
    );
  };

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-6">
      <div ref={boardRef} className="relative">
        <GameTable
          players={players}
          currentPlayerSeat={currentPlayerSeat}
          renderPlayer={renderPlayer}
          renderCenter={renderCenter}
          topLeftSlot={
            <TeamLevelDisplay 
              teamLevels={teamLevels} 
              currentLevel={currentLevel}
              currentPlayerSeat={currentPlayerSeat}
            />
          }
        />
        
        {flyingCards.map((fc) => (
          <CardFlyAnimation
            key={fc.id}
            card={fc.card}
            fromPosition={getPositionForSeat(fc.fromSeat, fc.fromPoolSlot)}
            toPosition={getPositionForSeat(fc.toSeat, fc.toPoolSlot)}
            onComplete={() => removeFlyingCard(fc.id)}
          />
        ))}
      </div>

      {status === TributeStatus.TRIBUTE_STATUS_RETURNING && canReturnTribute && (
        <>
          <PlayerHand
            cards={playerHand}
            selectedCards={selectedCards}
            onCardSelect={onCardSelect}
          />
          <div className="flex justify-center">
            <button
              onClick={handleReturnCard}
              disabled={selectedCards.length !== 1}
              className={`px-6 py-2 rounded-lg font-medium transition-colors ${
                selectedCards.length === 1
                  ? 'bg-primary text-primary-foreground hover:bg-primary/90'
                  : 'bg-muted text-muted-foreground cursor-not-allowed'
              }`}
            >
              确认还贡
            </button>
          </div>
        </>
      )}
    </div>
  );
};

export default TributeBoard;
