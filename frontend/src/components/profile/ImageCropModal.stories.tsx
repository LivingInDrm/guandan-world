import type { Story } from "@ladle/react";
import { useState } from "react";
import { ImageCropModal } from "./ImageCropModal";

const SAMPLE_IMAGE = "https://picsum.photos/800/600";

export const Default: Story = () => {
  const [result, setResult] = useState<string | null>(null);

  const handleConfirm = (file: File) => {
    const url = URL.createObjectURL(file);
    setResult(url);
  };

  if (result) {
    return (
      <div className="min-h-screen bg-primitive-neutral-900 flex flex-col items-center justify-center gap-4 p-8">
        <p className="text-fg-secondary">裁切结果：</p>
        <img
          src={result}
          alt="Cropped"
          className="w-32 h-32 rounded-lg border border-stroke"
        />
        <button
          onClick={() => setResult(null)}
          className="text-gold-400 underline"
        >
          重新裁切
        </button>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-primitive-neutral-900">
      <ImageCropModal
        imageSrc={SAMPLE_IMAGE}
        onConfirm={handleConfirm}
        onCancel={() => alert("取消")}
      />
    </div>
  );
};
