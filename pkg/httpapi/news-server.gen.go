// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Retrieves all posts
	// (GET /posts)
	GetAllPosts(w http.ResponseWriter, r *http.Request)
	// Creates a new post
	// (POST /posts)
	CreatePost(w http.ResponseWriter, r *http.Request)
	// Removes post by ID
	// (DELETE /posts/{postID})
	DeletePost(w http.ResponseWriter, r *http.Request, postID PostID)
	// Retrieves post by ID
	// (GET /posts/{postID})
	GetPost(w http.ResponseWriter, r *http.Request, postID PostID)
	// Updates post by ID
	// (PUT /posts/{postID})
	UpdatePost(w http.ResponseWriter, r *http.Request, postID PostID)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GetAllPosts operation middleware
func (siw *ServerInterfaceWrapper) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAllPosts(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// CreatePost operation middleware
func (siw *ServerInterfaceWrapper) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreatePost(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// DeletePost operation middleware
func (siw *ServerInterfaceWrapper) DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "postID" -------------
	var postID PostID

	err = runtime.BindStyledParameter("simple", false, "postID", chi.URLParam(r, "postID"), &postID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "postID", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeletePost(w, r, postID)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetPost operation middleware
func (siw *ServerInterfaceWrapper) GetPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "postID" -------------
	var postID PostID

	err = runtime.BindStyledParameter("simple", false, "postID", chi.URLParam(r, "postID"), &postID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "postID", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPost(w, r, postID)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// UpdatePost operation middleware
func (siw *ServerInterfaceWrapper) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "postID" -------------
	var postID PostID

	err = runtime.BindStyledParameter("simple", false, "postID", chi.URLParam(r, "postID"), &postID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "postID", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdatePost(w, r, postID)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/posts", wrapper.GetAllPosts)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/posts", wrapper.CreatePost)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/posts/{postID}", wrapper.DeletePost)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/posts/{postID}", wrapper.GetPost)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/posts/{postID}", wrapper.UpdatePost)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RW227jNhD9FWLaRybydlNgoTe39hZGDTtVkl6wWRiMNLa4kEgtSdkxAv17QVKybEu5",
	"AW6QvknkzHBmzpnLA8QyL6RAYTSED1AwxXI0qNzfpdRmMrJfCepY8cJwKSB052QyAgrc/hbMpEBBsBzt",
	"n1eioPB7yRUmEBpVIgUdp5gza20pVc4MhFCWPAEKZltYTW0UFyuoqsoq60IKjc6PES5ZmRn7GUthULhP",
	"VhQZj5n1KfimrWMPe2/8qHAJIfwQtAEG/lYHY6WkiuoX/HuHAd4IvC8wNpgQjWqNiqBVgYrCRKxZxpMI",
	"v5eo39Cl+l2i6ocrCjNpPstSJG/nxEwasnRP2rtazVr9VSEzaHnxi0y2jklKFqgM9wjuuddDpea2wwQK",
	"hpsMH1HzdxRydj9FsTIphD8PKORcNL8fu9zanci7bxi7PLbO72K3acyy+RLCL0+n7KZIjnUr2onePZAs",
	"WE8C/kpREFszZMM0qSWBtjVi7Z8ZnmNfenjyVHU+V2edXHytKDgS9CGY9ODghIm90+GtIOSM3MJsfr34",
	"PL+ZjW6BnBHRUIbcbQlPGqFo/MfN+Op6MZn9OZxOvGhNbbLkmCWapEwTXrPe+qkb3eE0Gg9H/yzGf0+u",
	"rq9qVVMqocnGZnMryzqRRMscTcrFimxSHqeEa8IyhSzZErzn2tgXiDZSsRU25iez63E0G04X4yiaR3vm",
	"vXXbEHiMzr2y7ROuQdyKPpBy1JqtHs1ec92HT9tDv3gEWms7rPZJe4gZNlA+W/Sdx7zq155qsfR6eX0c",
	"NQZbGy8Tb8upz4mesutEXzqZl1RdxrQhtfgLS69bO/aIi6W0j9VdC2a40WR4OSFXbooAhTUq7V34cD44",
	"H9hYZIGCFRxC+OiOqBuoLoTAuui+VuiisPG5zj5JIITf0Ayz7NLJHI3MnwaDV40FbjDXz1HFQd9GzpRi",
	"274xMf/dSl14F/oM7lwNjsapM7Ub9k+rNluBG0VlnjO1hRAiNIrjGm2hZ8Tnr6JuK+kmsKVbva9gO8BO",
	"MlGP6d/dbV4L1Gur532h4z3UhBGBG4eOE/A8Dx787lj5es3QYBeykTuvIdtfVx9pRK1IUK+ztqH0QXDC",
	"LF0MLp5X261vpyB9Li3lXUe729rpX9FHe8Z/kbyT8Nf3lxMz9u2xaBrQIRpF2YNGO8feJSB92+3/HB4f",
	"0iE4TsJNaJ/5o31ayaSM3Q+FUmUQQmpMocMgKHZX5wI3+pzLgBU8WH8I3KZzaGcqY5YdmAiDILOHqdQm",
	"/DT4NAjsCv5vAAAA//9qge8wlA8AAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
