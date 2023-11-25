package main

var (
	jsonDataSet0 = ` {
		"datasetID":"ds0",
		"keyGroupRegex":"^([^/]+)/.+([0-9]{4}).+M([0-9][0-9]).+$",
		"keyGroupIndex":"1,3,2",
		"sinkJob" :"job-0",

		"groupType":"H",
		"horizontalWidth":3
	}
	`

	dataset0 = &DataSet{
		DatasetID:       "ds0",
		KeyGroupRegex:   "^([^/]+)/.+([0-9]{4}).+M([0-9][0-9]).+$",
		KeyGroupIndex:   "1,3,2",
		SinkJob:         "job-0",
		GroupType:       "H",
		HorizontalWidth: 3,
	}

	jsonDataSet1 = ` {
		"datasetID":"ds1",
		"keyGroupRegex":"^([^/]+)/.+([0-9]{4}).+M([0-9][0-9]).+$",
		"keyGroupIndex":"1,3,2",
		"sinkJob" :"job-1",

		"groupType":"V",
		"verticalStart":1,
		"verticalHeight":10,
		"groupSize":4,
		"interleaved":false
	}
	`

	dataset1 = &DataSet{
		DatasetID:      "ds1",
		KeyGroupRegex:  "^([^/]+)/.+([0-9]{4}).+M([0-9][0-9]).+$",
		KeyGroupIndex:  "1,3,2",
		SinkJob:        "job-1",
		GroupType:      "V",
		VerticalStart:  1,
		VerticalHeight: 10,
		GroupSize:      4,
	}

	jsonDataSet2 = ` {
		"datasetID":"ds2",
		"keyGroupRegex":"^([^/]+)/.+([0-9]{4}).+M([0-9][0-9]).+$",
		"keyGroupIndex":"1,3,2",
		"sinkJob" :"job-2",

		"groupType":"V",
		"verticalStart":1,
		"verticalHeight":10,
		"groupSize":4,
		"interleaved":true
	}
	`

	dataset2 = &DataSet{
		DatasetID:      "ds2",
		KeyGroupRegex:  "^([^/]+)/.+([0-9]{4}).+M([0-9][0-9]).+$",
		KeyGroupIndex:  "1,3,2",
		SinkJob:        "job-2",
		GroupType:      "V",
		VerticalStart:  1,
		VerticalHeight: 10,
		GroupSize:      4,
		Interleaved:    true,
	}

	jsonDataSetMWA = ` {
		"datasetID":"mwa-dat:1257010784",
		"keyGroupRegex":"^.+/([0-9]{10})_([0-9]{10})_.+$",
		"keyGroupIndex":"1,3,2",
		"sinkJob" :"next-job",

		"groupType":"H",
		"horizontalWidth":24
		"verticalStart":  1257010786,
		"verticalHeight": 4797,
		"groupSize":    30
	}
	`

	datasetMWA = &DataSet{
		DatasetID:       "mwa-dat:1257010784",
		KeyGroupRegex:   "^.+/([0-9]{10})_([0-9]{10})_ch([0-9]{3}).+$",
		KeyGroupIndex:   "1,3,2",
		SinkJob:         "next-job",
		GroupType:       "H",
		HorizontalWidth: 24,
		VerticalStart:   1257010786,
		VerticalHeight:  4797,
		GroupSize:       30,
	}

	jsonDataSetCraftsFits = ` {
		"datasetID":"fits:Dec+6007_09_03/20221019",
		"keyGroupRegex":"^([^/]+/[^/]+)/.+-M([0-9]+)_([0-9]+).fits.*$",
		"keyGroupIndex":"1,2,3",
		"sinkJob":"next-job",

		"groupType":"V",
		"horizontalWidth":19,
		"verticalStart":  1,
		"verticalHeight": 4,
		"groupSize":    3,
		"interleaved":true
	}
	`
)
