package controller

import (
	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/controller/builder"
	"github.com/apecloud/kubeblocks/pkg/controller/kubebuilderx"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
	"github.com/free6om/db-operator-poc/api/v1alpha1"
)

type translateToKubeBlocksClusterReconciler struct{}

func (r *translateToKubeBlocksClusterReconciler) PreCondition(tree *kubebuilderx.ObjectTree) *kubebuilderx.CheckResult {
	if tree.GetRoot() == nil || model.IsObjectDeleting(tree.GetRoot()) {
		return kubebuilderx.ResultUnsatisfied
	}
	return kubebuilderx.ResultSatisfied
}

func (r *translateToKubeBlocksClusterReconciler) Reconcile(tree *kubebuilderx.ObjectTree) (*kubebuilderx.ObjectTree, error) {
	pgCluster, _ := tree.GetRoot().(*v1alpha1.PGCluster)
	kbCluster := builder.NewClusterBuilder(pgCluster.Namespace, pgCluster.Name).GetObject()
	object, err := tree.Get(kbCluster)
	if err != nil {
		return nil, err
	}
	if object != nil {
		kbCluster, _ = object.(*appsv1alpha1.Cluster)
	}
	kbCluster.Spec.ClusterDefRef = "postgresql"
	kbCluster.Spec.ClusterVersionRef = "postgresql-12.14.0"
	kbCluster.Spec.ComponentSpecs = []appsv1alpha1.ClusterComponentSpec{
		{
			Name:            "postgresql",
			ComponentDefRef: "postgresql",
			Replicas:        int32(pgCluster.Spec.Instances),
			Resources:       pgCluster.Spec.Resources,
			SwitchPolicy:    &appsv1alpha1.ClusterSwitchPolicy{Type: appsv1alpha1.Noop},
		},
	}
	kbCluster.Spec.TerminationPolicy = appsv1alpha1.Delete
	if err = tree.Update(kbCluster); err != nil {
		return nil, err
	}
	return tree, nil
}

func translateToKubeBlocksCluster() kubebuilderx.Reconciler {
	return &translateToKubeBlocksClusterReconciler{}
}

var _ kubebuilderx.Reconciler = &translateToKubeBlocksClusterReconciler{}
