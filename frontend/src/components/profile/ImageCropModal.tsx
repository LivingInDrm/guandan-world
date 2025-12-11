import React, { useState, useCallback } from 'react';
import Cropper from 'react-easy-crop';
import type { Area } from 'react-easy-crop';
import { Loader2 } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, Button, Slider, Label } from '../ui-next';
import { getCroppedImage } from './cropImage';

export interface ImageCropModalProps {
  imageSrc: string;
  onConfirm: (croppedFile: File) => void;
  onCancel: () => void;
}

export const ImageCropModal: React.FC<ImageCropModalProps> = ({
  imageSrc,
  onConfirm,
  onCancel,
}) => {
  const [crop, setCrop] = useState({ x: 0, y: 0 });
  const [zoom, setZoom] = useState(1);
  const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const onCropComplete = useCallback((_: Area, croppedPixels: Area) => {
    setCroppedAreaPixels(croppedPixels);
  }, []);

  const handleConfirm = async () => {
    if (!croppedAreaPixels) return;

    setIsProcessing(true);
    setError(null);
    try {
      const croppedFile = await getCroppedImage(imageSrc, croppedAreaPixels);
      onConfirm(croppedFile);
    } catch (err) {
      console.error('Crop failed:', err);
      setError('裁切失败，请重试');
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <Dialog open onOpenChange={(isOpen) => !isOpen && onCancel()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>裁切头像</DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-4">
          <div className="relative w-full h-64 bg-table-400 rounded-lg overflow-hidden">
            <Cropper
              image={imageSrc}
              crop={crop}
              zoom={zoom}
              aspect={1}
              onCropChange={setCrop}
              onZoomChange={setZoom}
              onCropComplete={onCropComplete}
            />
          </div>

          <div className="flex items-center gap-3">
            <Label className="shrink-0">缩放</Label>
            <Slider
              value={[zoom]}
              min={1}
              max={3}
              step={0.1}
              onValueChange={(values) => setZoom(values[0])}
              className="flex-1"
            />
          </div>

          {error && (
            <p className="text-sm text-destructive">{error}</p>
          )}

          <div className="flex gap-3 justify-end">
            <Button intent="neutral" onClick={onCancel} disabled={isProcessing}>
              取消
            </Button>
            <Button
              intent="primary"
              onClick={handleConfirm}
              disabled={isProcessing || !croppedAreaPixels}
            >
              {isProcessing && <Loader2 className="animate-spin" />}
              确认
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
