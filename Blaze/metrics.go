package main

import (
	"os"
)

func getFileAgeMetric(path string) (int64, error) {

	fileStats, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	time := fileStats.ModTime()

	mtime := time.Unix()

	return mtime, nil
}

func getFileSizeMetric(file string) (int64, error) {

	fileStats, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	fileSize := fileStats.Size()

	return fileSize, nil
}
