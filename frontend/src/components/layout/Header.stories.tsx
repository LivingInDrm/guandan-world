import type { Story } from "@ladle/react";
import { MemoryRouter, useLocation } from "react-router-dom";
import { useEffect } from "react";
import Header from "./Header";
import { useAuthStore } from "../../store/authStore";

const MockAuthProvider: React.FC<{ children: React.ReactNode; isAuthenticated?: boolean }> = ({ 
  children, 
  isAuthenticated = true 
}) => {
  const { login, logout } = useAuthStore();

  useEffect(() => {
    if (isAuthenticated) {
      login(
        { id: "1", username: "testuser", nickname: "测试用户", online: true },
        { access_token: "mock-token", refresh_token: "mock-refresh", expires_at: new Date(Date.now() + 3600000).toISOString(), user_id: "1" }
      );
    } else {
      logout();
    }
    return () => logout();
  }, [isAuthenticated, login, logout]);

  return <>{children}</>;
};

const LocationDisplay = () => {
  const location = useLocation();
  return (
    <div className="fixed bottom-4 left-4 bg-black/80 text-white px-3 py-2 rounded text-sm">
      当前路径: {location.pathname}
    </div>
  );
};

export const Default: Story = () => (
  <MemoryRouter initialEntries={["/lobby"]}>
    <MockAuthProvider>
      <Header />
      <LocationDisplay />
    </MockAuthProvider>
  </MemoryRouter>
);

export const WithLongNickname: Story = () => (
  <MemoryRouter initialEntries={["/lobby"]}>
    <MockAuthProvider>
      <Header />
      <LocationDisplay />
    </MockAuthProvider>
  </MemoryRouter>
);
WithLongNickname.decorators = [
  (Story) => {
    const { login } = useAuthStore();
    useEffect(() => {
      login(
        { id: "2", username: "longuser", nickname: "这是一个非常长的用户昵称测试", online: true },
        { access_token: "mock-token", refresh_token: "mock-refresh", expires_at: new Date(Date.now() + 3600000).toISOString(), user_id: "2" }
      );
    }, [login]);
    return <Story />;
  },
];

export const NotAuthenticated: Story = () => (
  <MemoryRouter>
    <MockAuthProvider isAuthenticated={false}>
      <Header />
    </MockAuthProvider>
  </MemoryRouter>
);
