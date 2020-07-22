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

package jsonschema_test

import (
	"encoding/json"

	"github.com/appvia/kore/pkg/utils/validation"

	"github.com/appvia/kore/pkg/utils/jsonschema"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const schema = `
{
	"$id": "https://appvia.io/kore/schemas/test.json",
	"type": "object",
	"properties": {
		"p1": {
			"type": "string"
		},
		"p2": {
			"type": "string",
			"immutable": true
		},
		"p3": {
			"type": "object",
			"properties": {
				"p31": {
					"type": "string"
				},
				"p32": {
					"type": "string",
					"immutable": true
				}
			}
		},
		"p4": {
			"type": "object",
			"immutable": true,
			"properties": {
				"p31": {
					"type": "string"
				},
				"p32": {
					"type": "string",
					"immutable": true
				}
			}
		},
		"p5": {
			"type": "array",
			"items": { "$ref": "#" }
		},
		"p6": {
			"type": "array",
			"items": {
				"type": "object",
				"immutable": true,
				"properties": {
					"p61": {
						"type": "string"
					}
				}
			}
		},
		"p7": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"p71": {
						"type": "string",
						"identifier": true
					},
					"p72": {
						"type": "string",
						"immutable": true
					}
				}
			}
		}
	}
}
`

type Object struct {
	P1 string   `json:"p1"`
	P2 string   `json:"p2"`
	P3 Child    `json:"p3"`
	P4 Child    `json:"p4"`
	P5 []Object `json:"p5"`
	P6 []P6     `json:"p6"`
	P7 []P7     `json:"p7"`
}

type Child struct {
	C1 string `json:"p31"`
	C2 string `json:"p32"`
}

type P6 struct {
	P61 string `json:"p61"`
}

type P7 struct {
	P71 string `json:"p71"`
	P72 string `json:"p72"`
}

func defaultObject() Object {
	return Object{
		P1: "p1v",
		P2: "p2v",
		P3: Child{C1: "p3c1v", C2: "p3c2v"},
		P4: Child{C1: "p4c1v", C2: "p4c2v"},
		P5: []Object{
			{P1: "p51p1v", P2: "p51p2v"},
			{P1: "p52p1v", P2: "p52p2v"},
		},
		P6: []P6{
			{P61: "p61v"},
			{P61: "p62v"},
		},
		P7: []P7{
			{P71: "id1", P72: "value1"},
			{P71: "id2", P72: "value2"},
		},
	}
}

func createJSON(o Object) []byte {
	res, _ := json.Marshal(o)
	return res
}

var _ = Describe("ValidateImmutableProperties", func() {

	It("should return no error if no change", func() {
		j1 := createJSON(defaultObject())
		j2 := createJSON(defaultObject())

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should return no error if a non-immutable field has changed", func() {
		j1 := createJSON(defaultObject())
		o2 := defaultObject()
		o2.P1 = "changed"
		j2 := createJSON(o2)

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should return an error if an immutable field has changed", func() {
		j1 := createJSON(defaultObject())
		o2 := defaultObject()
		o2.P2 = "changed"
		j2 := createJSON(o2)

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).To(HaveOccurred())
		Expect(err.(*validation.Error).FieldErrors).To(HaveLen(1))
		Expect(err.(*validation.Error).FieldErrors[0]).To(Equal(validation.FieldError{
			Field:   "prefix.p2",
			ErrCode: validation.ReadOnly,
			Message: "updating the field is not allowed",
		}))
	})

	It("should return no error if a non-immutable field has changed in an inner object", func() {
		j1 := createJSON(defaultObject())
		o2 := defaultObject()
		o2.P3.C1 = "changed"
		j2 := createJSON(o2)

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should return an error if an immutable field has changed in an inner object", func() {
		j1 := createJSON(defaultObject())
		o2 := defaultObject()
		o2.P3.C2 = "changed"
		j2 := createJSON(o2)

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).To(HaveOccurred())
		Expect(err.(*validation.Error).FieldErrors).To(HaveLen(1))
		Expect(err.(*validation.Error).FieldErrors[0]).To(Equal(validation.FieldError{
			Field:   "prefix.p3.p32",
			ErrCode: validation.ReadOnly,
			Message: "updating the field is not allowed",
		}))
	})

	It("should return an error if an immutable object has changed", func() {
		j1 := createJSON(defaultObject())
		o2 := defaultObject()
		o2.P4.C1 = "changed"
		j2 := createJSON(o2)

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).To(HaveOccurred())
		Expect(err.(*validation.Error).FieldErrors).To(HaveLen(1))
		Expect(err.(*validation.Error).FieldErrors[0]).To(Equal(validation.FieldError{
			Field:   "prefix.p4",
			ErrCode: validation.ReadOnly,
			Message: "updating the field is not allowed",
		}))
	})

	It("should return no error if a non-immutable field changed in an array element", func() {
		j1 := createJSON(defaultObject())
		o2 := defaultObject()
		o2.P5[0].P1 = "changed"
		j2 := createJSON(o2)

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should return an error if an immutable field changed in an array element", func() {
		j1 := createJSON(defaultObject())
		o2 := defaultObject()
		o2.P5[0].P2 = "changed"
		j2 := createJSON(o2)

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).To(HaveOccurred())
		Expect(err.(*validation.Error).FieldErrors).To(HaveLen(1))
		Expect(err.(*validation.Error).FieldErrors[0]).To(Equal(validation.FieldError{
			Field:   "prefix.p5.0.p2",
			ErrCode: validation.ReadOnly,
			Message: "updating the field is not allowed",
		}))
	})

	It("should return an error if the array order changed with immutable fields", func() {
		j1 := createJSON(defaultObject())
		o2 := defaultObject()
		o2.P5[0], o2.P5[1] = o2.P5[1], o2.P5[0]
		j2 := createJSON(o2)

		err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
		Expect(err).To(HaveOccurred())
		Expect(err.(*validation.Error).FieldErrors).To(HaveLen(2))
		Expect(err.(*validation.Error).FieldErrors[0]).To(Equal(validation.FieldError{
			Field:   "prefix.p5.0.p2",
			ErrCode: validation.ReadOnly,
			Message: "updating the field is not allowed",
		}))
		Expect(err.(*validation.Error).FieldErrors[1]).To(Equal(validation.FieldError{
			Field:   "prefix.p5.1.p2",
			ErrCode: validation.ReadOnly,
			Message: "updating the field is not allowed",
		}))
	})

	When("array objects are immutable", func() {
		It("should return an error if the array order changed", func() {
			j1 := createJSON(defaultObject())
			o2 := defaultObject()
			o2.P6[0], o2.P6[1] = o2.P6[1], o2.P6[0]
			j2 := createJSON(o2)

			err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
			Expect(err).To(HaveOccurred())
			Expect(err.(*validation.Error).FieldErrors).To(HaveLen(2))
			Expect(err.(*validation.Error).FieldErrors[0]).To(Equal(validation.FieldError{
				Field:   "prefix.p6.0",
				ErrCode: validation.ReadOnly,
				Message: "updating the field is not allowed",
			}))
			Expect(err.(*validation.Error).FieldErrors[1]).To(Equal(validation.FieldError{
				Field:   "prefix.p6.1",
				ErrCode: validation.ReadOnly,
				Message: "updating the field is not allowed",
			}))
		})

		It("should not return an error if a new object is added", func() {
			j1 := createJSON(defaultObject())
			o2 := defaultObject()
			o2.P6 = append(o2.P6, P6{P61: "p63v"})
			j2 := createJSON(o2)

			err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not return an error if an object is deleted", func() {
			j1 := createJSON(defaultObject())
			o2 := defaultObject()
			o2.P6 = o2.P6[0:1]
			j2 := createJSON(o2)

			err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return an error if any object changed", func() {
			j1 := createJSON(defaultObject())
			o2 := defaultObject()
			o2.P6[0].P61 = "changed"
			j2 := createJSON(o2)

			err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
			Expect(err).To(HaveOccurred())
			Expect(err.(*validation.Error).FieldErrors).To(HaveLen(1))
			Expect(err.(*validation.Error).FieldErrors[0]).To(Equal(validation.FieldError{
				Field:   "prefix.p6.0",
				ErrCode: validation.ReadOnly,
				Message: "updating the field is not allowed",
			}))
		})
	})

	When("an array of objects has the same identifier field", func() {
		It("should not return an error if the array order changed", func() {
			j1 := createJSON(defaultObject())
			o2 := defaultObject()
			o2.P7[0], o2.P7[1] = o2.P7[1], o2.P7[0]
			j2 := createJSON(o2)

			err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not return an error if a new object is added", func() {
			j1 := createJSON(defaultObject())
			o2 := defaultObject()
			o2.P7 = append(o2.P7, P7{P71: "id3", P72: "value3"})
			j2 := createJSON(o2)

			err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not return an error if an object is deleted", func() {
			j1 := createJSON(defaultObject())
			o2 := defaultObject()
			o2.P7 = o2.P7[0:1]
			j2 := createJSON(o2)

			err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return an error if an immutable field changed", func() {
			j1 := createJSON(defaultObject())
			o2 := defaultObject()
			o2.P7[1].P72 = "changed"
			j2 := createJSON(o2)

			err := jsonschema.ValidateImmutableProperties(schema, "test", "prefix", j1, j2)
			Expect(err).To(HaveOccurred())
			Expect(err.(*validation.Error).FieldErrors).To(HaveLen(1))
			Expect(err.(*validation.Error).FieldErrors[0]).To(Equal(validation.FieldError{
				Field:   "prefix.p7.1.p72",
				ErrCode: validation.ReadOnly,
				Message: "updating the field is not allowed",
			}))
		})
	})
})
