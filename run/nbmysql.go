package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wangjun861205/nbmysql"
)

func main() {
	var path string
	var createTab bool
	var rootPath string
	flag.StringVar(&path, "p", "./", "specific parse root path")
	flag.BoolVar(&createTab, "c", false, "if create all table")
	flag.Parse()
	if filepath.IsAbs(path) {
		rootPath = path
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		rootPath = filepath.Join(pwd, path)
	}
	infos, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range infos {
		if !info.IsDir() && info.Name()[len(info.Name())-5:] == ".nbdb" {
			db, err := nbmysql.ParseDatabase(filepath.Join(rootPath, info.Name()))
			if err != nil {
				log.Fatal(err)
			}
			err = nbmysql.Gen(db, filepath.Join(rootPath, strings.TrimRight(info.Name(), ".nbdb")+".model.go"))
			if err != nil {
				log.Fatal(err)
			}
			if createTab {
				for _, tab := range db.Tables {
					err := db.CreateTableIfNotExists(tab)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}
