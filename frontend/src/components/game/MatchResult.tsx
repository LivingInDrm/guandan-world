import React from 'react';
import { MatchResult as MatchResultType, Player } from '../../types';

interface MatchResultProps {
  matchResult: MatchResultType;
  onReturnToLobby: () => void;
}

const MatchResult: React.FC<MatchResultProps> = ({
  matchResult,
  onReturnToLobby
}) => {
  // Helper function to get team for player
  const getTeamForPlayer = (playerSeat: number): number => {
    return playerSeat % 2; // Team 0: seats 0,2; Team 1: seats 1,3
  };

  // Helper function to get player by seat
  const getPlayerBySeat = (seat: number): Player | undefined => {
    return matchResult.players.find(p => p.seat === seat);
  };

  // Group players by team
  const getTeamPlayers = () => {
    const teamPlayers: { [team: number]: Player[] } = {
      0: [],
      1: []
    };

    matchResult.players.forEach(player => {
      const team = getTeamForPlayer(player.seat);
      teamPlayers[team].push(player);
    });

    return teamPlayers;
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
    const hours = Math.floor(durationMs / 3600000);
    const minutes = Math.floor((durationMs % 3600000) / 60000);
    const seconds = Math.floor((durationMs % 60000) / 1000);
    
    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  const teamPlayers = getTeamPlayers();
  const winningTeam = matchResult.winner;
  const losingTeam = 1 - winningTeam;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-3xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-800 mb-4">
            🎉 比赛结束 🎉
          </h1>
          <div className="text-xl text-gray-600 mb-2">
            队伍{winningTeam + 1}获得最终胜利！
          </div>
          <div className="text-lg text-gray-500">
            恭喜率先达到A级！
          </div>
        </div>

        {/* Final Team Results */}
        <div className="grid grid-cols-2 gap-8 mb-8">
          {/* Winning Team */}
          <div className="bg-gradient-to-br from-yellow-50 to-yellow-100 border-2 border-yellow-300 rounded-lg p-6">
            <div className="text-center mb-4">
              <h2 className="text-2xl font-bold text-yellow-800 mb-2">
                🏆 队伍{winningTeam + 1}
              </h2>
              <div className="text-lg font-semibold text-yellow-700">
                冠军队伍
              </div>
            </div>
            <div className="space-y-3">
              {teamPlayers[winningTeam].map(player => (
                <div key={player.seat} className="flex justify-between items-center bg-white rounded-lg p-3 shadow-sm">
                  <span className="font-medium text-gray-800">{player.username}</span>
                  <span className="text-sm bg-yellow-100 text-yellow-800 px-3 py-1 rounded-full">
                    座位{player.seat + 1}
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-4 border-t border-yellow-200">
              <div className="flex justify-between items-center">
                <span className="text-lg font-semibold text-yellow-800">最终等级:</span>
                <span className="text-2xl font-bold text-yellow-800">
                  {getLevelText(matchResult.final_levels[winningTeam])}
                </span>
              </div>
            </div>
          </div>

          {/* Losing Team */}
          <div className="bg-gray-50 border border-gray-200 rounded-lg p-6">
            <div className="text-center mb-4">
              <h2 className="text-2xl font-bold text-gray-700 mb-2">
                队伍{losingTeam + 1}
              </h2>
              <div className="text-lg font-semibold text-gray-600">
                亚军队伍
              </div>
            </div>
            <div className="space-y-3">
              {teamPlayers[losingTeam].map(player => (
                <div key={player.seat} className="flex justify-between items-center bg-white rounded-lg p-3 shadow-sm">
                  <span className="font-medium text-gray-800">{player.username}</span>
                  <span className="text-sm bg-gray-100 text-gray-600 px-3 py-1 rounded-full">
                    座位{player.seat + 1}
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-4 pt-4 border-t border-gray-200">
              <div className="flex justify-between items-center">
                <span className="text-lg font-semibold text-gray-700">最终等级:</span>
                <span className="text-2xl font-bold text-gray-700">
                  {getLevelText(matchResult.final_levels[losingTeam])}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Match Statistics */}
        <div className="bg-gray-50 rounded-lg p-6 mb-8">
          <h3 className="text-xl font-semibold text-gray-800 mb-4 text-center">比赛统计</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-600 mb-1">
                {matchResult.statistics.total_deals}
              </div>
              <div className="text-sm text-gray-600">总局数</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600 mb-1">
                {formatDuration(matchResult.statistics.total_duration)}
              </div>
              <div className="text-sm text-gray-600">总时长</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-purple-600 mb-1">
                {getLevelText(matchResult.final_levels[0])} vs {getLevelText(matchResult.final_levels[1])}
              </div>
              <div className="text-sm text-gray-600">最终等级</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-orange-600 mb-1">
                队伍{matchResult.winner + 1}
              </div>
              <div className="text-sm text-gray-600">获胜队伍</div>
            </div>
          </div>
        </div>

        {/* Congratulations Message */}
        <div className="bg-gradient-to-r from-blue-50 to-purple-50 border border-blue-200 rounded-lg p-6 mb-6">
          <div className="text-center">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">
              感谢参与本次掼蛋对战！
            </h3>
            <p className="text-gray-600 mb-4">
              经过 {matchResult.statistics.total_deals} 局激烈的对战，
              队伍{winningTeam + 1} 成功率先达到A级，获得最终胜利！
            </p>
            <div className="flex justify-center items-center space-x-4 text-sm text-gray-500">
              <span>🎯 精彩对局</span>
              <span>•</span>
              <span>🤝 友谊第一</span>
              <span>•</span>
              <span>🏆 比赛第二</span>
            </div>
          </div>
        </div>

        {/* Action Button */}
        <div className="flex justify-center">
          <button
            onClick={onReturnToLobby}
            className="px-8 py-3 bg-blue-600 text-white text-lg font-semibold rounded-lg hover:bg-blue-700 transition-colors shadow-lg"
          >
            返回大厅
          </button>
        </div>
      </div>
    </div>
  );
};

export default MatchResult;