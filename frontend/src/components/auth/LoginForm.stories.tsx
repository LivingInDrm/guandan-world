import type { Story } from "@ladle/react";
import { MemoryRouter } from "react-router-dom";
import LoginForm from "./LoginForm";

export const Default: Story = () => (
  <MemoryRouter>
    <div className="min-h-screen bg-background flex items-center justify-center p-8">
      <div className="max-w-md w-full">
        <LoginForm />
      </div>
    </div>
  </MemoryRouter>
);

export const WithoutCard: Story = () => (
  <MemoryRouter>
    <div className="p-8 max-w-md bg-card rounded-ds-lg shadow-ds-elevation-2">
      <LoginForm showCard={false} />
    </div>
  </MemoryRouter>
);

export const CustomTitle: Story = () => (
  <MemoryRouter>
    <div className="min-h-screen bg-background flex items-center justify-center p-8">
      <div className="max-w-md w-full">
        <LoginForm title="欢迎回来" />
      </div>
    </div>
  </MemoryRouter>
);
