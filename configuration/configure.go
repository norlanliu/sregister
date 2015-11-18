//Copyright (c) 2015 Qi Liu AT ICT
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.
package configuration

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	usageLine = `usage:	sregister [flags]
	start SRegister

	sregister --version
	show the version of SRegister
	
	sregister -h | --help
	show the help of SRegister`

	flagLine = `member flags:`

	logDirFlag         = "log_dir"
	servicesDirFlag    = "services_dir"
	version            = `SFinder/SRegister 0.1`
	defaultServicesDir = `/etc/sfinder/sregister/services`
	defaultLogDir      = `/var/lib/sfinder/sregister/log`
)

type configure struct {
	fs             *flag.FlagSet
	serviceConfDir string
	logDir         string
	version        bool
}

func NewConfigure() *configure {
	cfg := &configure{}
	cfg.fs = flag.NewFlagSet("sregister", flag.ContinueOnError)
	fs := cfg.fs

	fs.Usage = func() {
		fmt.Println(usageLine)
	}
	//command flags
	fs.StringVar(&cfg.serviceConfDir, servicesDirFlag, defaultServicesDir, "Path to services configuration directory")
	fs.StringVar(&cfg.logDir, logDirFlag, defaultLogDir, "Path to log configuration directory")
	fs.BoolVar(&cfg.version, "version", false, "show the version")

	return cfg
}

func (cfg *configure) ParseConfigure(arguments []string) error {
	perr := cfg.fs.Parse(arguments)
	switch perr {
	case nil:
	case flag.ErrHelp:
		fmt.Println(flagLine)
		cfg.fs.PrintDefaults()
		os.Exit(0)
	default:
		os.Exit(2)
	}

	if len(cfg.fs.Args()) != 0 {
		fmt.Errorf("'%s' is not a valid flag", cfg.fs.Arg(0))
	}

	if cfg.version {
		fmt.Println(version)
	}

	err := parseFromEnv(cfg)

	if err != nil {
		fmt.Errorf("%v", err)
	}

	flag.Set(logDirFlag, cfg.logDir)

	flag.Parse()

	return err
}

func (cfg *configure) GetServiceConfDir() string {
	return cfg.serviceConfDir
}

func parseFromEnv(cfg *configure) error {

	fs := cfg.fs

	defaultConf := map[string]string{
		servicesDirFlag: defaultServicesDir,
		logDirFlag:      defaultLogDir,
	}

	var err error
	alreadySet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})
	fs.VisitAll(func(f *flag.Flag) {
		if !alreadySet[f.Name] {
			key := flagToEnv(f.Name)
			value := os.Getenv(key)
			if value != "" {
				if serr := fs.Set(f.Name, value); serr != nil {
					err = fmt.Errorf("invalid value %q for %s: %v", value, key, serr)
				}
			} else {
				if serr := fs.Set(f.Name, defaultConf[f.Name]); serr != nil {
					err = fmt.Errorf("invalid value %q for %s: %v", value, key, serr)
				}
			}
		}
	})
	return err
}

func flagToEnv(flag string) string {
	return "SREG_" + strings.ToUpper(flag)
}
