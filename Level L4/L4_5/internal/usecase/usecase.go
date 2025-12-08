package usecase

import (
	"api_optimization/internal/entity"
)

// UseCase - бизнес логика
type UseCase struct{}

func New() *UseCase {
	return &UseCase{}
}

// ReturnSum - вернуть сумму
func (u *UseCase) ReturnSum(nums entity.Input) *entity.Output {
	return &entity.Output{Sum: nums.A + nums.B}
}
