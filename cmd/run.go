package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {
	if e := verifyFlags(); e != nil {
		return fmt.Errorf("glazer: %w", e)
	}

	// keep usage silence after the provided flags are passing validation
	cmd.SilenceUsage = true
	writer := os.Stdout

	rootTemp := readTemplate("root", inputFile)
	md, err := os.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	// render faker functions within JSON
	fd, err := generateFakerData(md)
	if err != nil {
		return err
	}
	fmt.Println(fd)

	if err := rootTemp.ExecuteTemplate(writer, filepath.Base(inputFile), fd); err != nil {
		return err
	}

	return nil
}

func verifyFlags() error {
	if inputFile == "" {
		return errors.New("input file cannot be empty")
	}
	if jsonFile == "" {
		return errors.New("mock config file cannot be empty")
	}

	return nil
}

func readTemplate(name string, files ...string) *template.Template {
	return template.Must(template.New(name).ParseFiles(files...))
}
