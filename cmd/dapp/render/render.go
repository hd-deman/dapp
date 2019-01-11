package render

import (
	"fmt"

	"github.com/flant/dapp/cmd/dapp/common"
	"github.com/flant/dapp/pkg/dapp"
	"github.com/flant/dapp/pkg/deploy"
	"github.com/flant/dapp/pkg/lock"
	"github.com/flant/dapp/pkg/true_git"
	"github.com/spf13/cobra"
)

var CmdData struct {
	Values       []string
	SecretValues []string
	Set          []string
	SetString    []string
}

var CommonCmdData common.CmdData

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "render",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runRender()
			if err != nil {
				return fmt.Errorf("render failed: %s", err)
			}
			return nil
		},
	}

	common.SetupDir(&CommonCmdData, cmd)
	common.SetupTmpDir(&CommonCmdData, cmd)
	common.SetupHomeDir(&CommonCmdData, cmd)

	cmd.PersistentFlags().StringArrayVarP(&CmdData.Values, "values", "", []string{}, "Additional helm values")
	cmd.PersistentFlags().StringArrayVarP(&CmdData.SecretValues, "secret-values", "", []string{}, "Additional helm secret values")
	cmd.PersistentFlags().StringArrayVarP(&CmdData.Set, "set", "", []string{}, "Additional helm sets")
	cmd.PersistentFlags().StringArrayVarP(&CmdData.SetString, "set-string", "", []string{}, "Additional helm STRING sets")

	return cmd
}

func runRender() error {
	if err := dapp.Init(*CommonCmdData.TmpDir, *CommonCmdData.HomeDir); err != nil {
		return fmt.Errorf("initialization error: %s", err)
	}

	if err := lock.Init(); err != nil {
		return err
	}

	if err := true_git.Init(); err != nil {
		return err
	}

	if err := deploy.Init(); err != nil {
		return err
	}

	projectDir, err := common.GetProjectDir(&CommonCmdData)
	if err != nil {
		return fmt.Errorf("getting project dir failed: %s", err)
	}

	dappfile, err := common.GetDappfile(projectDir)
	if err != nil {
		return fmt.Errorf("dappfile parsing failed: %s", err)
	}

	return deploy.RunRender(projectDir, dappfile, deploy.RenderOptions{
		Values:       CmdData.Values,
		SecretValues: CmdData.SecretValues,
		Set:          CmdData.Set,
		SetString:    CmdData.SetString,
	})
}
