package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type DeltaSource struct {
	Stat           Stat  `json:"stat"`
	BytesChanged   int64 `json:"bytes_changed"`
	ContentChanged int64 `json:"content_changed"`
}

type Stat struct {
	Mode      string `json:"mode"`
	Mtime     string `json:"mtime"`
	Size      int64  `json:"size"`
	Atime     string `json:"atime"`
	Ctime     string `json:"ctime"`
	ACL       ACL    `json:"acl"`
	InodeNum  int64  `json:"inode_num"`
	DeviceNum int64  `json:"device_num"`
}

type ACL struct {
	Mode     string `json:"mode"`
	Uid      int64  `json:"uid"`
	Gid      int64  `json:"gid"`
	CIFSAttr int64  `json:"cifs_attr"`
}

type DeltaRead struct {
	Path        string
	AbsoluteMB  float64
	ReducedMB   float64
	IncreasedMB float64
	TotalSizeMB float64
}

var (
	re         = regexp.MustCompile(`(?m)Path.*\nInode.*.?`)
	_          = kingpin.Version("0.0.1b")
	fp         = kingpin.Flag("inputfile", "file to read for parsing").Required().String()
	searchPath = kingpin.Flag("searchpath", "folder path to look for in the input file").Required().String()
)

const (
	MB = 1048576
)

func main() {
	kingpin.Parse()
	var doneDat DeltaRead
	doneDat.Path = *searchPath
	fd, err := readAndPrepFile(*fp)
	if err != nil {
		panic(err)
	}
	for i := range fd {

		brokenOut := strings.Split(fd[i], "\n")
		parsedPath := strings.Split(brokenOut[0], " ")[1]
		if strings.HasPrefix(parsedPath, *searchPath) {
			fmt.Println("section", i, parsedPath)
			var j DeltaSource
			err = json.Unmarshal([]byte(strings.TrimLeft(brokenOut[1], "Inode: ")), &j)
			if err != nil {
				panic(err)
			}

			if j.BytesChanged != 0 {
				switch {
				case j.BytesChanged > 0:
					jMB := (float64(j.BytesChanged) / MB)
					doneDat.IncreasedMB = doneDat.IncreasedMB + jMB
				case j.BytesChanged < 0:
					jMB := ((float64(j.BytesChanged) * -1) / MB)
					doneDat.ReducedMB = doneDat.ReducedMB + jMB
				}
			}
		}
	}
	// format the output
	increase := fmt.Sprintf("%0.2f", doneDat.IncreasedMB)
	reduce := fmt.Sprintf("%0.2f", doneDat.ReducedMB)
	doneDat.IncreasedMB, err = strconv.ParseFloat(increase, 64)
	doneDat.ReducedMB, err = strconv.ParseFloat(reduce, 64)
	doneDat.AbsoluteMB = doneDat.IncreasedMB + doneDat.ReducedMB
	doneDat.TotalSizeMB = doneDat.IncreasedMB - doneDat.ReducedMB
	pretty, err := json.MarshalIndent(doneDat, "", "   ")
	fmt.Println(string(pretty))

}

//readAndPrepFile removes any split up json payloads
func readAndPrepFile(path string) (fileData []string, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	preprocess := re.FindAll([]byte(strings.ReplaceAll(string(data), ",\n", ", ")), -1)
	for i := range preprocess {
		fileData = append(fileData, string(preprocess[i]))
	}

	return
}
