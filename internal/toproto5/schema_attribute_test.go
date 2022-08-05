package toproto5_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto5"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestSchemaAttribute(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		attr        fwschema.Attribute
		path        *tftypes.AttributePath
		expected    *tfprotov5.SchemaAttribute
		expectedErr string
	}

	tests := map[string]testCase{
		"deprecated": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:               types.StringType,
				Optional:           true,
				DeprecationMessage: "deprecated, use new_string instead",
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:       "string",
				Type:       tftypes.String,
				Optional:   true,
				Deprecated: true,
			},
		},
		"description-plain": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:        types.StringType,
				Optional:    true,
				Description: "A string attribute",
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:            "string",
				Type:            tftypes.String,
				Optional:        true,
				Description:     "A string attribute",
				DescriptionKind: tfprotov5.StringKindPlain,
			},
		},
		"description-markdown": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:                types.StringType,
				Optional:            true,
				MarkdownDescription: "A string attribute",
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:            "string",
				Type:            tftypes.String,
				Optional:        true,
				Description:     "A string attribute",
				DescriptionKind: tfprotov5.StringKindMarkdown,
			},
		},
		"description-both": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:                types.StringType,
				Optional:            true,
				Description:         "A string attribute",
				MarkdownDescription: "A string attribute (markdown)",
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:            "string",
				Type:            tftypes.String,
				Optional:        true,
				Description:     "A string attribute (markdown)",
				DescriptionKind: tfprotov5.StringKindMarkdown,
			},
		},
		"attr-string": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:     types.StringType,
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "string",
				Type:     tftypes.String,
				Optional: true,
			},
		},
		"attr-bool": {
			name: "bool",
			attr: tfsdk.Attribute{
				Type:     types.BoolType,
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "bool",
				Type:     tftypes.Bool,
				Optional: true,
			},
		},
		"attr-number": {
			name: "number",
			attr: tfsdk.Attribute{
				Type:     types.NumberType,
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "number",
				Type:     tftypes.Number,
				Optional: true,
			},
		},
		"attr-list": {
			name: "list",
			attr: tfsdk.Attribute{
				Type:     types.ListType{ElemType: types.NumberType},
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "list",
				Type:     tftypes.List{ElementType: tftypes.Number},
				Optional: true,
			},
		},
		"attr-map": {
			name: "map",
			attr: tfsdk.Attribute{
				Type:     types.MapType{ElemType: types.StringType},
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "map",
				Type:     tftypes.Map{ElementType: tftypes.String},
				Optional: true,
			},
		},
		"attr-object": {
			name: "object",
			attr: tfsdk.Attribute{
				Type: types.ObjectType{AttrTypes: map[string]attr.Type{
					"foo": types.StringType,
					"bar": types.NumberType,
					"baz": types.BoolType,
				}},
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name: "object",
				Type: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
					"foo": tftypes.String,
					"bar": tftypes.Number,
					"baz": tftypes.Bool,
				}},
				Optional: true,
			},
		},
		"attr-set": {
			name: "set",
			attr: tfsdk.Attribute{
				Type:     types.SetType{ElemType: types.NumberType},
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "set",
				Type:     tftypes.Set{ElementType: tftypes.Number},
				Optional: true,
			},
		},
		// TODO: add tuple attribute when we support it
		"required": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:     types.StringType,
				Required: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "string",
				Type:     tftypes.String,
				Required: true,
			},
		},
		"optional": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:     types.StringType,
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "string",
				Type:     tftypes.String,
				Optional: true,
			},
		},
		"computed": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:     types.StringType,
				Computed: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "string",
				Type:     tftypes.String,
				Computed: true,
			},
		},
		"optional-computed": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:     types.StringType,
				Computed: true,
				Optional: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:     "string",
				Type:     tftypes.String,
				Computed: true,
				Optional: true,
			},
		},
		"sensitive": {
			name: "string",
			attr: tfsdk.Attribute{
				Type:      types.StringType,
				Optional:  true,
				Sensitive: true,
			},
			path: tftypes.NewAttributePath(),
			expected: &tfprotov5.SchemaAttribute{
				Name:      "string",
				Type:      tftypes.String,
				Optional:  true,
				Sensitive: true,
			},
		},
		"nested-attr-single": {
			name: "single_nested",
			attr: tfsdk.Attribute{
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"string": {
						Type:     types.StringType,
						Optional: true,
					},
					"computed": {
						Type:      types.NumberType,
						Computed:  true,
						Sensitive: true,
					},
				}),
				Optional: true,
			},
			path:        tftypes.NewAttributePath(),
			expectedErr: "protocol version 5 cannot have Attributes set",
		},
		"nested-attr-list": {
			name: "list_nested",
			attr: tfsdk.Attribute{
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"string": {
						Type:     types.StringType,
						Optional: true,
					},
					"computed": {
						Type:      types.NumberType,
						Computed:  true,
						Sensitive: true,
					},
				}),
				Optional: true,
			},
			path:        tftypes.NewAttributePath(),
			expectedErr: "protocol version 5 cannot have Attributes set",
		},
		"nested-attr-set": {
			name: "set_nested",
			attr: tfsdk.Attribute{
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"string": {
						Type:     types.StringType,
						Optional: true,
					},
					"computed": {
						Type:      types.NumberType,
						Computed:  true,
						Sensitive: true,
					},
				}),
				Optional: true,
			},
			path:        tftypes.NewAttributePath(),
			expectedErr: "protocol version 5 cannot have Attributes set",
		},
		"attr-and-nested-attr-set": {
			name: "whoops",
			attr: tfsdk.Attribute{
				Type: types.StringType,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"testing": {
						Type:     types.StringType,
						Optional: true,
					},
				}),
				Optional: true,
			},
			path:        tftypes.NewAttributePath(),
			expectedErr: "protocol version 5 cannot have Attributes set",
		},
		"attr-and-nested-attr-unset": {
			name: "whoops",
			attr: tfsdk.Attribute{
				Optional: true,
			},
			path:        tftypes.NewAttributePath(),
			expectedErr: "must have Type set",
		},
		"attr-and-nested-attr-empty": {
			name: "whoops",
			attr: tfsdk.Attribute{
				Optional:   true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{}),
			},
			path:        tftypes.NewAttributePath(),
			expectedErr: "must have Type set",
		},
		"missing-required-optional-and-computed": {
			name: "whoops",
			attr: tfsdk.Attribute{
				Type: types.StringType,
			},
			path:        tftypes.NewAttributePath(),
			expectedErr: "must have Required, Optional, or Computed set",
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := toproto5.SchemaAttribute(context.Background(), tc.name, tc.path, tc.attr)
			if err != nil {
				if tc.expectedErr == "" {
					t.Errorf("Unexpected error: %s", err)
					return
				}
				if err.Error() != tc.expectedErr {
					t.Errorf("Expected error to be %q, got %q", tc.expectedErr, err.Error())
					return
				}
				// got expected error
				return
			}
			if err == nil && tc.expectedErr != "" {
				t.Errorf("Expected error to be %q, got nil", tc.expectedErr)
				return
			}
			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("Unexpected diff (+wanted, -got): %s", diff)
				return
			}
		})
	}
}
