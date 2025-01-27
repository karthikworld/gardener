// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllerinstallation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	v1alpha1constants "github.com/gardener/gardener/pkg/apis/core/v1alpha1/constants"
	"github.com/gardener/gardener/pkg/apis/core/v1alpha1/helper"
	gardencorev1alpha1helper "github.com/gardener/gardener/pkg/apis/core/v1alpha1/helper"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/chartrenderer"
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/externalversions"
	gardencorelisters "github.com/gardener/gardener/pkg/client/core/listers/core/v1alpha1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/gardenlet/apis/config"
	"github.com/gardener/gardener/pkg/logger"
	seedpkg "github.com/gardener/gardener/pkg/operation/seed"
	"github.com/gardener/gardener/pkg/utils"
	"github.com/gardener/gardener/pkg/utils/flow"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"

	resourcesv1alpha1 "github.com/gardener/gardener-resource-manager/pkg/apis/resources/v1alpha1"
	"github.com/gardener/gardener-resource-manager/pkg/manager"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const installationTypeHelm = "helm"

func (c *Controller) controllerInstallationAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		logger.Logger.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	c.controllerInstallationQueue.Add(key)
}

func (c *Controller) controllerInstallationUpdate(oldObj, newObj interface{}) {
	old, ok1 := oldObj.(*gardencorev1alpha1.ControllerInstallation)
	new, ok2 := newObj.(*gardencorev1alpha1.ControllerInstallation)
	if !ok1 || !ok2 {
		return
	}

	if new.DeletionTimestamp == nil && old.Spec.RegistrationRef.ResourceVersion == new.Spec.RegistrationRef.ResourceVersion && old.Spec.SeedRef.ResourceVersion == new.Spec.SeedRef.ResourceVersion {
		return
	}

	c.controllerInstallationAdd(newObj)
}

func (c *Controller) controllerInstallationDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		logger.Logger.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	c.controllerInstallationQueue.Add(key)
}

func (c *Controller) reconcileControllerInstallationKey(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	controllerInstallation, err := c.controllerInstallationLister.Get(name)
	if apierrors.IsNotFound(err) {
		logger.Logger.Debugf("[CONTROLLERINSTALLATION RECONCILE] %s - skipping because ControllerInstallation has been deleted", key)
		return nil
	}
	if err != nil {
		logger.Logger.Infof("[CONTROLLERINSTALLATION RECONCILE] %s - unable to retrieve object from store: %v", key, err)
		return err
	}

	return c.controllerInstallationControl.Reconcile(controllerInstallation)
}

// ControlInterface implements the control logic for updating ControllerInstallations. It is implemented as an interface to allow
// for extensions that provide different semantics. Currently, there is only one implementation.
type ControlInterface interface {
	Reconcile(*gardencorev1alpha1.ControllerInstallation) error
}

// NewDefaultControllerInstallationControl returns a new instance of the default implementation ControlInterface that
// implements the documented semantics for ControllerInstallations. You should use an instance returned from
// NewDefaultControllerInstallationControl() for any scenario other than testing.
func NewDefaultControllerInstallationControl(k8sGardenClient kubernetes.Interface, k8sGardenCoreInformers gardencoreinformers.SharedInformerFactory, recorder record.EventRecorder, config *config.GardenletConfiguration, seedLister gardencorelisters.SeedLister, controllerRegistrationLister gardencorelisters.ControllerRegistrationLister, controllerInstallationLister gardencorelisters.ControllerInstallationLister, gardenNamespace *corev1.Namespace) ControlInterface {
	return &defaultControllerInstallationControl{k8sGardenClient, k8sGardenCoreInformers, recorder, config, seedLister, controllerRegistrationLister, controllerInstallationLister, gardenNamespace}
}

type defaultControllerInstallationControl struct {
	k8sGardenClient              kubernetes.Interface
	k8sGardenCoreInformers       gardencoreinformers.SharedInformerFactory
	recorder                     record.EventRecorder
	config                       *config.GardenletConfiguration
	seedLister                   gardencorelisters.SeedLister
	controllerRegistrationLister gardencorelisters.ControllerRegistrationLister
	controllerInstallationLister gardencorelisters.ControllerInstallationLister
	gardenNamespace              *corev1.Namespace
}

func (c *defaultControllerInstallationControl) Reconcile(obj *gardencorev1alpha1.ControllerInstallation) error {
	var (
		controllerInstallation = obj.DeepCopy()
		logger                 = logger.NewFieldLogger(logger.Logger, "controllerinstallation", controllerInstallation.Name)
	)

	if isResponsible, err := c.isResponsible(controllerInstallation); !isResponsible || err != nil {
		return err
	}

	if controllerInstallation.DeletionTimestamp != nil {
		return c.delete(controllerInstallation, logger)
	}
	return c.reconcile(controllerInstallation, logger)
}

func (c *defaultControllerInstallationControl) reconcile(controllerInstallation *gardencorev1alpha1.ControllerInstallation, logger logrus.FieldLogger) error {
	ctx := context.TODO()

	controllerInstallation, err := kutil.TryUpdateControllerInstallationWithEqualFunc(c.k8sGardenClient.GardenCore(), retry.DefaultBackoff, controllerInstallation.ObjectMeta, func(c *gardencorev1alpha1.ControllerInstallation) (*gardencorev1alpha1.ControllerInstallation, error) {
		if finalizers := sets.NewString(c.Finalizers...); !finalizers.Has(FinalizerName) {
			finalizers.Insert(FinalizerName)
			c.Finalizers = finalizers.UnsortedList()
		}
		return c, nil
	}, func(cur, updated *gardencorev1alpha1.ControllerInstallation) bool {
		return sets.NewString(cur.Finalizers...).Has(FinalizerName)
	})
	if err != nil {
		return err
	}

	var (
		newConditions      = helper.MergeConditions(controllerInstallation.Status.Conditions, helper.InitCondition(gardencorev1alpha1.ControllerInstallationValid), helper.InitCondition(gardencorev1alpha1.ControllerInstallationInstalled))
		conditionValid     = newConditions[0]
		conditionInstalled = newConditions[1]
	)

	defer func() {
		if _, err := c.updateConditions(controllerInstallation, conditionValid, conditionInstalled); err != nil {
			logger.Errorf("Failed to update the conditions : %+v", err)
		}
	}()

	controllerRegistration, err := c.controllerRegistrationLister.Get(controllerInstallation.Spec.RegistrationRef.Name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionFalse, "RegistrationNotFound", fmt.Sprintf("Referenced ControllerRegistration does not exist: %+v", err))
		} else {
			conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionUnknown, "RegistrationReadError", fmt.Sprintf("Referenced ControllerRegistration cannot be read: %+v", err))
		}
		return err
	}

	seed, err := c.seedLister.Get(controllerInstallation.Spec.SeedRef.Name)
	if err != nil {
		return err
	}

	k8sSeedClient, err := seedpkg.GetSeedClient(ctx, c.k8sGardenClient.Client(), c.config.SeedClientConnection.ClientConnectionConfiguration, c.config.SeedSelector == nil, seed.Name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionFalse, "SeedNotFound", fmt.Sprintf("Referenced Seed does not exist: %+v", err))
		} else {
			conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionUnknown, "SeedReadError", fmt.Sprintf("Referenced Seed cannot be read: %+v", err))
		}
		return err
	}
	chartRenderer, err := chartrenderer.NewForConfig(k8sSeedClient.RESTConfig())
	if err != nil {
		conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionUnknown, "ChartRendererCreationFailed", fmt.Sprintf("ChartRenderer cannot be recreated for referenced Seed: %+v", err))
		return err
	}

	var helmDeployment HelmDeployment
	if err := json.Unmarshal(controllerRegistration.Spec.Deployment.ProviderConfig.Raw, &helmDeployment); err != nil {
		conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionFalse, "ChartInformationInvalid", fmt.Sprintf("Chart Information cannot be unmarshalled: %+v", err))
		return err
	}

	namespace := getNamespaceForControllerInstallation(controllerInstallation)
	if err := kutil.CreateOrUpdate(ctx, k8sSeedClient.Client(), namespace, func() error {
		kutil.SetMetaDataLabel(&namespace.ObjectMeta, v1alpha1constants.GardenRole, v1alpha1constants.GardenRoleExtension)
		kutil.SetMetaDataLabel(&namespace.ObjectMeta, v1alpha1constants.LabelControllerRegistrationName, controllerRegistration.Name)
		return nil
	}); err != nil {
		return err
	}

	var (
		volumeProvider  string
		volumeProviders []gardencorev1alpha1.SeedVolumeProvider
	)
	if seed.Spec.Volume != nil {
		volumeProviders = seed.Spec.Volume.Providers
		if len(seed.Spec.Volume.Providers) > 0 {
			volumeProvider = seed.Spec.Volume.Providers[0].Name
		}
	}

	// Mix-in some standard values for seed.
	seedValues := map[string]interface{}{
		"gardener": map[string]interface{}{
			"garden": map[string]interface{}{
				"identity": c.gardenNamespace.UID,
			},
			"seed": map[string]interface{}{
				"identity":        seed.Name,
				"provider":        seed.Spec.Provider.Type,
				"volumeProvider":  volumeProvider,
				"volumeProviders": volumeProviders,
				"region":          seed.Spec.Provider.Region,
				"ingressDomain":   seed.Spec.DNS.IngressDomain,
				"blockCIDRs":      seed.Spec.BlockCIDRs,
				"protected":       gardencorev1alpha1helper.TaintsHave(seed.Spec.Taints, gardencorev1alpha1.SeedTaintProtected),
				"visible":         !gardencorev1alpha1helper.TaintsHave(seed.Spec.Taints, gardencorev1alpha1.SeedTaintInvisible),
				"taints":          seed.Spec.Taints,
				"networks":        seed.Spec.Networks,
			},
		},
	}

	release, err := chartRenderer.RenderArchive(helmDeployment.Chart, controllerRegistration.Name, namespace.Name, utils.MergeMaps(helmDeployment.Values, seedValues))
	if err != nil {
		conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionFalse, "ChartCannotBeRendered", fmt.Sprintf("Chart rendering process failed: %+v", err))
		return err
	}
	conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionTrue, "RegistrationValid", "Chart could be rendered successfully.")

	// Create secret
	data := make(map[string][]byte, len(release.Files()))
	for fileName, fileContent := range release.Files() {
		key := strings.Replace(fileName, "/", "_", -1)
		data[key] = []byte(fileContent)
	}

	var secretName = controllerInstallation.Name
	if err := manager.
		NewSecret(k8sSeedClient.Client()).
		WithNamespacedName(v1alpha1constants.GardenNamespace, secretName).
		WithKeyValues(data).
		Reconcile(ctx); err != nil {
		conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "InstallationFailed", fmt.Sprintf("Creation of ManagedResource secret %q failed: %+v", secretName, err))
		return err
	}

	if err := manager.
		NewManagedResource(k8sSeedClient.Client()).
		WithNamespacedName(v1alpha1constants.GardenNamespace, controllerInstallation.Name).
		WithSecretRef(secretName).
		WithClass(v1alpha1constants.SeedResourceManagerClass).
		Reconcile(ctx); err != nil {
		conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "InstallationFailed", fmt.Sprintf("Creation of ManagedResource %q failed: %+v", controllerInstallation.Name, err))
		return err
	}

	// TODO: Remove after several releases.
	controllerInstallation, err = kutil.TryUpdateControllerInstallationStatusWithEqualFunc(c.k8sGardenClient.GardenCore(), retry.DefaultBackoff, controllerInstallation.ObjectMeta,
		func(controllerInstallation *gardencorev1alpha1.ControllerInstallation) (*gardencorev1alpha1.ControllerInstallation, error) {
			var resources DeployedResources
			status, err := json.Marshal(resources)
			if err != nil {
				return nil, fmt.Errorf("could not marshal status for resources: %+v", err)
			}

			controllerInstallation.Status.ProviderStatus = &gardencorev1alpha1.ProviderConfig{
				RawExtension: runtime.RawExtension{
					Raw: status,
				},
			}
			return controllerInstallation, nil
		}, func(cur, updated *gardencorev1alpha1.ControllerInstallation) bool {
			return equality.Semantic.DeepEqual(cur.Status.ProviderStatus, updated.Status.ProviderStatus)
		},
	)
	if err != nil {
		conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "InstallationFailed", fmt.Sprintf("Could not write status for new resources: %+v", err))
		return err
	}

	conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionTrue, "InstallationSuccessful", "Installation of new resources succeeded.")
	return nil
}

func (c *defaultControllerInstallationControl) delete(controllerInstallation *gardencorev1alpha1.ControllerInstallation, logger logrus.FieldLogger) error {
	var (
		ctx                = context.TODO()
		newConditions      = helper.MergeConditions(controllerInstallation.Status.Conditions, helper.InitCondition(gardencorev1alpha1.ControllerInstallationValid), helper.InitCondition(gardencorev1alpha1.ControllerInstallationInstalled))
		conditionValid     = newConditions[0]
		conditionInstalled = newConditions[1]
	)

	defer func() {
		if _, err := c.updateConditions(controllerInstallation, conditionValid, conditionInstalled); client.IgnoreNotFound(err) != nil {
			logger.Errorf("Failed to update the conditions when trying to delete: %+v", err)
		}
	}()

	seed, err := c.seedLister.Get(controllerInstallation.Spec.SeedRef.Name)
	if err != nil {
		return err
	}

	k8sSeedClient, err := seedpkg.GetSeedClient(ctx, c.k8sGardenClient.Client(), c.config.SeedClientConnection.ClientConnectionConfiguration, c.config.SeedSelector == nil, seed.Name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionFalse, "SeedNotFound", fmt.Sprintf("Referenced Seed does not exist: %+v", err))
		} else {
			conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionUnknown, "SeedReadError", fmt.Sprintf("Referenced Seed cannot be read: %+v", err))
		}
		return err
	}

	controllerRegistration, err := c.controllerRegistrationLister.Get(controllerInstallation.Spec.RegistrationRef.Name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionFalse, "RegistrationNotFound", fmt.Sprintf("Referenced ControllerRegistration does not exist: %+v", err))
		} else {
			conditionValid = helper.UpdatedCondition(conditionValid, gardencorev1alpha1.ConditionUnknown, "RegistrationReadError", fmt.Sprintf("Referenced ControllerRegistration cannot be read: %+v", err))
		}
		return err
	}

	if err := c.cleanOldExtensions(ctx, k8sSeedClient.Client(), controllerRegistration); err != nil {
		if isDeletionInProgressError(err) {
			conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "DeletionPending", err.Error())
		} else {
			conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "DeletionFailed", fmt.Sprintf("Deletion of extension kinds failed: %+v", err))
		}
		return err
	}

	mr := &resourcesv1alpha1.ManagedResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      controllerInstallation.Name,
			Namespace: v1alpha1constants.GardenNamespace,
		},
	}
	err = k8sSeedClient.Client().Delete(ctx, mr)
	if err == nil {
		message := fmt.Sprintf("Deletion of ManagedResource %q is still pending.", controllerInstallation.Name)
		conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "DeletionPending", message)
		return errors.New(message)
	} else if !apierrors.IsNotFound(err) {
		conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "DeletionFailed", fmt.Sprintf("Deletion of ManagedResource %q failed: %+v", controllerInstallation.Name, err))
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      controllerInstallation.Name,
			Namespace: v1alpha1constants.GardenNamespace,
		},
	}
	if err := k8sSeedClient.Client().Delete(ctx, secret); client.IgnoreNotFound(err) != nil {
		conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "DeletionFailed", fmt.Sprintf("Deletion of ManagedResource secret %q failed: %+v", controllerInstallation.Name, err))
	}

	// TODO: Kept for backwards compatibility. Remove after several releases.
	if controllerInstallation.Status.ProviderStatus != nil {
		if err := c.cleanOldResources(k8sSeedClient, controllerInstallation, sets.NewString()); err != nil {
			if isDeletionInProgressError(err) {
				conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "DeletionPending", err.Error())
			} else {
				conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "DeletionFailed", fmt.Sprintf("Deletion of old resources failed: %+v", err))
			}
			return err
		}
	}

	if err := k8sSeedClient.Client().Delete(ctx, getNamespaceForControllerInstallation(controllerInstallation)); err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	conditionInstalled = helper.UpdatedCondition(conditionInstalled, gardencorev1alpha1.ConditionFalse, "DeletionSuccessful", "Deletion of old resources succeeded.")

	_, err = kutil.TryUpdateControllerInstallationWithEqualFunc(c.k8sGardenClient.GardenCore(), retry.DefaultBackoff, controllerInstallation.ObjectMeta, func(c *gardencorev1alpha1.ControllerInstallation) (*gardencorev1alpha1.ControllerInstallation, error) {
		finalizers := sets.NewString(c.Finalizers...)
		finalizers.Delete(FinalizerName)
		c.Finalizers = finalizers.UnsortedList()
		return c, nil
	}, func(cur, updated *gardencorev1alpha1.ControllerInstallation) bool {
		return !sets.NewString(cur.Finalizers...).Has(FinalizerName)
	})
	return err
}

func (c *defaultControllerInstallationControl) updateConditions(controllerInstallation *gardencorev1alpha1.ControllerInstallation, conditions ...gardencorev1alpha1.Condition) (*gardencorev1alpha1.ControllerInstallation, error) {
	return kutil.TryUpdateControllerInstallationStatusWithEqualFunc(c.k8sGardenClient.GardenCore(), retry.DefaultBackoff, controllerInstallation.ObjectMeta,
		func(controllerInstallation *gardencorev1alpha1.ControllerInstallation) (*gardencorev1alpha1.ControllerInstallation, error) {
			controllerInstallation.Status.Conditions = conditions
			return controllerInstallation, nil
		}, func(cur, updated *gardencorev1alpha1.ControllerInstallation) bool {
			return equality.Semantic.DeepEqual(cur.Status.Conditions, updated.Status.Conditions)
		},
	)
}

func (c *defaultControllerInstallationControl) isResponsible(controllerInstallation *gardencorev1alpha1.ControllerInstallation) (bool, error) {
	controllerRegistration, err := c.controllerRegistrationLister.Get(controllerInstallation.Spec.RegistrationRef.Name)
	if err != nil {
		return false, err
	}

	if deployment := controllerRegistration.Spec.Deployment; deployment != nil {
		return deployment.Type == installationTypeHelm, nil
	}
	return false, nil
}

func (c *defaultControllerInstallationControl) cleanOldExtensions(ctx context.Context, seedClient client.Client, controllerRegistration *gardencorev1alpha1.ControllerRegistration) error {
	var (
		fns               []flow.TaskFn
		relevantExtension []extensionsv1alpha1.Extension
		result            error
	)

	extensionList := &extensionsv1alpha1.ExtensionList{}
	if err := seedClient.List(ctx, extensionList); err != nil {
		return err
	}

	for _, res := range controllerRegistration.Spec.Resources {
		if res.Kind != extensionsv1alpha1.ExtensionResource {
			continue
		}

		for _, item := range extensionList.Items {
			if res.Type != item.Spec.Type {
				continue
			}

			relevantExtension = append(relevantExtension, item)
			del := &extensionsv1alpha1.Extension{
				ObjectMeta: metav1.ObjectMeta{
					Name:      item.GetName(),
					Namespace: item.GetNamespace(),
				},
			}
			fns = append(fns, func(ctx context.Context) error {
				return seedClient.Delete(ctx, del)
			})
		}
	}

	if errs := flow.Parallel(fns...)(ctx); errs != nil {
		multiErrs, ok := errs.(*multierror.Error)
		if !ok {
			return errs
		}

		for _, err := range multiErrs.WrappedErrors() {
			if !apierrors.IsNotFound(err) {
				result = multierror.Append(result, err)
			}
		}
	}

	if result != nil {
		return result
	}

	if len(relevantExtension) != 0 {
		return newDeletionInProgressError("deletion of extensions is still pending")
	}

	return nil
}

// cleanOldResources cleans ControllerInstallation resources which are not found in the given <newResourcesSet>.
//
// TODO: Kept for backwards compatibility. Remove after several releases.
func (c *defaultControllerInstallationControl) cleanOldResources(k8sSeedClient kubernetes.Interface, controllerInstallation *gardencorev1alpha1.ControllerInstallation, newResourcesSet sets.String) error {
	providerStatus := controllerInstallation.Status.ProviderStatus
	if providerStatus == nil {
		return nil
	}

	var oldResources DeployedResources
	if err := json.Unmarshal(providerStatus.Raw, &oldResources); err != nil {
		return err
	}

	var (
		deleted = true
		result  error
	)

	for _, oldResource := range oldResources.Resources {
		if !newResourcesSet.Has(objectReferenceToString(oldResource)) {
			reader := kubernetes.NewObjectReferenceReader(&oldResource)
			obj, err := reader.Read()
			if err != nil {
				result = multierror.Append(result, err)
				continue
			}

			if err := k8sSeedClient.Client().Delete(context.TODO(), obj); err != nil {
				if !apierrors.IsNotFound(err) {
					result = multierror.Append(result, err)
				}
				continue
			}
			deleted = false
		}
	}

	if result != nil {
		return result
	}

	if !deleted {
		return newDeletionInProgressError("deletion of old resources is still pending")
	}

	return nil
}

func objectReferenceToString(o corev1.ObjectReference) string {
	return fmt.Sprintf("%s/%s/%s/%s", o.APIVersion, o.Kind, o.Namespace, o.Name)
}

func getNamespaceForControllerInstallation(controllerInstallation *gardencorev1alpha1.ControllerInstallation) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("extension-%s", controllerInstallation.Name),
		},
	}
}

type deletionInProgressError struct {
	reason string
}

func newDeletionInProgressError(reason string) error {
	return &deletionInProgressError{
		reason: reason,
	}
}

func (e *deletionInProgressError) Error() string {
	return e.reason
}

func isDeletionInProgressError(err error) bool {
	_, ok := err.(*deletionInProgressError)
	return ok
}
