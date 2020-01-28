/**
 * Copyright (C) 2020 Rohith Jayawardene <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package hubctl

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/appvia/kore/pkg/utils"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// GetCaches is responsible for checking if are caches are up to date
func GetCaches(config Config) error {
	_, err := GetSwaggerCache(config)
	if err != nil {
		return err
	}

	return nil
}

//
func GetKubeConfig() (string, error) {
	path := func() string {
		p := os.ExpandEnv(os.Getenv("$KUBECONFIG"))
		if p != "" {
			return p
		}

		return os.ExpandEnv("${HOME}/.kube/config")
	}()

	found, err := utils.FileExists(path)
	if err != nil {
		return "", err
	}
	if !found {
		return "", errors.New("no kubeconfig found")
	}

	return path, nil
}

// GetSwaggerCache is responsible for updating the swagger cache
func GetSwaggerCache(config Config) ([]byte, error) {
	log.Debug("checking for cached resources file")

	current, err := config.GetSwagger()
	if err != nil {
		log.WithError(err).Error("failed read in cached swagger")

		return nil, err
	}
	// @step: we need to check if the swagger is up to date
	checksum, err := GetFileChecksum(config.GetSwaggerCachedFile())
	if err != nil {
		log.WithError(err).Warn("failed checking the swagger cache file")

		return current, nil
	}
	latest, err := config.GetSwaggerChecksum()
	if err != nil {
		log.WithError(err).Debug("failed retrieving latest swagger cache file")
	}
	if checksum == latest {
		return current, nil
	}

	// @step: else we need to download the latest
	update, err := config.GetSwaggerFromAPI()
	if err != nil {
		log.WithError(err).Debug("failed to retrieve swagger from apiserver")

		return current, nil
	}
	if err := config.UpdateSwaggerCache(update); err != nil {
		log.WithError(err).Error("failed to update the swagger cache")
	}

	return update, nil
}

// GetFileChecksum returns the checksum of a content of a file
func GetFileChecksum(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, _ = h.Write(content)

	return hex.EncodeToString(h.Sum(nil)), nil
}

// GetClientConfiguration is responsible for retrieving the client configuration
func GetClientConfiguration() (Config, error) {
	config := Config{}

	resolved := os.ExpandEnv(HubConfig)

	if err := os.MkdirAll(filepath.Dir(resolved), os.FileMode(0760)); err != nil {
		return config, err
	}

	// @step: read in the configuration
	content, err := ioutil.ReadFile(resolved)
	if err != nil {
		return config, err
	}
	if string(content) == "" {
		content = []byte("server: ''")
	}

	// @step: parse the configuration
	return config, yaml.NewDecoder(bytes.NewReader(content)).Decode(&config)
}
