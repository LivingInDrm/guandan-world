import type { Story } from "@ladle/react";
import TeamLevelDisplay from "./TeamLevelDisplay";

export const BothAtLevel2: Story = () => (
  <div className="p-8 bg-green-100 relative w-64 h-32">
    <TeamLevelDisplay
      teamLevels={[2, 2]}
      currentLevel={2}
      currentPlayerSeat={0}
    />
  </div>
);

export const MyTeamLeading: Story = () => (
  <div className="p-8 bg-green-100 relative w-64 h-32">
    <TeamLevelDisplay
      teamLevels={[7, 5]}
      currentLevel={7}
      currentPlayerSeat={0}
    />
  </div>
);

export const OpponentLeading: Story = () => (
  <div className="p-8 bg-green-100 relative w-64 h-32">
    <TeamLevelDisplay
      teamLevels={[5, 8]}
      currentLevel={8}
      currentPlayerSeat={0}
    />
  </div>
);

export const HighLevels: Story = () => (
  <div className="p-8 bg-green-100 relative w-64 h-32">
    <TeamLevelDisplay
      teamLevels={[14, 13]}
      currentLevel={14}
      currentPlayerSeat={0}
    />
  </div>
);

export const OpponentCurrentLevel: Story = () => (
  <div className="p-8 bg-green-100 relative w-64 h-32">
    <TeamLevelDisplay
      teamLevels={[10, 11]}
      currentLevel={11}
      currentPlayerSeat={0}
    />
  </div>
);

export const Seat1ViewMyTeamIsTeam1: Story = () => (
  <div className="p-8 bg-green-100 relative w-64 h-32">
    <TeamLevelDisplay
      teamLevels={[5, 8]}
      currentLevel={8}
      currentPlayerSeat={1}
    />
  </div>
);

export const Seat2ViewMyTeamIsTeam0: Story = () => (
  <div className="p-8 bg-green-100 relative w-64 h-32">
    <TeamLevelDisplay
      teamLevels={[11, 9]}
      currentLevel={11}
      currentPlayerSeat={2}
    />
  </div>
);

export const Seat3ViewMyTeamIsTeam1: Story = () => (
  <div className="p-8 bg-green-100 relative w-64 h-32">
    <TeamLevelDisplay
      teamLevels={[6, 12]}
      currentLevel={12}
      currentPlayerSeat={3}
    />
  </div>
);
