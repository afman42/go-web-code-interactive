package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"

	"os"
	"os/exec"
	"time"
)

const WindowsFolderTmp = "\\tmp\\"
const LinuxFolderTmp = "/tmp/"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func Shellout(language string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(language, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func StringWithCharset(length int) string {
	charset := "abcdefgABCDEFG123456"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// https://stackoverflow.com/questions/50740902/move-a-file-to-a-different-drive-with-go
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source file: %v", err)
	}
	return nil
}

func PathFileTemp(filename string) string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	os := runtime.GOOS
	dir = dir + LinuxFolderTmp
	if os == "windows" {
		dir = dir + WindowsFolderTmp
	}
	return filepath.Join(dir + filename)
}

func CheckIsNotData(slice []string, find string) bool {
	for _, v := range slice {
		if v == find {
			return true
		}
	}
	return false
}
