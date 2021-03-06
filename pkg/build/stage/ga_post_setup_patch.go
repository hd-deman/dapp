package stage

import (
	"github.com/flant/dapp/pkg/image"
	"github.com/flant/dapp/pkg/util"
)

const patchSizeStep = 1024 * 1024

func NewGAPostSetupPatchStage(gaPatchStageOptions *NewGaPatchStageOptions, baseStageOptions *NewBaseStageOptions) *GAPostSetupPatchStage {
	s := &GAPostSetupPatchStage{}
	s.GAPatchStage = newGAPatchStage(GAPostSetupPatch, gaPatchStageOptions, baseStageOptions)
	return s
}

type GAPostSetupPatchStage struct {
	*GAPatchStage
}

func (s *GAPostSetupPatchStage) GetDependencies(_ Conveyor, prevImage image.Image) (string, error) {
	var size int64
	for _, ga := range s.gitArtifacts {
		commit := ga.GetGACommitFromImageLabels(prevImage)
		if commit != "" {
			exist, err := ga.GitRepo().IsCommitExists(commit)
			if err != nil {
				return "", err
			}

			if exist {
				patchSize, err := ga.PatchSize(commit)
				if err != nil {
					return "", err
				}

				size += patchSize
			}
		}
	}

	return util.Sha256Hash(string(size / patchSizeStep)), nil
}
