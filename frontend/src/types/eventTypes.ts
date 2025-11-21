/**
 * Event Type Aliases
 * 
 * Simplified aliases for EventType enum to make code more readable
 * Usage: import { GameEventType } from '../types/eventTypes';
 */

import { EventType } from './generated/event';

// Simplified event type aliases for better readability
export const GameEventType = {
  // Match-level events
  MATCH_STARTED: EventType.EVENT_TYPE_MATCH_STARTED,
  MATCH_ENDED: EventType.EVENT_TYPE_MATCH_ENDED,
  
  // Deal-level events
  DEAL_STARTED: EventType.EVENT_TYPE_DEAL_STARTED,
  CARDS_DEALT: EventType.EVENT_TYPE_CARDS_DEALT,
  DEAL_ENDED: EventType.EVENT_TYPE_DEAL_ENDED,
  
  // Tribute phase events
  TRIBUTE_STARTED: EventType.EVENT_TYPE_TRIBUTE_STARTED,
  TRIBUTE_EXEMPTED: EventType.EVENT_TYPE_TRIBUTE_EXEMPTED,
  TRIBUTE_CARD_SUBMITTED: EventType.EVENT_TYPE_TRIBUTE_CARD_SUBMITTED,
  TRIBUTE_CARD_SELECTED: EventType.EVENT_TYPE_TRIBUTE_CARD_SELECTED,
  TRIBUTE_CARD_RETURNED: EventType.EVENT_TYPE_TRIBUTE_CARD_RETURNED,
  TRIBUTE_COMPLETED: EventType.EVENT_TYPE_TRIBUTE_COMPLETED,
  
  // Trick-level events
  TRICK_STARTED: EventType.EVENT_TYPE_TRICK_STARTED,
  TRICK_ENDED: EventType.EVENT_TYPE_TRICK_ENDED,
  
  // Player action events
  PLAYER_PLAYED: EventType.EVENT_TYPE_PLAYER_PLAYED,
  PLAYER_PASSED: EventType.EVENT_TYPE_PLAYER_PASSED,
  PLAYER_TIMEOUT: EventType.EVENT_TYPE_PLAYER_TIMEOUT,
  PLAYER_DISCONNECT: EventType.EVENT_TYPE_PLAYER_DISCONNECT,
  PLAYER_RECONNECT: EventType.EVENT_TYPE_PLAYER_RECONNECT,
} as const;

// Type for event type values
export type GameEventTypeName = typeof GameEventType[keyof typeof GameEventType];

// Helper function to get event name as string for logging
export function getEventTypeName(type: EventType): string {
  const entry = Object.entries(GameEventType).find(([_, value]) => value === type);
  return entry ? entry[0].toLowerCase() : `unknown(${type})`;
}
