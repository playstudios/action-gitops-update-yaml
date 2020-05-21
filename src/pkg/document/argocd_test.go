package document

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestArgoCD_Get(t *testing.T) {
	table := []struct {
		name         string
		testDataFile string
		path         string
		want         interface{}
		wantError    bool
	}{
		{
			name:         "helm value",
			testDataFile: "argocd.yaml",
			path:         "Spec.Source.Helm.Values.test.enabled",
			want:         false,
		},
		{
			name:         "targetRevision",
			testDataFile: "argocd.yaml",
			path:         "Spec.Source.TargetRevision",
			want:         "3.11.0",
		},
		{
			name:         "case insensitive",
			testDataFile: "argocd.yaml",
			path:         "spec.Source.targetRevision",
			want:         "3.11.0",
			wantError:    false,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(tt *testing.T) {
			b, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", test.testDataFile))
			if !assert.NoErrorf(tt, err, "error reading test data") {
				tt.FailNow()
			}
			doc, err := parse(bytes.NewReader(b), NewArgoCD)
			assert.NoError(tt, err, "error parsing yaml")
			assert.IsType(tt, &ArgoCD{}, doc, "unexpected Document type")

			got, err := doc.Get(test.path)
			if test.wantError {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.Equal(tt, test.want, got)
		})
	}
}

func TestArgoCD_Set(t *testing.T) {
	table := []struct {
		name         string
		testDataFile string
		path         string
		set          interface{}
		want         interface{}
		wantError    bool
	}{
		{
			name:         "helm value",
			testDataFile: "argocd.yaml",
			path:         "Spec.Source.Helm.Values.test.enabled",
			want:         "abc",
		},
		{
			name:         "target revision",
			testDataFile: "argocd.yaml",
			path:         "Spec.Source.targetRevision",
			want:         "abc",
		},
		// {
		// 	name:         "set valid but unsettable path",
		// 	testDataFile: "argocd.yaml",
		// 	path:         "ApiVersion",
		// 	want:         nil,
		// 	wantError:    true,
		// },
	}

	for _, test := range table {
		t.Run(test.name, func(tt *testing.T) {
			b, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", test.testDataFile))
			if !assert.NoErrorf(tt, err, "error reading test data") {
				tt.FailNow()
			}
			doc, err := parse(bytes.NewReader(b), NewArgoCD)
			assert.NoError(tt, err, "error parsing yaml")
			assert.IsType(tt, &ArgoCD{}, doc, "unexpected Document type")

			err = doc.Set(test.path, test.want)
			if test.wantError {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			got, err := doc.Get(test.path)
			assert.NoError(tt, err, "error getting new value")
			assert.EqualValues(tt, test.want, got)
		})
	}
}
