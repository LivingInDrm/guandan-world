import { create } from 'zustand';
import type { TributeView } from '../types/proto';

interface TributeState {
  tributeView: TributeView | null;
}

interface TributeActions {
  setTributeView: (view: TributeView | null) => void;
  reset: () => void;
}

type TributeStore = TributeState & TributeActions;

export const useTributeStore = create<TributeStore>((set) => ({
  tributeView: null,
  setTributeView: (tributeView) => set({ tributeView }),
  reset: () => set({ tributeView: null }),
}));
