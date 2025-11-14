// proto_adapter_comp.go - 牌型组合 Proto 适配器
//
// 职责:
// - CardComp 接口及其 11 种实现类型的转换（SDK ↔ Proto）
// - 牌型包括：Fold, IllegalComp, Single, Pair, Triple, FullHouse,
//   Straight, Plate, Tube, JokerBomb, NaiveBomb, StraightFlush
//
// 依赖:
// - proto_adapter_basic.go: ToProtoCards, FromProtoCards
// - proto_adapter_enums.go: ToProtoCompType, FromProtoCompType
//
// 被依赖:
// - proto_adapter_action.go: 使用 ToProtoCardComp, FromProtoCardComp
package sdk

import (
	pbgame "guandan-world/proto/gen/go/game"
)

// ==================== CardComp Adapters ====================

// ToProtoCardComp 转换 SDK CardComp 到 Proto CardComp
// 特殊处理：
// - 接口转换：根据具体类型使用对应的 oneof 字段
// - 所有牌型共享相同的 BaseComp 结构（Cards, Valid, Type）
func ToProtoCardComp(c CardComp) *pbgame.CardComp {
	if c == nil {
		return nil
	}

	switch comp := c.(type) {
	case *Fold:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_Fold{
				Fold: &pbgame.FoldComp{
					Cards: ToProtoCards(comp.Cards),
					Valid: comp.Valid,
					Type:  ToProtoCompType(comp.Type),
				},
			},
		}
	case *IllegalComp:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_Illegal{
				Illegal: &pbgame.IllegalComp{
					Cards: ToProtoCards(comp.Cards),
					Valid: comp.Valid,
					Type:  ToProtoCompType(comp.Type),
				},
			},
		}
	case *Single:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_Single{
				Single: &pbgame.SingleComp{
					Cards: ToProtoCards(comp.Cards),
					Valid: comp.Valid,
					Type:  ToProtoCompType(comp.Type),
				},
			},
		}
	case *Pair:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_Pair{
				Pair: &pbgame.PairComp{
					Cards:           ToProtoCards(comp.Cards),
					Valid:           comp.Valid,
					Type:            ToProtoCompType(comp.Type),
					NormalizedCards: ToProtoCards(comp.NormalizedCards),
				},
			},
		}
	case *Triple:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_Triple{
				Triple: &pbgame.TripleComp{
					Cards:           ToProtoCards(comp.Cards),
					Valid:           comp.Valid,
					Type:            ToProtoCompType(comp.Type),
					NormalizedCards: ToProtoCards(comp.NormalizedCards),
				},
			},
		}
	case *FullHouse:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_FullHouse{
				FullHouse: &pbgame.FullHouseComp{
					Cards:           ToProtoCards(comp.Cards),
					Valid:           comp.Valid,
					Type:            ToProtoCompType(comp.Type),
					NormalizedCards: ToProtoCards(comp.NormalizedCards),
				},
			},
		}
	case *Straight:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_Straight{
				Straight: &pbgame.StraightComp{
					Cards:           ToProtoCards(comp.Cards),
					Valid:           comp.Valid,
					Type:            ToProtoCompType(comp.Type),
					NormalizedCards: ToProtoCards(comp.NormalizedCards),
				},
			},
		}
	case *Plate:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_Plate{
				Plate: &pbgame.PlateComp{
					Cards:           ToProtoCards(comp.Cards),
					Valid:           comp.Valid,
					Type:            ToProtoCompType(comp.Type),
					NormalizedCards: ToProtoCards(comp.NormalizedCards),
				},
			},
		}
	case *Tube:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_Tube{
				Tube: &pbgame.TubeComp{
					Cards:           ToProtoCards(comp.Cards),
					Valid:           comp.Valid,
					Type:            ToProtoCompType(comp.Type),
					NormalizedCards: ToProtoCards(comp.NormalizedCards),
				},
			},
		}
	case *JokerBomb:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_JokerBomb{
				JokerBomb: &pbgame.JokerBombComp{
					Cards: ToProtoCards(comp.Cards),
					Valid: comp.Valid,
					Type:  ToProtoCompType(comp.Type),
				},
			},
		}
	case *NaiveBomb:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_NaiveBomb{
				NaiveBomb: &pbgame.NaiveBombComp{
					Cards:           ToProtoCards(comp.Cards),
					Valid:           comp.Valid,
					Type:            ToProtoCompType(comp.Type),
					NormalizedCards: ToProtoCards(comp.NormalizedCards),
				},
			},
		}
	case *StraightFlush:
		return &pbgame.CardComp{
			Comp: &pbgame.CardComp_StraightFlush{
				StraightFlush: &pbgame.StraightFlushComp{
					Cards:           ToProtoCards(comp.Cards),
					Valid:           comp.Valid,
					Type:            ToProtoCompType(comp.Type),
					NormalizedCards: ToProtoCards(comp.NormalizedCards),
				},
			},
		}
	default:
		return nil
	}
}

// FromProtoCardComp 转换 Proto CardComp 到 SDK CardComp
// 特殊处理：
// - oneof 转换：根据设置的字段返回对应的 Go 结构
// - 所有牌型共享相同的 BaseComp 结构（Cards, Valid, Type）
func FromProtoCardComp(pc *pbgame.CardComp) CardComp {
	if pc == nil {
		return nil
	}

	switch comp := pc.Comp.(type) {
	case *pbgame.CardComp_Fold:
		return &Fold{
			BaseComp: BaseComp{
				Cards: FromProtoCards(comp.Fold.Cards),
				Valid: comp.Fold.Valid,
				Type:  FromProtoCompType(comp.Fold.Type),
			},
		}
	case *pbgame.CardComp_Illegal:
		return &IllegalComp{
			BaseComp: BaseComp{
				Cards: FromProtoCards(comp.Illegal.Cards),
				Valid: comp.Illegal.Valid,
				Type:  FromProtoCompType(comp.Illegal.Type),
			},
		}
	case *pbgame.CardComp_Single:
		return &Single{
			BaseComp: BaseComp{
				Cards: FromProtoCards(comp.Single.Cards),
				Valid: comp.Single.Valid,
				Type:  FromProtoCompType(comp.Single.Type),
			},
		}
	case *pbgame.CardComp_Pair:
		return &Pair{
			BaseComp: BaseComp{
				Cards:           FromProtoCards(comp.Pair.Cards),
				NormalizedCards: FromProtoCards(comp.Pair.NormalizedCards),
				Valid:           comp.Pair.Valid,
				Type:            FromProtoCompType(comp.Pair.Type),
			},
		}
	case *pbgame.CardComp_Triple:
		return &Triple{
			BaseComp: BaseComp{
				Cards:           FromProtoCards(comp.Triple.Cards),
				NormalizedCards: FromProtoCards(comp.Triple.NormalizedCards),
				Valid:           comp.Triple.Valid,
				Type:            FromProtoCompType(comp.Triple.Type),
			},
		}
	case *pbgame.CardComp_FullHouse:
		return &FullHouse{
			BaseComp: BaseComp{
				Cards:           FromProtoCards(comp.FullHouse.Cards),
				NormalizedCards: FromProtoCards(comp.FullHouse.NormalizedCards),
				Valid:           comp.FullHouse.Valid,
				Type:            FromProtoCompType(comp.FullHouse.Type),
			},
		}
	case *pbgame.CardComp_Straight:
		return &Straight{
			BaseComp: BaseComp{
				Cards:           FromProtoCards(comp.Straight.Cards),
				NormalizedCards: FromProtoCards(comp.Straight.NormalizedCards),
				Valid:           comp.Straight.Valid,
				Type:            FromProtoCompType(comp.Straight.Type),
			},
		}
	case *pbgame.CardComp_Plate:
		return &Plate{
			BaseComp: BaseComp{
				Cards:           FromProtoCards(comp.Plate.Cards),
				NormalizedCards: FromProtoCards(comp.Plate.NormalizedCards),
				Valid:           comp.Plate.Valid,
				Type:            FromProtoCompType(comp.Plate.Type),
			},
		}
	case *pbgame.CardComp_Tube:
		return &Tube{
			BaseComp: BaseComp{
				Cards:           FromProtoCards(comp.Tube.Cards),
				NormalizedCards: FromProtoCards(comp.Tube.NormalizedCards),
				Valid:           comp.Tube.Valid,
				Type:            FromProtoCompType(comp.Tube.Type),
			},
		}
	case *pbgame.CardComp_JokerBomb:
		return &JokerBomb{
			BaseComp: BaseComp{
				Cards: FromProtoCards(comp.JokerBomb.Cards),
				Valid: comp.JokerBomb.Valid,
				Type:  FromProtoCompType(comp.JokerBomb.Type),
			},
		}
	case *pbgame.CardComp_NaiveBomb:
		return &NaiveBomb{
			BaseComp: BaseComp{
				Cards:           FromProtoCards(comp.NaiveBomb.Cards),
				NormalizedCards: FromProtoCards(comp.NaiveBomb.NormalizedCards),
				Valid:           comp.NaiveBomb.Valid,
				Type:            FromProtoCompType(comp.NaiveBomb.Type),
			},
		}
	case *pbgame.CardComp_StraightFlush:
		return &StraightFlush{
			BaseComp: BaseComp{
				Cards:           FromProtoCards(comp.StraightFlush.Cards),
				NormalizedCards: FromProtoCards(comp.StraightFlush.NormalizedCards),
				Valid:           comp.StraightFlush.Valid,
				Type:            FromProtoCompType(comp.StraightFlush.Type),
			},
		}
	default:
		return nil
	}
}
