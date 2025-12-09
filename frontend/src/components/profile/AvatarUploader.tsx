import React, { useRef } from 'react';

const MAX_AVATAR_SIZE = 5 * 1024 * 1024; // 5MB
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];

export interface AvatarUploaderProps {
  src?: string | null;
  size?: number;
  onChange?: (file: File) => void;
  onError?: (message: string) => void;
  disabled?: boolean;
}

export const AvatarUploader: React.FC<AvatarUploaderProps> = ({
  src,
  size = 96,
  onChange,
  onError,
  disabled = false,
}) => {
  const inputRef = useRef<HTMLInputElement>(null);

  const handleClick = () => {
    if (!disabled && inputRef.current) {
      inputRef.current.click();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > MAX_AVATAR_SIZE) {
        onError?.('图片大小不能超过 5MB');
        if (inputRef.current) inputRef.current.value = '';
        return;
      }
      if (!ALLOWED_TYPES.includes(file.type)) {
        onError?.('仅支持 JPEG、PNG、GIF、WebP 格式');
        if (inputRef.current) inputRef.current.value = '';
        return;
      }
      onChange?.(file);
    }
    if (inputRef.current) {
      inputRef.current.value = '';
    }
  };

  return (
    <div
      className={`
        relative overflow-hidden
        rounded-lg border border-table-300 shadow-card
        bg-table-100
        transition-all duration-200
        ${disabled ? 'cursor-not-allowed opacity-60' : 'cursor-pointer hover:shadow-card-hover'}
        group
      `}
      style={{ width: size, height: size }}
      onClick={handleClick}
    >
      {src ? (
        <img
          src={src}
          alt="头像"
          className="w-full h-full object-cover"
        />
      ) : (
        <div className="w-full h-full flex items-center justify-center bg-table-200">
          <svg
            className="w-1/2 h-1/2 text-table-300"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" />
          </svg>
        </div>
      )}

      {!disabled && (
        <div className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
          <svg
            className="w-8 h-8 text-white"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
            />
          </svg>
        </div>
      )}

      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={handleFileChange}
        disabled={disabled}
      />
    </div>
  );
};
