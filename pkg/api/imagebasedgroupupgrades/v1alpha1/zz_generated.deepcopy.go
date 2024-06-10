//go:build !ignore_autogenerated

/*
Copyright 2021.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterState) DeepCopyInto(out *ClusterState) {
	*out = *in
	if in.CurrentAction != nil {
		in, out := &in.CurrentAction, &out.CurrentAction
		*out = new(string)
		**out = **in
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterState.
func (in *ClusterState) DeepCopy() *ClusterState {
	if in == nil {
		return nil
	}
	out := new(ClusterState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBasedGroupUpgrade) DeepCopyInto(out *ImageBasedGroupUpgrade) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBasedGroupUpgrade.
func (in *ImageBasedGroupUpgrade) DeepCopy() *ImageBasedGroupUpgrade {
	if in == nil {
		return nil
	}
	out := new(ImageBasedGroupUpgrade)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ImageBasedGroupUpgrade) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBasedGroupUpgradeList) DeepCopyInto(out *ImageBasedGroupUpgradeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ImageBasedGroupUpgrade, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBasedGroupUpgradeList.
func (in *ImageBasedGroupUpgradeList) DeepCopy() *ImageBasedGroupUpgradeList {
	if in == nil {
		return nil
	}
	out := new(ImageBasedGroupUpgradeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ImageBasedGroupUpgradeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBasedGroupUpgradeSpec) DeepCopyInto(out *ImageBasedGroupUpgradeSpec) {
	*out = *in
	in.IBUSpec.DeepCopyInto(&out.IBUSpec)
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ClusterLabelSelectors != nil {
		in, out := &in.ClusterLabelSelectors, &out.ClusterLabelSelectors
		*out = make([]v1.LabelSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Actions != nil {
		in, out := &in.Actions, &out.Actions
		*out = make([]ImageBasedUpgradeAction, len(*in))
		copy(*out, *in)
	}
	out.RolloutStrategy = in.RolloutStrategy
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBasedGroupUpgradeSpec.
func (in *ImageBasedGroupUpgradeSpec) DeepCopy() *ImageBasedGroupUpgradeSpec {
	if in == nil {
		return nil
	}
	out := new(ImageBasedGroupUpgradeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBasedGroupUpgradeStatus) DeepCopyInto(out *ImageBasedGroupUpgradeStatus) {
	*out = *in
	in.StartedAt.DeepCopyInto(&out.StartedAt)
	in.CompletedAt.DeepCopyInto(&out.CompletedAt)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make([]ClusterState, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBasedGroupUpgradeStatus.
func (in *ImageBasedGroupUpgradeStatus) DeepCopy() *ImageBasedGroupUpgradeStatus {
	if in == nil {
		return nil
	}
	out := new(ImageBasedGroupUpgradeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageBasedUpgradeAction) DeepCopyInto(out *ImageBasedUpgradeAction) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageBasedUpgradeAction.
func (in *ImageBasedUpgradeAction) DeepCopy() *ImageBasedUpgradeAction {
	if in == nil {
		return nil
	}
	out := new(ImageBasedUpgradeAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RolloutStrategy) DeepCopyInto(out *RolloutStrategy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RolloutStrategy.
func (in *RolloutStrategy) DeepCopy() *RolloutStrategy {
	if in == nil {
		return nil
	}
	out := new(RolloutStrategy)
	in.DeepCopyInto(out)
	return out
}
