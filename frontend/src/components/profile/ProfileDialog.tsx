import React, { useState, useEffect } from 'react';
import { Loader2 } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';
import { getAvatarUrl } from '../../utils/avatar';
import {
  Input,
  Label,
  Button,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '../ui';
import { AvatarUploader } from './AvatarUploader';
import { ImageCropModal } from './ImageCropModal';

interface ProfileDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const ProfileDialog: React.FC<ProfileDialogProps> = ({ open, onOpenChange }) => {
  const { user, updateUser } = useAuthStore();

  const [nickname, setNickname] = useState('');
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [rawImageUrl, setRawImageUrl] = useState<string | null>(null);
  const [showCropModal, setShowCropModal] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (user && open) {
      setNickname(user.nickname || '');
      setAvatarFile(null);
      setAvatarPreview(null);
      setError(null);
    }
  }, [user, open]);

  useEffect(() => {
    if (avatarFile) {
      const url = URL.createObjectURL(avatarFile);
      setAvatarPreview(url);
      return () => URL.revokeObjectURL(url);
    } else {
      setAvatarPreview(null);
    }
  }, [avatarFile]);

  useEffect(() => {
    return () => {
      if (rawImageUrl) {
        URL.revokeObjectURL(rawImageUrl);
      }
    };
  }, [rawImageUrl]);

  const handleAvatarChange = (file: File) => {
    const url = URL.createObjectURL(file);
    setRawImageUrl(url);
    setShowCropModal(true);
    setError(null);
  };

  const handleCropConfirm = (croppedFile: File) => {
    setAvatarFile(croppedFile);
    setShowCropModal(false);
    if (rawImageUrl) {
      URL.revokeObjectURL(rawImageUrl);
      setRawImageUrl(null);
    }
  };

  const handleCropCancel = () => {
    setShowCropModal(false);
    if (rawImageUrl) {
      URL.revokeObjectURL(rawImageUrl);
      setRawImageUrl(null);
    }
  };

  const handleSave = async () => {
    if (!user) return;

    setIsLoading(true);
    setError(null);

    try {
      const updateRes = await apiClient.updateProfile({
        nickname: nickname || undefined,
        avatar: avatarFile || undefined,
      });

      if (updateRes.data) {
        updateUser(updateRes.data);
        onOpenChange(false);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '保存失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    onOpenChange(false);
  };

  const currentAvatarUrl = avatarPreview || getAvatarUrl(user?.avatar_key);

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle className="text-center">个人资料</DialogTitle>
            <DialogDescription className="sr-only">
              修改头像和昵称
            </DialogDescription>
          </DialogHeader>

          <div className="flex flex-col items-center gap-6">
            <AvatarUploader
              src={currentAvatarUrl}
              size={96}
              onChange={handleAvatarChange}
              onError={setError}
              disabled={isLoading}
            />

            <div className="w-full space-y-2">
              <Label htmlFor="nickname">昵称</Label>
              <Input
                id="nickname"
                value={nickname}
                onChange={(e) => setNickname(e.target.value)}
                placeholder="请输入昵称"
                disabled={isLoading}
                maxLength={50}
              />
            </div>

            {error && (
              <p className="text-sm text-action-danger">{error}</p>
            )}

            <div className="flex gap-4 w-full">
              <Button
                intent="neutral"
                onClick={handleCancel}
                disabled={isLoading}
                className="flex-1"
              >
                取消
              </Button>
              <Button
                intent="primary"
                onClick={handleSave}
                disabled={isLoading}
                className="flex-1"
              >
                {isLoading && <Loader2 className="animate-spin" />}
                保存
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      {showCropModal && rawImageUrl && (
        <ImageCropModal
          imageSrc={rawImageUrl}
          onConfirm={handleCropConfirm}
          onCancel={handleCropCancel}
        />
      )}
    </>
  );
};

export default ProfileDialog;
