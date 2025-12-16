import { RouterProvider } from 'react-router-dom';
import { router } from './router';
import { useGameService } from './hooks/useGameService';
import AuthProvider from './components/auth/AuthProvider';
import { OrientationGuard } from './components/ui';
import './App.css';

function App() {
  // Initialize game service
  useGameService();

  return (
    <AuthProvider>
      <OrientationGuard>
        <RouterProvider router={router} />
      </OrientationGuard>
    </AuthProvider>
  );
}

export default App;
