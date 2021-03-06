package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/clementauger/sqlg/parse"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cwd not found: %v", err)
	}
	// fmt.Println(inputFile)

	var clean bool
	var dryrun bool
	var engine string
	{
		c := flag.Bool("clean", false, "delete generated files")
		dry := flag.Bool("dry-run", false, "dry run, dont write files, just print the result")
		e := flag.String("engine", "", "target database engine")
		flag.Parse()
		clean = *c
		dryrun = *dry
		engine = *e
	}

	if clean {
		rmGeneratedFiles(flag.Args(), dryrun)
		return
	}
	generate(cwd, engine, dryrun)
}

func rmGeneratedFiles(files []string, dryrun bool) {

	var toRm []string
	for _, f := range files {
		s, err := os.Stat(f)
		if err != nil {
			log.Fatalf("stat: %v", err)
		}
		if s.IsDir() {
			sfiles, err := ioutil.ReadDir(f)
			if err != nil {
				log.Fatalf("read dir: %v", err)
			}
			for _, ff := range sfiles {
				if !ff.IsDir() && strings.HasSuffix(ff.Name(), "_gen.go") {
					toRm = append(toRm, filepath.Join(f, ff.Name()))
				}
			}
		} else if strings.HasSuffix(f, "_gen.go") {
			toRm = append(toRm, f)
		}
	}

	for _, file := range toRm {
		if ok, err := bodyHasGenComment(file); err != nil {
			log.Fatalf("read file: %v", err)
		} else if !ok {
			log.Println("skip ", file)
			continue
		}
		log.Println("rm   ", file)
		if !dryrun {
			err := os.Remove(file)
			if err != nil {
				log.Fatalf("rm: %v", err)
			}
		}
	}
}

func bodyHasGenComment(file string) (bool, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(b), "// Code generated by sqlg DO NOT EDIT"), nil
}

func generate(cwd, engine string, dryrun bool) {

	fileObjects, err := parse.FromPackagePath(cwd, engine, flag.Args())
	if err != nil {
		log.Fatalf("parsing failure: %v", err)
	}

	if len(fileObjects) < 1 {
		log.Println("no types to generate found")
	}
	var out map[string]string
	out, err = parse.Generate(fileObjects, engine)
	if err != nil {
		log.Fatalf("generation failure: %v", err)
	}

	for file, content := range out {
		if dryrun {
			fmt.Fprintln(os.Stderr, file)
			fmt.Fprintln(os.Stderr, content)
		} else {
			d := filepath.Dir(file)
			err = os.MkdirAll(d, os.ModePerm)
			if err != nil {
				log.Fatalf("mkdir failure: %v", err)
			}
			err = ioutil.WriteFile(file, []byte(content), os.ModePerm)
			if err != nil {
				log.Fatalf("write failure: %v", err)
			}
		}
	}
}
