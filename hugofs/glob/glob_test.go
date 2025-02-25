// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package glob

import (
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestResolveRootDir(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		input    string
		expected string
	}{
		{"data/foo.json", "data"},
		{"a/b/**/foo.json", "a/b"},
		{"dat?a/foo.json", ""},
		{"a/b[a-c]/foo.json", "a"},
	} {
		c.Assert(ResolveRootDir(test.input), qt.Equals, test.expected)
	}
}

func TestFilterGlobParts(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		input    []string
		expected []string
	}{
		{[]string{"a", "*", "c"}, []string{"a", "c"}},
	} {
		c.Assert(FilterGlobParts(test.input), qt.DeepEquals, test.expected)
	}
}

func TestNormalizePath(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		input    string
		expected string
	}{
		{filepath.FromSlash("data/FOO.json"), "data/foo.json"},
		{filepath.FromSlash("/data/FOO.json"), "data/foo.json"},
		{filepath.FromSlash("./FOO.json"), "foo.json"},
		{"//", ""},
	} {
		c.Assert(NormalizePath(test.input), qt.Equals, test.expected)
	}
}

func TestGetGlob(t *testing.T) {
	c := qt.New(t)
	g, err := GetGlob("**.JSON")
	c.Assert(err, qt.IsNil)
	c.Assert(g.Match("data/my.json"), qt.Equals, true)
}

func TestFilenameFilter(t *testing.T) {
	c := qt.New(t)

	excludeAlmostAllJSON, err := NewFilenameFilter([]string{"a/b/c/foo.json"}, []string{"**.json"})
	c.Assert(err, qt.IsNil)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("data/my.json")), qt.Equals, false)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("a/b/c/foo.json")), qt.Equals, true)
	c.Assert(excludeAlmostAllJSON.Match(filepath.FromSlash("a/b/c/foo.bar")), qt.Equals, false)

	nopFilter, err := NewFilenameFilter(nil, nil)
	c.Assert(err, qt.IsNil)
	c.Assert(nopFilter.Match("ab.txt"), qt.Equals, true)

	includeOnlyFilter, err := NewFilenameFilter([]string{"**.json", "**.jpg"}, nil)
	c.Assert(err, qt.IsNil)
	c.Assert(includeOnlyFilter.Match("ab.json"), qt.Equals, true)
	c.Assert(includeOnlyFilter.Match("ab.jpg"), qt.Equals, true)
	c.Assert(includeOnlyFilter.Match("ab.gif"), qt.Equals, false)

	exlcudeOnlyFilter, err := NewFilenameFilter(nil, []string{"**.json", "**.jpg"})
	c.Assert(err, qt.IsNil)
	c.Assert(exlcudeOnlyFilter.Match("ab.json"), qt.Equals, false)
	c.Assert(exlcudeOnlyFilter.Match("ab.jpg"), qt.Equals, false)
	c.Assert(exlcudeOnlyFilter.Match("ab.gif"), qt.Equals, true)

	var nilFilter *FilenameFilter
	c.Assert(nilFilter.Match("ab.gif"), qt.Equals, true)

	funcFilter := NewFilenameFilterForInclusionFunc(func(s string) bool { return strings.HasSuffix(s, ".json") })
	c.Assert(funcFilter.Match("ab.json"), qt.Equals, true)
	c.Assert(funcFilter.Match("ab.bson"), qt.Equals, false)

}

func BenchmarkGetGlob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetGlob("**/foo")
		if err != nil {
			b.Fatal(err)
		}
	}
}
