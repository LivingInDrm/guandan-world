import { useState, useCallback } from 'react';
import { CreatePlayerForm } from './components/CreatePlayerForm';
import { PlayerList } from './components/PlayerList';

export default function App() {
  const [refreshKey, setRefreshKey] = useState(0);

  const handlePlayerCreated = useCallback(() => {
    setRefreshKey((k) => k + 1);
  }, []);

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto py-6 px-4">
          <h1 className="text-3xl font-bold text-gray-900">AI Player Admin</h1>
        </div>
      </header>
      <main className="max-w-7xl mx-auto py-6 px-4 space-y-6">
        <CreatePlayerForm onCreated={handlePlayerCreated} />
        <PlayerList key={refreshKey} />
      </main>
    </div>
  );
}
