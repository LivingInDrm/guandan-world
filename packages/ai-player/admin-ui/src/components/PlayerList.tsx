import { useState, useEffect, useCallback } from 'react';
import { listPlayers, stopPlayer, removePlayer, type AIPlayerInfo } from '../api';

export function PlayerList() {
  const [players, setPlayers] = useState<AIPlayerInfo[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchPlayers = useCallback(async () => {
    try {
      const data = await listPlayers();
      setPlayers(data);
    } catch (err) {
      console.error('Failed to fetch players:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPlayers();
    const interval = setInterval(fetchPlayers, 2000);
    return () => clearInterval(interval);
  }, [fetchPlayers]);

  const handleStop = async (id: string) => {
    try {
      await stopPlayer(id);
      fetchPlayers();
    } catch (err) {
      console.error('Failed to stop player:', err);
    }
  };

  const handleRemove = async (id: string) => {
    try {
      await removePlayer(id);
      fetchPlayers();
    } catch (err) {
      console.error('Failed to remove player:', err);
    }
  };

  const getStateColor = (state: string) => {
    switch (state) {
      case 'PLAYING':
        return 'bg-green-100 text-green-800';
      case 'WAITING_IN_ROOM':
        return 'bg-yellow-100 text-yellow-800';
      case 'STOPPED':
        return 'bg-gray-100 text-gray-800';
      case 'IDLE':
        return 'bg-blue-100 text-blue-800';
      default:
        return 'bg-purple-100 text-purple-800';
    }
  };

  if (loading) {
    return (
      <div className="bg-white shadow rounded-lg p-6">
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow rounded-lg overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-200">
        <h2 className="text-xl font-semibold">AI Players ({players.length})</h2>
      </div>
      {players.length === 0 ? (
        <div className="p-6 text-gray-500 text-center">No AI players</div>
      ) : (
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Username</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Room</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">State</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {players.map((player) => (
              <tr key={player.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                  {player.id.slice(0, 8)}...
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {player.username || '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {player.roomId ? `${player.roomId.slice(0, 8)}... (seat ${player.playerSeat})` : '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStateColor(player.state)}`}>
                    {player.state}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                  {player.state !== 'STOPPED' && (
                    <button
                      onClick={() => handleStop(player.id)}
                      className="text-yellow-600 hover:text-yellow-900"
                    >
                      Stop
                    </button>
                  )}
                  <button
                    onClick={() => handleRemove(player.id)}
                    className="text-red-600 hover:text-red-900"
                  >
                    Remove
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
