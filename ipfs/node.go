package ipfs

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	lockfile "github.com/ipfs/go-ipfs/repo/fsrepo/lock"
)

type IpfsApi interface {
	Cat(string) ([]byte, error)
}

type IpfsCoreApi core.IpfsNode

const (
	nBitsForKeypairDefault = 2048
)

func StartIpfs(ctx context.Context, repoPath string) (*IpfsCoreApi, error) {
	if !fsrepo.IsInitialized(repoPath) {
		conf, err := config.Init(os.Stdout, nBitsForKeypairDefault)
		if err != nil {
			return nil, err
		}
		if err := fsrepo.Init(repoPath, conf); err != nil {
			return nil, err
		}
	}

	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	ncfg := &core.BuildCfg{
		Repo:      repo,
		Online:    true,
		Permament: true,
		Routing:   core.DHTOption,
	}

	node, err := core.NewNode(ctx, ncfg)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				closeIpfs(node, repoPath)
				return
			}
		}
	}()

	return (*IpfsCoreApi)(node), nil
}

func closeIpfs(node *core.IpfsNode, repoPath string) {
	repoLockFile := filepath.Join(repoPath, lockfile.LockFile)
	os.Remove(repoLockFile)
	node.Close()
}

func (ipfs *IpfsCoreApi) Cat(hash string) ([]byte, error) {
	node := ipfs.node()

	reader, err := coreunix.Cat(node.Context(), node, hash)
	if err != nil {
		return nil, err
	}

	res := make([]byte, reader.Size())
	_, err = reader.Read(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ipfs *IpfsCoreApi) node() *core.IpfsNode {
	return (*core.IpfsNode)(ipfs)
}
