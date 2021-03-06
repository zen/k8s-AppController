// Copyright 2016 Mirantis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"encoding/json"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/rest"
)

type Dependency struct {
	unversioned.TypeMeta `json:",inline"`

	// Standard object metadata
	api.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Parent string            `json:"parent"`
	Child  string            `json:"child"`
	Meta   map[string]string `json:"meta,omitempty"`
}

type DependencyList struct {
	unversioned.TypeMeta `json:",inline"`

	// Standard list metadata.
	unversioned.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Dependency `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type DependenciesInterface interface {
	List(opts api.ListOptions) (*DependencyList, error)
	Create(*Dependency) (*Dependency, error)
	Delete(name string, opts *api.DeleteOptions) error
}

type dependencies struct {
	rc *rest.RESTClient
}

func newDependencies(c rest.Config) (*dependencies, error) {
	rc, err := thirdPartyResourceRESTClient(&c)
	if err != nil {
		return nil, err
	}

	return &dependencies{rc}, nil
}

func (c dependencies) List(opts api.ListOptions) (*DependencyList, error) {
	resp, err := c.rc.Get().
		Namespace("default").
		Resource("dependencies").
		LabelsSelectorParam(opts.LabelSelector).
		DoRaw()

	if err != nil {
		return nil, err
	}

	result := &DependencyList{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c dependencies) Create(d *Dependency) (result *Dependency, err error) {
	result = &Dependency{}
	err = c.rc.Post().
		Namespace("default").
		Resource("Dependencies").
		Body(d).
		Do().
		Into(result)
	return
}

func (c *dependencies) Delete(name string, opts *api.DeleteOptions) error {
	return c.rc.Delete().
		Namespace("default").
		Resource("dependencies").
		Name(name).
		Body(opts).
		Do().
		Error()
}
