package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                           Copyright (c) 2006-2023 FUNBOX                           //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/knf"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"

	"github.com/essentialkaos/go-simpleyaml/v2"

	"github.com/funbox/init-exporter/procfile"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App props
const (
	APP  = "init-exporter-converter"
	VER  = "0.12.0"
	DESC = "Utility for converting procfiles from v1 to v2 format"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Supported arguments
const (
	OPT_CONFIG   = "c:config"
	OPT_IN_PLACE = "i:in-place"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// Config properies
const (
	MAIN_PREFIX               = "main:prefix"
	PATHS_WORKING_DIR         = "paths:working-dir"
	DEFAULTS_NPROC            = "defaults:nproc"
	DEFAULTS_NOFILE           = "defaults:nofile"
	DEFAULTS_RESPAWN          = "defaults:respawn"
	DEFAULTS_RESPAWN_COUNT    = "defaults:respawn-count"
	DEFAULTS_RESPAWN_INTERVAL = "defaults:respawn-interval"
	DEFAULTS_KILL_TIMEOUT     = "defaults:kill-timeout"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// DEFAULT_WORKING_DIR is path to default working dir
const DEFAULT_WORKING_DIR = "/tmp"

// ////////////////////////////////////////////////////////////////////////////////// //

type procData struct {
	Config               *procfile.Config
	Application          *procfile.Application
	HasCustomWorkingDirs bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_CONFIG:   {},
	OPT_IN_PLACE: {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.BOOL},

	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

var colorTagApp string
var colorTagVer string

var gitrev string

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	runtime.GOMAXPROCS(1)

	preConfigureUI()

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitrev).Print()
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	process(args.Strings())
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	term := os.Getenv("TERM")

	fmtc.DisableColors = true

	if term != "" {
		switch {
		case strings.Contains(term, "xterm"),
			strings.Contains(term, "color"),
			term == "screen":
			fmtc.DisableColors = false
		}
	}

	if !fsutil.IsCharacterDevice("/dev/stdout") && os.Getenv("FAKETTY") == "" {
		fmtc.DisableColors = true
	}

	if os.Getenv("NO_COLOR") != "" {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{#BCCF00}", "{#BCCF00}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{#148}", "{#148}"
	default:
		colorTagApp, colorTagVer = "{g}", "{g}"
	}
}

// process starts data processing
func process(files []string) {
	var err error
	var hasErrors bool

	if options.Has(OPT_CONFIG) {
		err = knf.Global(options.GetS(OPT_CONFIG))

		if err != nil {
			printErrorAndExit(err.Error())
		}
	}

	for _, file := range files {
		err = convert(file)

		if err != nil {
			hasErrors = true
			printError(err.Error())
		}
	}

	if hasErrors {
		os.Exit(1)
	}
}

// convert reads procfile in v1 format and print v2 data or save it to file
func convert(file string) error {
	var hasCustomWorkingDirs bool

	config := &procfile.Config{
		Name:             "",
		WorkingDir:       knf.GetS(PATHS_WORKING_DIR, "/tmp"),
		IsRespawnEnabled: knf.GetB(DEFAULTS_RESPAWN, true),
		RespawnInterval:  knf.GetI(DEFAULTS_RESPAWN_INTERVAL, 15),
		RespawnCount:     knf.GetI(DEFAULTS_RESPAWN_COUNT, 10),
		KillTimeout:      knf.GetI(DEFAULTS_KILL_TIMEOUT, 60),
		LimitFile:        knf.GetI(DEFAULTS_NOFILE, 10240),
		LimitProc:        knf.GetI(DEFAULTS_NPROC, 10240),
	}

	app, err := procfile.Read(file, config)

	if err != nil {
		return err
	}

	if app.ProcVersion != 1 {
		return fmt.Errorf("Given procfile already converted to v2 format")
	}

	config.WorkingDir, hasCustomWorkingDirs = getWorkingDir(app)

	err = validateApplication(app)

	if err != nil {
		return err
	}

	yamlData := renderProcfile(&procData{config, app, hasCustomWorkingDirs})

	err = validateYaml(yamlData)

	if err != nil {
		return fmt.Errorf("Can't convert given procfile to YAML: %v", err)
	}

	if !options.GetB(OPT_IN_PLACE) {
		fmt.Printf(yamlData)
		return nil
	}

	return writeData(file, yamlData)
}

// renderProcfile renders procfile
func renderProcfile(data *procData) string {
	var result string

	result += "version: 2\n"
	result += "\n"
	result += "start_on_runlevel: 2\n"
	result += "stop_on_runlevel: 5\n"
	result += "\n"

	if data.Config.IsRespawnEnabled {
		result += "respawn:\n"
		result += fmt.Sprintf("  count: %d\n", data.Config.RespawnCount)
		result += fmt.Sprintf("  interval: %d\n", data.Config.RespawnInterval)
		result += "\n"
	}

	result += "limits:\n"
	result += fmt.Sprintf("  nofile: %d\n", data.Config.LimitFile)
	result += fmt.Sprintf("  nproc: %d\n", data.Config.LimitProc)
	result += "\n"

	if !data.HasCustomWorkingDirs {
		result += "working_directory: " + data.Config.WorkingDir + "\n"
		result += "\n"
	}

	result += "commands:\n"

	for _, service := range data.Application.Services {
		result += "  " + service.Name + ":\n"

		if service.HasPreCmd() {
			result += "    pre: " + service.PreCmd + "\n"
		}

		result += "    command: " + service.Cmd + "\n"

		if service.HasPostCmd() {
			result += "    post: " + service.PostCmd + "\n"
		}

		if data.HasCustomWorkingDirs {
			result += "    working_directory: " + service.Options.WorkingDir + "\n"
		}

		if service.Options.IsCustomLogEnabled() {
			result += "    log: " + service.Options.LogFile + "\n"
		}

		if service.Options.IsEnvSet() {
			result += "    env:\n"
			for k, v := range service.Options.Env {
				result += fmt.Sprintf("      %s: %s\n", k, v)
			}
		}

		result += "\n"
	}

	return result
}

// getWorkingDir returns path to default working dir and flag
// if custom working dirs is used
func getWorkingDir(app *procfile.Application) (string, bool) {
	var dir = DEFAULT_WORKING_DIR

	for _, service := range app.Services {
		if dir == DEFAULT_WORKING_DIR {
			if service.Options.WorkingDir != "" {
				dir = service.Options.WorkingDir
			}

			continue
		}

		if dir != service.Options.WorkingDir {
			return DEFAULT_WORKING_DIR, true
		}
	}

	return dir, false
}

// validateApplication validates application and all services
func validateApplication(app *procfile.Application) error {
	errs := app.Validate()

	if len(errs) == 0 {
		return nil
	}

	return errs[0]
}

// validateYaml validates rendered yaml
func validateYaml(data string) error {
	_, err := simpleyaml.NewYaml([]byte(data))

	return err
}

// writeData writes procfile data to file
func writeData(file, data string) error {
	return ioutil.WriteFile(file, []byte(data), 0644)
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "init-exporter-converter"))
	case "fish":
		fmt.Printf(fish.Generate(info, "init-exporter-converter"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "init-exporter-converter"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "procfile…")

	info.AppNameColorTag = "{*}" + colorTagApp

	info.AddOption(OPT_CONFIG, "Path to init-exporter config", "file")
	info.AddOption(OPT_IN_PLACE, "Edit procfile in place")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"-i config/Procfile.production",
		"Convert Procfile.production to version 2 in-place",
	)

	info.AddExample(
		"-c /etc/init-exporter.conf Procfile.*",
		"Convert all procfiles to version 2 with defaults from init-exporter config and print result to console",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "FunBox",
		License: "MIT License",
		UpdateChecker: usage.UpdateChecker{
			"funbox/init-exporter-converter",
			update.GitHubChecker,
		},

		AppNameColorTag: "{*}" + colorTagApp,
		VersionColorTag: colorTagVer,
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
