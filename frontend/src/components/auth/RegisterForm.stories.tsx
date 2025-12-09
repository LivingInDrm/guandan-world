import type { Story } from "@ladle/react";
import { MemoryRouter } from "react-router-dom";
import RegisterForm from "./RegisterForm";

const Wrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <MemoryRouter>
    <div className="min-h-screen bg-table-900 flex items-center justify-center p-8">
      <div className="w-full max-w-md bg-table-800 rounded-sm p-6 shadow-lg">
        <h2 className="text-xl font-bold text-gold-400 mb-6 text-center">
          注册账号
        </h2>
        {children}
      </div>
    </div>
  </MemoryRouter>
);

export const Default: Story = () => (
  <Wrapper>
    <RegisterForm />
  </Wrapper>
);

export const FormOnly: Story = () => (
  <MemoryRouter>
    <div className="p-8 max-w-md">
      <RegisterForm />
    </div>
  </MemoryRouter>
);
