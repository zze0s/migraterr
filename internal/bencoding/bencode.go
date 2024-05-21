package bencoding

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeebo/bencode"
)

type Replacement struct {
	Old string
	New string
}

type BencodeFile map[string]any

func ProcessFile(reader io.Reader, filePath, exportPath string, replacements []Replacement, verbose, dry bool, currentIteration string) (bool, error) {
	if dry {
		if verbose {
			log.Printf("%s dry-run: process file %s\n", currentIteration, filePath)
		}
	}

	var decodedFile BencodeFile
	decoder := bencode.NewDecoder(reader)

	err := decoder.Decode(&decodedFile)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false, nil
		}
		return false, errors.Wrapf(err, "could not decode file: %s", filePath)
	}

	modified := false

	if len(replacements) > 0 {
		for k, v := range decodedFile {
			for _, replacement := range replacements {

				switch val := v.(type) {
				case string:
					if strings.Contains(val, replacement.Old) {
						if dry {
							if verbose {
								log.Printf("%s dry-run: replaced: %q with %q\n", currentIteration, replacement.Old, replacement.New)
							}
						} else {
							decodedFile[k] = strings.Replace(val, replacement.Old, replacement.New, -1)

							if verbose {
								log.Printf("%s replaced: %q with %q\n", currentIteration, replacement.Old, replacement.New)
							}
						}

						modified = true
					}
				default:
					continue
				}
			}
		}
	}

	if !modified {
		log.Printf("%s found no data to replace\n", currentIteration)
	}

	if dry {
		if verbose {
			log.Printf("%s dry-run: encode and write file %s\n", currentIteration, exportPath)
		}
	} else {
		exportFile, err := os.Create(exportPath)
		if err != nil {
			return modified, errors.Wrapf(err, "could not create file: %s", exportPath)
		}

		defer exportFile.Close()

		err = decodedFile.EncodeAndWriteFile(exportFile)
		if err != nil {
			return modified, errors.Wrapf(err, "could not write fastresume file %s", exportPath)
		}
	}

	return modified, nil
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

func (f BencodeFile) EncodeAndWriteFile(file *os.File) error {
	bufferedWriter := bufio.NewWriter(file)
	encoder := bencode.NewEncoder(bufferedWriter)
	if err := encoder.Encode(f); err != nil {
		return errors.Wrapf(err, "could not encode file: %s", file.Name())
	}

	if err := bufferedWriter.Flush(); err != nil {
		return errors.Wrapf(err, "could not flush buffered writer: %s", file.Name())
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

	var fastResume map[string]any
	if err := bencode.DecodeString(string(read), &fastResume); err != nil {
		log.Printf("could not decode bencode file %s\n", path)
	}

	_, fileName := filepath.Split(path)

	var builder strings.Builder

	fmt.Fprintf(&builder, "\nFilename: %s\n", fileName)

	for k, v := range fastResume {
		fmt.Fprintf(&builder, "%s: %v\n", k, v)
	}

	fmt.Print(builder.String())

	return nil
}
