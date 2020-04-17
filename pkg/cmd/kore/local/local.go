/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package local

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/appvia/kore/pkg/client/config"
	"github.com/appvia/kore/pkg/cmd/kore/assets"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/version"

	"github.com/spf13/cobra"
)

// CreateLocalOptions is used to provision a team
type CreateLocalOptions struct {
	cmdutils.Factory
	cmdutils.DefaultHandler
}

// NewCmdCreateLocal returns the create local command
func NewCmdCreateLocal(factory cmdutils.Factory) *cobra.Command {
	//o := &CreateLocalOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "local",
		Short:   "Used to configure and run a local instance of Kore",
		Example: "kore local configure",
		Run:     cmdutils.RunHelp,
	}

	command.AddCommand(
		NewCmdLocalConfigure(factory),
		NewCmdLocalStart(factory),
		NewCmdLocalStop(factory),
		NewCmdLocalLogs(factory),
	)

	return command
}

// Shared functions used by stop and start:

func localPath() string {
	return filepath.Dir(config.GetClientConfigurationPath()) + "/local"
}

// writeSupportFiles populates ~/.korectl/local/ with the supporting files to be mounted into containers used by kore local.
func writeSupportFiles() error {
	localPath := localPath()
	for k, v := range assets.LocalSupport {
		f := filepath.Join(localPath, k)
		_ = os.MkdirAll(filepath.Dir(f), os.ModePerm)
		err := ioutil.WriteFile(f, []byte(v), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func getComposeCmd(conf *config.Config, args ...string) (*exec.Cmd, error) {
	baseArgs := []string{"-f", "-", "-p", "korelocal"}
	cmd := exec.Command("docker-compose",
		append(baseArgs, args...)...,
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("KORE_TAG=%s", version.Release))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KORE_IDP_CLIENT_ID=%s", conf.AuthInfos[LocalProfileName].OIDC.ClientID))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KORE_IDP_CLIENT_SECRET=%s", conf.AuthInfos[LocalProfileName].OIDC.ClientSecret))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KORE_IDP_SERVER_URL=%s", conf.AuthInfos[LocalProfileName].OIDC.AuthorizeURL))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KORE_LOCAL_HOME=%s", localPath()))
	if err := pipeComposeToCmd(cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func pipeComposeToCmd(cmd *exec.Cmd) error {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()
	_, err = io.WriteString(stdin, assets.LocalCompose)
	if err != nil {
		return err
	}
	return nil
}

func startChecks(conf *config.Config) error {
	if !conf.HasProfile(LocalProfileName) {
		return errors.New("a 'local' profile has not been found in ~/.korectl/config - try running: kore local configure")
	}

	if !conf.HasAuthInfo(LocalProfileName) || !conf.IsOIDCProviderConfigured(LocalProfileName) {
		return errors.New("no OpenId provider was configured for your 'local' profile in ~/.korectl/config - try running: kore local configure")
	}
	return nil
}

func isKoreStarted(conf *config.Config) (bool, error) {
	cmd, err := getComposeCmd(conf, "ps")
	if err != nil {
		return false, err
	}
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdoutStderr))
		return false, err
	}

	return strings.Contains(string(stdoutStderr), "/kore-apiserver"), nil
}
