import { Router, Request, Response } from 'express';
import { PlayerManager } from './PlayerManager.js';
import type { AIPlayerConfig } from '../config.js';

export function createRouter(manager: PlayerManager): Router {
  const router = Router();

  router.get('/players', (_req: Request, res: Response) => {
    const players = manager.listPlayers();
    res.json(players);
  });

  router.post('/players', async (req: Request, res: Response) => {
    try {
      const config: AIPlayerConfig = req.body;
      if (!config.serverUrl || !config.roomCode) {
        res.status(400).json({ error: 'serverUrl and roomCode are required' });
        return;
      }
      if (!/^https?:\/\/.+/.test(config.serverUrl)) {
        res.status(400).json({ error: 'serverUrl must be a valid HTTP/HTTPS URL' });
        return;
      }
      if (config.autoRegister !== true && !config.username) {
        res.status(400).json({ error: 'username is required when autoRegister is not enabled' });
        return;
      }
      const id = await manager.createPlayer(config);
      const player = manager.getPlayer(id);
      res.status(201).json(player);
    } catch (error) {
      console.error('Failed to create player:', error);
      res.status(500).json({ error: 'Failed to create player' });
    }
  });

  router.get('/players/:id', (req: Request, res: Response) => {
    const player = manager.getPlayer(req.params.id);
    if (!player) {
      res.status(404).json({ error: 'Player not found' });
      return;
    }
    res.json(player);
  });

  router.delete('/players/:id', async (req: Request, res: Response) => {
    try {
      await manager.removePlayer(req.params.id);
      res.status(204).send();
    } catch (error) {
      if (error instanceof Error && error.message.includes('not found')) {
        res.status(404).json({ error: 'Player not found' });
        return;
      }
      console.error('Failed to remove player:', error);
      res.status(500).json({ error: 'Failed to remove player' });
    }
  });

  router.post('/players/:id/stop', async (req: Request, res: Response) => {
    try {
      await manager.stopPlayer(req.params.id);
      const player = manager.getPlayer(req.params.id);
      res.json(player);
    } catch (error) {
      if (error instanceof Error && error.message.includes('not found')) {
        res.status(404).json({ error: 'Player not found' });
        return;
      }
      console.error('Failed to stop player:', error);
      res.status(500).json({ error: 'Failed to stop player' });
    }
  });

  return router;
}
