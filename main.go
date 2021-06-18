package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

const (
	aurHost = "https://aur.archlinux.org"
)

type pkgDownloader struct {
	pkgNames []string
}

func newPkgDownloader() (pkgDownloader, error) {
	pkgNames, err := getForeignPkgNames()
	if err != nil {
		return pkgDownloader{}, err
	}

	return pkgDownloader{pkgNames: pkgNames}, nil
}

func main() {
	pkgDownloader, err := newPkgDownloader()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	for _, pkgName := range pkgDownloader.pkgNames {
		localVer, err := getLocalVer(pkgName)
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			if err := fetchPkg(pkgName); err != nil {
				log.Fatalf("%+v\n", err)
			}
			continue
		} else if err != nil {
			log.Fatalf("%+v\n", err)
		}
		remoteVer, err := fetchRemoteVer(pkgName)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
		if localVer != remoteVer {
			fmt.Printf("Download %s\n", pkgName)
			if err := fetchPkg(pkgName); err != nil {
				log.Fatalf("%+v\n", err)
			}
		}
	}
}
