/*
Copyright 2018 The Kubernetes Authors.

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

package validating_webhook

import (
	"context"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:admissionReviewVersions=v1,path=/validate-v1-pod,mutating=false,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=vwebhook.peerpods.io,sideEffects=None

// podValidator validates Pods
type PodValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// podValidator admits a pod if a specific runtimeClassName exists.
func (v *PodValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := v.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Spec.RuntimeClassName == nil || *pod.Spec.RuntimeClassName != "kata-remote-cc" {
		return admission.Denied("pod doesn't contain runtimeclass")
	}

	return admission.Allowed("")
}

// podValidator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *PodValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
