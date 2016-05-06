package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/monmohan/xferspdy/data"
	"os"
	"path/filepath"
)

var (
	fPath   = flag.String("file", "", "File path to create the fingerprint, mandatory")
	blockSz = flag.Uint64("blocksz", 2*1024, "Block Size, default block size is 2KB")
)

func main() {
	flag.Parse()
	if *fPath == "" {
		glog.Exitf("File path is required\n")
	}
	glog.V(2).Infof("File path %s , Block Size %d \n", *fPath, *blockSz)

	fgprt := data.NewFingerprint(*fPath, uint32(*blockSz))
	glog.V(4).Infof("Signature  %v \n", *fgprt)

	dir, fname := filepath.Split(*fPath)

	fname = filepath.Join(dir, fname+".fingerprint")

	fpfile, err := os.OpenFile(fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)

	if err != nil {
		glog.Fatalf("Error in creating finger print file %v \n, Error :", filepath.Join(dir, fname+".fingerprint"), err)
	}

	enc := gob.NewEncoder(fpfile)
	enc.Encode(*fgprt)
	fmt.Printf("Fingerprint for file: %v \nGenerated:  %v \n ", *fPath, fpfile.Name())
	fpfile.Close()

	fpfile, err = os.Open(fname)
	defer fpfile.Close()

	var fp data.Fingerprint
	dec := gob.NewDecoder(fpfile)
	err = dec.Decode(&fp)

	glog.V(4).Infof("Verifying signature , created %v\n decoded from file %v\n", *fgprt, fp)

	if err != nil || (len(fgprt.BlockMap) != len(fp.BlockMap)) {
		glog.Fatalf("Failed to decode finger print during verification %v\n", err)
	}

	glog.Flush()

}
