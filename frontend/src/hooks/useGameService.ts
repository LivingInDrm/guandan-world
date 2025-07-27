import { useEffect } from 'react';
import { gameService } from '../services/gameService';

export const useGameService = () => {
  useEffect(() => {
    gameService.initialize();
  }, []);

  return gameService;
};