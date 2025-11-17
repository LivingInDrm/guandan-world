package sdk

import (
	"sync/atomic"
	"time"

	eventpb "guandan-world/proto/event"
)

// EventMetadataProvider manages event metadata generation
type EventMetadataProvider struct {
	seqCounter int64
}

// NewEventMetadataProvider creates a new metadata provider
func NewEventMetadataProvider() *EventMetadataProvider {
	return &EventMetadataProvider{
		seqCounter: 0,
	}
}

// NextSeq returns the next sequence number
func (emp *EventMetadataProvider) NextSeq() int64 {
	return atomic.AddInt64(&emp.seqCounter, 1)
}

// FillMetadata fills event metadata fields
// Parameters:
//   - event: the proto event to fill
//   - match: current match (required for match_id)
//   - deal: current deal (optional, for deal_index)
//   - trick: current trick (optional, for trick_index)
//   - actorSeat: actor seat (-1 if not applicable)
func (emp *EventMetadataProvider) FillMetadata(
	event *eventpb.GameEvent,
	match *Match,
	deal *Deal,
	trick *Trick,
	actorSeat int,
) {
	if event == nil {
		return
	}

	// Fill match_id
	if match != nil {
		event.MatchId = match.ID
	}

	// Fill deal_index
	if deal != nil && match != nil {
		event.DealIndex = int32(len(match.DealHistory))
	} else {
		event.DealIndex = -1
	}

	// Fill trick_index
	if trick != nil && deal != nil {
		event.TrickIndex = int32(len(deal.TrickHistory))
	} else {
		event.TrickIndex = -1
	}

	// Fill actor_seat
	if actorSeat >= 0 && actorSeat < 4 {
		event.ActorSeat = int32(actorSeat)
	} else {
		event.ActorSeat = -1
	}

	// Fill seq and timestamp
	event.Seq = emp.NextSeq()
	event.CreatedAtMs = time.Now().UnixMilli()
}

// CreateBaseEvent creates a base event with type and metadata
func (emp *EventMetadataProvider) CreateBaseEvent(
	eventType eventpb.EventType,
	match *Match,
	deal *Deal,
	trick *Trick,
	actorSeat int,
) *eventpb.GameEvent {
	event := &eventpb.GameEvent{
		Type: eventType,
	}
	emp.FillMetadata(event, match, deal, trick, actorSeat)
	return event
}
