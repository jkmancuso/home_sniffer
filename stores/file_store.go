package stores

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type fileConfig struct {
	filename string
}

type fileStore struct {
	handle *os.File
	cfg    fileConfig
}

func NewFileStore() (fileStore, error) {
	fileCfg := newFileCfg()

	fStore := fileStore{
		cfg: fileCfg,
	}

	f, err := os.Create(fileCfg.filename)

	if err != nil {
		log.Errorf("Err: %v\ncould not output to file with params: %+v", err, fileCfg)
		return fStore, err
	}

	fStore.setHandle(f)

	log.Debug("Returning file store")

	return fStore, nil
}

func (store fileStore) Teardown() {
	log.Printf("Tearing down file store")
}

func (store *fileStore) setHandle(handle *os.File) {
	store.handle = handle
}

func newFileCfg() fileConfig {

	filepath := os.Getenv("FILEPATH")

	fileCfg := fileConfig{
		filename: filepath,
	}

	log.Debugf("Loading file cfg: %+v\n", fileCfg)

	return fileCfg

}

func (store fileStore) Send(data []string) error {

	var err error

	for _, payload := range data {
		_, err = store.handle.Write([]byte(payload + "\n"))

	}

	return err
}
