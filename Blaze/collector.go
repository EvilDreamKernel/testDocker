package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"path/filepath"
)

type FileExporterCollector struct {
	rootPath string
	fileAgeMetric           *prometheus.Desc
	fileSizeMetric          *prometheus.Desc
	directoryElementsNumber *prometheus.Desc
}

func (collector *FileExporterCollector) getFileList(root_path string) ([]string, error) {
	var files []string

	err := filepath.Walk(root_path, func(path string, info os.FileInfo, err error) error {
		fileStat, _ := os.Stat(path)
		if !fileStat.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (collector *FileExporterCollector) isRootPathExists() (bool, error) {
	var exists bool

	fileStat, err := os.Stat(collector.rootPath)
	if err == nil {
		if fileStat.IsDir() {
			exists = true
		}
	} else {
		exists = false
		if !os.IsNotExist(err) {
			return exists, err
		}
	}

	return exists, nil
}

func newFileExporterCollector(path string) *FileExporterCollector {
	return &FileExporterCollector{
		rootPath: path,
		fileAgeMetric: prometheus.NewDesc(
			"file_age_unix",
			"Last modification time of file",
			[]string{"path"}, nil),
		fileSizeMetric: prometheus.NewDesc(
			"file_size_bytes",
			"File size in bytes",
			[]string{"path"}, nil),
		directoryElementsNumber: prometheus.NewDesc(
			"directory_elements_number",
			"Elements number in directory",
			[]string{"path"}, nil),
	}
}

func (collector *FileExporterCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.fileSizeMetric
	ch <- collector.fileAgeMetric
	ch <- collector.directoryElementsNumber
}

func (collector *FileExporterCollector) Collect(ch chan<- prometheus.Metric) {
	isRootPathExists, err := collector.isRootPathExists()
	if err != nil {
		log.Fatal(err)
	}
	if  isRootPathExists {
		fileList, err := collector.getFileList(collector.rootPath)
		if err != nil {
			log.Fatal(err)
		}
		ch <- prometheus.MustNewConstMetric(
			collector.directoryElementsNumber,
			prometheus.GaugeValue,
			float64(len(fileList)),
			collector.rootPath)
		for _, file := range fileList {
			fileAge, err := getFileAgeMetric(file)
			if err != nil {
				log.Fatal(err)
			}
			fileSize, err := getFileSizeMetric(file)
			if err != nil {
				log.Fatal(err)
			}
			ch <- prometheus.MustNewConstMetric(
				collector.fileAgeMetric,
				prometheus.GaugeValue,
				float64(fileAge),
				file)
			ch <- prometheus.MustNewConstMetric(
				collector.fileSizeMetric,
				prometheus.GaugeValue,
				float64(fileSize),
				file)
		}
	}
}