import type { Story } from '@ladle/react';
import GameTable, { type PlayerPosition } from './GameTable';
import PlayerCard from './PlayerCard';
import type { Player } from '../../types';

const makePlayer = (id: string, username: string, seat: number): Player => ({
  id,
  username,
  seat,
  online: true,
  auto_play: false,
});

const mockPlayers: Player[] = [
  makePlayer('p1', '玩家A', 0),
  makePlayer('p2', '玩家B', 1),
  makePlayer('p3', '玩家C', 2),
  makePlayer('p4', '玩家D', 3),
];

const partialPlayers: (Player | null)[] = [
  mockPlayers[0],
  null,
  mockPlayers[2],
  null,
];

const emptyPlayers: (Player | null)[] = [null, null, null, null];

const renderPlayer = (
  player: Player | null,
  position: PlayerPosition,
  _seat: number,
  isHighlighted = false
) => (
  <PlayerCard player={player} position={position} isHighlighted={isHighlighted} />
);

const renderCenter = (content?: string) => () => (
  <div className="flex items-center justify-center h-full text-muted-foreground">
    {content || '中心区域'}
  </div>
);

export const EmptyTable: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={emptyPlayers}
      currentPlayerSeat={0}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={renderCenter()}
    />
  </div>
);
EmptyTable.meta = {
  description: '空桌子 - 无玩家加入',
};

export const PartialPlayers: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={partialPlayers}
      currentPlayerSeat={0}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={renderCenter()}
    />
  </div>
);
PartialPlayers.meta = {
  description: '部分玩家 - 只有2名玩家加入',
};

export const AllPlayersJoined: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={mockPlayers}
      currentPlayerSeat={0}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={renderCenter()}
    />
  </div>
);
AllPlayersJoined.meta = {
  description: '满员桌子 - 四名玩家都已加入',
};

export const ViewFromSeat1: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={mockPlayers}
      currentPlayerSeat={1}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={renderCenter()}
    />
  </div>
);
ViewFromSeat1.meta = {
  description: '座位1视角 - 玩家B在底部',
};

export const ViewFromSeat2: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={mockPlayers}
      currentPlayerSeat={2}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={renderCenter()}
    />
  </div>
);
ViewFromSeat2.meta = {
  description: '座位2视角 - 玩家C在底部',
};

export const ViewFromSeat3: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={mockPlayers}
      currentPlayerSeat={3}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={renderCenter()}
    />
  </div>
);
ViewFromSeat3.meta = {
  description: '座位3视角 - 玩家D在底部',
};

export const WithHighlightedPlayer: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={mockPlayers}
      currentPlayerSeat={0}
      renderPlayer={(player, position, seat) =>
        renderPlayer(player, position, seat, seat === 1)
      }
      renderCenter={renderCenter('玩家B回合')}
    />
  </div>
);
WithHighlightedPlayer.meta = {
  description: '高亮玩家 - 玩家B被高亮显示（当前回合）',
};

export const WithTopLeftSlot: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={mockPlayers}
      currentPlayerSeat={0}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={renderCenter()}
      topLeftSlot={
        <div className="absolute top-4 left-4 px-3 py-1.5 rounded-lg bg-card/80 backdrop-blur-sm text-sm text-foreground">
          房间: 12345
        </div>
      }
    />
  </div>
);
WithTopLeftSlot.meta = {
  description: '顶部左侧插槽 - 显示房间信息',
};

export const WithCenterContent: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={mockPlayers}
      currentPlayerSeat={0}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={() => (
        <div className="flex flex-col items-center justify-center h-full gap-2">
          <div className="text-lg font-semibold text-foreground">游戏进行中</div>
          <div className="text-sm text-muted-foreground">当前级别: 10</div>
          <div className="flex gap-2 mt-2">
            <div className="px-2 py-1 text-xs rounded bg-primary/20 text-primary">我方: 10级</div>
            <div className="px-2 py-1 text-xs rounded bg-destructive/20 text-destructive">对方: 8级</div>
          </div>
        </div>
      )}
    />
  </div>
);
WithCenterContent.meta = {
  description: '中心区域内容 - 显示游戏状态信息',
};

export const CustomClassName: Story = () => (
  <div className="p-4 bg-gray-800">
    <GameTable
      players={mockPlayers}
      currentPlayerSeat={0}
      renderPlayer={(player, position, seat) => renderPlayer(player, position, seat)}
      renderCenter={renderCenter()}
      className="h-[40rem]"
    />
  </div>
);
CustomClassName.meta = {
  description: '自定义高度 - 使用 className 调整尺寸',
};
