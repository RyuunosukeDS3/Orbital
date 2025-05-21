package argocd

import (
	"context"
	"fmt"
	"strings"

	"orbital/internal/config"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	appclient "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

// Change SetReplicaCount from a function to a variable of function type
var SetReplicaCount = setReplicaCountImpl

func setReplicaCountImpl(appName string, replicaCount string) error {
	clientOpts := apiclient.ClientOptions{
		ServerAddr: config.AppConfig.ArgoCDURL,
		AuthToken:  config.AppConfig.ArgoCDToken,
		Insecure:   true,
	}
	conn, appIf, err := apiclient.NewClientOrDie(&clientOpts).NewApplicationClient()
	if err != nil {
		return fmt.Errorf("failed to create ArgoCD client: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("warning: failed to close connection: %v\n", err)
		}
	}()

	ctx := context.Background()

	app, err := appIf.Get(ctx, &appclient.ApplicationQuery{Name: &appName})
	if err != nil {
		return fmt.Errorf("failed to get ArgoCD app: %w", err)
	}

	updatedSpec, err := UpdateReplicaCount(app.Spec.DeepCopy(), replicaCount)
	if err != nil {
		return err
	}

	_, err = appIf.UpdateSpec(ctx, &appclient.ApplicationUpdateSpecRequest{
		Name: &appName,
		Spec: updatedSpec,
	})
	if err != nil {
		return fmt.Errorf("failed to update replicaCount: %w", err)
	}

	return nil
}

// Rest of your functions stay unchanged

func UpdateReplicaCount(spec *v1alpha1.ApplicationSpec, replicaCount string) (*v1alpha1.ApplicationSpec, error) {
	for i := range spec.Sources {
		if UpdateSourceIfTargeted(&spec.Sources[i], replicaCount) {
			return spec, nil
		}
	}
	return nil, fmt.Errorf("could not find a source with a valueFile ending in 'values.yaml'")
}

func UpdateSourceIfTargeted(source *v1alpha1.ApplicationSource, replicaCount string) bool {
	if source.Helm == nil {
		return false
	}

	if !HasValuesFile(source.Helm.ValueFiles) {
		return false
	}

	UpdateOrAddParameter(source, replicaCount)
	return true
}

func HasValuesFile(files []string) bool {
	for _, vf := range files {
		if strings.HasSuffix(vf, "values.yaml") {
			return true
		}
	}
	return false
}

func UpdateOrAddParameter(source *v1alpha1.ApplicationSource, replicaCount string) {
	for i, p := range source.Helm.Parameters {
		if p.Name == "replicaCount" {
			source.Helm.Parameters[i].Value = replicaCount
			return
		}
	}
	source.Helm.Parameters = append(source.Helm.Parameters, v1alpha1.HelmParameter{
		Name:        "replicaCount",
		Value:       replicaCount,
		ForceString: true,
	})
}
