import React from 'react';
import { DealResult as DealResultType, Player, VictoryType } from '../../types';

interface DealResultProps {
  dealResult: DealResultType;
  players: Player[];
  teamLevels: [number, number];
  onContinue: () => void;
  onExit: () => void;
  isMatchFinished: boolean;
}

const DealResult: React.FC<DealResultProps> = ({
  dealResult,
  players,
  teamLevels,
  onContinue,
  onExit,
  isMatchFinished
}) => {
  // Helper function to get team for player
  const getTeamForPlayer = (playerSeat: number): number => {
    return playerSeat % 2; // Team 0: seats 0,2; Team 1: seats 1,3
  };

  // Helper function to get player by seat
  const getPlayerBySeat = (seat: number): Player | undefined => {
    return players.find(p => p.seat === seat);
  };

  // Group rankings by team
  const getTeamRankings = () => {
    const teamRankings: { [team: number]: Array<{ rank: number; player: Player }> } = {
      0: [],
      1: []
    };

    dealResult.rankings.forEach((playerSeat, index) => {
      const team = getTeamForPlayer(playerSeat);
      const player = getPlayerBySeat(playerSeat);
      if (player) {
        teamRankings[team].push({
          rank: index + 1, // Convert to 1-based ranking
          player
        });
      }
    });

    return teamRankings;
  };

  // Get victory type display text
  const getVictoryTypeText = (victoryType: VictoryType): string => {
    switch (victoryType) {
      case VictoryType.DOUBLE_DOWN:
        return '双下';
      case VictoryType.SINGLE_LAST:
        return '单贡';
      case VictoryType.PARTNER_LAST:
        return '对贡';
      default:
        return '胜利';
    }
  };

  // Get level display text
  const getLevelText = (level: number): string => {
    if (level <= 10) return level.toString();
    switch (level) {
      case 11: return 'J';
      case 12: return 'Q';
      case 13: return 'K';
      case 14: return 'A';
      default: return level.toString();
    }
  };

  // Format duration
  const formatDuration = (durationMs: number): string => {
    const minutes = Math.floor(durationMs / 60000);
    const seconds = Math.floor((durationMs % 60000) / 1000);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  const teamRankings = getTeamRankings();
  const winningTeam = dealResult.winning_team;
  const losingTeam = 1 - winningTeam;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div className="text-center mb-6">
          <h2 className="text-2xl font-bold text-gray-800 mb-2">
            {isMatchFinished ? '比赛结束' : '局结算'}
          </h2>
          <div className="text-lg text-gray-600">
            {getVictoryTypeText(dealResult.victory_type)} - 
            队伍{winningTeam + 1}获胜 (+{dealResult.upgrades[winningTeam]}级)
          </div>
        </div>

        {/* Team Rankings */}
        <div className="grid grid-cols-2 gap-6 mb-6">
          {/* Winning Team */}
          <div className="bg-green-50 border border-green-200 rounded-lg p-4">
            <h3 className="text-lg font-semibold text-green-800 mb-3 text-center">
              队伍{winningTeam + 1} (胜方)
            </h3>
            <div className="space-y-2">
              {teamRankings[winningTeam].map(({ rank, player }) => (
                <div key={player.seat} className="flex justify-between items-center">
                  <span className="font-medium">{player.username}</span>
                  <span className="text-sm bg-green-100 px-2 py-1 rounded">
                    第{rank}名
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-3 pt-3 border-t border-green-200">
              <div className="flex justify-between text-sm">
                <span>当前等级:</span>
                <span className="font-medium">
                  {getLevelText(teamLevels[winningTeam])}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span>升级:</span>
                <span className="font-medium text-green-600">
                  +{dealResult.upgrades[winningTeam]}级
                </span>
              </div>
            </div>
          </div>

          {/* Losing Team */}
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <h3 className="text-lg font-semibold text-red-800 mb-3 text-center">
              队伍{losingTeam + 1} (负方)
            </h3>
            <div className="space-y-2">
              {teamRankings[losingTeam].map(({ rank, player }) => (
                <div key={player.seat} className="flex justify-between items-center">
                  <span className="font-medium">{player.username}</span>
                  <span className="text-sm bg-red-100 px-2 py-1 rounded">
                    第{rank}名
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-3 pt-3 border-t border-red-200">
              <div className="flex justify-between text-sm">
                <span>当前等级:</span>
                <span className="font-medium">
                  {getLevelText(teamLevels[losingTeam])}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span>升级:</span>
                <span className="font-medium text-gray-500">
                  +{dealResult.upgrades[losingTeam]}级
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Deal Statistics */}
        <div className="bg-gray-50 rounded-lg p-4 mb-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-3">本局统计</h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div className="flex justify-between">
              <span>游戏时长:</span>
              <span className="font-medium">{formatDuration(dealResult.duration)}</span>
            </div>
            <div className="flex justify-between">
              <span>总轮次:</span>
              <span className="font-medium">{dealResult.trick_count}</span>
            </div>
            <div className="flex justify-between">
              <span>胜利类型:</span>
              <span className="font-medium">{getVictoryTypeText(dealResult.victory_type)}</span>
            </div>
            <div className="flex justify-between">
              <span>上贡情况:</span>
              <span className="font-medium">
                {dealResult.statistics.tribute_info.has_tribute ? '有上贡' : '无上贡'}
              </span>
            </div>
          </div>
        </div>

        {/* Player Statistics */}
        <div className="bg-gray-50 rounded-lg p-4 mb-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-3">玩家统计</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-300">
                  <th className="text-left py-2">玩家</th>
                  <th className="text-center py-2">排名</th>
                  <th className="text-center py-2">出牌次数</th>
                  <th className="text-center py-2">获胜轮次</th>
                  <th className="text-center py-2">过牌次数</th>
                  <th className="text-center py-2">超时次数</th>
                </tr>
              </thead>
              <tbody>
                {dealResult.statistics.player_stats.map((stats) => {
                  const player = getPlayerBySeat(stats.player_seat);
                  return (
                    <tr key={stats.player_seat} className="border-b border-gray-200">
                      <td className="py-2 font-medium">
                        {player?.username || `玩家${stats.player_seat + 1}`}
                      </td>
                      <td className="text-center py-2">第{stats.finish_rank}名</td>
                      <td className="text-center py-2">{stats.cards_played}</td>
                      <td className="text-center py-2">{stats.tricks_won}</td>
                      <td className="text-center py-2">{stats.pass_count}</td>
                      <td className="text-center py-2">{stats.timeout_count}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex justify-center space-x-4">
          {!isMatchFinished && (
            <button
              onClick={onContinue}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              继续游戏
            </button>
          )}
          <button
            onClick={onExit}
            className="px-6 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
          >
            {isMatchFinished ? '返回大厅' : '退出房间'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default DealResult;