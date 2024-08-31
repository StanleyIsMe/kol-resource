package config

import "fmt"

type MissingBaseConfigError struct {
	Err error
}

func (e *MissingBaseConfigError) Error() string {
	return fmt.Sprintf("missing base config: %v", e.Err)
}

type MissingEnvConfigError struct {
	Env string
	Err error
}

func (e *MissingEnvConfigError) Error() string {
	return fmt.Sprintf("missing config %s: %v", e.Env, e.Err)
}
