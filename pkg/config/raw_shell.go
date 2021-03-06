package config

type rawShell struct {
	BeforeInstall             interface{} `yaml:"beforeInstall,omitempty"`
	Install                   interface{} `yaml:"install,omitempty"`
	BeforeSetup               interface{} `yaml:"beforeSetup,omitempty"`
	Setup                     interface{} `yaml:"setup,omitempty"`
	CacheVersion              string      `yaml:"cacheVersion,omitempty"`
	BeforeInstallCacheVersion string      `yaml:"beforeInstallCacheVersion,omitempty"`
	InstallCacheVersion       string      `yaml:"installCacheVersion,omitempty"`
	BeforeSetupCacheVersion   string      `yaml:"beforeSetupCacheVersion,omitempty"`
	SetupCacheVersion         string      `yaml:"setupCacheVersion,omitempty"`

	rawDimg *rawDimg `yaml:"-"` // parent

	UnsupportedAttributes map[string]interface{} `yaml:",inline"`
}

func (c *rawShell) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if parent, ok := parentStack.Peek().(*rawDimg); ok {
		c.rawDimg = parent
	}

	type plain rawShell
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	if err := checkOverflow(c.UnsupportedAttributes, c, c.rawDimg.doc); err != nil {
		return err
	}

	return nil
}

func (c *rawShell) toDirective() (shell *Shell, err error) {
	shell = &Shell{}
	shell.CacheVersion = c.CacheVersion
	shell.BeforeInstallCacheVersion = c.BeforeInstallCacheVersion
	shell.InstallCacheVersion = c.InstallCacheVersion
	shell.BeforeSetupCacheVersion = c.BeforeSetupCacheVersion
	shell.SetupCacheVersion = c.SetupCacheVersion

	if beforeInstall, err := InterfaceToStringArray(c.BeforeInstall, c, c.rawDimg.doc); err != nil {
		return nil, err
	} else {
		shell.BeforeInstall = beforeInstall
	}

	if install, err := InterfaceToStringArray(c.Install, c, c.rawDimg.doc); err != nil {
		return nil, err
	} else {
		shell.Install = install
	}

	if beforeSetup, err := InterfaceToStringArray(c.BeforeSetup, c, c.rawDimg.doc); err != nil {
		return nil, err
	} else {
		shell.BeforeSetup = beforeSetup
	}

	if setup, err := InterfaceToStringArray(c.Setup, c, c.rawDimg.doc); err != nil {
		return nil, err
	} else {
		shell.Setup = setup
	}

	shell.raw = c

	if err := c.validateDirective(shell); err != nil {
		return nil, err
	}

	return shell, nil
}

func (c *rawShell) validateDirective(shell *Shell) error {
	if err := shell.validate(); err != nil {
		return err
	}

	return nil
}
