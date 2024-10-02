package abilities

type AbilityError struct {
	Msg string
}

func (e *AbilityError) Error() string {
	return e.Msg
}

func NewAbilityError(msg string) *AbilityError {
	return &AbilityError{
		Msg: msg,
	}
}

func IsAbilityError(err error) bool {
	if _, ok := err.(*AbilityError); ok {
		return true
	}

	return false
}

func ErrorAbilityOnCooldown() *AbilityError {
	return NewAbilityError("Ability is on cooldown")
}

func ErrorUnknownAbilityType() *AbilityError {
	return NewAbilityError("Unknown ability type")
}

func ErrorUnknownCostType() *AbilityError {
	return NewAbilityError("Unknown cost type")
}

func ErrorUnknownTargetType() *AbilityError {
	return NewAbilityError("Unknown target type")
}
