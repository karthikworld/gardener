// +build !ignore_autogenerated

/*
Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package config

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	klog "k8s.io/klog"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupBucketControllerConfiguration) DeepCopyInto(out *BackupBucketControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupBucketControllerConfiguration.
func (in *BackupBucketControllerConfiguration) DeepCopy() *BackupBucketControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(BackupBucketControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupEntryControllerConfiguration) DeepCopyInto(out *BackupEntryControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	if in.DeletionGracePeriodHours != nil {
		in, out := &in.DeletionGracePeriodHours, &out.DeletionGracePeriodHours
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupEntryControllerConfiguration.
func (in *BackupEntryControllerConfiguration) DeepCopy() *BackupEntryControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(BackupEntryControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConditionThreshold) DeepCopyInto(out *ConditionThreshold) {
	*out = *in
	if in.Duration != nil {
		in, out := &in.Duration, &out.Duration
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConditionThreshold.
func (in *ConditionThreshold) DeepCopy() *ConditionThreshold {
	if in == nil {
		return nil
	}
	out := new(ConditionThreshold)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControllerInstallationControllerConfiguration) DeepCopyInto(out *ControllerInstallationControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerInstallationControllerConfiguration.
func (in *ControllerInstallationControllerConfiguration) DeepCopy() *ControllerInstallationControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ControllerInstallationControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiscoveryConfiguration) DeepCopyInto(out *DiscoveryConfiguration) {
	*out = *in
	if in.DiscoveryCacheDir != nil {
		in, out := &in.DiscoveryCacheDir, &out.DiscoveryCacheDir
		*out = new(string)
		**out = **in
	}
	if in.HTTPCacheDir != nil {
		in, out := &in.HTTPCacheDir, &out.HTTPCacheDir
		*out = new(string)
		**out = **in
	}
	if in.TTL != nil {
		in, out := &in.TTL, &out.TTL
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiscoveryConfiguration.
func (in *DiscoveryConfiguration) DeepCopy() *DiscoveryConfiguration {
	if in == nil {
		return nil
	}
	out := new(DiscoveryConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenClientConnection) DeepCopyInto(out *GardenClientConnection) {
	*out = *in
	out.ClientConnectionConfiguration = in.ClientConnectionConfiguration
	if in.GardenClusterAddress != nil {
		in, out := &in.GardenClusterAddress, &out.GardenClusterAddress
		*out = new(string)
		**out = **in
	}
	if in.GardenClusterCACert != nil {
		in, out := &in.GardenClusterCACert, &out.GardenClusterCACert
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.BootstrapKubeconfig != nil {
		in, out := &in.BootstrapKubeconfig, &out.BootstrapKubeconfig
		*out = new(corev1.SecretReference)
		**out = **in
	}
	if in.KubeconfigSecret != nil {
		in, out := &in.KubeconfigSecret, &out.KubeconfigSecret
		*out = new(corev1.SecretReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenClientConnection.
func (in *GardenClientConnection) DeepCopy() *GardenClientConnection {
	if in == nil {
		return nil
	}
	out := new(GardenClientConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenletConfiguration) DeepCopyInto(out *GardenletConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.GardenClientConnection != nil {
		in, out := &in.GardenClientConnection, &out.GardenClientConnection
		*out = new(GardenClientConnection)
		(*in).DeepCopyInto(*out)
	}
	if in.SeedClientConnection != nil {
		in, out := &in.SeedClientConnection, &out.SeedClientConnection
		*out = new(SeedClientConnection)
		**out = **in
	}
	if in.ShootClientConnection != nil {
		in, out := &in.ShootClientConnection, &out.ShootClientConnection
		*out = new(ShootClientConnection)
		**out = **in
	}
	if in.Controllers != nil {
		in, out := &in.Controllers, &out.Controllers
		*out = new(GardenletControllerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.LeaderElection != nil {
		in, out := &in.LeaderElection, &out.LeaderElection
		*out = new(LeaderElectionConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.Discovery != nil {
		in, out := &in.Discovery, &out.Discovery
		*out = new(DiscoveryConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.LogLevel != nil {
		in, out := &in.LogLevel, &out.LogLevel
		*out = new(string)
		**out = **in
	}
	if in.KubernetesLogLevel != nil {
		in, out := &in.KubernetesLogLevel, &out.KubernetesLogLevel
		*out = new(klog.Level)
		**out = **in
	}
	if in.FeatureGates != nil {
		in, out := &in.FeatureGates, &out.FeatureGates
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SeedConfig != nil {
		in, out := &in.SeedConfig, &out.SeedConfig
		*out = new(SeedConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.SeedSelector != nil {
		in, out := &in.SeedSelector, &out.SeedSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenletConfiguration.
func (in *GardenletConfiguration) DeepCopy() *GardenletConfiguration {
	if in == nil {
		return nil
	}
	out := new(GardenletConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GardenletConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenletControllerConfiguration) DeepCopyInto(out *GardenletControllerConfiguration) {
	*out = *in
	if in.BackupBucket != nil {
		in, out := &in.BackupBucket, &out.BackupBucket
		*out = new(BackupBucketControllerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.BackupEntry != nil {
		in, out := &in.BackupEntry, &out.BackupEntry
		*out = new(BackupEntryControllerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.ControllerInstallation != nil {
		in, out := &in.ControllerInstallation, &out.ControllerInstallation
		*out = new(ControllerInstallationControllerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.Seed != nil {
		in, out := &in.Seed, &out.Seed
		*out = new(SeedControllerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.Shoot != nil {
		in, out := &in.Shoot, &out.Shoot
		*out = new(ShootControllerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.ShootCare != nil {
		in, out := &in.ShootCare, &out.ShootCare
		*out = new(ShootCareControllerConfiguration)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenletControllerConfiguration.
func (in *GardenletControllerConfiguration) DeepCopy() *GardenletControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(GardenletControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LeaderElectionConfiguration) DeepCopyInto(out *LeaderElectionConfiguration) {
	*out = *in
	out.LeaderElectionConfiguration = in.LeaderElectionConfiguration
	if in.LockObjectNamespace != nil {
		in, out := &in.LockObjectNamespace, &out.LockObjectNamespace
		*out = new(string)
		**out = **in
	}
	if in.LockObjectName != nil {
		in, out := &in.LockObjectName, &out.LockObjectName
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeaderElectionConfiguration.
func (in *LeaderElectionConfiguration) DeepCopy() *LeaderElectionConfiguration {
	if in == nil {
		return nil
	}
	out := new(LeaderElectionConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SeedClientConnection) DeepCopyInto(out *SeedClientConnection) {
	*out = *in
	out.ClientConnectionConfiguration = in.ClientConnectionConfiguration
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SeedClientConnection.
func (in *SeedClientConnection) DeepCopy() *SeedClientConnection {
	if in == nil {
		return nil
	}
	out := new(SeedClientConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SeedConfig) DeepCopyInto(out *SeedConfig) {
	*out = *in
	in.Seed.DeepCopyInto(&out.Seed)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SeedConfig.
func (in *SeedConfig) DeepCopy() *SeedConfig {
	if in == nil {
		return nil
	}
	out := new(SeedConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SeedControllerConfiguration) DeepCopyInto(out *SeedControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	if in.ReserveExcessCapacity != nil {
		in, out := &in.ReserveExcessCapacity, &out.ReserveExcessCapacity
		*out = new(bool)
		**out = **in
	}
	if in.SyncPeriod != nil {
		in, out := &in.SyncPeriod, &out.SyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SeedControllerConfiguration.
func (in *SeedControllerConfiguration) DeepCopy() *SeedControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(SeedControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ShootCareControllerConfiguration) DeepCopyInto(out *ShootCareControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	if in.SyncPeriod != nil {
		in, out := &in.SyncPeriod, &out.SyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ConditionThresholds != nil {
		in, out := &in.ConditionThresholds, &out.ConditionThresholds
		*out = make([]ConditionThreshold, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ShootCareControllerConfiguration.
func (in *ShootCareControllerConfiguration) DeepCopy() *ShootCareControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ShootCareControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ShootClientConnection) DeepCopyInto(out *ShootClientConnection) {
	*out = *in
	out.ClientConnectionConfiguration = in.ClientConnectionConfiguration
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ShootClientConnection.
func (in *ShootClientConnection) DeepCopy() *ShootClientConnection {
	if in == nil {
		return nil
	}
	out := new(ShootClientConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ShootControllerConfiguration) DeepCopyInto(out *ShootControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	if in.ReconcileInMaintenanceOnly != nil {
		in, out := &in.ReconcileInMaintenanceOnly, &out.ReconcileInMaintenanceOnly
		*out = new(bool)
		**out = **in
	}
	if in.RespectSyncPeriodOverwrite != nil {
		in, out := &in.RespectSyncPeriodOverwrite, &out.RespectSyncPeriodOverwrite
		*out = new(bool)
		**out = **in
	}
	if in.RetryDuration != nil {
		in, out := &in.RetryDuration, &out.RetryDuration
		*out = new(v1.Duration)
		**out = **in
	}
	if in.RetrySyncPeriod != nil {
		in, out := &in.RetrySyncPeriod, &out.RetrySyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.SyncPeriod != nil {
		in, out := &in.SyncPeriod, &out.SyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ShootControllerConfiguration.
func (in *ShootControllerConfiguration) DeepCopy() *ShootControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ShootControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}
