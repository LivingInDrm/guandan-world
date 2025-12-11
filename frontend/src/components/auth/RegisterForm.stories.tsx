import type { Story } from "@ladle/react";
import { MemoryRouter } from "react-router-dom";
import RegisterForm from "./RegisterForm";

export const Default: Story = () => (
  <MemoryRouter>
    <div className="min-h-screen bg-surface-base flex items-center justify-center p-8">
      <div className="max-w-md w-full">
        <RegisterForm />
      </div>
    </div>
  </MemoryRouter>
);

export const WithoutCard: Story = () => (
  <MemoryRouter>
    <div className="p-8 max-w-md bg-surface-elevated rounded-lg shadow-elevation-2">
      <RegisterForm showCard={false} />
    </div>
  </MemoryRouter>
);

export const CustomTitle: Story = () => (
  <MemoryRouter>
    <div className="min-h-screen bg-surface-base flex items-center justify-center p-8">
      <div className="max-w-md w-full">
        <RegisterForm title="创建账号" />
      </div>
    </div>
  </MemoryRouter>
);
