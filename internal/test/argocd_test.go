package test

import (
	"testing"

	"orbital/internal/argocd"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestHasValuesFile(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected bool
	}{
		{"Has values.yaml", []string{"other.yaml", "values.yaml"}, true},
		{"Has nested values.yaml", []string{"some/path/values.yaml"}, true},
		{"No match", []string{"test.yaml", "config.yml"}, false},
		{"Empty list", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, argocd.HasValuesFile(tt.files))
		})
	}
}

func TestUpdateOrAddParameter(t *testing.T) {
	source := &v1alpha1.ApplicationSource{
		Helm: &v1alpha1.ApplicationSourceHelm{
			Parameters: []v1alpha1.HelmParameter{
				{Name: "replicaCount", Value: "2"},
			},
		},
	}
	argocd.UpdateOrAddParameter(source, "5")

	assert.Equal(t, "5", source.Helm.Parameters[0].Value)

	source = &v1alpha1.ApplicationSource{
		Helm: &v1alpha1.ApplicationSourceHelm{},
	}
	argocd.UpdateOrAddParameter(source, "3")

	assert.Len(t, source.Helm.Parameters, 1)
	assert.Equal(t, "replicaCount", source.Helm.Parameters[0].Name)
	assert.Equal(t, "3", source.Helm.Parameters[0].Value)
	assert.True(t, source.Helm.Parameters[0].ForceString)
}

func TestUpdateSourceIfTargeted(t *testing.T) {
	source := &v1alpha1.ApplicationSource{
		Helm: &v1alpha1.ApplicationSourceHelm{
			ValueFiles: []string{"foo.yaml", "bar/values.yaml"},
		},
	}
	ok := argocd.UpdateSourceIfTargeted(source, "2")
	assert.True(t, ok)
	assert.Len(t, source.Helm.Parameters, 1)
	assert.Equal(t, "2", source.Helm.Parameters[0].Value)

	noHelm := &v1alpha1.ApplicationSource{}
	assert.False(t, argocd.UpdateSourceIfTargeted(noHelm, "1"))

	noValues := &v1alpha1.ApplicationSource{
		Helm: &v1alpha1.ApplicationSourceHelm{
			ValueFiles: []string{"foo.yaml"},
		},
	}
	assert.False(t, argocd.UpdateSourceIfTargeted(noValues, "1"))
}

func TestUpdateReplicaCount(t *testing.T) {
	spec := &v1alpha1.ApplicationSpec{
		Sources: []v1alpha1.ApplicationSource{
			{
				Helm: &v1alpha1.ApplicationSourceHelm{
					ValueFiles: []string{"values.yaml"},
				},
			},
		},
	}
	newSpec, err := argocd.UpdateReplicaCount(spec)
	assert.NoError(t, err)
	assert.Equal(t, "1", newSpec.Sources[0].Helm.Parameters[0].Value)

	spec2 := &v1alpha1.ApplicationSpec{
		Sources: []v1alpha1.ApplicationSource{
			{
				Helm: &v1alpha1.ApplicationSourceHelm{
					ValueFiles: []string{"not-it.yaml"},
				},
			},
		},
	}
	_, err = argocd.UpdateReplicaCount(spec2)
	assert.ErrorContains(t, err, "could not find a source with a valueFile ending in")
}
