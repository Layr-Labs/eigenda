package config

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigenda/inabox/utils"
)

type PkConfig struct {
	EcdsaMap map[string]KeyInfo `yaml:"ecdsaMap"`
	BlsMap   map[string]KeyInfo `yaml:"blsMap"`
}

type KeyInfo struct {
	// The private key (e.g. ECDSA or BLS) in string.
	PrivateKey string `yaml:"privateKey"`
	// The password used to encrypt the private key.
	Password string `yaml:"password"`
	// The file path to the encrypted private key.
	KeyFile string `yaml:"keyFile"`
}

type ConfigLock struct {
	RootPath string
	Path     string
	TestName string

	Config Config
	Pks    PkConfig `yaml:"privateKeys"`
	Envs
}

type Envs struct {
	Churner    ChurnerVars
	Dispersers []DisperserVars
	Batcher    BatcherVars
	Encoder    []EncoderVars
	Operators  []OperatorVars
	Stakers    []Staker
	Retriever  RetrieverVars
}

func (config *ConfigLock) GetKey(name string) (key, address string) {
	if keyInfo, ok := config.Config.Pks.EcdsaMap[name]; ok {
		key = keyInfo.PrivateKey
		address = utils.GetAddress(key)
		return
	}
	key = config.Pks.EcdsaMap[name].PrivateKey
	address = utils.GetAddress(key)
	return
}

func (config *ConfigLock) GetKeyString(name string) string {
	key, _ := config.GetKey(name)
	keyInt, ok := new(big.Int).SetString(key, 0)
	if !ok {
		log.Panicf("Error: could not parse key %s", key)
	}
	return keyInt.String()
}

// Constructs a mapping between service names/deployer names (e.g., 'dis0', 'opr1') and private keys. Order of priority: Map, List, File
func LoadPrivateKeys(env *Config, rootPath string) PkConfig {
	// construct full list of names
	names := make([]string, 0)
	for _, d := range env.Deployers {
		names = append(names, d.Name)
	}
	addNames := func(prefix string, num int) {
		for i := 0; i < num; i++ {
			names = append(names, fmt.Sprintf("%v%v", prefix, i))
		}
	}
	addNames("dis", env.Services.Counts.NumDis)
	addNames("opr", env.Services.Counts.NumOpr)
	addNames("staker", env.Services.Counts.NumOpr)

	// Collect private keys from file
	keyPath := filepath.Join(rootPath, "inabox", "secrets")

	// Read ECDSA private keys
	fileData := utils.MustReadFile(filepath.Join(keyPath, "ecdsa_keys/private_key_hex.txt"))
	ecdsaPks := strings.Split(string(fileData), "\n")
	// Read ECDSA passwords
	fileData = utils.MustReadFile(filepath.Join(keyPath, "ecdsa_keys/password.txt"))
	ecdsaPwds := strings.Split(string(fileData), "\n")
	// Read BLS private keys
	fileData = utils.MustReadFile(filepath.Join(keyPath, "bls_keys/private_key_hex.txt"))
	blsPks := strings.Split(string(fileData), "\n")
	// Read BLS passwords
	fileData = utils.MustReadFile(filepath.Join(keyPath, "bls_keys/password.txt"))
	blsPwds := strings.Split(string(fileData), "\n")

	if len(ecdsaPks) != len(blsPks) || len(blsPks) != len(ecdsaPwds) || len(ecdsaPwds) != len(blsPwds) {
		log.Panic("the number of keys and passwords for ECDSA and BLS must be the same")
	}

	pks := PkConfig{}
	// Add missing items to map
	pks.EcdsaMap = make(map[string]KeyInfo)
	pks.BlsMap = make(map[string]KeyInfo)

	ind := 0
	for _, name := range names {
		_, exists := pks.EcdsaMap[name]
		if !exists {

			if ind >= len(ecdsaPks) {
				log.Panic("not enough pks")
			}

			pks.EcdsaMap[name] = KeyInfo{
				PrivateKey: ecdsaPks[ind],
				Password:   ecdsaPwds[ind],
				KeyFile:    fmt.Sprintf("%s/ecdsa_keys/keys/%v.ecdsa.key.json", keyPath, ind+1),
			}
			pks.BlsMap[name] = KeyInfo{
				PrivateKey: blsPks[ind],
				Password:   blsPwds[ind],
				KeyFile:    fmt.Sprintf("%s/bls_keys/keys/%v.bls.key.json", keyPath, ind+1),
			}

			ind++
		}
	}

	return pks
}

func OpenConfigLock(rootPath, testName string) (testEnv *ConfigLock) {
	data := utils.MustReadFile(filepath.Join(rootPath, "inabox/testdata", testName, "config.lock"))
	var lock ConfigLock
	utils.MustUnmarshalYaml[ConfigLock](data, &lock)
	return &lock
}

func OpenCwdConfigLock() *ConfigLock {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicf("Couldn't get cwd: %v", err)
	}
	data := utils.MustReadFile(filepath.Join(cwd, "config.lock"))
	var testEnv ConfigLock
	utils.MustUnmarshalYaml(data, &testEnv)
	return &testEnv
}
