import { useState } from 'react';
import { createPlayer, type AIPlayerConfig } from '../api';

interface CreatePlayerFormProps {
  onCreated: () => void;
}

export function CreatePlayerForm({ onCreated }: CreatePlayerFormProps) {
  const [serverUrl, setServerUrl] = useState('http://localhost:8080');
  const [roomCode, setRoomCode] = useState('');
  const [autoRegister, setAutoRegister] = useState(true);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [nickname, setNickname] = useState('');
  const [level, setLevel] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const config: AIPlayerConfig = {
      serverUrl,
      roomCode,
      autoRegister,
      level,
    };
    if (nickname) config.nickname = nickname;
    if (!autoRegister) {
      config.username = username;
      if (password) config.password = password;
    }

    try {
      await createPlayer(config);
      setRoomCode('');
      setNickname('');
      setUsername('');
      setPassword('');
      onCreated();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h2 className="text-xl font-semibold mb-4">Create AI Player</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Server URL</label>
            <input
              type="text"
              value={serverUrl}
              onChange={(e) => setServerUrl(e.target.value)}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 border px-3 py-2"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Room Code</label>
            <input
              type="text"
              value={roomCode}
              onChange={(e) => setRoomCode(e.target.value)}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 border px-3 py-2"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Nickname</label>
            <input
              type="text"
              value={nickname}
              onChange={(e) => setNickname(e.target.value)}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 border px-3 py-2"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Level</label>
            <select
              value={level}
              onChange={(e) => setLevel(Number(e.target.value))}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 border px-3 py-2"
            >
              <option value={1}>1 - Easy</option>
              <option value={2}>2 - Medium</option>
              <option value={3}>3 - Hard</option>
            </select>
          </div>
        </div>
        <div className="flex items-center">
          <input
            type="checkbox"
            id="autoRegister"
            checked={autoRegister}
            onChange={(e) => setAutoRegister(e.target.checked)}
            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
          <label htmlFor="autoRegister" className="ml-2 block text-sm text-gray-900">
            Auto Register
          </label>
        </div>
        {!autoRegister && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4 bg-gray-50 rounded-md">
            <div>
              <label className="block text-sm font-medium text-gray-700">Username</label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 border px-3 py-2"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 border px-3 py-2"
              />
            </div>
          </div>
        )}
        {error && <div className="text-red-600 text-sm">{error}</div>}
        <button
          type="submit"
          disabled={loading}
          className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {loading ? 'Creating...' : 'Create Player'}
        </button>
      </form>
    </div>
  );
}
