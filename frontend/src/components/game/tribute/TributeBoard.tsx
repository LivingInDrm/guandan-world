import React, { useRef, useCallback, useState, useEffect } from 'react';
import type { Player, Card, TributePair } from '../../../types';
import { WS_MESSAGE_TYPES } from '../../../types';
import { TributeType } from '../../../types/generated/event';
import { TributeStatus } from '../../../types/generated/view';
import { EventType } from '../../../types/generated/event';
import { wsClient } from '../../../services/websocket';
import TributePool from './TributePool';
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
  playerHand: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onSelectTribute: (deckIndex: number) => void;
  onReturnTribute: (deckIndex: number) => void;
}

const getTributeTypeName = (type: TributeType): string => {
  switch (type) {
    case TributeType.TRIBUTE_TYPE_DOUBLE_DOWN:
      return '双下';
    case TributeType.TRIBUTE_TYPE_SINGLE_LAST:
      return '单下';
    case TributeType.TRIBUTE_TYPE_PARTNER_LAST:
      return '末游';
    case TributeType.TRIBUTE_TYPE_NONE:
      return '无需进贡';
    default:
      return '未知';
  }
};

interface TributePlayerAreaProps {
  player: Player | null;
  position: 'bottom' | 'left' | 'top' | 'right';
  seat: number;
  isGiver: boolean;
  isReceiver: boolean;
  hasSubmitted: boolean;
  hasReceived: boolean;
  isCurrentSelector: boolean;
}

const TributePlayerArea: React.FC<TributePlayerAreaProps> = ({
  player,
  position,
  isGiver,
  isReceiver,
  hasSubmitted,
  hasReceived,
  isCurrentSelector,
}) => {
  const getPositionClasses = () => {
    switch (position) {
      case 'bottom':
        return 'absolute bottom-4 left-1/2 transform -translate-x-1/2';
      case 'left':
        return 'absolute left-4 top-1/2 transform -translate-y-1/2';
      case 'top':
        return 'absolute top-4 left-1/2 transform -translate-x-1/2';
      case 'right':
        return 'absolute right-4 top-1/2 transform -translate-y-1/2';
      default:
        return '';
    }
  };

  const getRoleBadge = () => {
    if (isGiver) {
      return (
        <div className="text-xs px-2 py-1 rounded bg-red-200 text-red-800">
          {hasSubmitted ? '已上贡' : '待上贡'}
        </div>
      );
    }
    if (isReceiver) {
      return (
        <div className={`text-xs px-2 py-1 rounded ${
          isCurrentSelector 
            ? 'bg-yellow-200 text-yellow-800 animate-pulse' 
            : 'bg-green-200 text-green-800'
        }`}>
          {hasReceived ? '已收贡' : isCurrentSelector ? '选牌中' : '待收贡'}
        </div>
      );
    }
    return null;
  };

  if (!player) {
    return (
      <div className={`${getPositionClasses()} w-24 h-16`}>
        <div className="bg-gray-100 border-2 border-dashed border-gray-300 rounded-lg p-2 text-center">
          <div className="text-sm text-gray-400">空座位</div>
        </div>
      </div>
    );
  }

  return (
    <div className={`${getPositionClasses()} w-32`}>
      <div className={`border-2 rounded-lg p-2 text-center bg-white ${
        isCurrentSelector ? 'border-yellow-400 shadow-lg' : 'border-gray-300'
      }`}>
        <div className="text-sm font-medium truncate">{player.username}</div>
        {getRoleBadge()}
      </div>
    </div>
  );
};

const TributeBoard: React.FC<TributeBoardProps> = ({
  tributeData,
  players,
  currentPlayerSeat,
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
    isImmune,
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
      if (poolSlot !== undefined && poolSlotRefs.current[poolSlot]) {
        const slotRect = poolSlotRefs.current[poolSlot]!.getBoundingClientRect();
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
    const handleGameEvent = (message: { data: any }) => {
      const event = message.data;
      if (!event || !event.type) return;

      switch (event.type) {
        case EventType.EVENT_TYPE_TRIBUTE_CARD_SUBMITTED: {
          if (event.tributeCardSubmitted?.submittedCard && event.actorSeat !== undefined) {
            const card = event.tributeCardSubmitted.submittedCard as Card;
            const actorSeat = event.actorSeat as number;
            
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
            const card = event.tributeCardSelected.selectedCard as Card;
            const actorSeat = event.actorSeat as number;
            
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
            const card = event.tributeCardReturned.returnedCard as Card;
            const actorSeat = event.actorSeat as number;
            const targetPlayer = event.tributeCardReturned.targetPlayer as number;

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
  }, [tributeType, poolCards, flyingCards]);

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

  const getStepTitle = () => {
    if (isImmune) return '抗贡成功';
    switch (status) {
      case TributeStatus.TRIBUTE_STATUS_WAITING:
        return '等待进贡';
      case TributeStatus.TRIBUTE_STATUS_SELECTING:
        return currentSelectingSeat === currentPlayerSeat ? '请选择贡牌' : '等待选牌';
      case TributeStatus.TRIBUTE_STATUS_RETURNING:
        return '还贡阶段';
      case TributeStatus.TRIBUTE_STATUS_FINISHED:
        return '上贡完成';
      default:
        return '进贡阶段';
    }
  };

  const handleSlotRefReady = useCallback((slotIndex: number, element: HTMLDivElement | null) => {
    poolSlotRefs.current[slotIndex] = element;
  }, []);

  const renderPlayerArea = (relativeSeat: number) => {
    const seat = (currentPlayerSeat + relativeSeat) % 4;
    const positions: Array<'bottom' | 'left' | 'top' | 'right'> = ['bottom', 'left', 'top', 'right'];
    return (
      <TributePlayerArea
        key={seat}
        player={players[seat]}
        position={positions[relativeSeat]}
        seat={seat}
        isGiver={isGiver(seat)}
        isReceiver={isReceiver(seat)}
        hasSubmitted={hasSubmitted(seat)}
        hasReceived={hasReceived(seat)}
        isCurrentSelector={currentSelectingSeat === seat}
      />
    );
  };

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-6">
      <div ref={boardRef} className="relative w-full h-96 bg-green-100 border border-gray-300 rounded-lg">
        <div className="absolute top-4 left-4 bg-white border border-gray-300 rounded-lg p-3 shadow-sm z-10">
          <div className="text-sm font-medium mb-2">{getStepTitle()}</div>
          <div className="text-xs text-gray-600">
            类型: {getTributeTypeName(tributeType)}
          </div>
        </div>

        <div ref={poolRef} className="absolute inset-0 flex items-center justify-center">
          <TributePool
            poolCards={poolCards}
            canSelect={canSelectFromPool}
            onSelectCard={handlePoolCardSelect}
            tributeType={tributeType}
            messages={[]}
            onSlotRefReady={handleSlotRefReady}
          />
        </div>

        {[0, 1, 2, 3].map(renderPlayerArea)}

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
                  ? 'bg-blue-500 text-white hover:bg-blue-600'
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
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
