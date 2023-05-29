package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"migraterr/internal/bencoding"

	"github.com/spf13/cobra"
)

func RunBencode() *cobra.Command {

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

	command.AddCommand(RunBencodeEdit())
	command.AddCommand(RunBencodeInfo())

	return command
}

func RunBencodeEdit() *cobra.Command {

	command := &cobra.Command{
		Use:   "edit",
		Short: "Edit bencode files",
		Long:  `Edit bencoded files like .torrents and fastresume`,
		Example: `  migraterr bencode edit
  migraterr bencode edit --help`,
		SilenceUsage: true,
	}

	var (
		dry          bool
		verbose      bool
		export       string
		glob         string
		replacements []string
	)

	command.Flags().BoolVar(&dry, "dry-run", false, "Dry run, don't write changes")
	command.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	command.Flags().StringVar(&glob, "glob", "", "Glob to files: eg ./files/**/*.torrent.rtorrent")
	command.Flags().StringVar(&export, "export", "", "Export to directory. Will replace if not specified")
	command.Flags().StringSliceVar(&replacements, "replace", []string{}, "Replace: pattern|replace")

	command.Run = func(cmd *cobra.Command, args []string) {
		if glob == "" {
			log.Fatal("must have dir\n")
		}

		if len(replacements) == 0 {
			log.Fatal("must supply replacements\n")
		}

		processedFiles := 0

		files, err := filepath.Glob(glob)
		if err != nil {
			log.Fatal("could not open files\n")
		}

		for _, file := range files {
			_, fileName := filepath.Split(file)

			exportFile := ""
			if export != "" {
				exportFile = filepath.Join(export, fileName)
			}

			if err := bencoding.Process(file, exportFile, replacements, verbose); err != nil {
				log.Fatalf("error processing file: %q\n", err)
			}

			processedFiles++

			if verbose {
				fmt.Printf("[%d/%d] sucessfully processed file %s\n", len(files), processedFiles, fileName)
			}

		}

		fmt.Printf("migraterr bencode processed %d files\n", processedFiles)
	}

	return command
}

func RunBencodeInfo() *cobra.Command {

	command := &cobra.Command{
		Use:   "info",
		Short: "Info bencode files",
		Long:  `Info bencoded files like .torrents and fastresume`,
		Example: `  migraterr bencode info
  migraterr bencode info --help`,
		SilenceUsage: true,
	}

	var glob string

	command.Flags().StringVar(&glob, "glob", "", "Glob to files: eg ./files/**/*.torrent.rtorrent")

	command.Run = func(cmd *cobra.Command, args []string) {
		if glob == "" {
			log.Fatal("must have dir\n")
		}

		files, err := filepath.Glob(glob)
		if err != nil {
			log.Fatal("could not open files\n")
		}

		fmt.Printf("found %d files\n", len(files))

		for _, file := range files {
			if err := bencoding.Info(file); err != nil {
				log.Fatalf("error reading file: %s  %v", file, err)
			}
		}

	}

	return command
}
