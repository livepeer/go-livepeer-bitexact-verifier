package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
	"os/user"

	"github.com/livepeer/go-livepeer-bitexact-verifier/ipfs"
	"github.com/livepeer/go-livepeer-bitexact-verifier/verifier"
)

func verifierMain(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Insufficient number of arguments. 2 required and %d provided", len(args))
	}

	hash := args[0]
	transcodingOptions := args[1]

	glog.Infof("IPFS hash of segment data: %v", hash)
	glog.Infof("Transcoding options for segment data: %v", transcodingOptions)

	usr, err := user.Current()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup IPFS
	glog.Infof("Starting IPFS node...")

	ipfsPath := fmt.Sprintf("%v/.ipfs", usr.HomeDir)
	ipfsApi, err := ipfs.StartIpfs(ctx, ipfsPath)
	if err != nil {
		return err
	}

	glog.Infof("Fetching segment data from IPFS using hash %v...", hash)

	data, err := ipfsApi.Cat(hash)
	if err != nil {
		return err
	}

	glog.Infof("Retrieved segment data from IPFS for hash %v", hash)

	// Setup verifier
	workDir := fmt.Sprintf("%v/.lpData", usr.HomeDir)
	vf := verifier.NewVerifier(workDir)

	glog.Infof("Computing result hash...")

	resHash, err := vf.ComputeResultHash(data, transcodingOptions)
	if err != nil {
		return err
	}

	// Note: Oraclize needs the hex encoded hash to NOT be 0x prefixed in order to unhexlify
	fmt.Printf("%v\n", resHash[2:])

	return nil
}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	if err := verifierMain(os.Args[1:]); err != nil {
		glog.Error(err)

		os.Exit(1)
	}
}
