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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/appvia/kore/pkg/utils"
	yml "github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Document struct {
	// Endpoint is the rest storage endpoint
	Endpoint string
	// Object the resource to send
	Object *unstructured.Unstructured
}

// ParseDocument returns a collection of parsed documents and the api endpoints
func ParseDocument(src io.Reader, namespace string) ([]*Document, error) {
	var list []*Document

	// @step: read in the content of the file
	content, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	// @step: split the yaml documents up
	documents := strings.Split(string(content), "---")
	global := []string{"teams", "users", "plans"}

	for _, x := range documents {
		if x == "" {
			continue
		}

		doc, err := yml.YAMLToJSON([]byte(x))
		if err != nil {
			return nil, err
		}
		// @step: attempt to read the document into an unstructured
		u := &unstructured.Unstructured{}
		if err := u.UnmarshalJSON(doc); err != nil {
			return nil, err
		}

		// @checks
		// - ensure we have a name
		// - ensure we have a api kind
		if u.GetName() == "" {
			return nil, errors.New("resource must have names")
		}
		if u.GetKind() == "" {
			return nil, errors.New("resource must have an api kind")
		}
		if u.GetAPIVersion() == "" {
			return nil, errors.New("resource requires an api group")
		}

		// @step: we pluralize the kind and use that route the resource
		kind := strings.ToLower(utils.ToPlural(u.GetKind()))
		isGlobal := utils.Contains(kind, global)

		team := u.GetNamespace()
		if !isGlobal {
			if namespace != "" {
				if team != "" && team != namespace {
					return nil, errors.New("resource name and team selected are different")
				}
				team = namespace
			}
			if team == "" {
				return nil, errors.New("all resource must have a team namespace")
			}
		}
		team = strings.ToLower(team)
		name := strings.ToLower(u.GetName())

		remapping := map[string]string{
			"kubernetes": "clusters",
		}
		for k, v := range remapping {
			if k == kind {
				kind = v
			}
		}

		item := &Document{Object: u}
		switch isGlobal {
		case true:
			item.Endpoint = fmt.Sprintf("%s/%s", kind, name)
		default:
			item.Endpoint = fmt.Sprintf("%s/%s/%s", team, kind, name)
		}

		list = append(list, item)
	}

	return list, nil
}

// GetCaches is responsible for checking if are caches are up to date
func GetCaches(config *Config) error {
	/*
		content, err := GetSwaggerCache(*config)
		if err != nil {
			return err
		}
		_, err = fastjson.Parse(string(content))
		if err != nil {
			return err
		}
	*/

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

// BuildResourcesFromSwagger builds a list of global and namespaces resources
// from the swagger

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
func GetClientConfiguration() (*Config, error) {
	config := &Config{}

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
	return config, yaml.NewDecoder(bytes.NewReader(content)).Decode(config)
}
