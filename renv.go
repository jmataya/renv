// Package renv loads an environment file, parses it, and adds the values to
// the currently executing process using `os.Setenv()`.
//
// Here is an example .env file:
//
// 		MY_VAR=hi
//		ANOTHER_VAR=wat
//	  #COMMENTED_VAR=does-not-get-set
//
// Note the #COMMENTED_VAR will not be set.
//
// Usage:
//
//		import "renv"
//		renv.LoadEnv("/my/path/.env")
package renv

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// FindEnv locates an environment file anywhere between the path defined and
// all the way up the source tree. The searching functionality is only intended
// for use in development and test environments, so it will only search for
// files within $GOPATH, if path is part of $GOPATH.
func FindEnv(envPath string) (string, error) {
	goenv := os.Getenv("RENV")
	envFile := ".env"

	if len(goenv) > 0 && goenv != "development" {
		envFile = fmt.Sprintf("%s.%s", envFile, goenv)
	}

	envFilePath := path.Join(envPath, envFile)
	if _, err := os.Stat(envFilePath); err != nil {
		if os.IsNotExist(err) {
			gopath := os.Getenv("GOPATH")
			inGoPath := strings.HasPrefix(envPath, gopath)
			isGoPath, err := path.Match(gopath, envPath)
			if err != nil {
				return "", err
			}

			if inGoPath && !isGoPath {
				return FindEnv(path.Dir(envPath))
			}

			return "", fmt.Errorf("%s not found in project", envFile)
		}
		return "", err
	}

	return envFilePath, nil
}

// LoadEnv reads an environment file and uses `os.Setenv` to add the values to
// the current process.
func LoadEnv(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		kvp := strings.SplitN(line, "=", 2)
		if len(kvp) != 2 {
			return fmt.Errorf("Environment line %s is malformed - must be of format key=value", line)
		}

		os.Setenv(strings.TrimSpace(kvp[0]), strings.TrimSpace(kvp[1]))
	}

	return nil
}
