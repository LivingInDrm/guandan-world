export interface AIPlayerConfig {
  serverUrl: string;
  roomCode: string;
  username?: string;
  password?: string;
  autoRegister?: boolean;
  nickname?: string;
  level?: number;
}

export interface AIPlayerInfo {
  id: string;
  config: AIPlayerConfig;
  state: string;
  username: string | null;
  roomId: string | null;
  playerSeat: number | null;
  createdAt: string;
}

const API_BASE = '/api';

export async function listPlayers(): Promise<AIPlayerInfo[]> {
  const res = await fetch(`${API_BASE}/players`);
  if (!res.ok) throw new Error('Failed to list players');
  return res.json();
}

export async function createPlayer(config: AIPlayerConfig): Promise<AIPlayerInfo> {
  const res = await fetch(`${API_BASE}/players`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(config),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Failed to create player');
  }
  return res.json();
}

export async function stopPlayer(id: string): Promise<AIPlayerInfo> {
  const res = await fetch(`${API_BASE}/players/${id}/stop`, { method: 'POST' });
  if (!res.ok) throw new Error('Failed to stop player');
  return res.json();
}

export async function removePlayer(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/players/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to remove player');
}
