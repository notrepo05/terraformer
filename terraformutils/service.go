// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package terraformutils

import (
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils/providerwrapper"
)

type ServiceGenerator interface {
	InitResources() error
	GetResources() []Resource
	SetResources(resources []Resource)
	ParseFilter(rawFilter string) []ResourceFilter
	ParseFilters(rawFilters []string)
	PostConvertHook() error
	GetArgs() map[string]interface{}
	SetArgs(args map[string]interface{})
	SetName(name string)
	SetVerbose(bool)
	SetProviderName(name string)
	GetProviderName() string
	GetName() string
	InitialCleanup()
	PopulateIgnoreKeys(*providerwrapper.ProviderWrapper)
	PostRefreshCleanup()
}

type Service struct {
	Name         string
	Resources    []Resource
	ProviderName string
	Args         map[string]interface{}
	Filter       []ResourceFilter
	Verbose      bool
}

func (s *Service) SetProviderName(providerName string) {
	s.ProviderName = providerName
}

func (s *Service) GetProviderName() string {
	return s.ProviderName
}

func (s *Service) SetVerbose(verbose bool) {
	s.Verbose = verbose
}

func (s *Service) ParseFilters(rawFilters []string) {
	s.Filter = []ResourceFilter{}
	for _, rawFilter := range rawFilters {
		filters := s.ParseFilter(rawFilter)
		s.Filter = append(s.Filter, filters...)
	}
}

func (s *Service) ParseFilter(rawFilter string) []ResourceFilter {

	var filters []ResourceFilter
	filters = append(filters, ResourceFilter{
		ServiceName:      "",
		FieldPath:        strings.TrimPrefix("", "Name="),
		AcceptableValues: ParseFilterValues(strings.TrimPrefix("", "Value=")),
	})

	return filters
}

func (s *Service) SetName(name string) {
	s.Name = name
}
func (s *Service) GetName() string {
	return s.Name
}

func (s *Service) InitialCleanup() {
	FilterCleanup(s, true)
}

func (s *Service) PostRefreshCleanup() {
	if len(s.Filter) != 0 {
		FilterCleanup(s, false)
	}
}

func (s *Service) GetArgs() map[string]interface{} {
	return s.Args
}
func (s *Service) SetArgs(args map[string]interface{}) {
	s.Args = args
}

func (s *Service) GetResources() []Resource {
	return s.Resources
}
func (s *Service) SetResources(resources []Resource) {
	s.Resources = resources
}

func (s *Service) InitResources() error {
	panic("implement me")
}

func (s *Service) PostConvertHook() error {
	return nil
}

func (s *Service) PopulateIgnoreKeys(providerWrapper *providerwrapper.ProviderWrapper) {
	var resourcesTypes []string
	for _, r := range s.Resources {
		resourcesTypes = append(resourcesTypes, r.InstanceInfo.Type)
	}
	keys := IgnoreKeys(resourcesTypes, providerWrapper)
	for k, v := range keys {
		for i := range s.Resources {
			if s.Resources[i].InstanceInfo.Type == k {
				s.Resources[i].IgnoreKeys = append(s.Resources[i].IgnoreKeys, v...)
			}
		}
	}
}
