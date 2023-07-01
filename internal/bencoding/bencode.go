package bencoding

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeebo/bencode"
)

type BencodeFile map[string]interface{}

func Process(filePath, exportDir string, replacements []string, verbose, dry bool) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("error reading file: %s err: %q\n", filePath, err)
		return errors.Wrapf(err, "could not read file: %s", filePath)
	}

	if verbose {
		log.Printf("reading file %s\n", filePath)
	}

	exportPath := filePath
	if exportDir != "" {
		_, file := filepath.Split(filePath)
		exportPath = filepath.Join(exportDir, file)
	}

	if dry {
		if verbose {
			log.Printf("dry-run: process file %s\n", filePath)
			log.Printf("dry-run: encode and write file %s\n", exportPath)
		}
	} else {
		var decodedFile BencodeFile
		if err := bencode.DecodeBytes(file, &decodedFile); err != nil {
			log.Printf("could not decode bencode file: %s\n", filePath)
			return errors.Wrapf(err, "could not decode file: %s", filePath)
		}

		if len(replacements) > 0 {
			for k, v := range decodedFile {
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
							decodedFile[k] = strings.Replace(val, parts[0], parts[1], -1)
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

		if err := decodedFile.EncodeAndWrite(exportPath); err != nil {
			log.Printf("could not write fastresume file %s error: %q\n", exportPath, err)
			return err
		}
	}

	if verbose {
		log.Printf("sucessfully processed file %s\n", exportPath)
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

func (f BencodeFile) EncodeAndWrite(path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Printf("os create error: %q\n", err)
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		log.Printf("os create error: %q\n", err)
		return err
	}

	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	encoder := bencode.NewEncoder(bufferedWriter)
	if err := encoder.Encode(f); err != nil {
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
