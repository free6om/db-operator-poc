/*
Copyright 2024.

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

	"github.com/apecloud/kubeblocks/pkg/controller/kubebuilderx"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/free6om/db-operator-poc/api/v1alpha1"
)

// PGClusterReconciler reconciles a PGCluster object
type PGClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	record.EventRecorder
}

//+kubebuilder:rbac:groups=apps.pg.dboperator.io,resources=pgclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.pg.dboperator.io,resources=pgclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.pg.dboperator.io,resources=pgclusters/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps.kubeblocks.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PGCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PGClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("pg-cluster")

	err := kubebuilderx.NewController(ctx, r.Client, req, r.EventRecorder, logger).
		Prepare(objectTree()).
		Do(translateToKubeBlocksCluster()).
		Commit()

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *PGClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.PGCluster{}).
		Complete(r)
}
