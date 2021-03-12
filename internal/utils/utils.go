package utils

import (
	"fmt"
	"log"
	"os"
)

// GetRelativeDirectory tries to find the path to the relative path given two different ways. First it checks if they're
// near the executable. If it doesn't find them near the executable it will search the working directory. If they're not found in either
// location it returns an error.
func GetRelativeDirectory(relativePath string) (string, error) {
	// Check for frontend relative to the executable.
	execDir, err := os.Executable()
	if err != nil {
		log.Println(err)
	} else {
		dir := fmt.Sprintf("%s/%s", execDir, relativePath)
		_, err := os.Stat(dir)
		if err == nil {
			return dir, nil
		}
	}
	// Check for frontend files relative to the working directory.
	workDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	} else {
		dir := fmt.Sprintf("%s/%s", workDir, relativePath)
		_, err := os.Stat(dir)
		if err == nil {
			return dir, nil
		}
	}
	return "", fmt.Errorf("directory not found")
}
