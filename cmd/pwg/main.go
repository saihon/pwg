package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/saihon/pwg/data"

	flag "github.com/saihon/flags"
	"github.com/saihon/pwg"
)

var (
	Name     = "pwg"
	Version  = "v0.0.1"
	options  *pwg.Options
	flagUser *flag.FlagSet
)

func init() {
	options = new(pwg.Options)

	flag.CommandLine.Init(Name, flag.ExitOnError, true)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [subcommand] [options...] [arguments...]\n\n", Name)
		flag.PrintCustom()
		fmt.Fprint(os.Stderr, `
Command:

  user, username
	Generate a string  like a username
	show more details. username --help

`)
	}

	flag.Bool("version", 'v', false, "Output version information", func(_ flag.Getter) error {
		fmt.Fprintf(os.Stderr, "%s %s\n", Name, Version)
		return flag.ErrHelp
	})

	flag.BoolVar(&options.All, "all", 'a', false, "Use all kinds of character\n( same as -lLns )", nil)
	flag.BoolVar(&options.Number, "number", 'n', false, "Use a numbers", nil)
	flag.BoolVar(&options.Symbol, "symbol", 's', false, "Use a symbols", nil)
	flag.BoolVar(&options.UpperCase, "uppercase", 'L', false, "Use upper case letters", nil)
	flag.BoolVar(&options.LowerCase, "lowercase", 'l', false, "Use lower case letters", nil)
	flag.IntVar(&options.Generate, "generate", 'g', 1, "Number of passwords to generate", nil)
	flag.IntVar(&options.Length, "length", 'd', 0, "Password length. default 6", nil)
	flag.BoolVar(&options.Evenly, "evenly", 'e', false, "Use character types as evenly\nas possible. default false", nil)
	flag.StringVar(&options.Custom, "custom", 'c', "", "Specify any character string\nto be used for password", nil)

	flagUser = flag.NewFlagSet(Name+" username", flag.ExitOnError, false)
	flagUser.Bool("version", 'v', false, "Output version information", func(_ flag.Getter) error {
		fmt.Fprintf(os.Stderr, "%s %s\n", Name, Version)
		return flag.ErrHelp
	})

	flagUser.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [options...] [arguments...]\n\n", flagUser.Name())
		flagUser.PrintCustom()
		fmt.Fprint(os.Stderr, "\n")
	}
	flagUser.BoolVar(&options.Random, "random", 'r', false, "Random username length", nil)
	flagUser.IntVar(&options.Length, "length", 'd', 0, "Username length. default 5", nil)
	flagUser.IntVar(&options.Generate, "generate", 'g', 1, "Number of username to generate", nil)
	flagUser.BoolVar(&options.Capitalize, "capitalize", 'c', false, "Capitalize the beginning of a username", nil)
}

func flagparse() ([]string, bool) {
	flag.Parse()
	if flag.NArg() == 0 {
		return flag.Args(), false
	}

	switch flag.Arg(0) {
	case "user", "username":
		flagUser.Parse(flag.Args()[1:])
		return flagUser.Args(), true
	default:
		flag.CommandLine.StopImmediate = false
		flag.CommandLine.Parse(flag.Args())
		return flag.Args(), false
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: recover: %v\n", err)
			os.Exit(1)
		}
	}()

	os.Exit(_main())
}

func _main() int {
	argv, username := flagparse()

	if options.Length < 1 {
		options.Length = 6
		if len(argv) > 0 {
			n, err := strconv.Atoi(argv[0])
			if err == nil {
				options.Length = n
			}
		}
	}

	p := pwg.New(options)

	switch {
	case username:
		for b := range p.Users(data.Data) {
			os.Stdout.Write(append(b, '\n'))
		}
	default:
		for b := range p.Gen() {
			os.Stdout.Write(append(b, '\n'))
		}
	}

	return 0
}
