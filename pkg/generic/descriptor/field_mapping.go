/*
 * Copyright 2021 CloudWeGo Authors
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

package descriptor

import (
	"reflect"
)

// FiledMapping mapping handle for filed descriptor
type FiledMapping interface {
	Handle(field *FieldDescriptor)
}

// NewFieldMapping FiledMapping creator
type NewFieldMapping func(value string) FiledMapping

// GoTagAnnatition go.tag annatation define
var GoTagAnnatition = NewBAMAnnotation("go.tag", NewGoTag)

// go.tag = 'json:"xx"'
type goTag struct {
	tag reflect.StructTag
}

// NewGoTag go.tag annotation creator
var NewGoTag NewFieldMapping = func(value string) FiledMapping {
	return &goTag{reflect.StructTag(value)}
}

func (m *goTag) Handle(field *FieldDescriptor) {
	field.Alias = m.tag.Get("json")
}
