package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type ConfType struct {
	ChangeInterval string `json:"changeInterval"`
}

var (
	// Default config when someone starts the application for the first time.
	DefaultConf = &ConfType{}
	// The config
	conf = &ConfType{}
)

// Get path to the config file, the path is always to a file that exists, if err != nil
func (c *ConfType) getConfFile() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := path.Join(dir, ".wallpaperconf.json")

	// veryify that config exists
	if _, err := os.Stat(filePath); err == nil {
		// file exists - do nothing
	} else if os.IsNotExist(err) {
		// does not exist - create it with defaults
		data, err := json.Marshal(DefaultConf)
		if err != nil {
			return "", err
		}
		err = os.WriteFile(filePath, data, 0)
		if err != nil {
			return "", err
		}
	} else {
		return "", err
	}

	return filePath, nil
}

// Force reload the config from file.
func (c *ConfType) Reload() error {
	return c.load()
}

// Load the config.
func (c *ConfType) load() error {
	filePath, err := c.getConfFile()
	if err != nil {
		return err
	}

	// The file exists
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, conf)
	return err

}

// Get the config
func Get() *ConfType {
	return conf
}

func init() {
	// Load the config on start
	err := conf.load()
	if err != nil {
		log.Printf("!! Could not load config file.")
		log.Print(err)
	}
}
