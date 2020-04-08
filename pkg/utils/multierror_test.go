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

package utils_test

import (
	"errors"
	"testing"

	"github.com/appvia/kore/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAppendMultiErrorWithNilErrorsReturnsNil(t *testing.T) {
	var err1, err2 error
	me := utils.AppendMultiError(err1, err2)
	assert.NoError(t, me)
}

func TestAppendMultiErrorReturnsFirstIfSecondIsNil(t *testing.T) {
	var err1, err2 error
	err1 = errors.New("first error")
	me := utils.AppendMultiError(err1, err2)
	assert.Error(t, me)
	assert.Equal(t, err1, me)
}

func TestAppendMultiErrorReturnsSecondIfFirstIsNil(t *testing.T) {
	var err1, err2 error
	err2 = errors.New("second error")
	me := utils.AppendMultiError(err1, err2)
	assert.Error(t, me)
	assert.Equal(t, err2, me)
}

func TestAppendMultiErrorJoinsErrors(t *testing.T) {
	err1 := errors.New("first error")
	err2 := errors.New("second error")
	me := utils.AppendMultiError(err1, err2)
	assert.EqualError(t, me, "first error, second error")
}

func TestAppendMultiErrorFirstIsAMultiError(t *testing.T) {
	err1 := errors.New("first error")
	err2 := utils.AppendMultiError(err1, errors.New("second error"))
	err3 := errors.New("third error")
	me := utils.AppendMultiError(err2, err3)
	assert.EqualError(t, me, "first error, second error, third error")
}

func TestAppendMultiErrorSecondIsAMultiError(t *testing.T) {
	err1 := errors.New("first error")
	err2 := errors.New("second error")
	err3 := utils.AppendMultiError(err2, errors.New("third error"))
	me := utils.AppendMultiError(err1, err3)
	assert.EqualError(t, me, "first error, second error, third error")
}

func TestAppendMultiErrorBothAreAMultiError(t *testing.T) {
	err1 := errors.New("first error")
	err2 := utils.AppendMultiError(err1, errors.New("second error"))
	err3 := errors.New("third error")
	err4 := utils.AppendMultiError(err3, errors.New("forth error"))
	me := utils.AppendMultiError(err2, err4)
	assert.EqualError(t, me, "first error, second error, third error, forth error")
}
