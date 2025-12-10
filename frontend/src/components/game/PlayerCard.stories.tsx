import type { Story } from '@ladle/react';
import PlayerCard from './PlayerCard';
import type { Player } from '../../types';

function createMockPlayer(username: string, seat: number): Player {
  return {
    id: `player-${seat}`,
    username,
    seat,
    online: true,
    auto_play: false,
  };
}

const ContainerWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <div className="relative w-[600px] h-[400px] bg-muted/30 rounded-lg border border-dashed border-border">
    {children}
  </div>
);

export const Default: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={createMockPlayer('张三', 0)}
      position="bottom"
    />
  </ContainerWrapper>
);
Default.meta = {
  description: '默认状态 - 底部位置玩家',
};

export const EmptySeat: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={null}
      position="bottom"
    />
  </ContainerWrapper>
);
EmptySeat.meta = {
  description: '空座位状态',
};

export const BottomPosition: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={createMockPlayer('底部玩家', 0)}
      position="bottom"
    />
  </ContainerWrapper>
);
BottomPosition.meta = {
  description: '底部位置',
};

export const LeftPosition: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={createMockPlayer('左侧玩家', 1)}
      position="left"
    />
  </ContainerWrapper>
);
LeftPosition.meta = {
  description: '左侧位置',
};

export const TopPosition: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={createMockPlayer('顶部玩家', 2)}
      position="top"
    />
  </ContainerWrapper>
);
TopPosition.meta = {
  description: '顶部位置',
};

export const RightPosition: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={createMockPlayer('右侧玩家', 3)}
      position="right"
    />
  </ContainerWrapper>
);
RightPosition.meta = {
  description: '右侧位置',
};

export const AllPositions: Story = () => (
  <ContainerWrapper>
    <PlayerCard player={createMockPlayer('底部', 0)} position="bottom" />
    <PlayerCard player={createMockPlayer('左侧', 1)} position="left" />
    <PlayerCard player={createMockPlayer('顶部', 2)} position="top" />
    <PlayerCard player={createMockPlayer('右侧', 3)} position="right" />
  </ContainerWrapper>
);
AllPositions.meta = {
  description: '四个位置同时展示',
};

export const Highlighted: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={createMockPlayer('当前玩家', 0)}
      position="bottom"
      isHighlighted={true}
    />
  </ContainerWrapper>
);
Highlighted.meta = {
  description: '高亮状态 - 当前玩家回合',
};

export const WithStatusSlotLeft: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={createMockPlayer('左侧玩家', 1)}
      position="left"
      statusSlot={
        <div className="px-2 py-1 bg-accent text-accent-foreground rounded text-xs">
          出牌中...
        </div>
      }
    />
  </ContainerWrapper>
);
WithStatusSlotLeft.meta = {
  description: '左侧位置带状态槽',
};

export const WithStatusSlotRight: Story = () => (
  <ContainerWrapper>
    <PlayerCard
      player={createMockPlayer('右侧玩家', 3)}
      position="right"
      statusSlot={
        <div className="px-2 py-1 bg-warning text-warning-foreground rounded text-xs">
          10s
        </div>
      }
    />
  </ContainerWrapper>
);
WithStatusSlotRight.meta = {
  description: '右侧位置带状态槽（倒计时）',
};

export const DifferentPlayers: Story = () => (
  <div className="flex gap-8 p-4">
    <div className="relative w-[200px] h-[150px] bg-muted/30 rounded-lg border border-dashed border-border">
      <PlayerCard player={createMockPlayer('Alice', 0)} position="bottom" />
    </div>
    <div className="relative w-[200px] h-[150px] bg-muted/30 rounded-lg border border-dashed border-border">
      <PlayerCard player={createMockPlayer('Bob', 1)} position="bottom" />
    </div>
    <div className="relative w-[200px] h-[150px] bg-muted/30 rounded-lg border border-dashed border-border">
      <PlayerCard player={createMockPlayer('王五', 2)} position="bottom" />
    </div>
  </div>
);
DifferentPlayers.meta = {
  description: '不同用户名展示不同头像',
};
