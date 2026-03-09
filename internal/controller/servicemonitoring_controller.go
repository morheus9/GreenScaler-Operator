package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/morheus9/GreenScaler-Operator/api/v1alpha1"
)

// ServiceMonitoringReconciler reconciles a ServiceMonitoring object
type ServiceMonitoringReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=app.example.com,resources=servicemonitorings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.example.com,resources=servicemonitorings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.example.com,resources=servicemonitorings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceMonitoring object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *ServiceMonitoringReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)
	logger.Info("Start Reconciler")
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceMonitoringReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.ServiceMonitoring{}).
		Named("servicemonitoring").
		Complete(r)
}
