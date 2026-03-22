/*
Copyright 2026.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: JSON tags are required for API serialization. After changing types, run
// `make generate` and `make manifests`.

// GreenScalerServiceSpec defines the desired state of GreenScalerService.
type GreenScalerServiceSpec struct {
	// timeZone is an IANA time zone name (e.g. "UTC", "Europe/Moscow") used to
	// evaluate schedule windows. Defaults to UTC if empty.
	// +optional
	TimeZone string `json:"timeZone,omitempty"`

	// targets lists workloads (Deployments or StatefulSets) to scale.
	// +kubebuilder:validation:MinItems=1
	Targets []ScaleTarget `json:"targets"`

	// schedule defines time windows and the replica count for each window.
	// The first window that matches the current local time (in timeZone) wins.
	// +kubebuilder:validation:MinItems=1
	Schedule []ScaleWindow `json:"schedule"`
}

// ScaleTarget identifies a workload to scale.
type ScaleTarget struct {
	// kind is the workload API kind. Supported values: Deployment, StatefulSet.
	// +kubebuilder:validation:Enum=Deployment;StatefulSet
	Kind string `json:"kind"`
	// name is the target resource name.
	Name string `json:"name"`
	// namespace is the target namespace. If empty, the GreenScalerService's namespace is used.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// ScaleWindow is a half-open [from, to) interval in local wall-clock time (HH:MM, 24h).
// If from equals to, the window is treated as the full 24-hour day.
type ScaleWindow struct {
	// from is the inclusive start of the window (HH:MM).
	// +kubebuilder:validation:Pattern=`^([01][0-9]|2[0-3]):[0-5][0-9]$`
	From string `json:"from"`
	// to is the exclusive end of the window (HH:MM), except when from == to (full day).
	// +kubebuilder:validation:Pattern=`^([01][0-9]|2[0-3]):[0-5][0-9]$`
	To string `json:"to"`
	// replicas is the desired replica count while this window is active.
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas"`
}

// GreenScalerServiceStatus defines the observed state of GreenScalerService.
type GreenScalerServiceStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the GreenScalerService resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// lastAppliedReplicas is the replica count last written to target workloads.
	// +optional
	LastAppliedReplicas *int32 `json:"lastAppliedReplicas,omitempty"`

	// lastReconcileTime is when status was last updated successfully (controller clock).
	// +optional
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GreenScalerService is the Schema for the greenscalerservices API
type GreenScalerService struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of GreenScalerService
	// +required
	Spec GreenScalerServiceSpec `json:"spec"`

	// status defines the observed state of GreenScalerService
	// +optional
	Status GreenScalerServiceStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// GreenScalerServiceList contains a list of GreenScalerService
type GreenScalerServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []GreenScalerService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GreenScalerService{}, &GreenScalerServiceList{})
}
