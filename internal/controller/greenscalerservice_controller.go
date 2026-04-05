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

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/morheus9/GreenScaler-Operator/api/v1alpha1"
)

// GreenScalerServiceReconciler reconciles a GreenScalerService object
type GreenScalerServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=app.example.com,resources=greenscalerservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.example.com,resources=greenscalerservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.example.com,resources=greenscalerservices/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;update;patch

// Reconcile applies time-based replica counts from the GreenScalerService spec to
// each target workload. It requeues every minute so schedule boundaries are picked
// up without watching the clock. Errors surface as reconcile retries (with backoff).
//
// See: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile
func (r *GreenScalerServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)
	logger.V(1).Info("reconciling GreenScalerService", "namespacedName", req.NamespacedName)
	logger.Info("Reconciliation complete")

	var scaler appv1alpha1.GreenScalerService
	if err := r.Get(ctx, req.NamespacedName, &scaler); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	loc := time.UTC
	if scaler.Spec.TimeZone != "" {
		loaded, err := time.LoadLocation(scaler.Spec.TimeZone)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("invalid spec.timeZone %q: %w", scaler.Spec.TimeZone, err)
		}
		loc = loaded
	}

	now := time.Now().In(loc)
	replicas, ok, err := desiredReplicas(now, scaler.Spec.Schedule)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !ok {
		// No schedule window matches current time; requeue to catch the next boundary.
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	for _, t := range scaler.Spec.Targets {
		ns := t.Namespace
		if ns == "" {
			ns = scaler.Namespace
		}

		key := types.NamespacedName{Name: t.Name, Namespace: ns}
		switch t.Kind {
		case "Deployment":
			var obj appsv1.Deployment
			if err := r.Get(ctx, key, &obj); err != nil {
				return ctrl.Result{}, err
			}
			if obj.Spec.Replicas == nil || *obj.Spec.Replicas != replicas {
				obj.Spec.Replicas = ptrInt32(replicas)
				if err := r.Update(ctx, &obj); err != nil {
					return ctrl.Result{}, err
				}
			}
		case "StatefulSet":
			var obj appsv1.StatefulSet
			if err := r.Get(ctx, key, &obj); err != nil {
				return ctrl.Result{}, err
			}
			if obj.Spec.Replicas == nil || *obj.Spec.Replicas != replicas {
				obj.Spec.Replicas = ptrInt32(replicas)
				if err := r.Update(ctx, &obj); err != nil {
					return ctrl.Result{}, err
				}
			}
		default:
			return ctrl.Result{}, fmt.Errorf("unsupported target kind: %s", t.Kind)
		}
	}

	scaler.Status.LastAppliedReplicas = ptrInt32(replicas)
	nowMeta := metav1.NewTime(time.Now())
	scaler.Status.LastReconcileTime = &nowMeta
	if err := r.Status().Update(ctx, &scaler); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GreenScalerServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.GreenScalerService{}).
		Named("greenscalerservice").
		Complete(r)
}

// desiredReplicas returns the replica count for the first schedule window that matches now.
// ok is false when no window matches (operator should requeue until a window applies).
func desiredReplicas(now time.Time, schedule []appv1alpha1.ScaleWindow) (int32, bool, error) {
	minutesNow := now.Hour()*60 + now.Minute()

	for _, w := range schedule {
		start, err := hhmmToMinutes(w.From)
		if err != nil {
			return 0, false, fmt.Errorf("invalid schedule.from %q: %w", w.From, err)
		}
		end, err := hhmmToMinutes(w.To)
		if err != nil {
			return 0, false, fmt.Errorf("invalid schedule.to %q: %w", w.To, err)
		}

		// from == to means the window always matches (full day).
		if start == end {
			return w.Replicas, true, nil
		}

		if start < end {
			if minutesNow >= start && minutesNow < end {
				return w.Replicas, true, nil
			}
			continue
		}

		// Window wraps midnight (e.g. 22:00–06:00): active if after start OR before end.
		if minutesNow >= start || minutesNow < end {
			return w.Replicas, true, nil
		}
	}

	return 0, false, nil
}

// hhmmToMinutes parses "HH:MM" into minutes since midnight.
func hhmmToMinutes(hhmm string) (int, error) {
	t, err := time.Parse("15:04", hhmm)
	if err != nil {
		return 0, err
	}
	return t.Hour()*60 + t.Minute(), nil
}

// ptrInt32 returns a pointer to v (used for optional *int32 API fields).
func ptrInt32(v int32) *int32 {
	return &v
}
