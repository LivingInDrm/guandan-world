import React from 'react';
import type { TributeFlowProps } from './types';
import CardDisplay from '../CardDisplay';

export const StartContent: React.FC<TributeFlowProps> = () => {
  return (
    <div className="text-center p-8">
      <div className="text-4xl mb-4">🎴</div>
      <h2 className="text-2xl font-bold mb-2">上贡阶段开始</h2>
      <p className="text-gray-400">根据上局排名进行进贡</p>
    </div>
  );
};

export const ImmunityCheckContent: React.FC<TributeFlowProps> = ({ tributePhase }) => {
  const isImmune = tributePhase.is_immune;
  return (
    <div className="text-center p-8">
      <div className="text-4xl mb-4">{isImmune ? '🛡️' : '⚔️'}</div>
      <h2 className="text-2xl font-bold mb-2">
        {isImmune ? '抗贡成功' : '未触发抗贡'}
      </h2>
      {isImmune && <p className="text-green-400">双大王在同一方，免除上贡！</p>}
      {!isImmune && <p className="text-gray-400">即将开始进贡流程...</p>}
    </div>
  );
};

export const SubmittingContent: React.FC<TributeFlowProps> = () => {
  return (
    <div className="text-center p-8">
      <h2 className="text-xl font-bold mb-4">进贡中...</h2>
      <div className="flex justify-center items-center gap-8">
        <div className="animate-pulse">等待贡牌提交...</div>
      </div>
    </div>
  );
};

export const SelectingContent: React.FC<TributeFlowProps> = ({ 
  tributePhase,
  currentPlayerSeat,
  onSelectTribute
}) => {
  const isSelecting = tributePhase.selecting_player === currentPlayerSeat;
  const poolCards = tributePhase.pool_cards || [];

  return (
    <div className="text-center p-4">
      <h2 className="text-xl font-bold mb-4">
        {isSelecting ? '请选择一张贡牌' : '等待对方选牌'}
      </h2>
      
      <div className="flex justify-center gap-4 mt-8">
        {poolCards.map((card) => (
          <button
            key={card.deckIndex}
            type="button"
            onClick={() => isSelecting && onSelectTribute(card.deckIndex)}
            className={`transform transition-all bg-transparent border-0 p-0 ${isSelecting ? 'cursor-pointer hover:-translate-y-2' : 'cursor-default'}`}
            aria-label={`选择贡牌 ${card.rank}`}
            disabled={!isSelecting}
          >
            <CardDisplay card={card} />
          </button>
        ))}
      </div>
    </div>
  );
};

export const ReturningContent: React.FC<TributeFlowProps> = ({ 
  tributePhase,
  currentPlayerSeat,
  selectedCards,
  onReturnTribute
}) => {
  // Check if we are a receiver
  const isReceiver = Object.values(tributePhase.tribute_map || {}).includes(currentPlayerSeat);
  
  // Check if we have already returned (if return_cards has an entry for us)
  const hasReturned = tributePhase.return_cards && (currentPlayerSeat in tributePhase.return_cards);

  const canReturn = isReceiver && !hasReturned && selectedCards.length === 1;

  return (
    <div className="text-center p-4">
      <h2 className="text-xl font-bold mb-4">还贡阶段</h2>
      
      {isReceiver && !hasReturned ? (
        <>
          <p className="mb-4">请从手牌中点击选择一张牌，然后点击确认</p>
          <div className="max-w-2xl mx-auto bg-gray-800/50 p-4 rounded-lg flex flex-col items-center gap-4">
            <div className="text-sm text-gray-400" aria-live="polite">
              {selectedCards.length === 0 ? "请选择一张手牌" : 
               selectedCards.length === 1 ? "点击下方按钮确认" : "只能选择一张牌"}
            </div>
            
            <button
              type="button"
              onClick={() => {
                if (canReturn) {
                  onReturnTribute(selectedCards[0].deckIndex);
                }
              }}
              disabled={!canReturn}
              className={`px-6 py-2 rounded font-bold transition-colors ${
                canReturn 
                  ? 'bg-blue-500 hover:bg-blue-600 text-white' 
                  : 'bg-gray-600 text-gray-400 cursor-not-allowed'
              }`}
            >
              确认还贡
            </button>
          </div>
        </>
      ) : (
        <p className="text-gray-400">
          {hasReturned ? "已还贡，等待其他人..." : "等待进贡者还贡..."}
        </p>
      )}
    </div>
  );
};

export const FinishedContent: React.FC<TributeFlowProps> = () => {
  return (
    <div className="text-center p-8">
      <h2 className="text-2xl font-bold mb-4">上贡完成</h2>
      <div className="grid grid-cols-1 gap-4 max-w-md mx-auto">
        <div className="bg-gray-800 p-4 rounded">
           准备开始游戏...
        </div>
      </div>
    </div>
  );
};
