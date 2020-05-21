package document

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParse(t *testing.T) {
	table := []struct {
		name         string
		testDataFile string
		peekFuncs    []peekFunc
		wantType     interface{}
		wantError    bool
	}{
		{
			name: "simple argocd",
			testDataFile: "argocd.yaml",
			peekFuncs: []peekFunc{NewArgoCD},
			wantType: &ArgoCD{},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(tt *testing.T) {
			b, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", test.testDataFile))
			if !assert.NoErrorf(tt, err, "error reading test data") {
				tt.FailNow()
			}
			got, err := parse(bytes.NewReader(b), test.peekFuncs...)
			if test.wantError {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.IsType(tt, test.wantType, got, "unexpected Document type")
		})
	}
}
