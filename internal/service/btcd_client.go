package service

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"io/ioutil"
	"path/filepath"
)

func getBtcdClient() (*rpcclient.Client, error) {
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		return nil, err
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18334",
		Endpoint:     "ws",
		User:         "user",
		Pass:         "12345",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, nil)
	return client, err
}
