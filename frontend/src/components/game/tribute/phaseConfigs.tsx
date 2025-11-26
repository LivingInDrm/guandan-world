import type { UIPhase, PhaseConfig } from './types';
import type { TributePhase } from '../../../types';
import { 
  StartContent, 
  ImmunityCheckContent, 
  SubmittingContent, 
  SelectingContent, 
  ReturningContent, 
  FinishedContent 
} from './contents';

export const getPhaseConfig = (phase: UIPhase, tributePhase: TributePhase): PhaseConfig => {
  const configs: Record<UIPhase, PhaseConfig> = {
    START: {
      title: '上贡阶段开始',
      icon: '🎴',
      duration: 3000,
      renderContent: (props) => <StartContent {...props} />
    },
    IMMUNITY_CHECK: {
      title: '抗贡检测',
      icon: '🛡️',
      duration: tributePhase.is_immune ? 3000 : 2000,
      renderContent: (props) => <ImmunityCheckContent {...props} />
    },
    SUBMITTING: {
      title: '进贡',
      icon: '⬆️',
      duration: 2000, // Allow time for "flying cards" animation
      renderContent: (props) => <SubmittingContent {...props} />
    },
    SELECTING: {
      title: '选贡',
      icon: '👆',
      // No fixed duration for selecting phase, unless it's auto-selection which we might handle in the component or Flow
      renderContent: (props) => <SelectingContent {...props} />
    },
    RETURNING: {
      title: '还贡',
      icon: '⬇️',
      renderContent: (props) => <ReturningContent {...props} />
    },
    FINISHED: {
      title: '上贡完成',
      icon: '✅',
      duration: 3000,
      renderContent: (props) => <FinishedContent {...props} />
    }
  };

  return configs[phase];
};
