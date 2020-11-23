package config

import (
	"path/filepath"

	"github.com/triasteam/go-streamnet/types"
)

var (
	DefaultStreamNetGoDir = ".sng"
	defaultConfigDir      = "config"
	defaultDataDir        = "data"

	defaultConfigFileName  = "config.toml"
	defaultGenesisJSONName = "genesis.json"
	defaultAddrBookName    = "addrbook.json"

	defaultConfigFilePath  = filepath.Join(defaultConfigDir, defaultConfigFileName)
	defaultGenesisJSONPath = filepath.Join(defaultConfigDir, defaultGenesisJSONName)
	defaultAddrBookPath    = filepath.Join(defaultConfigDir, defaultAddrBookName)

	// GenesisTrunk is the trunk of genesis
	GenesisTrunk = types.Sha256([]byte("StreamNet_Trunk"))

	// GenesisBranch is the branck of genesis
	GenesisBranch = types.Sha256([]byte("StreamNet_Branch"))
)

// Config defines the top level configuration for a StreamNet node
type TODO_Config struct {
	RPC       *RPCConfig
	P2P       *P2PConfig
	Consensus *ConsensusConfig
}

// DefaultConfig returns a default configuration for a StreamNet node
func DefaultConfig() *Config {
	return &Config{}
}

type RPCConfig struct {
	// TCP or UNIX socket address for the RPC server to listen on
	ListenAddress string
}

type P2PConfig struct {
	// Address to listen for incoming connections
	ListenAddress string

	// Path to address book
	AddrBook string
}

type ConsensusConfig struct {
	mwm int
}
