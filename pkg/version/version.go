package version

import "os"

type Version struct {
	Version string
	Build   string
}

func GetVersion() (*Version, error) {
	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	version, err := os.ReadFile(wd + "/VERSION")

	if err != nil {
		return nil, err
	}

	return &Version{
		Version: string(version),
		Build:   "dev",
	}, nil
}
