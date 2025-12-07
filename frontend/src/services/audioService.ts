import { CompType } from '../types/generated/common';

import singleSound from '../assets/audio/single.m4a';
import pairSound from '../assets/audio/pair.m4a';
import tripleSound from '../assets/audio/triple.m4a';
import fullhouseSound from '../assets/audio/fullhouse.m4a';
import straightSound from '../assets/audio/straight.m4a';
import plateSound from '../assets/audio/plate.m4a';
import tubeSound from '../assets/audio/tube.m4a';
import jokerbombSound from '../assets/audio/jokerbomb.m4a';
import naivebombSound from '../assets/audio/naivebomb.m4a';
import flushSound from '../assets/audio/flush.m4a';
import passSound from '../assets/audio/pass.m4a';

const compTypeToSound: Record<number, string> = {
  [CompType.COMP_TYPE_SINGLE]: singleSound,
  [CompType.COMP_TYPE_PAIR]: pairSound,
  [CompType.COMP_TYPE_TRIPLE]: tripleSound,
  [CompType.COMP_TYPE_FULL_HOUSE]: fullhouseSound,
  [CompType.COMP_TYPE_STRAIGHT]: straightSound,
  [CompType.COMP_TYPE_PLATE]: plateSound,
  [CompType.COMP_TYPE_TUBE]: tubeSound,
  [CompType.COMP_TYPE_JOKER_BOMB]: jokerbombSound,
  [CompType.COMP_TYPE_NAIVE_BOMB]: naivebombSound,
  [CompType.COMP_TYPE_STRAIGHT_FLUSH]: flushSound,
};

class AudioService {
  private audioCache: Map<string, HTMLAudioElement> = new Map();
  private volume: number = 1.0;
  private muted: boolean = false;

  constructor() {
    this.preload();
  }

  private preload(): void {
    const allSounds = [...Object.values(compTypeToSound), passSound];
    allSounds.forEach((src) => {
      const audio = new Audio(src);
      audio.preload = 'auto';
      this.audioCache.set(src, audio);
    });
  }

  private play(src: string): void {
    if (this.muted) return;

    const cachedAudio = this.audioCache.get(src);
    if (cachedAudio) {
      const audio = cachedAudio.cloneNode(true) as HTMLAudioElement;
      audio.volume = this.volume;
      audio.play().catch((err) => {
        console.warn('Audio play failed:', err);
      });
    } else {
      const audio = new Audio(src);
      audio.volume = this.volume;
      audio.play().catch((err) => {
        console.warn('Audio play failed:', err);
      });
    }
  }

  playCardSound(compType: CompType): void {
    const src = compTypeToSound[compType];
    if (src) {
      this.play(src);
    }
  }

  playPassSound(): void {
    this.play(passSound);
  }

  setVolume(volume: number): void {
    this.volume = Math.max(0, Math.min(1, volume));
  }

  setMuted(muted: boolean): void {
    this.muted = muted;
  }
}

export const audioService = new AudioService();
