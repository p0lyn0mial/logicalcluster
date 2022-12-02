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

import (
	"path"
	"regexp"
	"strings"
)

var (
	// Wildcard is the Path indicating a requests that spans many logical clusters.
	Wildcard = Path{value: "*"}
)

const (
	separator = ":"
)

// Path represents a colon separated list of words describing a path in a logical cluster hierarchy,
// like a file path in a file-system.
//
// For instance, in the following hierarchy:
//
// root/                    (62208dab)
// ├── accounting           (c8a942c5)
// │   └── us-west          (33bab531)
// │       └── invoices     (f5865fce)
// └── management           (e7e08986)
//
//	└── us-west-invoices    (f5865fce)
//
// the following would all be valid paths for the `project` logical cluster:
//   - root:accounting:us-west:invoices
//   - 62208dab:accounting:us-west:invoices
//   - c8a942c5:us-west:invoices
//   - 33bab531:invoices
//   - f5865fce
//   - root:management:us-west-invoices
//   - 62208dab:management:us-west-invoices
//   - e7e08986:us-west-invoices
type Path struct {
	value string
}

// New returns a new Path.
func New(value string) Path {
	return Path{value}
}

// NewValidated returns a Path and whether it is valid.
func NewValidated(value string) (Path, bool) {
	n := Path{value}
	return n, n.IsValid()
}

// Empty returns true if the stored path is unset.
// It is a convenience method for checking against an empty value.
func (n Path) Empty() bool {
	return n.value == ""
}

// Name return a new Name object from the stored path and whether it can be created.
// A convenience method for working with methods which accept a Name type.
func (n Path) Name() (Name, bool) {
	if _, hasParent := n.Parent(); hasParent {
		return "", false
	}
	return Name(n.value), true
}

// RequestPath returns a URL path segment used to access API for the stored path.
func (n Path) RequestPath() string {
	return path.Join("/clusters", n.value)
}

// String returns string representation of the stored value.
// Satisfies the Stringer interface.
func (n Path) String() string {
	return n.value
}

// Parent returns a new path with all but the last element of the stored path.
func (n Path) Parent() (Path, bool) {
	parent, _ := n.Split()
	return parent, parent.value != ""
}

// Split splits the path immediately following the final colon,
// separating it into a new path and a logical cluster name component.
// If there is no colon in the path,
// Split returns an empty path and a name set to the path.
func (n Path) Split() (parent Path, name string) {
	i := strings.LastIndex(n.value, separator)
	if i < 0 {
		return Path{}, n.value
	}
	return Path{n.value[:i]}, n.value[i+1:]
}

// Base returns the last element of the path.
func (n Path) Base() string {
	_, name := n.Split()
	return name
}

// Join returns a new path by adding the given path segment
// into already existing path and separating it with a colon.
func (n Path) Join(name string) Path {
	if n.value == "" {
		return Path{name}
	}
	return Path{n.value + separator + name}
}

func (n Path) HasPrefix(other Path) bool {
	return strings.HasPrefix(n.value, other.value)
}

const lclusterNameFmt string = "[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?"

var lclusterRegExp = regexp.MustCompile("^" + lclusterNameFmt + "(:" + lclusterNameFmt + ")*$")

// IsValid returns true if the path is a Wildcard or a colon separated list of words where each word
// starts with a lower-case letter and contains only lower-case letters, digits and hyphens.
func (n Path) IsValid() bool {
	return n == Wildcard || lclusterRegExp.MatchString(n.value)
}
