package bencoding

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeebo/bencode"
)

func Process(path, export string, replacements []string, verbose bool) error {
	read, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading file: %s err: %q\n", path, err)
	}

	var fastResume map[string]interface{}
	if err := bencode.DecodeString(string(read), &fastResume); err != nil {
		log.Printf("could not decode bencode file: %s\n", path)
	}

	if len(replacements) > 0 {
		for k, v := range fastResume {
			for _, replacement := range replacements {
				if !strings.Contains(replacement, "|") {
					continue
				}

				parts := strings.Split(replacement, "|")

				if len(parts) == 0 || len(parts) > 2 {
					continue
				}

				switch val := v.(type) {
				case string:
					if strings.Contains(val, parts[0]) {
						fastResume[k] = strings.Replace(val, parts[0], parts[1], -1)
					}
				default:
					continue
				}

				if verbose {
					fmt.Printf("replaced: '%s' with '%s'\n", parts[0], parts[1])
				}
			}
		}
	}

	if export != "" {
		if err = Encode(export, fastResume); err != nil {
			log.Printf("could not export fastresume file %s error: %q\n", path, err)
			return err
		}
	} else {
		if err = Encode(path, fastResume); err != nil {
			log.Printf("could not write fastresume file %s error: %q\n", path, err)
			return err
		}
	}

	if verbose {
		log.Printf("sucessfully processed file %s\n", path)
	}

	return nil
}

func Encode(path string, data any) error {
	file, err := os.Create(path)
	if err != nil {
		log.Printf("os create error: %q\n", err)
		return err
	}

	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	encoder := bencode.NewEncoder(bufferedWriter)
	if err := encoder.Encode(data); err != nil {
		log.Printf("encode error: %q\n", err)
		return err
	}

	if err := bufferedWriter.Flush(); err != nil {
		return err
	}

	return nil
}

func Info(path string) error {
	read, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading file: %s err: %q\n", path, err)
	}

	var fastResume map[string]interface{}
	if err := bencode.DecodeString(string(read), &fastResume); err != nil {
		log.Printf("could not decode bencode file %s\n", path)
	}

	_, fileName := filepath.Split(path)

	fmt.Printf("\nFilename: %s\n", fileName)
	for k, v := range fastResume {
		fmt.Printf("%s: %v\n", k, v)
	}
	fmt.Printf("\n")

	return nil
}
