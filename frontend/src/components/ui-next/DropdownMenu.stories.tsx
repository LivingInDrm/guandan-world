import type { Story } from "@ladle/react";
import { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuCheckboxItem,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from "./DropdownMenu";
import { Button } from "./Button";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center min-h-[300px]">
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button>打开菜单</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem>个人资料</DropdownMenuItem>
        <DropdownMenuItem>游戏设置</DropdownMenuItem>
        <DropdownMenuItem>帮助中心</DropdownMenuItem>
        <DropdownMenuItem>退出登录</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
);

export const WithSeparator: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center min-h-[300px]">
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button>用户菜单</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-48">
        <DropdownMenuLabel>我的账户</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem>个人资料</DropdownMenuItem>
        <DropdownMenuItem>游戏记录</DropdownMenuItem>
        <DropdownMenuItem>成就徽章</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem>设置</DropdownMenuItem>
        <DropdownMenuItem>帮助</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem>退出登录</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
);

export const WithCheckbox: Story = () => {
  const [showStatus, setShowStatus] = useState(true);
  const [showNotification, setShowNotification] = useState(false);
  const [autoPlay, setAutoPlay] = useState(true);

  return (
    <div className="p-8 bg-gray-800 flex items-center justify-center min-h-[300px]">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button>显示设置</Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-48">
          <DropdownMenuLabel>显示选项</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuCheckboxItem
            checked={showStatus}
            onCheckedChange={setShowStatus}
          >
            显示在线状态
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            checked={showNotification}
            onCheckedChange={setShowNotification}
          >
            消息通知
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            checked={autoPlay}
            onCheckedChange={setAutoPlay}
          >
            自动出牌
          </DropdownMenuCheckboxItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};

export const WithRadio: Story = () => {
  const [theme, setTheme] = useState("dark");

  return (
    <div className="p-8 bg-gray-800 flex items-center justify-center min-h-[300px]">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button>主题设置</Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-48">
          <DropdownMenuLabel>选择主题</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuRadioGroup value={theme} onValueChange={setTheme}>
            <DropdownMenuRadioItem value="light">浅色模式</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="dark">深色模式</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="system">跟随系统</DropdownMenuRadioItem>
          </DropdownMenuRadioGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};
