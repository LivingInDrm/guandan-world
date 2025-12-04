import { useTributeStore } from '../store/tributeStore';

export function useTributeData() {
  const tributeView = useTributeStore(s => s.tributeView);
  if (!tributeView) return null;

  return {
    status: tributeView.status,
    tributeType: tributeView.tributeType,
    givers: tributeView.givers,
    receivers: tributeView.receivers,
    tributePairs: tributeView.tributePairs,
    poolCards: tributeView.poolCards,
    isImmune: tributeView.isImmune,
  };
}
