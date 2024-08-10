package task

import "github.com/DiamondDrakeVentures/patchver/common"

type Task interface {
	ID() string
	SetID(id string)
	Type() string
	Name() string
	SetName(name string)

	Execute(logger common.Logger) error
}
