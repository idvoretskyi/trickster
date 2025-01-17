/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
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

// Package reverseproxy provides the HTTP Reverse Proxy (no caching) Backend provider
package reverseproxy

import (
	"net/http"

	"github.com/tricksterproxy/trickster/pkg/backends"
	bo "github.com/tricksterproxy/trickster/pkg/backends/options"
)

var _ backends.Backend = (*Client)(nil)

// Client Implements the Proxy Client Interface
type Client struct {
	backends.Backend
}

// NewClient returns a new Client Instance
func NewClient(name string, o *bo.Options, router http.Handler) (backends.Backend, error) {
	c := &Client{}
	b, err := backends.New(name, o, c.RegisterHandlers, router, nil)
	c.Backend = b
	return c, err
}
