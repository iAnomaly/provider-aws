/*
Copyright 2021 The Crossplane Authors.

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

// Code generated by ack-generate. DO NOT EDIT.

package domain

import (
	"context"

	svcapi "github.com/aws/aws-sdk-go/service/opensearchservice"
	svcsdk "github.com/aws/aws-sdk-go/service/opensearchservice"
	svcsdkapi "github.com/aws/aws-sdk-go/service/opensearchservice/opensearchserviceiface"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	cpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	svcapitypes "github.com/crossplane-contrib/provider-aws/apis/opensearchservice/v1alpha1"
	awsclient "github.com/crossplane-contrib/provider-aws/pkg/clients"
)

const (
	errUnexpectedObject = "managed resource is not an Domain resource"

	errCreateSession = "cannot create a new session"
	errCreate        = "cannot create Domain in AWS"
	errUpdate        = "cannot update Domain in AWS"
	errDescribe      = "failed to describe Domain"
	errDelete        = "failed to delete Domain"
)

type connector struct {
	kube client.Client
	opts []option
}

func (c *connector) Connect(ctx context.Context, mg cpresource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*svcapitypes.Domain)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	sess, err := awsclient.GetConfigV1(ctx, c.kube, mg, cr.Spec.ForProvider.Region)
	if err != nil {
		return nil, errors.Wrap(err, errCreateSession)
	}
	return newExternal(c.kube, svcapi.New(sess), c.opts), nil
}

func (e *external) Observe(ctx context.Context, mg cpresource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*svcapitypes.Domain)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}
	if meta.GetExternalName(cr) == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}
	input := GenerateDescribeDomainInput(cr)
	if err := e.preObserve(ctx, cr, input); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "pre-observe failed")
	}
	resp, err := e.client.DescribeDomainWithContext(ctx, input)
	if err != nil {
		return managed.ExternalObservation{ResourceExists: false}, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDescribe)
	}
	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err := e.lateInitialize(&cr.Spec.ForProvider, resp); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "late-init failed")
	}
	GenerateDomain(resp).Status.AtProvider.DeepCopyInto(&cr.Status.AtProvider)

	upToDate, err := e.isUpToDate(cr, resp)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "isUpToDate check failed")
	}
	return e.postObserve(ctx, cr, resp, managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        upToDate,
		ResourceLateInitialized: !cmp.Equal(&cr.Spec.ForProvider, currentSpec),
	}, nil)
}

func (e *external) Create(ctx context.Context, mg cpresource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*svcapitypes.Domain)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Creating())
	input := GenerateCreateDomainInput(cr)
	if err := e.preCreate(ctx, cr, input); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "pre-create failed")
	}
	resp, err := e.client.CreateDomainWithContext(ctx, input)
	if err != nil {
		return managed.ExternalCreation{}, awsclient.Wrap(err, errCreate)
	}

	if resp.DomainStatus.ARN != nil {
		cr.Status.AtProvider.ARN = resp.DomainStatus.ARN
	} else {
		cr.Status.AtProvider.ARN = nil
	}
	if resp.DomainStatus.AccessPolicies != nil {
		cr.Spec.ForProvider.AccessPolicies = resp.DomainStatus.AccessPolicies
	} else {
		cr.Spec.ForProvider.AccessPolicies = nil
	}
	if resp.DomainStatus.AdvancedOptions != nil {
		f2 := map[string]*string{}
		for f2key, f2valiter := range resp.DomainStatus.AdvancedOptions {
			var f2val string
			f2val = *f2valiter
			f2[f2key] = &f2val
		}
		cr.Spec.ForProvider.AdvancedOptions = f2
	} else {
		cr.Spec.ForProvider.AdvancedOptions = nil
	}
	if resp.DomainStatus.AdvancedSecurityOptions != nil {
		f3 := &svcapitypes.AdvancedSecurityOptionsInput{}
		if resp.DomainStatus.AdvancedSecurityOptions.AnonymousAuthEnabled != nil {
			f3.AnonymousAuthEnabled = resp.DomainStatus.AdvancedSecurityOptions.AnonymousAuthEnabled
		}
		if resp.DomainStatus.AdvancedSecurityOptions.Enabled != nil {
			f3.Enabled = resp.DomainStatus.AdvancedSecurityOptions.Enabled
		}
		if resp.DomainStatus.AdvancedSecurityOptions.InternalUserDatabaseEnabled != nil {
			f3.InternalUserDatabaseEnabled = resp.DomainStatus.AdvancedSecurityOptions.InternalUserDatabaseEnabled
		}
		if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions != nil {
			f3f4 := &svcapitypes.SAMLOptionsInput{}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Enabled != nil {
				f3f4.Enabled = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Enabled
			}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp != nil {
				f3f4f1 := &svcapitypes.SAMLIDp{}
				if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp.EntityId != nil {
					f3f4f1.EntityID = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp.EntityId
				}
				if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp.MetadataContent != nil {
					f3f4f1.MetadataContent = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.Idp.MetadataContent
				}
				f3f4.IDp = f3f4f1
			}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.RolesKey != nil {
				f3f4.RolesKey = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.RolesKey
			}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes != nil {
				f3f4.SessionTimeoutMinutes = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes
			}
			if resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.SubjectKey != nil {
				f3f4.SubjectKey = resp.DomainStatus.AdvancedSecurityOptions.SAMLOptions.SubjectKey
			}
			f3.SAMLOptions = f3f4
		}
		cr.Spec.ForProvider.AdvancedSecurityOptions = f3
	} else {
		cr.Spec.ForProvider.AdvancedSecurityOptions = nil
	}
	if resp.DomainStatus.AutoTuneOptions != nil {
		f4 := &svcapitypes.AutoTuneOptionsInput{}
		cr.Spec.ForProvider.AutoTuneOptions = f4
	} else {
		cr.Spec.ForProvider.AutoTuneOptions = nil
	}
	if resp.DomainStatus.ChangeProgressDetails != nil {
		f5 := &svcapitypes.ChangeProgressDetails{}
		if resp.DomainStatus.ChangeProgressDetails.ChangeId != nil {
			f5.ChangeID = resp.DomainStatus.ChangeProgressDetails.ChangeId
		}
		if resp.DomainStatus.ChangeProgressDetails.Message != nil {
			f5.Message = resp.DomainStatus.ChangeProgressDetails.Message
		}
		cr.Status.AtProvider.ChangeProgressDetails = f5
	} else {
		cr.Status.AtProvider.ChangeProgressDetails = nil
	}
	if resp.DomainStatus.ClusterConfig != nil {
		f6 := &svcapitypes.ClusterConfig{}
		if resp.DomainStatus.ClusterConfig.ColdStorageOptions != nil {
			f6f0 := &svcapitypes.ColdStorageOptions{}
			if resp.DomainStatus.ClusterConfig.ColdStorageOptions.Enabled != nil {
				f6f0.Enabled = resp.DomainStatus.ClusterConfig.ColdStorageOptions.Enabled
			}
			f6.ColdStorageOptions = f6f0
		}
		if resp.DomainStatus.ClusterConfig.DedicatedMasterCount != nil {
			f6.DedicatedMasterCount = resp.DomainStatus.ClusterConfig.DedicatedMasterCount
		}
		if resp.DomainStatus.ClusterConfig.DedicatedMasterEnabled != nil {
			f6.DedicatedMasterEnabled = resp.DomainStatus.ClusterConfig.DedicatedMasterEnabled
		}
		if resp.DomainStatus.ClusterConfig.DedicatedMasterType != nil {
			f6.DedicatedMasterType = resp.DomainStatus.ClusterConfig.DedicatedMasterType
		}
		if resp.DomainStatus.ClusterConfig.InstanceCount != nil {
			f6.InstanceCount = resp.DomainStatus.ClusterConfig.InstanceCount
		}
		if resp.DomainStatus.ClusterConfig.InstanceType != nil {
			f6.InstanceType = resp.DomainStatus.ClusterConfig.InstanceType
		}
		if resp.DomainStatus.ClusterConfig.WarmCount != nil {
			f6.WarmCount = resp.DomainStatus.ClusterConfig.WarmCount
		}
		if resp.DomainStatus.ClusterConfig.WarmEnabled != nil {
			f6.WarmEnabled = resp.DomainStatus.ClusterConfig.WarmEnabled
		}
		if resp.DomainStatus.ClusterConfig.WarmType != nil {
			f6.WarmType = resp.DomainStatus.ClusterConfig.WarmType
		}
		if resp.DomainStatus.ClusterConfig.ZoneAwarenessConfig != nil {
			f6f9 := &svcapitypes.ZoneAwarenessConfig{}
			if resp.DomainStatus.ClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount != nil {
				f6f9.AvailabilityZoneCount = resp.DomainStatus.ClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount
			}
			f6.ZoneAwarenessConfig = f6f9
		}
		if resp.DomainStatus.ClusterConfig.ZoneAwarenessEnabled != nil {
			f6.ZoneAwarenessEnabled = resp.DomainStatus.ClusterConfig.ZoneAwarenessEnabled
		}
		cr.Spec.ForProvider.ClusterConfig = f6
	} else {
		cr.Spec.ForProvider.ClusterConfig = nil
	}
	if resp.DomainStatus.CognitoOptions != nil {
		f7 := &svcapitypes.CognitoOptions{}
		if resp.DomainStatus.CognitoOptions.Enabled != nil {
			f7.Enabled = resp.DomainStatus.CognitoOptions.Enabled
		}
		if resp.DomainStatus.CognitoOptions.IdentityPoolId != nil {
			f7.IdentityPoolID = resp.DomainStatus.CognitoOptions.IdentityPoolId
		}
		if resp.DomainStatus.CognitoOptions.RoleArn != nil {
			f7.RoleARN = resp.DomainStatus.CognitoOptions.RoleArn
		}
		if resp.DomainStatus.CognitoOptions.UserPoolId != nil {
			f7.UserPoolID = resp.DomainStatus.CognitoOptions.UserPoolId
		}
		cr.Spec.ForProvider.CognitoOptions = f7
	} else {
		cr.Spec.ForProvider.CognitoOptions = nil
	}
	if resp.DomainStatus.Created != nil {
		cr.Status.AtProvider.Created = resp.DomainStatus.Created
	} else {
		cr.Status.AtProvider.Created = nil
	}
	if resp.DomainStatus.Deleted != nil {
		cr.Status.AtProvider.Deleted = resp.DomainStatus.Deleted
	} else {
		cr.Status.AtProvider.Deleted = nil
	}
	if resp.DomainStatus.DomainEndpointOptions != nil {
		f10 := &svcapitypes.DomainEndpointOptions{}
		if resp.DomainStatus.DomainEndpointOptions.CustomEndpoint != nil {
			f10.CustomEndpoint = resp.DomainStatus.DomainEndpointOptions.CustomEndpoint
		}
		if resp.DomainStatus.DomainEndpointOptions.CustomEndpointCertificateArn != nil {
			f10.CustomEndpointCertificateARN = resp.DomainStatus.DomainEndpointOptions.CustomEndpointCertificateArn
		}
		if resp.DomainStatus.DomainEndpointOptions.CustomEndpointEnabled != nil {
			f10.CustomEndpointEnabled = resp.DomainStatus.DomainEndpointOptions.CustomEndpointEnabled
		}
		if resp.DomainStatus.DomainEndpointOptions.EnforceHTTPS != nil {
			f10.EnforceHTTPS = resp.DomainStatus.DomainEndpointOptions.EnforceHTTPS
		}
		if resp.DomainStatus.DomainEndpointOptions.TLSSecurityPolicy != nil {
			f10.TLSSecurityPolicy = resp.DomainStatus.DomainEndpointOptions.TLSSecurityPolicy
		}
		cr.Spec.ForProvider.DomainEndpointOptions = f10
	} else {
		cr.Spec.ForProvider.DomainEndpointOptions = nil
	}
	if resp.DomainStatus.DomainId != nil {
		cr.Status.AtProvider.DomainID = resp.DomainStatus.DomainId
	} else {
		cr.Status.AtProvider.DomainID = nil
	}
	if resp.DomainStatus.DomainName != nil {
		cr.Status.AtProvider.DomainName = resp.DomainStatus.DomainName
	} else {
		cr.Status.AtProvider.DomainName = nil
	}
	if resp.DomainStatus.EBSOptions != nil {
		f13 := &svcapitypes.EBSOptions{}
		if resp.DomainStatus.EBSOptions.EBSEnabled != nil {
			f13.EBSEnabled = resp.DomainStatus.EBSOptions.EBSEnabled
		}
		if resp.DomainStatus.EBSOptions.Iops != nil {
			f13.IOPS = resp.DomainStatus.EBSOptions.Iops
		}
		if resp.DomainStatus.EBSOptions.Throughput != nil {
			f13.Throughput = resp.DomainStatus.EBSOptions.Throughput
		}
		if resp.DomainStatus.EBSOptions.VolumeSize != nil {
			f13.VolumeSize = resp.DomainStatus.EBSOptions.VolumeSize
		}
		if resp.DomainStatus.EBSOptions.VolumeType != nil {
			f13.VolumeType = resp.DomainStatus.EBSOptions.VolumeType
		}
		cr.Spec.ForProvider.EBSOptions = f13
	} else {
		cr.Spec.ForProvider.EBSOptions = nil
	}
	if resp.DomainStatus.EncryptionAtRestOptions != nil {
		f14 := &svcapitypes.EncryptionAtRestOptions{}
		if resp.DomainStatus.EncryptionAtRestOptions.Enabled != nil {
			f14.Enabled = resp.DomainStatus.EncryptionAtRestOptions.Enabled
		}
		if resp.DomainStatus.EncryptionAtRestOptions.KmsKeyId != nil {
			f14.KMSKeyID = resp.DomainStatus.EncryptionAtRestOptions.KmsKeyId
		}
		cr.Status.AtProvider.EncryptionAtRestOptions = f14
	} else {
		cr.Status.AtProvider.EncryptionAtRestOptions = nil
	}
	if resp.DomainStatus.Endpoint != nil {
		cr.Status.AtProvider.Endpoint = resp.DomainStatus.Endpoint
	} else {
		cr.Status.AtProvider.Endpoint = nil
	}
	if resp.DomainStatus.Endpoints != nil {
		f16 := map[string]*string{}
		for f16key, f16valiter := range resp.DomainStatus.Endpoints {
			var f16val string
			f16val = *f16valiter
			f16[f16key] = &f16val
		}
		cr.Status.AtProvider.Endpoints = f16
	} else {
		cr.Status.AtProvider.Endpoints = nil
	}
	if resp.DomainStatus.EngineVersion != nil {
		cr.Spec.ForProvider.EngineVersion = resp.DomainStatus.EngineVersion
	} else {
		cr.Spec.ForProvider.EngineVersion = nil
	}
	if resp.DomainStatus.LogPublishingOptions != nil {
		f18 := map[string]*svcapitypes.LogPublishingOption{}
		for f18key, f18valiter := range resp.DomainStatus.LogPublishingOptions {
			f18val := &svcapitypes.LogPublishingOption{}
			if f18valiter.CloudWatchLogsLogGroupArn != nil {
				f18val.CloudWatchLogsLogGroupARN = f18valiter.CloudWatchLogsLogGroupArn
			}
			if f18valiter.Enabled != nil {
				f18val.Enabled = f18valiter.Enabled
			}
			f18[f18key] = f18val
		}
		cr.Spec.ForProvider.LogPublishingOptions = f18
	} else {
		cr.Spec.ForProvider.LogPublishingOptions = nil
	}
	if resp.DomainStatus.NodeToNodeEncryptionOptions != nil {
		f19 := &svcapitypes.NodeToNodeEncryptionOptions{}
		if resp.DomainStatus.NodeToNodeEncryptionOptions.Enabled != nil {
			f19.Enabled = resp.DomainStatus.NodeToNodeEncryptionOptions.Enabled
		}
		cr.Spec.ForProvider.NodeToNodeEncryptionOptions = f19
	} else {
		cr.Spec.ForProvider.NodeToNodeEncryptionOptions = nil
	}
	if resp.DomainStatus.Processing != nil {
		cr.Status.AtProvider.Processing = resp.DomainStatus.Processing
	} else {
		cr.Status.AtProvider.Processing = nil
	}
	if resp.DomainStatus.ServiceSoftwareOptions != nil {
		f21 := &svcapitypes.ServiceSoftwareOptions{}
		if resp.DomainStatus.ServiceSoftwareOptions.AutomatedUpdateDate != nil {
			f21.AutomatedUpdateDate = &metav1.Time{*resp.DomainStatus.ServiceSoftwareOptions.AutomatedUpdateDate}
		}
		if resp.DomainStatus.ServiceSoftwareOptions.Cancellable != nil {
			f21.Cancellable = resp.DomainStatus.ServiceSoftwareOptions.Cancellable
		}
		if resp.DomainStatus.ServiceSoftwareOptions.CurrentVersion != nil {
			f21.CurrentVersion = resp.DomainStatus.ServiceSoftwareOptions.CurrentVersion
		}
		if resp.DomainStatus.ServiceSoftwareOptions.Description != nil {
			f21.Description = resp.DomainStatus.ServiceSoftwareOptions.Description
		}
		if resp.DomainStatus.ServiceSoftwareOptions.NewVersion != nil {
			f21.NewVersion = resp.DomainStatus.ServiceSoftwareOptions.NewVersion
		}
		if resp.DomainStatus.ServiceSoftwareOptions.OptionalDeployment != nil {
			f21.OptionalDeployment = resp.DomainStatus.ServiceSoftwareOptions.OptionalDeployment
		}
		if resp.DomainStatus.ServiceSoftwareOptions.UpdateAvailable != nil {
			f21.UpdateAvailable = resp.DomainStatus.ServiceSoftwareOptions.UpdateAvailable
		}
		if resp.DomainStatus.ServiceSoftwareOptions.UpdateStatus != nil {
			f21.UpdateStatus = resp.DomainStatus.ServiceSoftwareOptions.UpdateStatus
		}
		cr.Status.AtProvider.ServiceSoftwareOptions = f21
	} else {
		cr.Status.AtProvider.ServiceSoftwareOptions = nil
	}
	if resp.DomainStatus.SnapshotOptions != nil {
		f22 := &svcapitypes.SnapshotOptions{}
		if resp.DomainStatus.SnapshotOptions.AutomatedSnapshotStartHour != nil {
			f22.AutomatedSnapshotStartHour = resp.DomainStatus.SnapshotOptions.AutomatedSnapshotStartHour
		}
		cr.Status.AtProvider.SnapshotOptions = f22
	} else {
		cr.Status.AtProvider.SnapshotOptions = nil
	}
	if resp.DomainStatus.UpgradeProcessing != nil {
		cr.Status.AtProvider.UpgradeProcessing = resp.DomainStatus.UpgradeProcessing
	} else {
		cr.Status.AtProvider.UpgradeProcessing = nil
	}
	if resp.DomainStatus.VPCOptions != nil {
		f24 := &svcapitypes.VPCDerivedInfo{}
		if resp.DomainStatus.VPCOptions.AvailabilityZones != nil {
			f24f0 := []*string{}
			for _, f24f0iter := range resp.DomainStatus.VPCOptions.AvailabilityZones {
				var f24f0elem string
				f24f0elem = *f24f0iter
				f24f0 = append(f24f0, &f24f0elem)
			}
			f24.AvailabilityZones = f24f0
		}
		if resp.DomainStatus.VPCOptions.SecurityGroupIds != nil {
			f24f1 := []*string{}
			for _, f24f1iter := range resp.DomainStatus.VPCOptions.SecurityGroupIds {
				var f24f1elem string
				f24f1elem = *f24f1iter
				f24f1 = append(f24f1, &f24f1elem)
			}
			f24.SecurityGroupIDs = f24f1
		}
		if resp.DomainStatus.VPCOptions.SubnetIds != nil {
			f24f2 := []*string{}
			for _, f24f2iter := range resp.DomainStatus.VPCOptions.SubnetIds {
				var f24f2elem string
				f24f2elem = *f24f2iter
				f24f2 = append(f24f2, &f24f2elem)
			}
			f24.SubnetIDs = f24f2
		}
		if resp.DomainStatus.VPCOptions.VPCId != nil {
			f24.VPCID = resp.DomainStatus.VPCOptions.VPCId
		}
		cr.Status.AtProvider.VPCOptions = f24
	} else {
		cr.Status.AtProvider.VPCOptions = nil
	}

	return e.postCreate(ctx, cr, resp, managed.ExternalCreation{}, err)
}

func (e *external) Update(ctx context.Context, mg cpresource.Managed) (managed.ExternalUpdate, error) {
	return e.update(ctx, mg)

}

func (e *external) Delete(ctx context.Context, mg cpresource.Managed) error {
	cr, ok := mg.(*svcapitypes.Domain)
	if !ok {
		return errors.New(errUnexpectedObject)
	}
	cr.Status.SetConditions(xpv1.Deleting())
	input := GenerateDeleteDomainInput(cr)
	ignore, err := e.preDelete(ctx, cr, input)
	if err != nil {
		return errors.Wrap(err, "pre-delete failed")
	}
	if ignore {
		return nil
	}
	resp, err := e.client.DeleteDomainWithContext(ctx, input)
	return e.postDelete(ctx, cr, resp, awsclient.Wrap(cpresource.Ignore(IsNotFound, err), errDelete))
}

type option func(*external)

func newExternal(kube client.Client, client svcsdkapi.OpenSearchServiceAPI, opts []option) *external {
	e := &external{
		kube:           kube,
		client:         client,
		preObserve:     nopPreObserve,
		postObserve:    nopPostObserve,
		lateInitialize: nopLateInitialize,
		isUpToDate:     alwaysUpToDate,
		preCreate:      nopPreCreate,
		postCreate:     nopPostCreate,
		preDelete:      nopPreDelete,
		postDelete:     nopPostDelete,
		update:         nopUpdate,
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

type external struct {
	kube           client.Client
	client         svcsdkapi.OpenSearchServiceAPI
	preObserve     func(context.Context, *svcapitypes.Domain, *svcsdk.DescribeDomainInput) error
	postObserve    func(context.Context, *svcapitypes.Domain, *svcsdk.DescribeDomainOutput, managed.ExternalObservation, error) (managed.ExternalObservation, error)
	lateInitialize func(*svcapitypes.DomainParameters, *svcsdk.DescribeDomainOutput) error
	isUpToDate     func(*svcapitypes.Domain, *svcsdk.DescribeDomainOutput) (bool, error)
	preCreate      func(context.Context, *svcapitypes.Domain, *svcsdk.CreateDomainInput) error
	postCreate     func(context.Context, *svcapitypes.Domain, *svcsdk.CreateDomainOutput, managed.ExternalCreation, error) (managed.ExternalCreation, error)
	preDelete      func(context.Context, *svcapitypes.Domain, *svcsdk.DeleteDomainInput) (bool, error)
	postDelete     func(context.Context, *svcapitypes.Domain, *svcsdk.DeleteDomainOutput, error) error
	update         func(context.Context, cpresource.Managed) (managed.ExternalUpdate, error)
}

func nopPreObserve(context.Context, *svcapitypes.Domain, *svcsdk.DescribeDomainInput) error {
	return nil
}

func nopPostObserve(_ context.Context, _ *svcapitypes.Domain, _ *svcsdk.DescribeDomainOutput, obs managed.ExternalObservation, err error) (managed.ExternalObservation, error) {
	return obs, err
}
func nopLateInitialize(*svcapitypes.DomainParameters, *svcsdk.DescribeDomainOutput) error {
	return nil
}
func alwaysUpToDate(*svcapitypes.Domain, *svcsdk.DescribeDomainOutput) (bool, error) {
	return true, nil
}

func nopPreCreate(context.Context, *svcapitypes.Domain, *svcsdk.CreateDomainInput) error {
	return nil
}
func nopPostCreate(_ context.Context, _ *svcapitypes.Domain, _ *svcsdk.CreateDomainOutput, cre managed.ExternalCreation, err error) (managed.ExternalCreation, error) {
	return cre, err
}
func nopPreDelete(context.Context, *svcapitypes.Domain, *svcsdk.DeleteDomainInput) (bool, error) {
	return false, nil
}
func nopPostDelete(_ context.Context, _ *svcapitypes.Domain, _ *svcsdk.DeleteDomainOutput, err error) error {
	return err
}
func nopUpdate(context.Context, cpresource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}
