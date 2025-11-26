import React, { useState, useEffect } from 'react';
import type { TributeFlowProps, UIPhase } from './types';
import TributePhaseContent from './TributePhaseContent';
import { TributeStatus } from '../../../types';

const TributeFlow: React.FC<TributeFlowProps> = (props) => {
  const { tributePhase } = props;
  const [uiPhase, setUIPhase] = useState<UIPhase>('START');

  // Initial State Check: If we load/mount and the game is already deep in tribute phase
  useEffect(() => {
    if (tributePhase.status === TributeStatus.RETURNING) {
      setUIPhase('RETURNING');
    } else if (tributePhase.status === TributeStatus.FINISHED) {
      setUIPhase('FINISHED');
    } else if (tributePhase.status === TributeStatus.SELECTING) {
      setUIPhase('SELECTING');
    }
  }, [tributePhase.status]);

  // Automatic Transitions
  useEffect(() => {
    if (uiPhase === 'START') {
      const timer = setTimeout(() => setUIPhase('IMMUNITY_CHECK'), 3000);
      return () => clearTimeout(timer);
    }
  }, [uiPhase]);

  useEffect(() => {
    if (uiPhase === 'IMMUNITY_CHECK') {
      const delay = tributePhase.is_immune ? 3000 : 2000;
      const timer = setTimeout(() => {
        setUIPhase(tributePhase.is_immune ? 'FINISHED' : 'SUBMITTING');
      }, delay);
      return () => clearTimeout(timer);
    }
  }, [uiPhase, tributePhase.is_immune]);

  useEffect(() => {
    if (uiPhase === 'SUBMITTING') {
      // Transition to SELECTING after animation
      const timer = setTimeout(() => setUIPhase('SELECTING'), 2000);
      return () => clearTimeout(timer);
    }
  }, [uiPhase]);

  // Backend Status Synchronization (State Jumping)
  useEffect(() => {
    // If backend is RETURNING, jump to RETURNING immediately (skip others)
    if (tributePhase.status === TributeStatus.RETURNING && 
        uiPhase !== 'RETURNING' && 
        uiPhase !== 'FINISHED') {
      setUIPhase('RETURNING');
    }

    // If backend is FINISHED, jump to FINISHED immediately
    if (tributePhase.status === TributeStatus.FINISHED && 
        uiPhase !== 'FINISHED') {
      setUIPhase('FINISHED');
    }
    
    // If backend is SELECTING and we are in SELECTING (already there), do nothing.
    // If backend is SELECTING and we are in START/IMMUNITY/SUBMITTING, we let the timers play out
    // so the user sees the flow.
  }, [tributePhase.status, uiPhase]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 pointer-events-none">
      <div className="w-full max-w-4xl h-[600px] pointer-events-auto">
        <TributePhaseContent 
          phase={uiPhase} 
          {...props} 
        />
      </div>
    </div>
  );
};

export default TributeFlow;
