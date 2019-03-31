package file

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/shibukawa/configdir"
)

// CreateFullPath make parent directories as needed
func CreateFullPath(path string) error {
	return os.MkdirAll(path, 0750)
}

// CreateFullPathForFile take the directory path in which this file should exist and create this path
func CreateFullPathForFile(filePath string) (err error) {
	dir, _ := filepath.Split(filePath)
	return CreateFullPath(dir)
}

// IsADirectoryPath get info about the path string not about a physical resource in the filesystem
// return true if the path is a directory path
func IsADirectoryPath(path string) bool {
	dir, file := filepath.Split(path)
	return (len(dir) != 0 && len(file) == 0)
}

// ExistFile return true if this file exist
func ExistFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// HomeConfigFolder return the default config folder path
func HomeConfigFolder(systemVendorName, systemProgramName string) (*configdir.Config, error) {
	configDirs := configdir.New(systemVendorName, systemProgramName)
	folders := configDirs.QueryFolders(configdir.Global)
	if len(folders) <= 0 {
		return nil, errors.New("some strange error was happen, can not determine the home config folder")
	}
	return folders[0], nil
}

// DefaultConfigPath return the path for file `filename` under the default config folder
func DefaultConfigPath(fileName string, homeConf *configdir.Config) (path string) {
	return filepath.Join(homeConf.Path, fileName)
}
