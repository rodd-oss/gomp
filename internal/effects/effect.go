package effects

type EffectType int

const (
	TypeDamage EffectType = iota
	TypeHeal
)

type Effect struct {
	Name         string
	Type         EffectType
	Damage       int
	FriendlyFire bool // if true, the effect will be applied to friendly units
}
