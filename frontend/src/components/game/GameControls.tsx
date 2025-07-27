import React, { useState, useEffect, useCallback } from 'react';
import { Card } from '../../types';

interface GameControlsProps {
  selectedCards: Card[];
  canPlay: boolean;
  isMyTurn: boolean;
  turnTimeoutSeconds: number;
  onPlayCards: (cards: Card[]) => void;
  onPass: () => void;
  disabled?: boolean;
}

interface CountdownTimerProps {
  seconds: number;
  isActive: boolean;
}

const CountdownTimer: React.FC<CountdownTimerProps> = ({ seconds, isActive }) => {
  const [timeLeft, setTimeLeft] = useState(seconds);

  useEffect(() => {
    setTimeLeft(seconds);
  }, [seconds]);

  useEffect(() => {
    if (!isActive || timeLeft <= 0) return;

    const timer = setInterval(() => {
      setTimeLeft(prev => {
        if (prev <= 1) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [isActive, timeLeft]);

  if (!isActive) return null;

  const getTimerColor = () => {
    if (timeLeft <= 5) return 'text-red-600 bg-red-100';
    if (timeLeft <= 10) return 'text-orange-600 bg-orange-100';
    return 'text-blue-600 bg-blue-100';
  };

  return (
    <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${getTimerColor()}`}>
      <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" />
      </svg>
      {timeLeft}秒
    </div>
  );
};

interface PlayValidationResult {
  isValid: boolean;
  error?: string;
  cardType?: string;
}

const GameControls: React.FC<GameControlsProps> = ({
  selectedCards,
  canPlay,
  isMyTurn,
  turnTimeoutSeconds,
  onPlayCards,
  onPass,
  disabled = false
}) => {
  const [validationResult, setValidationResult] = useState<PlayValidationResult>({ isValid: true });

  // Simple card type validation (this would be more complex in a real implementation)
  const validateCards = useCallback((cards: Card[]): PlayValidationResult => {
    if (cards.length === 0) {
      return { isValid: false, error: '请选择要出的牌' };
    }

    if (cards.length === 1) {
      return { isValid: true, cardType: '单牌' };
    }

    if (cards.length === 2) {
      // Check if it's a pair
      if (cards[0].rank === cards[1].rank) {
        return { isValid: true, cardType: '对子' };
      }
      return { isValid: false, error: '两张牌必须是对子' };
    }

    if (cards.length === 3) {
      // Check if it's three of a kind
      if (cards.every(card => card.rank === cards[0].rank)) {
        return { isValid: true, cardType: '三张' };
      }
      return { isValid: false, error: '三张牌必须是同点数' };
    }

    // For more complex combinations, we would need more sophisticated validation
    if (cards.length >= 5) {
      return { isValid: true, cardType: `${cards.length}张牌型` };
    }

    return { isValid: false, error: '无效的牌型组合' };
  }, []);

  useEffect(() => {
    if (selectedCards.length > 0) {
      const result = validateCards(selectedCards);
      setValidationResult(result);
    } else {
      setValidationResult({ isValid: true });
    }
  }, [selectedCards, validateCards]);

  const handlePlayCards = () => {
    if (selectedCards.length === 0) return;
    
    const validation = validateCards(selectedCards);
    if (!validation.isValid) {
      setValidationResult(validation);
      return;
    }

    onPlayCards(selectedCards);
  };

  const handlePass = () => {
    onPass();
  };

  const isPlayDisabled = disabled || !canPlay || !isMyTurn || selectedCards.length === 0 || !validationResult.isValid;
  const isPassDisabled = disabled || !canPlay || !isMyTurn;

  return (
    <div className="bg-white border border-gray-300 rounded-lg p-4">
      {/* Turn Status and Timer */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-3">
          {isMyTurn ? (
            <div className="flex items-center">
              <div className="w-3 h-3 bg-green-500 rounded-full mr-2 animate-pulse"></div>
              <span className="text-green-700 font-medium">轮到您出牌</span>
            </div>
          ) : (
            <div className="flex items-center">
              <div className="w-3 h-3 bg-gray-400 rounded-full mr-2"></div>
              <span className="text-gray-600">等待其他玩家</span>
            </div>
          )}
        </div>
        
        <CountdownTimer 
          seconds={turnTimeoutSeconds} 
          isActive={isMyTurn && canPlay && !disabled} 
        />
      </div>

      {/* Card Validation Info */}
      {selectedCards.length > 0 && (
        <div className="mb-4">
          {validationResult.isValid ? (
            <div className="flex items-center text-green-700 bg-green-50 px-3 py-2 rounded">
              <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
              </svg>
              <span className="text-sm">
                {validationResult.cardType} - 可以出牌
              </span>
            </div>
          ) : (
            <div className="flex items-center text-red-700 bg-red-50 px-3 py-2 rounded">
              <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
              </svg>
              <span className="text-sm">{validationResult.error}</span>
            </div>
          )}
        </div>
      )}

      {/* Action Buttons */}
      <div className="flex space-x-3">
        <button
          onClick={handlePlayCards}
          disabled={isPlayDisabled}
          className={`
            flex-1 py-3 px-4 rounded-lg font-medium transition-all duration-200
            ${isPlayDisabled 
              ? 'bg-gray-200 text-gray-500 cursor-not-allowed' 
              : 'bg-blue-600 text-white hover:bg-blue-700 active:bg-blue-800 shadow-md hover:shadow-lg'
            }
          `}
        >
          <div className="flex items-center justify-center">
            <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
            </svg>
            出牌 {selectedCards.length > 0 && `(${selectedCards.length}张)`}
          </div>
        </button>

        <button
          onClick={handlePass}
          disabled={isPassDisabled}
          className={`
            flex-1 py-3 px-4 rounded-lg font-medium transition-all duration-200
            ${isPassDisabled 
              ? 'bg-gray-200 text-gray-500 cursor-not-allowed' 
              : 'bg-gray-600 text-white hover:bg-gray-700 active:bg-gray-800 shadow-md hover:shadow-lg'
            }
          `}
        >
          <div className="flex items-center justify-center">
            <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
            不出
          </div>
        </button>
      </div>

      {/* Help Text */}
      {!isMyTurn && (
        <div className="mt-3 text-center text-sm text-gray-500">
          等待其他玩家操作...
        </div>
      )}
      
      {isMyTurn && selectedCards.length === 0 && (
        <div className="mt-3 text-center text-sm text-gray-500">
          请选择要出的牌，或点击"不出"跳过
        </div>
      )}
    </div>
  );
};

export default GameControls;