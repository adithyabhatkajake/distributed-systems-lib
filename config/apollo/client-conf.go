package apollo

import (
	"github.com/adithyabhatkajake/libe2c/crypto"
	pb "github.com/golang/protobuf/proto"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
)

// ClientConfig is an aggregation of all configs for the client in one place
type ClientConfig struct {
	// All Configs
	Config *ClientDataConfig
	// PKI Algorithm
	alg crypto.PKIAlgo
	// My private key
	PvtKey crypto.PrivKey
	// Mapping between nodes and their Public keys
	NodeKeyMap map[uint64]crypto.PubKey
}

// MarshalBinary implements Serializable interface
func (cc ClientConfig) MarshalBinary() ([]byte, error) {
	return pb.Marshal(cc.Config)
}

// UnmarshalBinary implements Deserializable interface
func (cc *ClientConfig) UnmarshalBinary(inpBytes []byte) error {
	data := &ClientDataConfig{}
	err := pb.Unmarshal(inpBytes, data)
	if err != nil {
		return nil
	}
	cc.Config = data
	cc.init()
	return err
}

// MarshalJSON implements JSON Marshaller interface
// https://talks.golang.org/2015/json.slide#19
func (cc ClientConfig) MarshalJSON() ([]byte, error) {
	return jsonpb.Marshal(cc.Config)
}

// UnmarshalJSON implements Unmarshaller JSON interface
// https://talks.golang.org/2015/json.slide#19
func (cc *ClientConfig) UnmarshalJSON(bytes []byte) error {
	err := jsonpb.Unmarshal(bytes, cc.Config)
	if err != nil {
		return err
	}
	cc.init()
	return err
}

// init function initializes the structure assuming the config has been set
func (cc *ClientConfig) init() {
	alg, exists := crypto.GetAlgo(cc.Config.CryptoCon.KeyType)
	if exists == false {
		panic("Unknown key type")
	}
	cc.alg = alg
	cc.PvtKey = alg.PrivKeyFromBytes(cc.Config.CryptoCon.PvtKey)
	cc.NodeKeyMap = make(map[uint64]crypto.PubKey)
	for idx, pubkeyBytes := range cc.Config.CryptoCon.NodeKeyMap {
		cc.NodeKeyMap[idx] = alg.PubKeyFromBytes(pubkeyBytes)
	}
}

// NewClientConfig creates a ClientConfig object from ClientDataConfig
func NewClientConfig(con *ClientDataConfig) *ClientConfig {
	cc := &ClientConfig{}
	cc.Config = con
	cc.init()
	return cc
}
