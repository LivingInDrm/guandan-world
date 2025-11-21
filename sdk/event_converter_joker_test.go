package sdk

import (
	"testing"
)

func TestJokerConversion(t *testing.T) {
	// Test Small Joker
	smallJoker, _ := NewCard(15, "Joker", 2)
	smallJoker.DeckIndex = 104

	// Test Big Joker
	bigJoker, _ := NewCard(16, "Joker", 2)
	bigJoker.DeckIndex = 105

	// Test SDK -> Proto -> SDK for Small Joker
	protoSmall := ConvertCardToProto(smallJoker)
	if protoSmall.Suit != -1 {
		t.Errorf("Small Joker proto suit should be -1, got %d", protoSmall.Suit)
	}
	if protoSmall.Rank != 15 {
		t.Errorf("Small Joker proto rank should be 15, got %d", protoSmall.Rank)
	}

	sdkSmall := ConvertCardFromProto(protoSmall)
	if sdkSmall.Color != "Joker" {
		t.Errorf("Small Joker SDK color should be 'Joker', got '%s'", sdkSmall.Color)
	}
	if sdkSmall.Number != 15 {
		t.Errorf("Small Joker SDK number should be 15, got %d", sdkSmall.Number)
	}

	// Test SDK -> Proto -> SDK for Big Joker
	protoBig := ConvertCardToProto(bigJoker)
	if protoBig.Suit != -1 {
		t.Errorf("Big Joker proto suit should be -1, got %d", protoBig.Suit)
	}
	if protoBig.Rank != 16 {
		t.Errorf("Big Joker proto rank should be 16, got %d", protoBig.Rank)
	}

	sdkBig := ConvertCardFromProto(protoBig)
	if sdkBig.Color != "Joker" {
		t.Errorf("Big Joker SDK color should be 'Joker', got '%s'", sdkBig.Color)
	}
	if sdkBig.Number != 16 {
		t.Errorf("Big Joker SDK number should be 16, got %d", sdkBig.Number)
	}

	// Test validator key generation
	validator := &PlayValidator{}
	smallKey := validator.getCardKey(sdkSmall)
	bigKey := validator.getCardKey(sdkBig)

	expectedSmallKey := "15_Joker"
	expectedBigKey := "16_Joker"

	if smallKey != expectedSmallKey {
		t.Errorf("Small Joker key should be '%s', got '%s'", expectedSmallKey, smallKey)
	}
	if bigKey != expectedBigKey {
		t.Errorf("Big Joker key should be '%s', got '%s'", expectedBigKey, bigKey)
	}
}
