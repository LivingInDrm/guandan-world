import React, { useState } from 'react';

interface CreateRoomModalProps {
  onClose: () => void;
  onConfirm: () => void;
}

const CreateRoomModal: React.FC<CreateRoomModalProps> = ({ onClose, onConfirm }) => {
  const [isCreating, setIsCreating] = useState(false);

  const handleConfirm = async () => {
    setIsCreating(true);
    try {
      await onConfirm();
    } finally {
      setIsCreating(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <div className="p-6">
          {/* Header */}
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-semibold text-gray-900">创建新房间</h3>
            <button
              onClick={onClose}
              disabled={isCreating}
              className="text-gray-400 hover:text-gray-600 disabled:opacity-50"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {/* Content */}
          <div className="mb-6">
            <p className="text-gray-600 mb-4">
              确认创建新房间？您将成为房主，负责管理房间和开始游戏。
            </p>
            
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h4 className="font-medium text-blue-900 mb-2">房间规则</h4>
              <ul className="text-sm text-blue-800 space-y-1">
                <li>• 房间最多容纳4名玩家</li>
                <li>• 房主可以在人数满足时开始游戏</li>
                <li>• 房主离开时会自动转移给其他玩家</li>
                <li>• 所有玩家离开后房间自动关闭</li>
              </ul>
            </div>
          </div>

          {/* Actions */}
          <div className="flex justify-end space-x-3">
            <button
              onClick={onClose}
              disabled={isCreating}
              className="px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors disabled:opacity-50"
            >
              取消
            </button>
            <button
              onClick={handleConfirm}
              disabled={isCreating}
              className={`px-4 py-2 text-white rounded-lg font-medium transition-colors ${
                isCreating
                  ? 'bg-gray-400 cursor-not-allowed'
                  : 'bg-blue-600 hover:bg-blue-700'
              }`}
            >
              {isCreating ? (
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  创建中...
                </div>
              ) : (
                '确认创建'
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateRoomModal;