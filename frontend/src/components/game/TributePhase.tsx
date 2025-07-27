import React, { useState, useEffect } from 'react';
import { 
  TributePhase as TributePhaseType, 
  TributeStatus, 
  TributeAction, 
  TributeActionType,
  Card, 
  Player 
} from '../../types';

interface TributePhaseProps {
  tributePhase: TributePhaseType;
  players: (Player | null)[];
  currentPlayerSeat: number;
  onSelectTribute: (cardId: string) => void;
  onReturnTribute: (cardId: string) => void;
}

interface CountdownTimerProps {
  targetTime: string;
  onTimeout: () => void;
}

const CountdownTimer: React.FC<CountdownTimerProps> = ({ targetTime, onTimeout }) => {
  const [timeLeft, setTimeLeft] = useState<number>(0);

  useEffect(() => {
    const updateTimer = () => {
      const now = new Date().getTime();
      const target = new Date(targetTime).getTime();
      const difference = Math.max(0, Math.floor((target - now) / 1000));
      
      setTimeLeft(difference);
      
      if (difference === 0) {
        onTimeout();
      }
    };

    updateTimer();
    const interval = setInterval(updateTimer, 1000);

    return () => clearInterval(interval);
  }, [targetTime, onTimeout]);

  return (
    <div className="text-lg font-bold text-red-600">
      {timeLeft}秒
    </div>
  );
};

interface CardDisplayProps {
  card: Card;
  onClick?: () => void;
  selectable?: boolean;
  selected?: boolean;
  size?: 'small' | 'medium' | 'large';
}

const CardDisplay: React.FC<CardDisplayProps> = ({ 
  card, 
  onClick, 
  selectable = false, 
  selected = false,
  size = 'medium'
}) => {
  const getSuitSymbol = (suit: number) => {
    switch (suit) {
      case 0: return '♠';
      case 1: return '♥';
      case 2: return '♣';
      case 3: return '♦';
      default: return '';
    }
  };

  const getRankText = (rank: number) => {
    if (rank === 15) return '小王';
    if (rank === 16) return '大王';
    if (rank <= 10) return rank.toString();
    switch (rank) {
      case 11: return 'J';
      case 12: return 'Q';
      case 13: return 'K';
      case 14: return 'A';
      default: return rank.toString();
    }
  };

  const getSuitColor = (suit: number) => {
    return suit === 1 || suit === 3 ? 'text-red-600' : 'text-black';
  };

  const getSizeClasses = () => {
    switch (size) {
      case 'small':
        return 'w-8 h-12 text-xs';
      case 'large':
        return 'w-16 h-24 text-lg';
      default:
        return 'w-12 h-18 text-sm';
    }
  };

  return (
    <div
      className={`
        ${getSizeClasses()}
        bg-white border-2 rounded-lg flex flex-col items-center justify-center
        ${selectable ? 'cursor-pointer hover:shadow-md' : ''}
        ${selected ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}
        ${selectable && !selected ? 'hover:border-gray-400' : ''}
      `}
      onClick={selectable ? onClick : undefined}
    >
      {card.is_joker ? (
        <div className="text-center">
          <div className="font-bold text-red-600">{getRankText(card.rank)}</div>
        </div>
      ) : (
        <div className="text-center">
          <div className={`font-bold ${getSuitColor(card.suit)}`}>
            {getRankText(card.rank)}
          </div>
          <div className={`${getSuitColor(card.suit)}`}>
            {getSuitSymbol(card.suit)}
          </div>
        </div>
      )}
    </div>
  );
};

interface TributeInfoDisplayProps {
  tributePhase: TributePhaseType;
  players: (Player | null)[];
}

const TributeInfoDisplay: React.FC<TributeInfoDisplayProps> = ({ tributePhase, players }) => {
  const getPlayerName = (seat: number) => {
    return players[seat]?.username || `玩家${seat + 1}`;
  };

  const renderTributeInfo = () => {
    if (tributePhase.is_immune) {
      return (
        <div className="bg-yellow-100 border border-yellow-400 rounded-lg p-3 mb-4">
          <div className="text-center text-yellow-800 font-medium">
            🛡️ 抗贡生效 - 败方持有2张及以上大王，免于上贡
          </div>
        </div>
      );
    }

    const tributeInfos: string[] = [];
    
    // Check if this is double down scenario
    const isDoubleDown = Object.values(tributePhase.tribute_map).some(receiver => receiver === -1);
    
    if (isDoubleDown) {
      tributeInfos.push('双下情况：败方贡牌到池，胜方选择');
      Object.entries(tributePhase.tribute_map).forEach(([giver, receiver]) => {
        if (receiver === -1) {
          tributeInfos.push(`${getPlayerName(parseInt(giver))} 贡牌到池`);
        }
      });
    } else {
      Object.entries(tributePhase.tribute_map).forEach(([giver, receiver]) => {
        if (receiver !== -1) {
          tributeInfos.push(`${getPlayerName(parseInt(giver))} → ${getPlayerName(receiver)}`);
        }
      });
    }

    return (
      <div className="bg-blue-100 border border-blue-400 rounded-lg p-3 mb-4">
        <div className="text-blue-800 font-medium mb-2">上贡信息</div>
        <div className="space-y-1">
          {tributeInfos.map((info, index) => (
            <div key={index} className="text-sm text-blue-700">{info}</div>
          ))}
        </div>
      </div>
    );
  };

  return renderTributeInfo();
};

interface PoolSelectionProps {
  poolCards: Card[];
  selectingPlayer: number;
  currentPlayerSeat: number;
  players: (Player | null)[];
  selectTimeout: string;
  onSelectCard: (cardId: string) => void;
  onTimeout: () => void;
}

const PoolSelection: React.FC<PoolSelectionProps> = ({
  poolCards,
  selectingPlayer,
  currentPlayerSeat,
  players,
  selectTimeout,
  onSelectCard,
  onTimeout
}) => {
  const [selectedCard, setSelectedCard] = useState<string | null>(null);
  const isMyTurn = selectingPlayer === currentPlayerSeat;

  const handleCardClick = (card: Card) => {
    if (!isMyTurn) return;
    setSelectedCard(card.id);
  };

  const handleConfirmSelection = () => {
    if (selectedCard) {
      onSelectCard(selectedCard);
      setSelectedCard(null);
    }
  };

  return (
    <div className="bg-white border border-gray-300 rounded-lg p-4">
      <div className="text-center mb-4">
        <div className="text-lg font-medium mb-2">
          {isMyTurn ? '请选择一张贡牌' : `等待 ${players[selectingPlayer]?.username || `玩家${selectingPlayer + 1}`} 选择贡牌`}
        </div>
        <CountdownTimer targetTime={selectTimeout} onTimeout={onTimeout} />
      </div>

      <div className="flex justify-center gap-4 mb-4">
        {poolCards.map((card) => (
          <CardDisplay
            key={card.id}
            card={card}
            selectable={isMyTurn}
            selected={selectedCard === card.id}
            onClick={() => handleCardClick(card)}
            size="large"
          />
        ))}
      </div>

      {isMyTurn && selectedCard && (
        <div className="text-center">
          <button
            onClick={handleConfirmSelection}
            className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-2 rounded-lg font-medium"
          >
            确认选择
          </button>
        </div>
      )}
    </div>
  );
};

interface ReturnTributeProps {
  tributePhase: TributePhaseType;
  currentPlayerSeat: number;
  players: (Player | null)[];
  playerHand: Card[];
  onReturnCard: (cardId: string) => void;
}

const ReturnTribute: React.FC<ReturnTributeProps> = ({
  tributePhase,
  currentPlayerSeat,
  players,
  playerHand,
  onReturnCard
}) => {
  const [selectedCard, setSelectedCard] = useState<string | null>(null);

  // Check if current player needs to return tribute
  const needsReturn = Object.values(tributePhase.tribute_map).includes(currentPlayerSeat) &&
                     !tributePhase.return_cards[currentPlayerSeat];

  if (!needsReturn) {
    return (
      <div className="bg-gray-100 border border-gray-300 rounded-lg p-4 text-center">
        <div className="text-gray-600">等待其他玩家还贡...</div>
      </div>
    );
  }

  const handleCardClick = (card: Card) => {
    setSelectedCard(card.id);
  };

  const handleConfirmReturn = () => {
    if (selectedCard) {
      onReturnCard(selectedCard);
      setSelectedCard(null);
    }
  };

  return (
    <div className="bg-white border border-gray-300 rounded-lg p-4">
      <div className="text-center mb-4">
        <div className="text-lg font-medium mb-2">请选择一张牌还贡</div>
        <div className="text-sm text-gray-600">选择一张最小的牌还给上贡者</div>
      </div>

      <div className="flex flex-wrap justify-center gap-2 mb-4 max-h-32 overflow-y-auto">
        {playerHand.map((card) => (
          <CardDisplay
            key={card.id}
            card={card}
            selectable={true}
            selected={selectedCard === card.id}
            onClick={() => handleCardClick(card)}
            size="small"
          />
        ))}
      </div>

      {selectedCard && (
        <div className="text-center">
          <button
            onClick={handleConfirmReturn}
            className="bg-green-500 hover:bg-green-600 text-white px-6 py-2 rounded-lg font-medium"
          >
            确认还贡
          </button>
        </div>
      )}
    </div>
  );
};

interface TributeResultDisplayProps {
  tributePhase: TributePhaseType;
  players: (Player | null)[];
}

const TributeResultDisplay: React.FC<TributeResultDisplayProps> = ({ tributePhase, players }) => {
  const getPlayerName = (seat: number) => {
    return players[seat]?.username || `玩家${seat + 1}`;
  };

  const renderTributeResults = () => {
    const results: JSX.Element[] = [];

    // Show tribute cards
    Object.entries(tributePhase.tribute_cards).forEach(([giver, card]) => {
      const giverSeat = parseInt(giver);
      const receiver = tributePhase.tribute_map[giverSeat];
      
      results.push(
        <div key={`tribute-${giver}`} className="flex items-center justify-between p-2 bg-blue-50 rounded">
          <div className="flex items-center gap-2">
            <span className="text-sm">{getPlayerName(giverSeat)} 上贡:</span>
            <CardDisplay card={card} size="small" />
          </div>
          <div className="text-sm text-gray-600">
            → {receiver === -1 ? '贡池' : getPlayerName(receiver)}
          </div>
        </div>
      );
    });

    // Show return cards
    Object.entries(tributePhase.return_cards).forEach(([receiver, card]) => {
      const receiverSeat = parseInt(receiver);
      
      results.push(
        <div key={`return-${receiver}`} className="flex items-center justify-between p-2 bg-green-50 rounded">
          <div className="flex items-center gap-2">
            <span className="text-sm">{getPlayerName(receiverSeat)} 还贡:</span>
            <CardDisplay card={card} size="small" />
          </div>
        </div>
      );
    });

    return results;
  };

  return (
    <div className="bg-white border border-gray-300 rounded-lg p-4">
      <div className="text-lg font-medium mb-3 text-center">贡牌结果</div>
      <div className="space-y-2">
        {renderTributeResults()}
      </div>
      <div className="text-center mt-4 text-sm text-gray-600">
        3秒后自动进入出牌阶段...
      </div>
    </div>
  );
};

const TributePhase: React.FC<TributePhaseProps> = ({
  tributePhase,
  players,
  currentPlayerSeat,
  onSelectTribute,
  onReturnTribute
}) => {
  const handleTimeout = () => {
    // Handle timeout - this would typically be handled by the parent component
    // or through WebSocket communication
    console.log('Tribute phase timeout');
  };

  const renderPhaseContent = () => {
    switch (tributePhase.status) {
      case TributeStatus.WAITING:
        return (
          <div className="text-center py-8">
            <div className="text-lg text-gray-600">准备上贡阶段...</div>
          </div>
        );

      case TributeStatus.SELECTING:
        return (
          <PoolSelection
            poolCards={tributePhase.pool_cards}
            selectingPlayer={tributePhase.selecting_player}
            currentPlayerSeat={currentPlayerSeat}
            players={players}
            selectTimeout={tributePhase.select_timeout}
            onSelectCard={onSelectTribute}
            onTimeout={handleTimeout}
          />
        );

      case TributeStatus.RETURNING:
        return (
          <ReturnTribute
            tributePhase={tributePhase}
            currentPlayerSeat={currentPlayerSeat}
            players={players}
            playerHand={[]} // This would need to be passed from parent
            onReturnCard={onReturnTribute}
          />
        );

      case TributeStatus.FINISHED:
        return (
          <TributeResultDisplay
            tributePhase={tributePhase}
            players={players}
          />
        );

      default:
        return null;
    }
  };

  return (
    <div className="w-full max-w-4xl mx-auto p-4">
      <div className="bg-gradient-to-r from-blue-500 to-purple-600 text-white text-center py-3 rounded-t-lg">
        <h2 className="text-xl font-bold">上贡阶段</h2>
      </div>
      
      <div className="bg-gray-50 p-4 space-y-4">
        <TributeInfoDisplay tributePhase={tributePhase} players={players} />
        {renderPhaseContent()}
      </div>
    </div>
  );
};

export default TributePhase;