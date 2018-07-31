package autoload

import (
	"os"

	"github.com/jmataya/renv"
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envFile, err := renv.FindEnv(cwd)
	if err != nil {
		panic(err)
	}

	if err := renv.LoadEnv(envFile); err != nil {
		panic(err)
	}
}
