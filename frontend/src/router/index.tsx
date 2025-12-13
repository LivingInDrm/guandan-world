import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from '../components/layout/Layout';
import LoginPage from '../components/auth/LoginPage';
import ProtectedRoute from '../components/auth/ProtectedRoute';
import RoomLobby from '../components/lobby/RoomLobby';
import GamePage from '../components/game/GamePage';
import JoinByInvite from '../components/lobby/JoinByInvite';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <Navigate to="/lobby" replace />
      },
      {
        path: 'login',
        element: <LoginPage />
      },
      {
        path: 'join',
        element: <JoinByInvite />
      },
      {
        path: 'lobby',
        element: (
          <ProtectedRoute>
            <RoomLobby />
          </ProtectedRoute>
        )
      },
      {
        path: 'game/:roomId',
        element: (
          <ProtectedRoute>
            <GamePage />
          </ProtectedRoute>
        )
      }
    ]
  }
]);

export default router;
