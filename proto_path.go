package grpcdump

import (
	"os"
	"path/filepath"
	"strings"
	"path"
	"log"
)

type ArrayFlags []string

func (arrayFlags *ArrayFlags) String() string {
	return "Specify the directory in which to search for imports"
}

func (arrayFlags *ArrayFlags) Set(value string) error {
	*arrayFlags = append(*arrayFlags, value)
	return nil
}

func (arrayFlags *ArrayFlags) ParseDir(suffix string) []string {
	ret := make([]string, 0)

	for _, file := range *arrayFlags {
		fileinfo, _ := os.Stat(file)

		if fileinfo.IsDir() {
			filepath.Walk(file, func(p string, info os.FileInfo, err error) error {
				if strings.HasSuffix(info.Name(), suffix) {
					f := path.Join(file, info.Name())
					log.Printf("proto files: %s", f)
					ret = append(ret, f)
				}

				return nil
			})
		} else {
			log.Printf("proto files: %s", file)
			ret = append(ret, file)
		}
	}

	return ret
}