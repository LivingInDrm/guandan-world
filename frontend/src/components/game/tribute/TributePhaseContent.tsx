import React from 'react';
import type { PhaseContentProps } from './types';
import { getPhaseConfig } from './phaseConfigs';

const TributePhaseContent: React.FC<PhaseContentProps> = ({ 
  phase, 
  tributePhase, 
  players, 
  ...props 
}) => {
  const config = getPhaseConfig(phase, tributePhase);
  
  return (
    <div className="tribute-phase-content w-full h-full flex flex-col items-center justify-center bg-black/50 backdrop-blur-sm rounded-xl border border-white/10 shadow-2xl overflow-hidden transition-all duration-500">
      {/* 头部 */}
      <div className="phase-header flex items-center gap-3 p-4 border-b border-white/10 w-full bg-white/5">
        {config.icon && <span className="text-2xl" aria-hidden="true">{config.icon}</span>}
        <h2 className="text-xl font-bold text-white/90">{config.title}</h2>
      </div>
      
      {/* 主体内容 */}
      <div className="phase-body flex-1 w-full p-4 overflow-y-auto">
        {config.renderContent({ tributePhase, players, ...props })}
      </div>
    </div>
  );
};

export default TributePhaseContent;
