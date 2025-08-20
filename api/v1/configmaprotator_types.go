/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigMapRotatorSpec defines the desired state of ConfigMapRotator
type ConfigMapRotatorSpec struct {
	// ConfigMapName is the name of the ConfigMap to rotate
	ConfigMapName string `json:"configMapName"`

	// RotationIntervalHours defines how often to rotate (in hours)
	RotationIntervalHours int `json:"rotationIntervalHours"`

	// DataTemplate defines the structure of data to generate
	DataTemplate map[string]string `json:"dataTemplate"`
}

// ConfigMapRotatorStatus defines the observed state of ConfigMapRotator.
type ConfigMapRotatorStatus struct {
	// LastRotationTime tracks when the ConfigMap was last rotated
	LastRotationTime *metav1.Time `json:"lastRotationTime,omitempty"`

	// CurrentGeneration tracks the current version
	CurrentGeneration int64 `json:"currentGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ConfigMapRotator is the Schema for the configmaprotators API
type ConfigMapRotator struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ConfigMapRotator
	// +required
	Spec ConfigMapRotatorSpec `json:"spec"`

	// status defines the observed state of ConfigMapRotator
	// +optional
	Status ConfigMapRotatorStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ConfigMapRotatorList contains a list of ConfigMapRotator
type ConfigMapRotatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigMapRotator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigMapRotator{}, &ConfigMapRotatorList{})
}
