package autoload

import (
	"fmt"
	"os"
	"strings"

	"github.com/jmataya/renv"
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envFile, err := renv.FindEnv(cwd)
	if err != nil {
		if strings.Contains(err.Error(), "not found in project") {
			fmt.Printf("warning: %s\n", err.Error())
			return
		}

		panic(err)
	}

	if err := renv.LoadEnv(envFile); err != nil {
		panic(err)
	}
}
