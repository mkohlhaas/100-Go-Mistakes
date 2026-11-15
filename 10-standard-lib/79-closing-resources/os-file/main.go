package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func listing1(filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %v\n", err)
		}
	}()

	return nil
}

func writeToFile1(filename string, content []byte) (err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	_, err = f.Write(content)
	return
}

func writeToFile2(filename string, content []byte) (err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return f.Sync()
}

func writeToFile(filename string, content []byte) (err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return
	}

	defer func() {
		if e := f.Close(); e != nil {
			if err == nil {
				err = e
			}
		}
	}()

	if _, err = f.Write(content); err != nil {
		return
	}

	if err = f.Sync(); err != nil {
		return
	}

	return
}

func writeToFileJoined(filename string, content []byte) (err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return
	}

	defer func() {
		if e := f.Close(); e != nil {
			if err == nil {
				err = e
			} else {
				err = errors.Join(err, e)
			}
		}
	}()

	if _, err = f.Write(content); err != nil {
		return
	}

	if err = f.Sync(); err != nil {
		return
	}

	return
}

func main() {
	content := []byte("hello\n")

	if err := listing1("/tmp/opencode/listing1.txt"); err != nil {
		log.Printf("listing1: %v", err)
	} else {
		fmt.Println("listing1: ok")
	}

	if err := writeToFile1("/tmp/opencode/file1.txt", content); err != nil {
		log.Printf("writeToFile1: %v", err)
	} else {
		fmt.Println("writeToFile1: ok")
	}

	if err := writeToFile2("/tmp/opencode/file2.txt", content); err != nil {
		log.Printf("writeToFile2: %v", err)
	} else {
		fmt.Println("writeToFile2: ok")
	}

	if err := writeToFile("/tmp/opencode/file_fixed.txt", content); err != nil {
		log.Printf("writeToFile: %v", err)
	} else {
		fmt.Println("writeToFile: ok")
	}

	if err := writeToFileJoined("/tmp/opencode/file_joined.txt", content); err != nil {
		log.Printf("writeToFileJoined: %v", err)
	} else {
		fmt.Println("writeToFileJoined: ok")
	}
}
