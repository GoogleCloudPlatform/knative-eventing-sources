/*
Copyright 2021 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by injection-gen. DO NOT EDIT.

package filteredFactory

import (
	context "context"

	externalversions "github.com/google/knative-gcp/pkg/client/informers/externalversions"
	client "github.com/google/knative-gcp/pkg/client/injection/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformerFactory(withInformerFactory)
}

// Key is used as the key for associating information with a context.Context.
type Key struct {
	Selector string
}

type LabelKey struct{}

func WithSelectors(ctx context.Context, selector ...string) context.Context {
	return context.WithValue(ctx, LabelKey{}, selector)
}

func withInformerFactory(ctx context.Context) context.Context {
	c := client.Get(ctx)
	opts := []externalversions.SharedInformerOption{}
	if injection.HasNamespaceScope(ctx) {
		opts = append(opts, externalversions.WithNamespace(injection.GetNamespaceScope(ctx)))
	}
	untyped := ctx.Value(LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	for _, selector := range labelSelectors {
		thisOpts := append(opts, externalversions.WithTweakListOptions(func(l *v1.ListOptions) {
			l.LabelSelector = selector
		}))
		ctx = context.WithValue(ctx, Key{Selector: selector},
			externalversions.NewSharedInformerFactoryWithOptions(c, controller.GetResyncPeriod(ctx), thisOpts...))
	}
	return ctx
}

// Get extracts the InformerFactory from the context.
func Get(ctx context.Context, selector string) externalversions.SharedInformerFactory {
	untyped := ctx.Value(Key{Selector: selector})
	if untyped == nil {
		logging.FromContext(ctx).Panicf(
			"Unable to fetch github.com/google/knative-gcp/pkg/client/informers/externalversions.SharedInformerFactory with selector %s from context.", selector)
	}
	return untyped.(externalversions.SharedInformerFactory)
}
