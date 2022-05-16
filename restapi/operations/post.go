// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PostHandlerFunc turns a function with the right signature into a post handler
type PostHandlerFunc func(PostParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostHandlerFunc) Handle(params PostParams) middleware.Responder {
	return fn(params)
}

// PostHandler interface for that can handle valid post params
type PostHandler interface {
	Handle(PostParams) middleware.Responder
}

// NewPost creates a new http.Handler for the post operation
func NewPost(ctx *middleware.Context, handler PostHandler) *Post {
	return &Post{Context: ctx, Handler: handler}
}

/* Post swagger:route POST / post

Post post API

*/
type Post struct {
	Context *middleware.Context
	Handler PostHandler
}

func (o *Post) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostBody post body
//
// swagger:model PostBody
type PostBody struct {

	// Commands to execute in order.
	// Example: ["terraform -chdir=out/workflow/tfbase.tar.gz plan"]
	Commands []string `json:"commands"`

	// If set to true all commands are getting executed and errors ignored.
	// Example: true
	Continue *bool `json:"continue,omitempty"`

	// Variables set for all commands. This translatyes into TF_VAR_* environment variables.
	// Example: [{"name":"instance_name","value":"myinstance"}]
	Variables []*PostParamsBodyVariablesItems0 `json:"variables"`
}

// Validate validates this post body
func (o *PostBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateVariables(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostBody) validateVariables(formats strfmt.Registry) error {
	if swag.IsZero(o.Variables) { // not required
		return nil
	}

	for i := 0; i < len(o.Variables); i++ {
		if swag.IsZero(o.Variables[i]) { // not required
			continue
		}

		if o.Variables[i] != nil {
			if err := o.Variables[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("body" + "." + "variables" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("body" + "." + "variables" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this post body based on the context it is used
func (o *PostBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateVariables(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostBody) contextValidateVariables(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Variables); i++ {

		if o.Variables[i] != nil {
			if err := o.Variables[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("body" + "." + "variables" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("body" + "." + "variables" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostBody) UnmarshalBinary(b []byte) error {
	var res PostBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostParamsBodyVariablesItems0 post params body variables items0
//
// swagger:model PostParamsBodyVariablesItems0
type PostParamsBodyVariablesItems0 struct {

	// Name of the variable.
	Name string `json:"name,omitempty"`

	// Value of the variable.
	Value string `json:"value,omitempty"`
}

// Validate validates this post params body variables items0
func (o *PostParamsBodyVariablesItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post params body variables items0 based on context it is used
func (o *PostParamsBodyVariablesItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostParamsBodyVariablesItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostParamsBodyVariablesItems0) UnmarshalBinary(b []byte) error {
	var res PostParamsBodyVariablesItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
