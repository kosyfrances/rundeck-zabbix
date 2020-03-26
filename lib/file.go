package lib

import (
	"fmt"
	"os"
)

// DumpToFile dumps data, a list of bytes to file given.
// If the filePath doesn't exist, it creates it, or appends to the file
func DumpToFile(filePath string, data []byte) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file. %v", err)
	}

	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("cannot write to file. %v", err)
	}
	return nil
}
