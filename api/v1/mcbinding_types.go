/*
Copyright 2022.

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

// McbindingSpec defines the desired state of Mcbinding
type McbindingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//机器归属人
	Person string `json:"person,omitempty"`
	//需不需要直接node中获取消息  但是查询想通过这个bingding去查
	//操作系统
	Os string `json:"os,omitempty"`
	//机器地点
	Place string `json:"place,omitempty"`
	//项目信息
	Project string `json:"project,omitempty"`
	//板块
	Unit  string `json:"unit ,omitempty"`
	//部门
	Department   string `json:"department  ,omitempty"`
	//集群信息和集群类型 是否要存储  暂定不需要  global按照cluster 去查询信息时候  已经带上了cluster信息
	//global 查集体信息的时候  getList   存储类型就很重要了   拿到具体的item再来子集群 检索


	// Foo is an example field of Mcbinding. Edit mcbinding_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
}

// McbindingStatus defines the observed state of Mcbinding
type McbindingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}


//标记可以使用Status
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Mcbinding is the Schema for the mcbindings API
type Mcbinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   McbindingSpec   `json:"spec,omitempty"`
	Status McbindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// McbindingList contains a list of Mcbinding
type McbindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mcbinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mcbinding{}, &McbindingList{})
}
