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
		datasetID := datasetPrefix + ":" + ss[0]
		name := ss[1]
		fmt.Println("dataset-id:", datasetID)
		// The same set of datasets has the same regex
		for k, v := range mapDataset {
			fmt.Printf("k:%s,v:%v\n", k, v)
			if strings.HasPrefix(k, datasetID) {
				dataset = v
				break
			}
		}

		// dataset = mapDataset[datasetID]
		logrus.Printf("dataset-id:%s,dataset:%v, dataset-map:%v", datasetID, dataset, mapDataset)
		if dataset == nil {
			fmt.Fprintf(os.Stderr, "dataset %s not found\n", datasetID)
			os.Exit(1)
		}
		entity := dataset.parseEntity(name)
		entity.datasetID = datasetID + ":" + entity.datasetID

		datasetID = entity.datasetID
		// set the real dataset
		dataset = mapDataset[datasetID]
		dataset.addEntity(entity)
		logrus.Printf("entity:%v\n", entity)

		if dataset.checkNewGroupAvailable(entity) {
			logrus.Println("new-group added.")
			dataset.outputNewGroups(entity)
		}
	}

	os.Exit(code)
}
