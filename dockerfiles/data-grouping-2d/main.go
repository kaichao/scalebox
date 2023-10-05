package main

import (
	"fmt"
	"os"
	"strings"

	scalebox "github.com/kaichao/scalebox/golang/misc"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		logger.Fatalf("cmdline param needed\n")
		os.Exit(1)
	}

	keyText := os.Args[1]
	dataset := parseDataSet(keyText)
	var code int
	if dataset != nil {
		// new dataset
		mapDataset[dataset.DatasetID] = dataset
		fmt.Printf("new added dataset:%v\n", dataset)
		scalebox.AppendToFile(datasetFile, keyText)
		logrus.Println("dataset-map:", mapDataset)
	} else {
		// new entity
		ss := strings.Split(keyText, ",")
		datasetID := ss[0]
		name := ss[1]
		dataset = mapDataset[datasetID]
		fmt.Printf("dataset-id:%s,dataset:%v\n", datasetID, dataset)
		logrus.Println("dataset-map:", mapDataset)
		if dataset == nil {
			fmt.Fprintf(os.Stderr, "dataset %s not found\n", datasetID)
			os.Exit(1)
		}
		entity := dataset.parseEntity(name)
		dataset.addEntity(entity)
		logrus.Printf("504,%v", entity)
		if dataset.checkNewGroupAvailable(entity) {
			logrus.Println("505:new-group")
			dataset.outputNewGroups(entity)
		}
	}

	os.Exit(code)
}
