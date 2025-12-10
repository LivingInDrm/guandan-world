import avatar001 from '../assets/avatar_set/avatar_001.jpg';
import avatar002 from '../assets/avatar_set/avatar_002.jpg';
import avatar003 from '../assets/avatar_set/avatar_003.jpg';
import avatar004 from '../assets/avatar_set/avatar_004.jpg';
import avatar005 from '../assets/avatar_set/avatar_005.jpg';
import avatar006 from '../assets/avatar_set/avatar_006.jpg';
import avatar007 from '../assets/avatar_set/avatar_007.jpg';
import avatar008 from '../assets/avatar_set/avatar_008.jpg';
import avatar009 from '../assets/avatar_set/avatar_009.jpg';
import avatar010 from '../assets/avatar_set/avatar_010.jpg';
import avatar011 from '../assets/avatar_set/avatar_011.jpg';
import avatar012 from '../assets/avatar_set/avatar_012.jpg';
import avatar013 from '../assets/avatar_set/avatar_013.jpg';
import avatar014 from '../assets/avatar_set/avatar_014.jpg';
import avatar015 from '../assets/avatar_set/avatar_015.jpg';
import avatar016 from '../assets/avatar_set/avatar_016.jpg';
import avatar017 from '../assets/avatar_set/avatar_017.jpg';
import avatar018 from '../assets/avatar_set/avatar_018.jpg';
import avatar019 from '../assets/avatar_set/avatar_019.jpg';
import avatar020 from '../assets/avatar_set/avatar_020.jpg';

export const LOCAL_AVATARS = [
  avatar001, avatar002, avatar003, avatar004, avatar005,
  avatar006, avatar007, avatar008, avatar009, avatar010,
  avatar011, avatar012, avatar013, avatar014, avatar015,
  avatar016, avatar017, avatar018, avatar019, avatar020,
];

export function getAvatarByUsername(username: string): string {
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    hash = ((hash << 5) - hash) + username.charCodeAt(i);
    hash = hash & hash;
  }
  return LOCAL_AVATARS[Math.abs(hash) % LOCAL_AVATARS.length];
}

const IMG_BASE_URL = import.meta.env.VITE_IMG_BASE_URL || '';

export function getAvatarUrl(avatarKey?: string | null): string | null {
  if (!avatarKey) return null;
  const base = IMG_BASE_URL.replace(/\/+$/, '');
  return `${base}/${avatarKey}`;
}
