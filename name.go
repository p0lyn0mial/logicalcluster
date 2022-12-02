/*
Copyright 2022 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logicalcluster

import "regexp"

var (
	clusterNameRegExp = regexp.MustCompile(clusterNameString)
)

const (
	clusterNameString string = "^[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"
)

// Name holds a value that uniquely identifies a logical cluster, which:
//  1. can be used to access a cluster via `/clusters/$value`
//  2. is part of an etcd key path
type Name string

// Path creates a new Path object from the stored value.
// A convenience method for working with methods which accept a Path type.
func (n Name) Path() Path {
	return New(string(n))
}

// String returns string representation of the stored value.
// Satisfies the Stringer interface.
func (n Name) String() string {
	return string(n)
}

// IsValid returns true if the stored value matches a defined format.
// A convenience method that could be used for enforcing a well-known structure of a logical cluster name.
//
// As of today a valid value starts with a lower-case letter or digit
// and contains only lower-case letters, digits and hyphens.
func (n Name) IsValid() bool {
	return clusterNameRegExp.MatchString(string(n))
}

// Empty returns true if the stored value is unset.
// It is a convenience method for checking against an empty value.
func (n Name) Empty() bool {
	return n == ""
}

// Object is a local interface representation of the Kubernetes metav1.Object, to avoid dependencies on k8s.io/apimachinery.
type Object interface {
	GetAnnotations() map[string]string
}

// AnnotationKey is the name of the annotation key used to denote an object's logical cluster.
const AnnotationKey = "kcp.dev/cluster"

// From returns the logical cluster from the given object.
func From(obj Object) Name {
	return Name(obj.GetAnnotations()[AnnotationKey])
}
