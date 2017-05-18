package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                       Copyright (c) 2006-2017 FB GROUP LLC                         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/knf"
	"pkg.re/essentialkaos/ek.v9/options"
	"pkg.re/essentialkaos/ek.v9/usage"

	"pkg.re/essentialkaos/go-simpleyaml.v1"

	"github.com/funbox/init-exporter/procfile"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App props
const (
	APP  = "init-exporter-converter"
	VER  = "0.8.0"
	DESC = "Utility for converting procfiles from v1 to v2 format"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Supported arguments
const (
	OPT_CONFIG    = "c:config"
	OPT_IN_PLACE  = "i:in-place"
	OPT_NO_COLORS = "nc:no-colors"
	OPT_HELP      = "h:help"
	OPT_VERSION   = "v:version"
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
	OPT_CONFIG:    {},
	OPT_IN_PLACE:  {Type: options.BOOL},
	OPT_NO_COLORS: {Type: options.BOOL},
	OPT_HELP:      {Type: options.BOOL},
	OPT_VERSION:   {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	runtime.GOMAXPROCS(1)

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		fmtc.Println("Error while options parsing:")

		for _, err := range errs {
			fmtc.Printf("  %v\n", err)
		}

		os.Exit(1)
	}

	if options.GetB(OPT_NO_COLORS) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VERSION) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	process(args[0])
}

// process start data processing
func process(file string) {
	var err error

	if options.Has(OPT_CONFIG) {
		err = knf.Global(options.GetS(OPT_CONFIG))

		if err != nil {
			printErrorAndExit(err.Error())
		}
	}

	err = convert(file)

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// convert read procfile in v1 format and print v2 data or save it to file
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
		printErrorAndExit("Given procfile already converted to v2 format")
	}

	config.WorkingDir, hasCustomWorkingDirs = getWorkingDir(app)

	validateApplication(app)

	yamlData := renderProcfile(&procData{config, app, hasCustomWorkingDirs})

	err = validateYaml(yamlData)

	if err != nil {
		printErrorAndExit("Can't convert given procfile to YAML: %v", err)
	}

	if !options.GetB(OPT_IN_PLACE) {
		fmt.Printf(yamlData)
		return nil
	}

	return writeData(file, yamlData)
}

// renderProcfile render procfile
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

// getWorkingDir return path to default working dir and flag
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

// validateApplication validate application and all services
func validateApplication(app *procfile.Application) {
	errs := app.Validate()

	if len(errs) == 0 {
		return
	}

	printError("Errors while application validation:")

	for _, err := range errs {
		printError("  - %v", err)
	}

	os.Exit(1)
}

// validateYaml validate rendered yaml
func validateYaml(data string) error {
	_, err := simpleyaml.NewYaml([]byte(data))

	return err
}

// writeData write procfile data to file
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

// showUsage print usage info to console
func showUsage() {
	info := usage.NewInfo("", "procfile")

	info.AddOption(OPT_CONFIG, "Path to init-exporter config", "file")
	info.AddOption(OPT_IN_PLACE, "Edit procfile in place")
	info.AddOption(OPT_NO_COLORS, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VERSION, "Show version")

	info.AddExample(
		"-i config/Procfile.production",
		"Convert Procfile.production to version 2 in-place",
	)

	info.AddExample(
		"config/Procfile.production -c /etc/init-exporter.conf Procfile.production",
		"Convert Procfile.production to version 2 with defaults from init-exporter config and print result to console",
	)

	info.Render()
}

// showAbout print version info to console
func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "FB Group",
		License: "MIT License",
	}

	about.Render()
}
