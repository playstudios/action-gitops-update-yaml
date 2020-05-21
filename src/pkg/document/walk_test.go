package document

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestWalk(t *testing.T) {
	d := valueMap{
		"a":  "b",
		"a1": 1,
		"a2": 123.456,
		"c": []string{
			"123",
			"456",
			"789",
		},
		"c1": []int{
			963,
			852,
			741,
		},
		"d": map[string]interface{}{
			"e": "f",
			"g": 1.2,
			"h": 5,
			"i": []int{
				9,
				8,
				7,
			},
			"i1": []float32{
				1.23,
				4.56,
				7.89,
			},
			"j": map[string]interface{}{
				"k": map[string]interface{}{
					"l": map[string]interface{}{
						"m": "n",
					},
				},
			},
		},
	}

	table := []struct {
		name      string
		path      string
		want      interface{}
		wantError bool
	}{
		{
			name: "root level string",
			path: "a",
			want: "b",
		},
		{
			name: "root level int",
			path: "a1",
			want: 1,
		},
		{
			name: "root level float",
			path: "a2",
			want: 123.456,
		},
		{
			name:      "root level invalid key",
			path:      "a200",
			want:      nil,
			wantError: true,
		},

		{
			name: "root string slice idx 0",
			path: "c.0",
			want: "123",
		},
		{
			name: "root string slice idx 1",
			path: "c.1",
			want: "456",
		},
		{
			name: "root string slice idx 2",
			path: "c.2",
			want: "789",
		},
		{
			name:      "root string slice idx out of bounds",
			path:      "c.10",
			want:      nil,
			wantError: true,
		},

		{
			name: "nested 1 string",
			path: "d.e",
			want: "f",
		},
		{
			name: "nested 1 float",
			path: "d.g",
			want: 1.2,
		},

		{
			name: "nested 1 int slice",
			path: "d.i.0",
			want: 9,
		},
		{
			name:      "nested 1 int slice out of bounds",
			path:      "d.i.10",
			want:      nil,
			wantError: true,
		},

		{
			name: "nested 1 float slice",
			path: "d.i1.1",
			want: 4.56,
		},

		{
			name: "nested 2 string",
			path: "d.j.k.l.m",
			want: "n",
		},
	}

	for _, test := range table {
		t.Run(test.name, func(tt *testing.T) {
			got, err := walk(strings.Split(test.path, "."), &d, nil)
			if test.wantError {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.EqualValues(tt, test.want, got)
		})
	}
}
