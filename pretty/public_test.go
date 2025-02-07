// Copyright 2013 Google Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pretty

import (
	"testing"
	"time"
)

func TestDiff(t *testing.T) {
	type example struct {
		Name    string
		Age     int
		Friends []string
	}

	tests := []struct {
		desc      string
		got, want interface{}
		diff      string
	}{
		{
			desc: "basic struct",
			got: example{
				Name: "Zaphd",
				Age:  42,
				Friends: []string{
					"Ford Prefect",
					"Trillian",
					"Marvin",
				},
			},
			want: example{
				Name: "Zaphod",
				Age:  42,
				Friends: []string{
					"Ford Prefect",
					"Trillian",
				},
			},
			diff: ` {
- Name: "Zaphd",
+ Name: "Zaphod",
  Age: 42,
  Friends: [
   "Ford Prefect",
   "Trillian",
-  "Marvin",
  ],
 }`,
		},
	}

	for _, test := range tests {
		got, _ := Compare(test.got, test.want)
		if want := test.diff; got != want {
			t.Errorf("%s:", test.desc)
			t.Errorf("  got:  %q", got)
			t.Errorf("  want: %q", want)
		}
	}
}

func TestSkipZeroFields(t *testing.T) {
	type example struct {
		Name    string
		Species string
		Age     int
		Friends []string
	}

	tests := []struct {
		desc      string
		got, want interface{}
		diff      string
	}{
		{
			desc: "basic struct",
			got: example{
				Name:    "Zaphd",
				Species: "Betelgeusian",
				Age:     42,
			},
			want: example{
				Name:    "Zaphod",
				Species: "Betelgeusian",
				Age:     42,
				Friends: []string{
					"Ford Prefect",
					"Trillian",
					"",
				},
			},
			diff: ` {
- Name: "Zaphd",
+ Name: "Zaphod",
  Species: "Betelgeusian",
  Age: 42,
+ Friends: [
+  "Ford Prefect",
+  "Trillian",
+  "",
+ ],
 }`,
		},
	}

	cfg := *CompareConfig
	cfg.SkipZeroFields = true

	for _, test := range tests {
		got, _ := cfg.Compare(test.got, test.want)
		if want := test.diff; got != want {
			t.Errorf("%s:", test.desc)
			t.Errorf("  got:  %q", got)
			t.Errorf("  want: %q", want)
		}
	}
}

func TestRegressions(t *testing.T) {
	tests := []struct {
		issue  string
		config *Config
		value  interface{}
		want   string
	}{
		{
			issue: "kylelemons/godebug#13",
			config: &Config{
				PrintStringers: true,
			},
			value: struct{ Day *time.Weekday }{},
			want:  "{Day: nil}",
		},
	}

	for _, test := range tests {
		t.Run(test.issue, func(t *testing.T) {
			if got, want := test.config.Sprint(test.value), test.want; got != want {
				t.Errorf("%#v.Sprint(%#v) = %q, want %q", test.config, test.value, got, want)
			}
		})
	}
}

func TestNilvsEmptyStruct(t *testing.T) {
	type Bar struct {
		Value string
	}

	type Foo struct {
		Bar *Bar
	}

	a := Foo{
		Bar: nil,
	}

	b := Foo{
		Bar: &Bar{
			Value: "",
		},
	}

	want := ""
	wantPercentage := 1.0

	got, gotPercentage := Compare(a, b)
	if got != want || wantPercentage != gotPercentage {
		t.Errorf("GOT\n%#v\n%f", got, gotPercentage)
		t.Errorf("WANT\n%#v\n%f", want, wantPercentage)
	}
}

func TestIgnoreMoneyFormatDifferences(t *testing.T) {
	type example struct {
		Price string
	}

	tests := []struct {
		desc      string
		got, want interface{}
		diff      string
	}{
		{
			desc: "basic struct",
			got: example{
				Price: "3456.00",
			},
			want: example{
				Price: "$3,456.00",
			},
		},
	}

	cfg := *CompareConfig
	cfg.IgnoreMoneyFormatDifferences = true

	for _, test := range tests {
		got, _ := cfg.Compare(test.got, test.want)
		if want := test.diff; got != want {
			t.Errorf("%s:", test.desc)
			t.Errorf("  got:  %q", got)
			t.Errorf("  want: %q", want)
		}
	}
}
