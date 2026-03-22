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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GreenScalerServiceSpec defines the desired state of GreenScalerService
type GreenScalerServiceSpec struct {
	// timeZone specifies the IANA time zone, such as "UTC" or "Europe/Moscow".
	// If not specified, UTC is used.
	// +optional
	TimeZone string `json:"timeZone,omitempty"`

	// schedule - the schedule used to determine the required number of replicas.
	// The first suitable window is used.
	Targets []ScaleTarget `json:"targets"`

	// schedule - a schedule that determines the required number of replicas.
	// The first matching window is applied.
	// +kubebuilder:validation:MinItems=1
	Schedule []ScaleWindow `json:"schedule"`
}

// ScaleTarget describes the target workload to scale.
type ScaleTarget struct {
	// name - the name of the target resource.
	// +kubebuilder:validation:Enum=Deployment;StatefulSet
	Kind string `json:"kind"`
	// name - the name of the target resource.
	Name string `json:"name"`
	// namespace - the namespace of the target resource. If empty, namespace CR is taken.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// ScaleWindow describes the time window and the target size.
type ScaleWindow struct {
	// from - the beginning of the window in the HH:MM format.
	// +kubebuilder:validation:Pattern=`^([01][0-9]|2[0-3]):[0-5][0-9]$`
	From string `json:"from"`
	// to - the end of the window in the HH:MM format.
	// +kubebuilder:validation:Pattern=`^([01][0-9]|2[0-3]):[0-5][0-9]$`
	To string `json:"to"`
	// replicas - the desired number of replicas in the window.
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas"`
}

// GreenScalerServiceStatus defines the observed state of GreenScalerService.
type GreenScalerServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

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

	// lastAppliedReplicas - the last number of replicas that the operator tried to apply.
	// +optional
	LastAppliedReplicas *int32 `json:"lastAppliedReplicas,omitempty"`

	// lastReconcileTime - the time of the last successful CR processing.
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
