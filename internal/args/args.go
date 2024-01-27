package args

import (
	"errors"
	"fmt"
	"os"
)

const (
	SOFT = 0
	HARD = 1
)

type AppFlags struct {
	flagValueMap map[string]string
}

func NewAppFlags() AppFlags {
	fMap := make(map[string]string)
	return AppFlags{
		flagValueMap: fMap,
	}
}

func (f AppFlags) WithFlag(flag string, initialValue string) AppFlags {
	f.flagValueMap[flag] = initialValue
	return AppFlags{
		flagValueMap: f.flagValueMap,
	}
}

type AppArgs struct {
	flags AppFlags
	mode  int
}

func NewAppArgs(flags AppFlags, mode int) (AppArgs, error) {
	args := AppArgs{
		flags: flags,
		mode:  mode,
	}
	if mode != HARD && mode != SOFT {
		return args, errors.New("Invalid Arg mode provided")

	}
	err := args.initialise()
	if err != nil {
		return args, err
	}

	return args, nil
}

func (a *AppArgs) initialise() error {
	// NOTE: Skipping Program name
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		_, ok := a.flags.flagValueMap[args[i]]
		if ok && (i != len(args)-1) {
			_, ok := a.flags.flagValueMap[args[i+1]]
			if !ok {
				a.flags.flagValueMap[args[i]] = args[i+1]
				i += 1
			} else {
				if a.mode == HARD {
					return errors.New(fmt.Sprintf("No value provided for flag %s", args[i]))
				}
			}
		} else {
			if a.mode == HARD {
				return errors.New(fmt.Sprintf("No value provided for flag %s", args[i]))
			}
		}
	}

	// NOTE: This is to make sure all of the registered flags have a non empty value if in `HARD` Mode
	if a.mode == HARD {
		for k, v := range a.flags.flagValueMap {
			if v == "" {
				return errors.New(fmt.Sprintf("No value provided for flag %s", k))
			}
		}
	}

	return nil
}

func (a *AppArgs) GetFlagValue(flag string) (string, error) {
	val, ok := a.flags.flagValueMap[flag]
	if ok {
		return val, nil
	}
	return "", errors.New("Invalid flag requested")
}
