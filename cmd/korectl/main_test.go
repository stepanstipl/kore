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

package main

import (
	"bytes"
	"os"

	"github.com/urfave/cli/v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const TestAppName = "korectl-test"

var _ = Describe("Main", func() {
	var args []string
	var exitCode int
	var exitErr error
	var stdout *bytes.Buffer
	var stderr *bytes.Buffer

	BeforeSuite(func() {
		cli.OsExiter = func(code int) {
			exitCode = code
		}
	})

	AfterSuite(func() {
		cli.OsExiter = os.Exit
	})

	BeforeEach(func() {
		exitCode = 0
		exitErr = nil
		stdout = bytes.NewBuffer([]byte{})
		stderr = bytes.NewBuffer([]byte{})
	})

	JustBeforeEach(func() {
		exitCode, exitErr = Main(append([]string{TestAppName}, args...), stdout, stderr)
	})

	When("no arguments are passed", func() {
		It("should return with a non-zero exit code", func() {
			Expect(exitCode).ToNot(Equal(0))
		})

		It("should show the help", func() {
			Expect(stdout).To(ContainSubstring("USAGE:"))
		})
	})

	When("an unknown command is passed", func() {
		BeforeEach(func() {
			args = []string{"unknown-command"}
		})

		It("should return with a non-zero exit code", func() {
			Expect(exitCode).ToNot(Equal(0))
		})

		It("should return an error", func() {
			Expect(exitErr).To(HaveOccurred())
			Expect(exitErr.Error()).To(ContainSubstring("unknown command"))
		})
	})

})
