import React, { useState, useEffect } from 'react';
import { Volume2, VolumeX, Sun, Moon, Palette } from 'lucide-react';
import { Label } from '../ui';
import { Slider } from '../ui/Slider';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '../ui/Dialog';
import { audioService } from '../../services/audioService';
import { useThemeStore, type Theme } from '../../store/themeStore';

interface SettingsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const themeOptions: { value: Theme; label: string; icon: React.ReactNode }[] = [
  { value: 'light', label: '浅色', icon: <Sun className="w-4 h-4" /> },
  { value: 'dark', label: '深色', icon: <Moon className="w-4 h-4" /> },
  { value: 'chinese', label: '中国风', icon: <Palette className="w-4 h-4" /> },
];

const SettingsDialog: React.FC<SettingsDialogProps> = ({ open, onOpenChange }) => {
  const [volume, setVolume] = useState(() => Math.round(audioService.getVolume() * 100));
  const [muted, setMuted] = useState(() => audioService.getMuted());
  const { theme, setTheme } = useThemeStore();

  useEffect(() => {
    if (open) {
      setVolume(Math.round(audioService.getVolume() * 100));
      setMuted(audioService.getMuted());
    }
  }, [open]);

  const handleVolumeChange = (value: number[]) => {
    const newVolume = value[0];
    setVolume(newVolume);
    audioService.setVolume(newVolume / 100);
  };

  const handleVolumeCommit = () => {
    audioService.playPassSound();
  };

  const handleMutedToggle = () => {
    const newMuted = !muted;
    setMuted(newMuted);
    audioService.setMuted(newMuted);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="text-center">设置</DialogTitle>
          <DialogDescription className="sr-only">
            调整游戏音量和其他设置
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <Label>主题</Label>
          </div>
          <div className="grid grid-cols-3 gap-2">
            {themeOptions.map((option) => (
              <button
                key={option.value}
                type="button"
                onClick={() => setTheme(option.value)}
                className={`flex flex-col items-center gap-1.5 p-3 rounded-lg border transition-all ${
                  theme === option.value
                    ? 'border-primary bg-primary/10 text-primary'
                    : 'border-border bg-card text-muted-foreground hover:bg-muted'
                }`}
              >
                {option.icon}
                <span className="text-xs font-medium">{option.label}</span>
              </button>
            ))}
          </div>
        </div>

        <div className="border-t border-border" />

        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <Label htmlFor="volume-slider">音量</Label>
            <div className="flex items-center gap-2">
              <span className="text-sm text-table-400">{volume}%</span>
              <button
                type="button"
                onClick={handleMutedToggle}
                className="p-1 rounded hover:bg-table-200 transition-colors"
                title={muted ? '取消静音' : '静音'}
              >
                {muted ? (
                  <VolumeX className="w-5 h-5 text-table-400" />
                ) : (
                  <Volume2 className="w-5 h-5 text-table-400" />
                )}
              </button>
            </div>
          </div>
          <Slider
            id="volume-slider"
            value={[volume]}
            onValueChange={handleVolumeChange}
            onValueCommit={handleVolumeCommit}
            min={0}
            max={100}
            step={1}
            disabled={muted}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default SettingsDialog;
