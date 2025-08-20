/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	rotationv1 "github.com/akmanon/configmap-rotator-operator/api/v1"
)

// ConfigMapRotatorReconciler reconciles a ConfigMapRotator object
type ConfigMapRotatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=rotation.my.domain,resources=configmaprotators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rotation.my.domain,resources=configmaprotators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=rotation.my.domain,resources=configmaprotators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigMapRotator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *ConfigMapRotatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// Fetch the ConfigMapRotator instance
	var rotator rotationv1.ConfigMapRotator
	if err := r.Get(ctx, req.NamespacedName, &rotator); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if rotation is needed
	now := time.Now()
	rotationInterval := time.Duration(rotator.Spec.RotationIntervalHours) * time.Hour

	needsRotation := rotator.Status.LastRotationTime == nil ||
		now.Sub(rotator.Status.LastRotationTime.Time) >= rotationInterval

	if needsRotation {
		// Generate new ConfigMap data
		newData := make(map[string]string)
		for key, template := range rotator.Spec.DataTemplate {
			// Simple example: append timestamp for rotation
			newData[key] = fmt.Sprintf("%s-%d", template, now.Unix())
		}

		// Create or update ConfigMap
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rotator.Spec.ConfigMapName,
				Namespace: rotator.Namespace,
			},
		}

		// Use CreateOrUpdate from controllerutil
		result, err := ctrl.CreateOrUpdate(ctx, r.Client, configMap, func() error {
			configMap.Data = newData
			return ctrl.SetControllerReference(&rotator, configMap, r.Scheme)
		})

		if err != nil {
			logger.Error(err, "Failed to create or update ConfigMap")
			return ctrl.Result{}, err
		}

		logger.Info("ConfigMap operation completed",
			"configmap", rotator.Spec.ConfigMapName,
			"operation", result)

		// Update status
		rotator.Status.LastRotationTime = &metav1.Time{Time: now}
		rotator.Status.CurrentGeneration++

		if err := r.Status().Update(ctx, &rotator); err != nil {
			logger.Error(err, "Failed to update ConfigMapRotator status")
			return ctrl.Result{}, err
		}

		logger.Info("ConfigMap rotated", "configmap", rotator.Spec.ConfigMapName)
	}

	// Schedule next reconciliation
	nextCheck := rotationInterval
	if rotator.Status.LastRotationTime != nil {
		elapsed := now.Sub(rotator.Status.LastRotationTime.Time)
		if elapsed < rotationInterval {
			nextCheck = rotationInterval - elapsed
		}
	}

	return ctrl.Result{RequeueAfter: nextCheck}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapRotatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rotationv1.ConfigMapRotator{}).
		Named("configmaprotator").
		Complete(r)
}
