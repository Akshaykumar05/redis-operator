/*
Copyright 2020 Opstree Solutions.

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

package v1beta2

import (
	common "github.com/OT-CONTAINER-KIT/redis-operator/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RedisClusterSpec defines the desired state of RedisCluster
type RedisClusterSpec struct {
	Size             *int32           `json:"clusterSize"`
	KubernetesConfig KubernetesConfig `json:"kubernetesConfig"`
	// +kubebuilder:default:=v7
	ClusterVersion *string `json:"clusterVersion,omitempty"`
	// +kubebuilder:default:={livenessProbe:{initialDelaySeconds: 1, timeoutSeconds: 1, periodSeconds: 10, successThreshold: 1, failureThreshold:3}, readinessProbe:{initialDelaySeconds: 1, timeoutSeconds: 1, periodSeconds: 10, successThreshold: 1, failureThreshold:3}}
	RedisLeader RedisLeader `json:"redisLeader,omitempty"`
	// +kubebuilder:default:={livenessProbe:{initialDelaySeconds: 1, timeoutSeconds: 1, periodSeconds: 10, successThreshold: 1, failureThreshold:3}, readinessProbe:{initialDelaySeconds: 1, timeoutSeconds: 1, periodSeconds: 10, successThreshold: 1, failureThreshold:3}}
	RedisFollower      RedisFollower                `json:"redisFollower,omitempty"`
	RedisExporter      *RedisExporter               `json:"redisExporter,omitempty"`
	Storage            *ClusterStorage              `json:"storage,omitempty"`
	PodSecurityContext *corev1.PodSecurityContext   `json:"podSecurityContext,omitempty"`
	PriorityClassName  string                       `json:"priorityClassName,omitempty"`
	Resources          *corev1.ResourceRequirements `json:"resources,omitempty"`
	TLS                *TLSConfig                   `json:"TLS,omitempty"`
	ACL                *ACLConfig                   `json:"acl,omitempty"`
	InitContainer      *InitContainer               `json:"initContainer,omitempty"`
	Sidecars           *[]Sidecar                   `json:"sidecars,omitempty"`
	ServiceAccountName *string                      `json:"serviceAccountName,omitempty"`
	PersistenceEnabled *bool                        `json:"persistenceEnabled,omitempty"`
}

func (cr *RedisClusterSpec) GetReplicaCounts(t string) int32 {
	replica := cr.Size
	if t == "leader" && cr.RedisLeader.CommonAttributes.Replicas != nil {
		replica = cr.RedisLeader.CommonAttributes.Replicas
	} else if t == "follower" && cr.RedisFollower.CommonAttributes.Replicas != nil {
		replica = cr.RedisFollower.CommonAttributes.Replicas
	}
	return *replica
}

// RedisLeader interface will have the redis leader configuration
type RedisLeader struct {
	CommonAttributes              common.RedisLeader      `json:",inline"`
	SecurityContext               *corev1.SecurityContext `json:"securityContext,omitempty"`
	TerminationGracePeriodSeconds *int64                  `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,4,opt,name=terminationGracePeriodSeconds"`
}

// RedisFollower interface will have the redis follower configuration
type RedisFollower struct {
	CommonAttributes              common.RedisFollower    `json:",inline"`
	SecurityContext               *corev1.SecurityContext `json:"securityContext,omitempty"`
	TerminationGracePeriodSeconds *int64                  `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,4,opt,name=terminationGracePeriodSeconds"`
}

// RedisClusterStatus defines the observed state of RedisCluster
type RedisClusterStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="ClusterSize",type=integer,JSONPath=`.spec.clusterSize`,description=Current cluster node count
// +kubebuilder:printcolumn:name="LeaderReplicas",type=integer,JSONPath=`.spec.redisLeader.CommonAttributes.Replicas`,description=Overridden Leader replica count
// +kubebuilder:printcolumn:name="FollowerReplicas",type=integer,JSONPath=`.spec.redisFollower.CommonAttributes.Replicas`,description=Overridden Follower replica count
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description=Age of Cluster
// RedisCluster is the Schema for the redisclusters API
type RedisCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisClusterSpec   `json:"spec"`
	Status RedisClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RedisClusterList contains a list of RedisCluster
type RedisClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RedisCluster{}, &RedisClusterList{})
}
