import { RouterProvider } from 'react-router-dom';
import { router } from './router';
import { useGameService } from './hooks/useGameService';
import AuthProvider from './components/auth/AuthProvider';
import './App.css';

function App() {
  useGameService();

  return (
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  );
}

export default App;
