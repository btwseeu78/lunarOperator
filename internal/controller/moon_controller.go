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
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MoonReconciler reconciles a Moon object
type MoonReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lunar.arpan.io,resources=moons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lunar.arpan.io,resources=moons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lunar.arpan.io,resources=moons/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=list;watch;get;patch;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Moon object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *MoonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Moon", req.Name)
	log.Info("Starting Reconcile of the Controller")

	// find the objects first
	var moon lunarv1.Moon
	err := r.Get(ctx, req.NamespacedName, &moon)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Object Might Be Getting Deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to get Moon kind", "Name", req.Name, "NameSpace", req.Namespace)
		return ctrl.Result{}, err
	}

	//Now we create the deployment for moon kind

	found := &appsv1.Deployment{}
	actual, err := r.deplyomentForMoon(moon)

	if err != nil {
		log.Error(err, "Unable to generate Deployment Spec", "Name", req.Name, "Namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	// Get deployment ByName
	err = r.Get(ctx, types.NamespacedName{Namespace: moon.Namespace, Name: moon.Name}, found)

	if err != nil {
		//if not found create it
		if errors.IsNotFound(err) {
			log.Info("Deployment Does Not Exist Creating deployment", "Name", moon.Name, "NameSpace", moon.Namespace)
			err = r.Create(ctx, actual)
			if err != nil {
				log.Error(err, "Unable To Create the deployment")
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil

		} else {
			log.Error(err, "Unable to create deployment", "Name", moon.Name, "NameSpace", moon.Namespace)
			return ctrl.Result{}, err
		}
	}

	// If the details are there check if it has different values from actual

	if !equality.Semantic.DeepDerivative(actual.Spec.Template, found.Spec.Template) {
		log.Info("Updating deployment for", "Name", found.Name, "Namespace", found.Namespace)
		found = actual
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "failed to update the deployments", "Name", found.Name, "NameSpace", found.Namespace)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// check if there is a mismatch of size
	log.Info("Getting Resource version", "resourceversion", found.ObjectMeta.ResourceVersion)
	replicas := moon.Spec.Replicas
	if found.Spec.Replicas != replicas {
		found.Spec.Replicas = replicas
		err := r.Update(ctx, found)
		log.Info("Getting Resource version", "resourceversion", found.ObjectMeta.ResourceVersion)
		if err != nil {
			log.Error(err, "unable to update the replicas")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

func (r *MoonReconciler) deplyomentForMoon(moon lunarv1.Moon) (*appsv1.Deployment, error) {
	depl := &appsv1.Deployment{}
	depl = &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      moon.Name,
			Namespace: moon.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: moon.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"deployment/moon": moon.Name},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"deployment/moon": moon.Name},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "moon-container",
							Image: "gcr.io/google-samples/gb-frontend:v4",
							Env: []v1.EnvVar{
								{
									Name:  "SUN",
									Value: moon.Spec.SunName,
								},
							},
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      "TCP",
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	err := ctrl.SetControllerReference(&moon, depl, r.Scheme)
	return depl, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *MoonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lunarv1.Moon{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
