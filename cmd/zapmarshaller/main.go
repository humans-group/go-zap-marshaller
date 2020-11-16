package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mkorolyov/astparser"

	zapmarshaller "github.com/humans-group/go-zap-marshaller"
)

var (
	inputDir         = flag.String("in", "", "directory with go files to be parsed")
	excludeRegexpStr = flag.String("exclude", "", "exclude regexp to skip files")
	includeRegexpStr = flag.String("include", "", "include regexp to limit input files")
	outputDir        = flag.String("out", "", "output directory for generated files, without package")
	outPackage       = flag.String("package", "", "package for generated zap marshallers files")
)

func main() {
	flag.Parse()

	// load golang sources
	cfg := astparser.Config{InputDir: *inputDir}
	if *excludeRegexpStr != "" {
		cfg.ExcludeRegexp = *excludeRegexpStr
	}

	if *includeRegexpStr != "" {
		cfg.IncludeRegexp = *includeRegexpStr
	}

	sources, err := astparser.Load(cfg)
	if err != nil {
		log.Fatalf("failed to load sources from %s excluding %s: %v", *inputDir, *excludeRegexpStr, err)
	}

	generator := zapmarshaller.Generator{Cfg: zapmarshaller.Config{OutPackage: *outPackage}}
	if *outputDir != "" {
		zapFiles := generator.Generate(sources)
		// save
		for f, body := range zapFiles {
			filePath := *outputDir + "/" + strings.TrimSuffix(f, filepath.Ext(f)) + ".zap.go"

			if err := ioutil.WriteFile(filePath, body, 0600); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to save generated zap marshaller file %s: %v\n", filePath, err)
			}
		}
	}
}
