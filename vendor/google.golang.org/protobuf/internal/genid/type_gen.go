// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by generate-protos. DO NOT EDIT.

package genid

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

const File_google_protobuf_type_proto = "google/protobuf/type.proto"

// Full and short names for google.protobuf.Syntax.
const (
	Syntax_enum_fullname = "google.protobuf.Syntax"
	Syntax_enum_name     = "Syntax"
)

// Enum values for google.protobuf.Syntax.
const (
	Syntax_SYNTAX_PROTO2_enum_value   = 0
	Syntax_SYNTAX_PROTO3_enum_value   = 1
	Syntax_SYNTAX_EDITIONS_enum_value = 2
)

// Names for google.protobuf.Type.
const (
	Type_message_name     protoreflect.Name     = "Type"
	Type_message_fullname protoreflect.FullName = "google.protobuf.Type"
)

// Field names for google.protobuf.Type.
const (
	Type_Name_field_name          protoreflect.Name = "name"
	Type_Fields_field_name        protoreflect.Name = "fields"
	Type_Oneofs_field_name        protoreflect.Name = "oneofs"
	Type_Options_field_name       protoreflect.Name = "options"
	Type_SourceContext_field_name protoreflect.Name = "source_context"
	Type_Syntax_field_name        protoreflect.Name = "syntax"
	Type_Edition_field_name       protoreflect.Name = "edition"

	Type_Name_field_fullname          protoreflect.FullName = "google.protobuf.Type.name"
	Type_Fields_field_fullname        protoreflect.FullName = "google.protobuf.Type.fields"
	Type_Oneofs_field_fullname        protoreflect.FullName = "google.protobuf.Type.oneofs"
	Type_Options_field_fullname       protoreflect.FullName = "google.protobuf.Type.options"
	Type_SourceContext_field_fullname protoreflect.FullName = "google.protobuf.Type.source_context"
	Type_Syntax_field_fullname        protoreflect.FullName = "google.protobuf.Type.syntax"
	Type_Edition_field_fullname       protoreflect.FullName = "google.protobuf.Type.edition"
)

// Field numbers for google.protobuf.Type.
const (
	Type_Name_field_number          protoreflect.FieldNumber = 1
	Type_Fields_field_number        protoreflect.FieldNumber = 2
	Type_Oneofs_field_number        protoreflect.FieldNumber = 3
	Type_Options_field_number       protoreflect.FieldNumber = 4
	Type_SourceContext_field_number protoreflect.FieldNumber = 5
	Type_Syntax_field_number        protoreflect.FieldNumber = 6
	Type_Edition_field_number       protoreflect.FieldNumber = 7
)

// Names for google.protobuf.Field.
const (
	Field_message_name     protoreflect.Name     = "Field"
	Field_message_fullname protoreflect.FullName = "google.protobuf.Field"
)

// Field names for google.protobuf.Field.
const (
	Field_Kind_field_name         protoreflect.Name = "kind"
	Field_Cardinality_field_name  protoreflect.Name = "cardinality"
	Field_Number_field_name       protoreflect.Name = "number"
	Field_Name_field_name         protoreflect.Name = "name"
	Field_TypeUrl_field_name      protoreflect.Name = "type_url"
	Field_OneofIndex_field_name   protoreflect.Name = "oneof_index"
	Field_Packed_field_name       protoreflect.Name = "packed"
	Field_Options_field_name      protoreflect.Name = "options"
	Field_JsonName_field_name     protoreflect.Name = "json_name"
	Field_DefaultValue_field_name protoreflect.Name = "default_value"

	Field_Kind_field_fullname         protoreflect.FullName = "google.protobuf.Field.kind"
	Field_Cardinality_field_fullname  protoreflect.FullName = "google.protobuf.Field.cardinality"
	Field_Number_field_fullname       protoreflect.FullName = "google.protobuf.Field.number"
	Field_Name_field_fullname         protoreflect.FullName = "google.protobuf.Field.name"
	Field_TypeUrl_field_fullname      protoreflect.FullName = "google.protobuf.Field.type_url"
	Field_OneofIndex_field_fullname   protoreflect.FullName = "google.protobuf.Field.oneof_index"
	Field_Packed_field_fullname       protoreflect.FullName = "google.protobuf.Field.packed"
	Field_Options_field_fullname      protoreflect.FullName = "google.protobuf.Field.options"
	Field_JsonName_field_fullname     protoreflect.FullName = "google.protobuf.Field.json_name"
	Field_DefaultValue_field_fullname protoreflect.FullName = "google.protobuf.Field.default_value"
)

// Field numbers for google.protobuf.Field.
const (
	Field_Kind_field_number         protoreflect.FieldNumber = 1
	Field_Cardinality_field_number  protoreflect.FieldNumber = 2
	Field_Number_field_number       protoreflect.FieldNumber = 3
	Field_Name_field_number         protoreflect.FieldNumber = 4
	Field_TypeUrl_field_number      protoreflect.FieldNumber = 6
	Field_OneofIndex_field_number   protoreflect.FieldNumber = 7
	Field_Packed_field_number       protoreflect.FieldNumber = 8
	Field_Options_field_number      protoreflect.FieldNumber = 9
	Field_JsonName_field_number     protoreflect.FieldNumber = 10
	Field_DefaultValue_field_number protoreflect.FieldNumber = 11
)

// Full and short names for google.protobuf.Field.Kind.
const (
	Field_Kind_enum_fullname = "google.protobuf.Field.Kind"
	Field_Kind_enum_name     = "Kind"
)

// Enum values for google.protobuf.Field.Kind.
const (
	Field_TYPE_UNKNOWN_enum_value  = 0
	Field_TYPE_DOUBLE_enum_value   = 1
	Field_TYPE_FLOAT_enum_value    = 2
	Field_TYPE_INT64_enum_value    = 3
	Field_TYPE_UINT64_enum_value   = 4
	Field_TYPE_INT32_enum_value    = 5
	Field_TYPE_FIXED64_enum_value  = 6
	Field_TYPE_FIXED32_enum_value  = 7
	Field_TYPE_BOOL_enum_value     = 8
	Field_TYPE_STRING_enum_value   = 9
	Field_TYPE_GROUP_enum_value    = 10
	Field_TYPE_MESSAGE_enum_value  = 11
	Field_TYPE_BYTES_enum_value    = 12
	Field_TYPE_UINT32_enum_value   = 13
	Field_TYPE_ENUM_enum_value     = 14
	Field_TYPE_SFIXED32_enum_value = 15
	Field_TYPE_SFIXED64_enum_value = 16
	Field_TYPE_SINT32_enum_value   = 17
	Field_TYPE_SINT64_enum_value   = 18
)

// Full and short names for google.protobuf.Field.Cardinality.
const (
	Field_Cardinality_enum_fullname = "google.protobuf.Field.Cardinality"
	Field_Cardinality_enum_name     = "Cardinality"
)

// Enum values for google.protobuf.Field.Cardinality.
const (
	Field_CARDINALITY_UNKNOWN_enum_value  = 0
	Field_CARDINALITY_OPTIONAL_enum_value = 1
	Field_CARDINALITY_REQUIRED_enum_value = 2
	Field_CARDINALITY_REPEATED_enum_value = 3
)

// Names for google.protobuf.Enum.
const (
	Enum_message_name     protoreflect.Name     = "Enum"
	Enum_message_fullname protoreflect.FullName = "google.protobuf.Enum"
)

// Field names for google.protobuf.Enum.
const (
	Enum_Name_field_name          protoreflect.Name = "name"
	Enum_Enumvalue_field_name     protoreflect.Name = "enumvalue"
	Enum_Options_field_name       protoreflect.Name = "options"
	Enum_SourceContext_field_name protoreflect.Name = "source_context"
	Enum_Syntax_field_name        protoreflect.Name = "syntax"
	Enum_Edition_field_name       protoreflect.Name = "edition"

	Enum_Name_field_fullname          protoreflect.FullName = "google.protobuf.Enum.name"
	Enum_Enumvalue_field_fullname     protoreflect.FullName = "google.protobuf.Enum.enumvalue"
	Enum_Options_field_fullname       protoreflect.FullName = "google.protobuf.Enum.options"
	Enum_SourceContext_field_fullname protoreflect.FullName = "google.protobuf.Enum.source_context"
	Enum_Syntax_field_fullname        protoreflect.FullName = "google.protobuf.Enum.syntax"
	Enum_Edition_field_fullname       protoreflect.FullName = "google.protobuf.Enum.edition"
)

// Field numbers for google.protobuf.Enum.
const (
	Enum_Name_field_number          protoreflect.FieldNumber = 1
	Enum_Enumvalue_field_number     protoreflect.FieldNumber = 2
	Enum_Options_field_number       protoreflect.FieldNumber = 3
	Enum_SourceContext_field_number protoreflect.FieldNumber = 4
	Enum_Syntax_field_number        protoreflect.FieldNumber = 5
	Enum_Edition_field_number       protoreflect.FieldNumber = 6
)

// Names for google.protobuf.EnumValue.
const (
	EnumValue_message_name     protoreflect.Name     = "EnumValue"
	EnumValue_message_fullname protoreflect.FullName = "google.protobuf.EnumValue"
)

// Field names for google.protobuf.EnumValue.
const (
	EnumValue_Name_field_name    protoreflect.Name = "name"
	EnumValue_Number_field_name  protoreflect.Name = "number"
	EnumValue_Options_field_name protoreflect.Name = "options"

	EnumValue_Name_field_fullname    protoreflect.FullName = "google.protobuf.EnumValue.name"
	EnumValue_Number_field_fullname  protoreflect.FullName = "google.protobuf.EnumValue.number"
	EnumValue_Options_field_fullname protoreflect.FullName = "google.protobuf.EnumValue.options"
)

// Field numbers for google.protobuf.EnumValue.
const (
	EnumValue_Name_field_number    protoreflect.FieldNumber = 1
	EnumValue_Number_field_number  protoreflect.FieldNumber = 2
	EnumValue_Options_field_number protoreflect.FieldNumber = 3
)

// Names for google.protobuf.Option.
const (
	Option_message_name     protoreflect.Name     = "Option"
	Option_message_fullname protoreflect.FullName = "google.protobuf.Option"
)

// Field names for google.protobuf.Option.
const (
	Option_Name_field_name  protoreflect.Name = "name"
	Option_Value_field_name protoreflect.Name = "value"

	Option_Name_field_fullname  protoreflect.FullName = "google.protobuf.Option.name"
	Option_Value_field_fullname protoreflect.FullName = "google.protobuf.Option.value"
)

// Field numbers for google.protobuf.Option.
const (
	Option_Name_field_number  protoreflect.FieldNumber = 1
	Option_Value_field_number protoreflect.FieldNumber = 2
)