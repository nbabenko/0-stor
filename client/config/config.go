package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"

	"github.com/zero-os/0-stor/client/fullreadwrite"
)

const (
	chunkerStr      = "chunker"
	compressStr     = "compress"
	distributionStr = "distribution"
	encryptStr      = "encrypt"
	hashStr         = "hash"
	replicationStr  = "replication"
)

var (
	validPipes = map[string]struct{}{
		chunkerStr:      struct{}{},
		compressStr:     struct{}{},
		encryptStr:      struct{}{},
		hashStr:         struct{}{},
		replicationStr:  struct{}{},
		distributionStr: struct{}{},
	}
)

// Config defines 0-stor client config
type Config struct {
	Organization string   `yaml:"organization" validate:"nonzero"`
	Namespace    string   `yaml:"namespace" validate:"nonzero"`
	IyoClientID  string   `yaml:"iyo_client_id" validate:"nonzero"`
	IyoSecret    string   `yaml:"iyo_secret" validate:"nonzero"`
	Shards       []string `yaml:"shards" validate:"nonzero"` // 0-stor shards
	MetaShards   []string `yaml:"meta_shards"`
	Pipes        []Pipe   `yaml:"pipes" validate:"nonzero"`
}

// NewFromReader creates Config from a reader
func NewFromReader(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var conf Config

	// unmarshall
	if err := yaml.Unmarshal(b, &conf); err != nil {
		return nil, err
	}

	// validate
	if err := validator.Validate(conf); err != nil {
		return nil, err
	}

	// do post processing to each pipe
	for i, pipe := range conf.Pipes {
		if err := pipe.validate(); err != nil {
			return nil, err
		}

		if err := pipe.setConfigType(); err != nil {
			return nil, err
		}
		conf.Pipes[i] = pipe
	}
	return &conf, nil
}

// Write writes the config to the given writer
func (conf *Config) Write(w io.Writer) error {
	b, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

func (conf *Config) CreatePipeWriter(finalWriter fullreadwrite.Writer) (fullreadwrite.Writer, error) {
	nextWriter := finalWriter

	for i := len(conf.Pipes) - 1; i >= 0; i-- {
		pipe := conf.Pipes[i]
		w, err := pipe.CreateWriter(nextWriter, conf.Shards, conf.Organization, conf.Namespace)
		if err != nil {
			return nil, err
		}
		nextWriter = w
	}
	return nextWriter, nil
}

func (conf *Config) CreateAllReaders() ([]fullreadwrite.Reader, error) {
	var readers []fullreadwrite.Reader
	for _, pipe := range conf.Pipes {
		ar, err := pipe.CreateReader(nil, conf.Shards, conf.Organization, conf.Namespace)
		if err != nil {
			return nil, err
		}
		readers = append([]fullreadwrite.Reader{ar}, readers...)
	}
	return readers, nil
}