/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and cert-manager-cop contributors
SPDX-License-Identifier: Apache-2.0
*/

package generator

import (
	"context"
	"fmt"
	"io/fs"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	helmgenerator "github.com/sap/component-operator-runtime/pkg/manifests/helm"
	componentoperatorruntimetypes "github.com/sap/component-operator-runtime/pkg/types"
)

type ResourceGenerator struct {
	generator *helmgenerator.HelmGenerator
}

func NewResourceGenerator(fsys fs.FS, chartPath string, client client.Client) (*ResourceGenerator, error) {
	generator, err := helmgenerator.NewHelmGenerator(fsys, chartPath, client)
	if err != nil {
		return nil, err
	}
	return &ResourceGenerator{generator: generator}, nil
}

func (g *ResourceGenerator) Generate(ctx context.Context, namespace string, name string, parameters componentoperatorruntimetypes.Unstructurable) ([]client.Object, error) {
	values := parameters.ToUnstructured()

	values["fullnameOverride"] = name

	delete(values, "namespace")
	delete(values, "name")

	values["installCRDs"] = true

	var additionalResources []client.Object
	if v, ok := values["additionalResources"]; ok {
		v, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("invalid parameter found (expected array): .additionalResources")
		}
		for i, object := range v {
			object, ok := object.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid parameter found (expected object): .additionalResources[%d]", i)
			}
			additionalResources = append(additionalResources, &unstructured.Unstructured{Object: object})
		}
		delete(values, "additionalResources")
	}

	resources, err := g.generator.Generate(ctx, namespace, name, componentoperatorruntimetypes.UnstructurableMap(values))
	if err != nil {
		return nil, err
	}

	return append(resources, additionalResources...), nil
}
