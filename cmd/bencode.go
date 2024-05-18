package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"migraterr/internal/bencoding"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func BencodeCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "bencode",
		Short: "bencoded tools",
		Long:  `bencoded tools for editing and showing files like .torrents and fastresume`,
		Example: `  migraterr bencode 
  migraterr bencode --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(BencodeEditCommand())
	command.AddCommand(BencodeInfoCommand())

	return command
}

func BencodeEditCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "edit",
		Short: "Edit bencode files",
		Long:  `Edit bencoded files like .torrents and fastresume`,
		Example: `  migraterr bencode edit ./files/**/*.fastresume --replace "/old/path/|/new/path"
  migraterr bencode edit --help`,
		Args: cobra.MinimumNArgs(1),
	}

	var (
		dry            bool
		verbose        bool
		exportDir      string
		replaceStrings []string
		replacements   []bencoding.Replacement
	)

	command.Flags().BoolVar(&dry, "dry-run", false, "Dry run, don't write changes")
	command.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	command.Flags().StringVar(&exportDir, "export", "", "Export to directory. Will replace if not specified")
	command.Flags().StringSliceVar(&replaceStrings, "replace", []string{}, "Replace: pattern|replace")

	command.MarkFlagRequired("replace")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		path := args[0]

		if strings.HasPrefix(path, "~/") {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, path[2:])
		}

		if len(replaceStrings) == 0 {
			log.Fatal("must supply replacements\n")
		}

		// validate replaceStrings
		if len(replaceStrings) > 0 {
			for _, replacement := range replaceStrings {
				if !strings.Contains(replacement, "|") {
					return errors.Errorf("invalid replacement: %q", replacement)
				}

				parts := strings.Split(replacement, "|")
				if len(parts) > 2 {
					return errors.Errorf("invalid replacement: %q - contains multiple |", replacement)
				}

				replacements = append(replacements, bencoding.Replacement{
					Old: parts[0],
					New: parts[1],
				})
			}
		}

		processedFiles := 0

		files, err := filepath.Glob(path)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			return errors.Errorf("no files found matching pattern: %q", path)
		}

		if exportDir != "" {
			dir := filepath.Dir(exportDir)
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "could not create export dir: %q", exportDir)
			}
		}

		for idx, file := range files {
			_, fileName := filepath.Split(file)

			exportPath := file
			if exportDir != "" {
				exportPath = filepath.Join(exportDir, fileName)
			}

			if verbose {
				log.Printf("[%d/%d] reading file %s\n", idx, len(files), fileName)
			}

			err = openAndProcessFile(file, exportPath, replacements, verbose, dry)
			if err != nil {
				return errors.Wrapf(err, "could not open and process file: %q", file)
			}

			processedFiles++

			if verbose {
				log.Printf("[%d/%d] sucessfully processed file %s\n", processedFiles, len(files), fileName)
			}
		}

		log.Printf("migraterr bencode processed %d files\n", processedFiles)

		return nil
	}

	return command
}

func openAndProcessFile(filePath string, exportPath string, replacements []bencoding.Replacement, verbose bool, dry bool) error {
	fileReader, err := os.Open(filePath)
	if err != nil {
		return errors.Wrapf(err, "could not read filePath: %s", filePath)
	}

	defer fileReader.Close()

	err = bencoding.ProcessFile(fileReader, filePath, exportPath, replacements, verbose, dry)
	if err != nil {
		return errors.Wrapf(err, "error processing filePath: %s", filePath)
	}

	return nil
}

func BencodeInfoCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "info",
		Short: "Info bencode files",
		Long:  `Info bencoded files like .torrents and fastresume`,
		Example: `  migraterr bencode info ./files/**/*.fastresume
  migraterr bencode info --help`,
		Args: cobra.MinimumNArgs(1),
	}

	command.Run = func(cmd *cobra.Command, args []string) {
		path := args[0]

		if strings.HasPrefix(path, "~/") {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, path[2:])
		}

		files, err := filepath.Glob(path)
		if err != nil {
			log.Fatal("could not open files\n")
		}

		fmt.Printf("Found %d files matching %q\n", len(files), path)

		for _, file := range files {
			if err := bencoding.Info(file); err != nil {
				log.Fatalf("error reading file: %s  %v", file, err)
			}
		}

	}

	return command
}
