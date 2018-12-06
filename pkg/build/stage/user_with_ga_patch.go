package stage

import (
	"github.com/flant/dapp/pkg/build/builder"
	"github.com/flant/dapp/pkg/image"
)

func newUserWithGAPatchStage(builder builder.Builder, name StageName, baseStageOptions *NewBaseStageOptions) *UserWithGAPatchStage {
	s := &UserWithGAPatchStage{}
	s.UserStage = newUserStage(builder, name, baseStageOptions)
	s.GAPatchStage = newGAPatchStage(name, baseStageOptions)
	s.GAPatchStage.BaseStage = s.BaseStage
	return s
}

type UserWithGAPatchStage struct {
	*UserStage
	GAPatchStage *GAPatchStage
}

func (s *UserWithGAPatchStage) PrepareImage(c Conveyor, prevBuiltImage, image image.Image) error {
	if err := s.BaseStage.PrepareImage(c, prevBuiltImage, image); err != nil {
		return err
	}

	if err := s.GAPatchStage.PrepareImage(c, prevBuiltImage, image); err != nil {
		return nil
	}

	return nil
}
