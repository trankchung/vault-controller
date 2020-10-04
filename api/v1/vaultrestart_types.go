/*
Copyright 2020 Chung Tran <chung.k.tran@gmail.com>.

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

// +kubebuilder:validation:Required
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VaultRestartSpec defines the desired state of VaultRestart
type VaultRestartSpec struct {
	// Labels to match with pods to delete
	MatchingLabels client.MatchingLabels `json:"matchingLabels,omitempty"`

	// Reconcile interval in ns, s, m, or h units, i.e.: 60s
	// +optional
	// +kubebuilder:validation:Pattern=^\d+(ns|s|m|h)$
	PollingInterval string `json:"pollingInterval,omitempty"`
}

// VaultRestartStatus defines the observed state of VaultRestart
type VaultRestartStatus struct {
	// The Vault index saved from last run
	Index uint64 `json:"index,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VaultRestart is the Schema for the vaultrestarts API
// +kubebuilder:resource:shortName=vrst;vrsts,singular=vaultrestart
type VaultRestart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VaultRestartSpec   `json:"spec,omitempty"`
	Status VaultRestartStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VaultRestartList contains a list of VaultRestart
type VaultRestartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VaultRestart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VaultRestart{}, &VaultRestartList{})
}
