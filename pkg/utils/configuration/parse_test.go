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

package configuration_test

import (
	"context"
	"encoding/base64"
	"reflect"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/sirupsen/logrus"

	"github.com/appvia/kore/pkg/utils/configuration"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers/controllersfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type testObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              testObjectSpec `json:"spec,omitempty"`
}

func (t testObject) GetConfiguration() *apiextv1.JSON {
	return t.Spec.Configuration
}

func (t testObject) GetConfigurationFrom() []corev1.ConfigurationFromSource {
	return t.Spec.ConfigurationFrom
}

type testObjectSpec struct {
	Configuration     *apiextv1.JSON                   `json:"configuration,omitempty"`
	ConfigurationFrom []corev1.ConfigurationFromSource `json:"configurationFrom,omitempty"`
}

type testConfig struct {
	Param1     string                   `json:"param1"`
	Param2     string                   `json:"param2"`
	ParamMap   map[string]interface{}   `json:"param_map"`
	ParamArray []map[string]interface{} `json:"param_array"`
	ParamBool  bool                     `json:"param_bool"`
	ParamInt   int                      `json:"param_int"`
	ParamFloat float64                  `json:"param_float"`
}

var _ = Describe("ParseObjectConfiguration", func() {
	var client *controllersfakes.FakeClient
	var v *testConfig
	var o *testObject
	var parseErr error
	var secretData map[string]string

	BeforeEach(func() {
		client = &controllersfakes.FakeClient{}

		secretData = map[string]string{
			"secretkey1": "secretvalue1",
			"secretkey2": "secretvalue2",
		}

		v = &testConfig{}
		o = &testObject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testobject",
				Namespace: "testobjectns",
			},
			Spec: testObjectSpec{
				Configuration: &apiextv1.JSON{Raw: []byte(`{"param1":"value1"}`)},
				ConfigurationFrom: []corev1.ConfigurationFromSource{
					{
						Path: "param2",
						SecretKeyRef: &corev1.OptionalSecretKeySelector{
							SecretKeySelector: corev1.SecretKeySelector{
								Name:      "testsecret",
								Namespace: "testsecretns",
								Key:       "secretkey2",
							},
							Optional: false,
						},
					},
				},
			},
		}

	})

	JustBeforeEach(func() {
		client.GetStub = func(ctx context.Context, name types.NamespacedName, object runtime.Object) error {
			if name.Namespace == "testsecretns" && name.Name == "testsecret" {
				secret := &configv1.Secret{
					Spec: configv1.SecretSpec{
						Data: map[string]string{},
					},
				}
				for k, v := range secretData {
					secret.Spec.Data[k] = base64.StdEncoding.EncodeToString([]byte(v))
				}
				reflect.ValueOf(object).Elem().Set(reflect.ValueOf(secret).Elem())
				return nil
			}

			return errors.NewNotFound(schema.GroupResource{Resource: "Secret"}, name.Name)
		}

		logger := logrus.New()
		logger.Out = GinkgoWriter
		parseErr = configuration.ParseObjectConfiguration(context.Background(), client, o, v)
	})

	It("should parse the configuration from the secrets", func() {
		Expect(parseErr).ToNot(HaveOccurred())
		Expect(v).To(Equal(&testConfig{
			Param1: "value1",
			Param2: "secretvalue2",
		}))
	})

	When("the secret value is quoted string", func() {
		BeforeEach(func() {
			secretData["secretkey2"] = `"123"`
		})

		It("should set the string value", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1: "value1",
				Param2: "123",
			}))
		})
	})

	When("it loads multiple values from the same secret", func() {
		BeforeEach(func() {
			o.Spec.ConfigurationFrom = []corev1.ConfigurationFromSource{
				{
					Path: "param1",
					SecretKeyRef: &corev1.OptionalSecretKeySelector{
						SecretKeySelector: corev1.SecretKeySelector{
							Name:      "testsecret",
							Namespace: "testsecretns",
							Key:       "secretkey1",
						},
						Optional: false,
					},
				},
				{
					Path: "param2",
					SecretKeyRef: &corev1.OptionalSecretKeySelector{
						SecretKeySelector: corev1.SecretKeySelector{
							Name:      "testsecret",
							Namespace: "testsecretns",
							Key:       "secretkey2",
						},
						Optional: false,
					},
				},
			}
		})

		It("should cache the already loaded secrets", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1: "secretvalue1",
				Param2: "secretvalue2",
			}))
			Expect(client.GetCallCount()).To(Equal(1))
		})
	})

	When("name defines a parameter on a map parameter", func() {
		BeforeEach(func() {
			o.Spec.ConfigurationFrom[0].Path = "param_map.value"
		})

		It("should set the value on the map", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1: "value1",
				ParamMap: map[string]interface{}{
					"value": "secretvalue2",
				},
			}))
		})
	})

	When("name defines a parameter on an array parameter", func() {
		BeforeEach(func() {
			o.Spec.ConfigurationFrom[0].Path = "param_array.0.value"
		})

		It("should set the value on the map", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1: "value1",
				ParamArray: []map[string]interface{}{
					{
						"value": "secretvalue2",
					},
				},
			}))
		})
	})

	When("the secret value is bool", func() {
		BeforeEach(func() {
			secretData["param_bool"] = `true`
			o.Spec.ConfigurationFrom = []corev1.ConfigurationFromSource{
				{
					Path: "param_bool",
					SecretKeyRef: &corev1.OptionalSecretKeySelector{
						SecretKeySelector: corev1.SecretKeySelector{
							Name:      "testsecret",
							Namespace: "testsecretns",
							Key:       "param_bool",
						},
					},
				},
			}
		})

		It("should set the bool value", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1:    "value1",
				ParamBool: true,
			}))
		})
	})

	When("the secret value is integer", func() {
		BeforeEach(func() {
			secretData["param_int"] = `123`
			o.Spec.ConfigurationFrom = []corev1.ConfigurationFromSource{
				{
					Path: "param_int",
					SecretKeyRef: &corev1.OptionalSecretKeySelector{
						SecretKeySelector: corev1.SecretKeySelector{
							Name:      "testsecret",
							Namespace: "testsecretns",
							Key:       "param_int",
						},
					},
				},
			}
		})

		It("should set the integer value", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1:   "value1",
				ParamInt: 123,
			}))
		})
	})

	When("the secret value is float", func() {
		BeforeEach(func() {
			secretData["param_float"] = `123.456`
			o.Spec.ConfigurationFrom = []corev1.ConfigurationFromSource{
				{
					Path: "param_float",
					SecretKeyRef: &corev1.OptionalSecretKeySelector{
						SecretKeySelector: corev1.SecretKeySelector{
							Name:      "testsecret",
							Namespace: "testsecretns",
							Key:       "param_float",
						},
					},
				},
			}
		})

		It("should set the float value", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1:     "value1",
				ParamFloat: 123.456,
			}))
		})
	})

	When("the secret value is a JSON object", func() {
		BeforeEach(func() {
			secretData["param_map"] = `{"foo": "bar"}`
			o.Spec.ConfigurationFrom = []corev1.ConfigurationFromSource{
				{
					Path: "param_map",
					SecretKeyRef: &corev1.OptionalSecretKeySelector{
						SecretKeySelector: corev1.SecretKeySelector{
							Name:      "testsecret",
							Namespace: "testsecretns",
							Key:       "param_map",
						},
					},
				},
			}
		})

		It("should set the map value", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1: "value1",
				ParamMap: map[string]interface{}{
					"foo": "bar",
				},
			}))
		})
	})

	When("the secret value is an array", func() {
		BeforeEach(func() {
			secretData["param_array"] = `[{"foo": "bar"}]`
			o.Spec.ConfigurationFrom = []corev1.ConfigurationFromSource{
				{
					Path: "param_array",
					SecretKeyRef: &corev1.OptionalSecretKeySelector{
						SecretKeySelector: corev1.SecretKeySelector{
							Name:      "testsecret",
							Namespace: "testsecretns",
							Key:       "param_array",
						},
					},
				},
			}
		})

		It("should set the array value", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(v).To(Equal(&testConfig{
				Param1: "value1",
				ParamArray: []map[string]interface{}{
					{
						"foo": "bar",
					},
				},
			}))
		})
	})

	When("name defines a parameter on a non-map, non-array parameter", func() {
		BeforeEach(func() {
			o.Spec.ConfigurationFrom[0].Path = "param1.value"
		})

		It("should throw an error", func() {
			Expect(parseErr).To(HaveOccurred())
		})
	})

	When("the secret does not exist", func() {
		BeforeEach(func() {
			o.Spec.ConfigurationFrom[0].SecretKeyRef.Name = "nonexisting"
		})

		It("should throw an error", func() {
			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr.Error()).To(ContainSubstring("does not exist"))
		})

		When("the value is optional", func() {
			BeforeEach(func() {
				o.Spec.ConfigurationFrom[0].SecretKeyRef.Optional = true
			})

			It("should not error", func() {
				Expect(parseErr).ToNot(HaveOccurred())
				Expect(v).To(Equal(&testConfig{
					Param1: "value1",
				}))
			})
		})
	})

	When("the secret key does not exist", func() {
		BeforeEach(func() {
			o.Spec.ConfigurationFrom[0].SecretKeyRef.Key = "nonexisting"
		})

		It("should throw an error", func() {
			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr.Error()).To(ContainSubstring("does not exist"))
		})

		When("the value is optional", func() {
			BeforeEach(func() {
				o.Spec.ConfigurationFrom[0].SecretKeyRef.Optional = true
			})

			It("should not error", func() {
				Expect(parseErr).ToNot(HaveOccurred())
				Expect(v).To(Equal(&testConfig{
					Param1: "value1",
				}))
			})
		})
	})

	When("the namespace is not set", func() {
		BeforeEach(func() {
			o.Spec.ConfigurationFrom[0].SecretKeyRef.Namespace = ""
		})

		It("should default to the objects namespace", func() {
			Expect(parseErr).To(MatchError("failed to load secret testobjectns/testsecret: does not exist"))
		})
	})
})
