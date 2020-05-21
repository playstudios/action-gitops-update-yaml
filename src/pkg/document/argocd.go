package document

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

type argoApplication struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`

	Spec argoApplicationSpec `yaml:"spec"`
}

type argoApplicationSpec struct {
	Source *argoApplicationSource `yaml:"source"`
}

type argoApplicationSource struct {
	TargetRevision string              `yaml:"targetRevision"`
	Helm           argoApplicationHelm `yaml:"helm"`
}

type argoApplicationHelm struct {
	Values map[string]interface{} `yaml:"values"`
}

func (h *argoApplicationHelm) UnmarshalYAML(value *yaml.Node) error {
	m := map[string]string{}
	if err := value.Decode(&m); err != nil {
		return err
	}
	v := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(m["values"]), &v); err != nil {
		return err
	}
	h.Values = v
	return nil
}

type ArgoCD struct {
	manifest argoApplication
}

var _ Document = &ArgoCD{}

func NewArgoCD(r io.Reader) (Document, error) {
	d := &ArgoCD{
		manifest: argoApplication{},
	}
	if err := yaml.NewDecoder(r).Decode(&d.manifest); err != nil {
		return nil, fmt.Errorf("error decoding yaml: %w", err)
	}

	if d.manifest.ApiVersion == "argoproj.io/v1alpha1" && d.manifest.Kind == "Application" {
		return d, nil
	}
	return nil, nil
}

func (a *ArgoCD) Get(path string) (interface{}, error) {
	i, err := walk(strings.Split(path, "."), a.manifest, nil)
	if err != nil {
		return nil, fmt.Errorf("error walking values: %w", err)
	}
	return i, nil
}

func (a *ArgoCD) Set(path string, v interface{}) error {
	_, err := walk(strings.Split(path, "."), a.manifest, setValue(v))
	return err
}

func (a *ArgoCD) Write(w io.WriteCloser) error {
	// TODO argoApplicationHelm.Values needs to be encoded as a string, not a map
	return yaml.NewEncoder(w).Encode(a.manifest)
}
