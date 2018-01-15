package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
)

var (
	set bool
)

func init() {
	log.SetFlags(0)
	flag.BoolVar(&set, "set", false, "set access and modification time to exif DateTime")
}

func main() {
	flag.Parse()

	root := flag.Arg(0)
	if len(root) == 0 {
		log.Fatalf("usage: %s PATH\n", os.Args[0])
	}

	info, err := os.Stat(root)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}
	if !info.IsDir() {
		log.Fatalf("%s in not a directory\n", root)
	}

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		cd, err := getCreateDate(path)
		if err != nil {
			log.Println(errors.Wrapf(err, "ERROR %s", path))
			return nil
		}
		msg := fmt.Sprintf("tag DateTime is %s for %s", cd, path)
		if set {
			if err := os.Chtimes(path, cd, cd); err != nil {
				log.Println(errors.Wrapf(err, "ERROR %s: change access/modification time", path))
				return nil
			}
			msg = fmt.Sprintf("set atime/mtime to %s on %s", cd, path)
		}
		fmt.Println(msg)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func getCreateDate(path string) (time.Time, error) {
	var cd time.Time
	f, err := os.Open(path)
	if err != nil {
		return cd, errors.Wrap(err, "open file")
	}
	x, err := exif.Decode(f)
	if err != nil {
		return cd, errors.Wrap(err, "decode exif data")
	}
	d, err := x.DateTime()
	if err != nil {
		return cd, errors.Wrap(err, "read DateTime tag")
	}
	return d, nil
}
