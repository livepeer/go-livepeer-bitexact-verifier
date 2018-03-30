package verifier

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ffmpeg "github.com/livepeer/lpms/ffmpeg"
	"github.com/livepeer/lpms/transcoder"
)

type Verifier struct {
	WorkDir string
}

func NewVerifier(wd string) *Verifier {
	return &Verifier{WorkDir: wd}
}

func (v *Verifier) ComputeResultHash(data []byte, transcodingOptions string) (string, error) {
	profiles, err := parseTranscodingOptions(transcodingOptions)
	if err != nil {
		return "", err
	}

	ffmpeg.InitFFmpeg()
	defer ffmpeg.DeinitFFmpeg()

	tr := transcoder.NewFFMpegSegmentTranscoder(profiles, "", v.WorkDir)

	fname, err := v.createDataFile(data)
	if err != nil {
		return "", err
	}

	defer os.Remove(fname)

	tData, err := tr.Transcode(fname)
	if err != nil {
		return "", fmt.Errorf("Failed to transcode seg %v: %v", fname, err)
	}

	tDataHashes := make([][]byte, len(tData))
	for _, td := range tData {
		tDataHashes = append(tDataHashes, crypto.Keccak256(td))
	}

	dataHash := crypto.Keccak256(data)
	concatTDataHash := crypto.Keccak256(tDataHashes...)

	return common.ToHex(crypto.Keccak256(dataHash, concatTDataHash)), nil
}

func (v *Verifier) createDataFile(data []byte) (string, error) {
	inName := randName()
	if _, err := os.Stat(v.WorkDir); os.IsNotExist(err) {
		err := os.Mkdir(v.WorkDir, 0700)
		if err != nil {
			return "", fmt.Errorf("Verifier could not create work dir: %v", err)
		}
	}

	fname := path.Join(v.WorkDir, inName)
	if err := ioutil.WriteFile(fname, data, 0644); err != nil {
		return "", fmt.Errorf("Verifier could not write file: %v", err)
	}

	return fname, nil
}

func parseTranscodingOptions(transcodingOptions string) ([]ffmpeg.VideoProfile, error) {
	profiles := make([]ffmpeg.VideoProfile, 0)

	for i := 0; i < len(transcodingOptions); i += VideoProfileIDSize {
		opt := transcodingOptions[i : i+VideoProfileIDSize]

		pName, ok := VideoProfileNameLookup[opt]
		if !ok {
			return nil, fmt.Errorf("Missing video profile name for id: %v", opt)
		}

		p, ok := ffmpeg.VideoProfileLookup[pName]
		if !ok {
			return nil, fmt.Errorf("Missing video profile for name: %v", pName)
		}

		profiles = append(profiles, p)
	}

	return profiles, nil
}

func randName() string {
	rand.Seed(time.Now().UnixNano())
	x := make([]byte, 10, 10)
	for i := 0; i < len(x); i++ {
		x[i] = byte(rand.Uint32())
	}

	return fmt.Sprintf("%x.ts", x)
}
