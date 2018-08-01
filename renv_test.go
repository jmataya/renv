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
	os.Setenv("RENV", "development")

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

func TestFindEnvNestedDir(t *testing.T) {
	os.Setenv("RENV", "development")

	createFile(".env", "MYVAR=Donkey")

	cwd, _ := os.Getwd()
	testPath := path.Join(cwd, "some/path")
	os.MkdirAll(testPath, os.ModePerm)

	want := path.Join(cwd, ".env")
	envPath, err := FindEnv(testPath)
	if err != nil || envPath != want {
		t.Errorf("FindEnv(%s) = (%s, %v), want (%s, <nil>)", cwd, envPath, err, want)
	}

	deleteFile(".env")
	deleteFile("some")
}

func TestFindEnvTestFile(t *testing.T) {
	os.Setenv("RENV", "test")
	createFile(".env.test", "TESTVAR=Donkey")

	cwd, _ := os.Getwd()
	testPath := path.Join(cwd, "some/path")
	os.MkdirAll(testPath, os.ModePerm)

	want := path.Join(cwd, ".env.test")
	envPath, err := FindEnv(testPath)
	if err != nil || envPath != want {
		t.Errorf("FindEnv(%s) = (%s, %v), want (%s, <nil>)", cwd, envPath, err, want)
	}

	deleteFile(".env.test")
	deleteFile("some")
}

func TestFindEnvNoFile(t *testing.T) {
	os.Setenv("RENV", "development")

	cwd, _ := os.Getwd()
	want := ".env not found in project"

	_, err := FindEnv(cwd)
	if err == nil || err.Error() != want {
		t.Errorf("FindEnv(%s) = (_, %v), want (_, %v)", cwd, err, want)
	}
}

func TestFindEnvNoTestFile(t *testing.T) {
	os.Setenv("RENV", "test")
	createFile(".env", "TESTVAR=Donkey")

	cwd, _ := os.Getwd()
	os.MkdirAll(cwd, os.ModePerm)

	want := ".env.test not found in project"
	if _, err := FindEnv(cwd); err == nil || err.Error() != want {
		t.Errorf("FindEnv(%s) = (_, %v), want (_, %s)", cwd, err, want)
	}

	deleteFile(".env")
}

func TestLoadEnv(t *testing.T) {
	os.Setenv("RENV", "development")

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

func TestLoadEnvCommented(t *testing.T) {
	env := `
		FOO=test
		#BAR=var
	`

	os.Setenv("RENV", "development")

	createFile(".env", env)
	cwd, _ := os.Getwd()
	envPath, _ := FindEnv(cwd)

	if err := LoadEnv(envPath); err != nil {
		t.Errorf("LoadEnv(%s) = %v, want <nil>", envPath, err)
	} else if os.Getenv("FOO") != "test" {
		t.Errorf("os.Getenv(\"FOO\") = %s, want test", os.Getenv("FOO"))
	} else if os.Getenv("BAR") != "" {
		t.Errorf("os.Getenv(\"BAR\") = %s, want \"\"", os.Getenv("BAR"))
	}

	deleteFile(".env")
}

func TestLoadEnvInvalid(t *testing.T) {
	env := `
		FOO=test
		BAR
	`

	os.Setenv("RENV", "development")

	createFile(".env", env)
	cwd, _ := os.Getwd()
	envPath, _ := FindEnv(cwd)

	want := "Environment line BAR is malformed - must be of format key=value"

	if err := LoadEnv(envPath); err == nil || err.Error() != want {
		t.Errorf("LoadEnv(%s) = %v, want %s", envPath, err, want)
	} else if os.Getenv("FOO") != "test" {
		t.Errorf("os.Getenv(\"FOO\") = %s, want test", os.Getenv("FOO"))
	}

	deleteFile(".env")
}
