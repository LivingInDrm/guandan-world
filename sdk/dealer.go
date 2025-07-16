package sdk

import (
	"fmt"
	"math/rand"
	"time"
)

// Dealer handles all card dealing operations independently within the SDK
type Dealer struct {
	level int
	deck  []*Card
}

// NewDealer creates a new dealer for the specified level
func NewDealer(level int) (*Dealer, error) {
	if level < 2 || level > 14 {
		return nil, fmt.Errorf("invalid level: %d, must be between 2 and 14", level)
	}
	
	return &Dealer{
		level: level,
	}, nil
}

// CreateFullDeck creates a complete deck of 108 cards for Guandan
func (d *Dealer) CreateFullDeck() []*Card {
	deck := make([]*Card, 0, 108)
	
	// Add regular cards (2-A) for each suit, 2 copies each
	// Total: 4 suits × 13 cards × 2 copies = 104 cards
	for _, color := range Colors {
		for number := 2; number <= 14; number++ {
			for copy := 0; copy < 2; copy++ {
				card, err := NewCard(number, color, d.level)
				if err != nil {
					// This should never happen with valid inputs
					continue
				}
				deck = append(deck, card)
			}
		}
	}
	
	// Add jokers (2 small jokers + 2 big jokers)
	// Total: 4 jokers
	for copy := 0; copy < 2; copy++ {
		smallJoker, _ := NewCard(15, "Joker", d.level) // Black Joker
		bigJoker, _ := NewCard(16, "Joker", d.level)   // Red Joker
		deck = append(deck, smallJoker, bigJoker)
	}
	
	d.deck = deck
	return deck
}

// ShuffleDeck shuffles the deck using Fisher-Yates algorithm
func (d *Dealer) ShuffleDeck() {
	if d.deck == nil {
		d.CreateFullDeck()
	}
	
	// Seed random number generator with current time
	rand.Seed(time.Now().UnixNano())
	
	// Fisher-Yates shuffle
	for i := len(d.deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d.deck[i], d.deck[j] = d.deck[j], d.deck[i]
	}
}

// DealCards deals 27 cards to each of the 4 players and returns the hands
// Returns [4][]*Card representing each player's hand
// In standard Guandan, all 108 cards are dealt to players (no bottom cards)
func (d *Dealer) DealCards() ([4][]*Card, error) {
	// Create and shuffle deck
	d.CreateFullDeck()
	d.ShuffleDeck()
	
	if len(d.deck) != 108 {
		return [4][]*Card{}, fmt.Errorf("invalid deck size: %d, expected 108", len(d.deck))
	}
	
	var playerHands [4][]*Card
	
	// Deal 27 cards to each player
	for player := 0; player < 4; player++ {
		playerHands[player] = make([]*Card, 27)
		for card := 0; card < 27; card++ {
			cardIndex := player*27 + card
			if cardIndex >= len(d.deck) {
				return [4][]*Card{}, fmt.Errorf("not enough cards in deck")
			}
			playerHands[player][card] = d.deck[cardIndex]
		}
		
		// Sort each player's hand automatically
		playerHands[player] = d.SortPlayerHand(playerHands[player])
	}
	
	return playerHands, nil
}

// SortPlayerHand sorts a player's hand according to Guandan rules
func (d *Dealer) SortPlayerHand(cards []*Card) []*Card {
	return sortCards(cards)
}

// ValidateDeck validates that the deck contains exactly 108 cards with correct distribution
func (d *Dealer) ValidateDeck() error {
	if d.deck == nil {
		return fmt.Errorf("deck is nil")
	}
	
	if len(d.deck) != 108 {
		return fmt.Errorf("invalid deck size: %d, expected 108", len(d.deck))
	}
	
	// Count cards by type
	cardCounts := make(map[string]int)
	
	for _, card := range d.deck {
		key := fmt.Sprintf("%d-%s", card.Number, card.Color)
		cardCounts[key]++
	}
	
	// Validate regular cards (should have 2 of each)
	for _, color := range Colors {
		for number := 2; number <= 14; number++ {
			key := fmt.Sprintf("%d-%s", number, color)
			if cardCounts[key] != 2 {
				return fmt.Errorf("invalid count for %s: %d, expected 2", key, cardCounts[key])
			}
		}
	}
	
	// Validate jokers (should have 2 of each)
	if cardCounts["15-Joker"] != 2 {
		return fmt.Errorf("invalid count for Black Joker: %d, expected 2", cardCounts["15-Joker"])
	}
	if cardCounts["16-Joker"] != 2 {
		return fmt.Errorf("invalid count for Red Joker: %d, expected 2", cardCounts["16-Joker"])
	}
	
	return nil
}

// GetDeck returns the current deck (for testing purposes)
func (d *Dealer) GetDeck() []*Card {
	return d.deck
}

// GetLevel returns the current level
func (d *Dealer) GetLevel() int {
	return d.level
}

