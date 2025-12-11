import type { Story } from "@ladle/react";
import { MemoryRouter } from "react-router-dom";
import { useAuthStore } from "../../store/authStore";
import ProfileDialog from "./ProfileDialog";

const Wrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <MemoryRouter>
    <div className="min-h-screen bg-primitive-neutral-900 flex items-center justify-center p-8">
      {children}
    </div>
  </MemoryRouter>
);

export const Default: Story = () => {
  useAuthStore.setState({
    user: {
      id: "user-1",
      username: "testuser",
      nickname: "测试玩家",
      avatar_key: "avatars/default.png",
      online: true,
    },
    isAuthenticated: true,
  });

  return (
    <Wrapper>
      <ProfileDialog open={true} onOpenChange={() => {}} />
    </Wrapper>
  );
};

export const WithoutAvatar: Story = () => {
  useAuthStore.setState({
    user: {
      id: "user-2",
      username: "player2",
      nickname: "无头像用户",
      avatar_key: undefined,
      online: true,
    },
    isAuthenticated: true,
  });

  return (
    <Wrapper>
      <ProfileDialog open={true} onOpenChange={() => {}} />
    </Wrapper>
  );
};

export const WithoutNickname: Story = () => {
  useAuthStore.setState({
    user: {
      id: "user-3",
      username: "player3",
      nickname: undefined,
      avatar_key: "avatars/sample.png",
      online: true,
    },
    isAuthenticated: true,
  });

  return (
    <Wrapper>
      <ProfileDialog open={true} onOpenChange={() => {}} />
    </Wrapper>
  );
};

export const NewUser: Story = () => {
  useAuthStore.setState({
    user: {
      id: "user-4",
      username: "newplayer",
      nickname: undefined,
      avatar_key: undefined,
      online: true,
    },
    isAuthenticated: true,
  });

  return (
    <Wrapper>
      <ProfileDialog open={true} onOpenChange={() => {}} />
    </Wrapper>
  );
};
