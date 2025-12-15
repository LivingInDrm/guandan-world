import { AIPlayer, AIState } from '../AIPlayer.js';
import type { AIPlayerInfo } from '../AIPlayer.js';
import type { AIPlayerConfig } from '../config.js';

export class PlayerManager {
  private players: Map<string, AIPlayer> = new Map();

  async createPlayer(config: AIPlayerConfig): Promise<string> {
    const player = new AIPlayer(config);
    const id = player.getId();
    this.players.set(id, player);
    player.start().catch((error) => {
      console.error(`Player ${id} failed to start:`, error);
    });
    return id;
  }

  async stopPlayer(id: string): Promise<void> {
    const player = this.players.get(id);
    if (!player) {
      throw new Error(`Player ${id} not found`);
    }
    await player.stop();
  }

  async removePlayer(id: string): Promise<void> {
    const player = this.players.get(id);
    if (!player) {
      throw new Error(`Player ${id} not found`);
    }
    if (player.getState() !== AIState.STOPPED) {
      await player.stop();
    }
    this.players.delete(id);
  }

  getPlayer(id: string): AIPlayerInfo | undefined {
    const player = this.players.get(id);
    return player?.getInfo();
  }

  listPlayers(): AIPlayerInfo[] {
    return Array.from(this.players.values()).map((player) => player.getInfo());
  }

  async stopAll(): Promise<void> {
    const stopPromises = Array.from(this.players.values()).map((player) =>
      player.stop().catch((error) => {
        console.error(`Failed to stop player ${player.getId()}:`, error);
      })
    );
    await Promise.all(stopPromises);
  }
}
