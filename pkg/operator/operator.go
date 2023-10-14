/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and cert-manager-cop contributors
SPDX-License-Identifier: Apache-2.0
*/

package operator

import (
	"embed"
	"flag"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sap/component-operator-runtime/pkg/component"
	"github.com/sap/component-operator-runtime/pkg/operator"

	operatorv1alpha1 "github.com/sap/cert-manager-cop/api/v1alpha1"
	"github.com/sap/cert-manager-cop/internal/generator"
)

const Name = "cert-manager-cop.cs.sap.com"

//go:embed all:data
var data embed.FS

type Options struct {
	Name       string
	FlagPrefix string
}

type Operator struct {
	options Options
}

var defaultOperator operator.Operator = New()

func GetName() string {
	return defaultOperator.GetName()
}

func InitScheme(scheme *runtime.Scheme) {
	defaultOperator.InitScheme(scheme)
}

func InitFlags(flagset *flag.FlagSet) {
	defaultOperator.InitFlags(flagset)
}

func ValidateFlags() error {
	return defaultOperator.ValidateFlags()
}

func GetUncacheableTypes() []client.Object {
	return defaultOperator.GetUncacheableTypes()
}

func Setup(mgr ctrl.Manager, discoveryClient discovery.DiscoveryInterface) error {
	return defaultOperator.Setup(mgr, discoveryClient)
}

func New() *Operator {
	return NewWithOptions(Options{})
}

func NewWithOptions(options Options) *Operator {
	operator := &Operator{options: options}
	if operator.options.Name == "" {
		operator.options.Name = Name
	}
	return operator
}

func (o *Operator) GetName() string {
	return o.options.Name
}

func (o *Operator) InitScheme(scheme *runtime.Scheme) {
	utilruntime.Must(operatorv1alpha1.AddToScheme(scheme))
}

func (o *Operator) InitFlags(flagset *flag.FlagSet) {
}

func (o *Operator) ValidateFlags() error {
	return nil
}

func (o *Operator) GetUncacheableTypes() []client.Object {
	return []client.Object{&operatorv1alpha1.CertManager{}}
}

func (o *Operator) Setup(mgr ctrl.Manager, discoveryClient discovery.DiscoveryInterface) error {
	resourceGenerator, err := generator.NewResourceGenerator(
		Name,
		data,
		"data/charts/cert-manager",
		mgr.GetClient(),
		discoveryClient,
	)
	if err != nil {
		return errors.Wrap(err, "error initializing resource generator")
	}

	if err := component.NewReconciler[*operatorv1alpha1.CertManager](
		o.options.Name,
		mgr.GetClient(),
		discoveryClient,
		mgr.GetEventRecorderFor(o.options.Name),
		mgr.GetScheme(),
		resourceGenerator,
	).SetupWithManager(mgr); err != nil {
		return errors.Wrapf(err, "unable to create controller")
	}

	return nil
}