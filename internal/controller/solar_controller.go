/*
Copyright 2023.

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

	lunarv1 "github.com/btwseeu78/chasing-sun/api/v1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SolarReconciler reconciles a Solar object
type SolarReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lunar.arpan.io,resources=solars,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lunar.arpan.io,resources=solars/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lunar.arpan.io,resources=solars/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Solar object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *SolarReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Name", req.Name, "NameSpace", req.Namespace)

	solartype := &lunarv1.Solar{}

	err := r.Get(ctx, req.NamespacedName, solartype)

	if err != nil && errors.IsNotFound(err) {
		log.Info("The Controller probablyGetting Deleted", "Name", req.Name, "NameSpace", req.Namespace)
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Unable To find the object", "Name", req.Name, "Namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	cmap := &corev1.ConfigMap{}
	actual := &corev1.ConfigMap{}
	// find the ConfigMap

	err = r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, cmap)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("The ConfigMap Created is missing")
			actual, err = r.configMapForOperator(solartype)

			if err != nil {
				log.Error(err, "Unable to generate manifest", "Name", solartype.Spec.MyName, "NameSpace", solartype.Namespace)
			}

			err = r.Create(ctx, actual)

			if err != nil {
				log.Error(err, "Unable to Create Controller object", "Name", actual.Name, "NameSpace", actual.Namespace)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil

		}

	}

	//check for update the cache

	if !equality.Semantic.DeepDerivative(actual, cmap) {
		log.Info("Updating The diff in Configmap", "Name", req.Name, "NameSpece", req.Namespace)
		cmap = actual
		err = r.Update(ctx, cmap)
		if err != nil {
			log.Error(err, "Unable to Update", "Name", req.Name, "NameSpace", req.Namespace)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

func (r *SolarReconciler) configMapForOperator(solar *lunarv1.Solar) (*corev1.ConfigMap, error) {
	cmap := &corev1.ConfigMap{}
	configmapData := make(map[string]string, 0)
	sunProperties := `
	color.sun=red
	temp.sun=1300
	`
	configmapData["sun.properties"] = sunProperties
	cmap = &corev1.ConfigMap{
		Data: configmapData,
		ObjectMeta: metav1.ObjectMeta{
			Name:      solar.Spec.MyName,
			Namespace: solar.Namespace,
		},
	}
	err := ctrl.SetControllerReference(solar, cmap, r.Scheme)
	return cmap, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SolarReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lunarv1.Solar{}).
		Complete(r)
}
