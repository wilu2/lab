package openapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/legacy"
)

type Validator struct {
	router *legacy.Router
}

func NewValidator(spec string) (*Validator, error) {
	ctx := context.Background()
	loader := openapi3.Loader{Context: ctx}
	doc, err := loader.LoadFromData([]byte(spec))
	if err != nil {
		return nil, err
	}
	router, err := legacy.NewRouter(doc)
	if err != nil {
		return nil, err
	}
	return &Validator{router: router.(*legacy.Router)}, nil
}

func (v *Validator) newRequestValidationInput(req *http.Request) (*openapi3filter.RequestValidationInput, error) {
	route, pathParams, err := v.router.FindRoute(req)
	if err != nil {
		return nil, err
	}
	return &openapi3filter.RequestValidationInput{
		Request:    req,
		Route:      route,
		PathParams: pathParams,
	}, nil
}

func (v *Validator) ValidateRequest(ctx context.Context, req *http.Request) error {
	requestValidationInput, err := v.newRequestValidationInput(req)
	if err != nil {
		return err
	}

	return openapi3filter.ValidateRequest(ctx, requestValidationInput)
}

func (v *Validator) ValidateResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
	requestValidationInput, err := v.newRequestValidationInput(req)
	if err != nil {
		return fmt.Errorf("new request validation input: %w", err)
	}

	// restore the body after validation
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)
	defer func() { resp.Body = ioutil.NopCloser(&buf) }()
	defer ioutil.ReadAll(tee)

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 resp.StatusCode,
		Header:                 resp.Header,
		Body:                   ioutil.NopCloser(tee),
	}

	err = openapi3filter.ValidateResponse(ctx, responseValidationInput)
	if err != nil {
		return fmt.Errorf("%s validate response: %w", string(buf.Bytes()), err)
	}
	return nil
}
