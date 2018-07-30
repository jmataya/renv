# renv

Recursively find .env files in your current project and use them to populate
you environment.

The package works very similarly to [go-env](https://github.com/jpfuentes2/go-env),
where environment variables can be loaded into the application's runtime
environment through the definition of a `.env` or `.env.*` file. However, this
package supports a recursive autoload up to the calling project's root.
Recursive autoload allows you to, for example, define an `.env.test` at the
root and a test file at `<PROJECT_ROOT>/pkg/some_test.go` load the environment
variables correctly.

## Installation

`go get github.com/jmataya/dotenv`

## Usage

### Specify an Environment File

The most straightforward way to get started is to exactly specify the
environment file you want to load.

```golang
package main

import "github.com/jmataya/renv"

func main() {
        err := renv.LoadEnv("/path/to/env/file")
        if err != nil {
                // Do something...
        }
}
```

## Authors

* Jeff Mataya - [@jmataya](https://github.com/jmataya)

## Credit

Jacques Fuentes [@jpfuentes2](https://github.com/jpfuentes2) has an excellent
package called [go-env](https://github.com/jpfuentes2/go-env) that I've been
referencing for years.

_renv_ grew out of a desire to take that same functionality and allow
executables that may not be at a project root to leverage a project-wide .env
file.

## License

MIT
