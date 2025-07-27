import React from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from '../components/layout/Layout';
import LoginPage from '../components/auth/LoginPage';
import ProtectedRoute from '../components/auth/ProtectedRoute';
import RoomLobby from '../components/lobby/RoomLobby';
import RoomWaiting from '../components/room/RoomWaiting';

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
        path: 'lobby',
        element: (
          <ProtectedRoute>
            <RoomLobby />
          </ProtectedRoute>
        )
      },
      {
        path: 'room/:roomId',
        element: (
          <ProtectedRoute>
            <RoomWaiting />
          </ProtectedRoute>
        )
      }
    ]
  }
]);

export default router;