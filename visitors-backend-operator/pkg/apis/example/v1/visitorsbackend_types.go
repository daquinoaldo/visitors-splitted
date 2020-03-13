package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VisitorsBackendSpec defines the desired state of VisitorsBackend
type VisitorsBackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	DatabaseSize int32 `json:"database-size"`
	BackendSize  int32 `json:"backend-size"`
}

// VisitorsBackendStatus defines the observed state of VisitorsBackend
type VisitorsBackendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	BackendImage string `json:"backendImage"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VisitorsBackend is the Schema for the visitorsbackends API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=visitorsbackends,scope=Namespaced
type VisitorsBackend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VisitorsBackendSpec   `json:"spec,omitempty"`
	Status VisitorsBackendStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VisitorsBackendList contains a list of VisitorsBackend
type VisitorsBackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VisitorsBackend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VisitorsBackend{}, &VisitorsBackendList{})
}
