const IMG_BASE_URL = import.meta.env.VITE_IMG_BASE_URL || '';

export function getAvatarUrl(avatarKey?: string | null): string | null {
  if (!avatarKey) return null;
  const base = IMG_BASE_URL.replace(/\/+$/, '');
  return `${base}/${avatarKey}`;
}
