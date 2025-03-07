package v1alpha1

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Freight represents a collection of versioned artifacts.
type Freight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// ID is a system-assigned value that is derived deterministically from the
	// contents of the Freight. i.e. Two pieces of Freight can be compared for
	// equality by comparing their IDs.
	ID string `json:"id,omitempty"`
	// Commits describes specific Git repository commits.
	Commits []GitCommit `json:"commits,omitempty"`
	// Images describes specific versions of specific container images.
	Images []Image `json:"images,omitempty"`
	// Charts describes specific versions of specific Helm charts.
	Charts []Chart `json:"charts,omitempty"`
	// Status describes the current status of this Freight.
	Status FreightStatus `json:"status,omitempty"`
}

func (f *Freight) GetStatus() *FreightStatus {
	return &f.Status
}

// UpdateID deterministically calculates a piece of Freight's ID based on its
// contents and assigns it to the ID field.
func (f *Freight) UpdateID() {
	size := len(f.Commits) + len(f.Images) + len(f.Charts)
	artifacts := make([]string, 0, size)
	for _, commit := range f.Commits {
		artifacts = append(
			artifacts,
			fmt.Sprintf("%s:%s", commit.RepoURL, commit.ID),
		)
	}
	for _, image := range f.Images {
		artifacts = append(
			artifacts,
			fmt.Sprintf("%s:%s", image.RepoURL, image.Tag),
		)
	}
	for _, chart := range f.Charts {
		artifacts = append(
			artifacts,
			fmt.Sprintf("%s/%s:%s", chart.RegistryURL, chart.Name, chart.Version),
		)
	}
	sort.Strings(artifacts)
	f.ID = fmt.Sprintf(
		"%x",
		sha1.Sum([]byte(strings.Join(artifacts, "|"))),
	)
}

// GitCommit describes a specific commit from a specific Git repository.
type GitCommit struct {
	// RepoURL is the URL of a Git repository.
	RepoURL string `json:"repoURL,omitempty"`
	// ID is the ID of a specific commit in the Git repository specified by
	// RepoURL.
	ID string `json:"id,omitempty"`
	// Branch denotes the branch of the repository where this commit was found.
	Branch string `json:"branch,omitempty"`
	// HealthCheckCommit is the ID of a specific commit. When specified,
	// assessments of Stage health will used this value (instead of ID) when
	// determining if applicable sources of Argo CD Application resources
	// associated with the Stage are or are not synced to this commit. Note that
	// there are cases (as in that of Kargo Render being utilized as a promotion
	// mechanism) wherein the value of this field may differ from the commit ID
	// found in the ID field.
	HealthCheckCommit string `json:"healthCheckCommit,omitempty"`
	// Message is the git commit message
	Message string `json:"message,omitempty"`
	// Author is the git commit author
	Author string `json:"author,omitempty"`
}

// FreightStatus describes a piece of Freight's most recently observed state.
type FreightStatus struct {
	// Qualifications describes the Stages for which this Freight has been
	// qualified.
	Qualifications map[string]Qualification `json:"qualifications,omitempty"`
}

// Qualification describes a Freight's qualification for a Stage.
type Qualification struct{}

//+kubebuilder:object:root=true

// FreightList is a list of Freight resources.
type FreightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Freight `json:"items"`
}
