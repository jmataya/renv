package renv

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func createFile(name, contents string) error {
	cwd, _ := os.Getwd()
	file, err := os.Create(path.Join(cwd, name))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s\n", contents)); err != nil {
		return err
	}

	file.Sync()
	return nil
}

func deleteFile(name string) error {
	cwd, _ := os.Getwd()
	return os.RemoveAll(path.Join(cwd, name))
}

func TestFindEnvCurrentDir(t *testing.T) {
	if err := createFile(".env", "MYVAR=Donkey"); err != nil {
		t.Errorf("createFile(...) = %v", err)
		return
	}

	cwd, _ := os.Getwd()
	want := path.Join(cwd, ".env")

	envPath, err := FindEnv(cwd)
	if err != nil || envPath != want {
		t.Errorf("FindEnv(%s) = (%s, %v), want (%s, <nil>)", cwd, envPath, err, want)
	}

	deleteFile(".env")
}

func TestFindEnvNoFile(t *testing.T) {
	cwd, _ := os.Getwd()
	want := ".env not found in project"

	_, err := FindEnv(cwd)
	if err == nil || err.Error() != want {
		t.Errorf("FindEnv(%s) = (_, %v), want (_, %v)", cwd, err, want)
	}
}

func TestLoadEnv(t *testing.T) {
	createFile(".env", "MYVAR=Donkey")
	cwd, _ := os.Getwd()
	envPath, _ := FindEnv(cwd)

	if err := LoadEnv(envPath); err != nil {
		t.Errorf("LoadEnv(%s) = %v, want <nil>", envPath, err)
	} else if os.Getenv("MYVAR") != "Donkey" {
		t.Errorf("os.Getenv(\"MYVAR\") = %s, want Donkey", os.Getenv("MYVAR"))
	}

	deleteFile(".env")
}
