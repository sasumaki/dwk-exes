/*


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

package controllers

import (
	"context"
	"fmt"

	stabledwkv1 "controller/api/v1"

	core "k8s.io/api/core/v1"

	kbatch "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"github.com/labstack/gommon/log"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DummySiteReconciler reconciles a DummySite object
type DummySiteReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.dwk.my.domain,resources=dummysites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.dwk.my.domain,resources=dummysites/status,verbs=get;update;patch

func (r *DummySiteReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("dummysite", req.NamespacedName)

	var dummySite stabledwkv1.DummySite

	// create the workqueue
	if err := r.Get(ctx, req.NamespacedName, &dummySite); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Println("name", dummySite.Name)

	fmt.Println("status: ", dummySite.Status, dummySite.Status.Status)
	fmt.Println("Spec", dummySite.Spec.Website_url)

	job := *buildJob(dummySite)
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: dummySite.Namespace, Name: dummySite.Name + "dummysite-job"}, &job)

	if apierrors.IsNotFound(err) {
		log.Info("could not find existing Job for DummySite, creating one...")

		if err := r.Client.Create(ctx, &job); err != nil {
			log.Error(err, "failed to create Job resource")
			return ctrl.Result{}, err
		}
	}
	log.Info("existing Job resource already exists for DummySite")
	fmt.Println(job)

	return ctrl.Result{}, nil
}

// apiVersion: batch/v1
// kind: Job
// metadata:
//   name: {{{ job_name }}}
//   labels:
//     countdown: "{{ countdown_name }}"
//     delay: "{{{ delay }}}"
//     length: "{{{ length }}}"
// spec:
//   template:
//     spec:
//       containers:
//       - name: {{{ container_name }}}
//         image: {{{ image }}}
//         args: ["{{{ length }}}"]
//       restartPolicy: Never

func buildJob(dummySite stabledwkv1.DummySite) *kbatch.Job {
	job := kbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:            dummySite.Name + "dummysite-job",
			Namespace:       dummySite.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&dummySite, stabledwkv1.GroupVersion.WithKind("DummySite"))},
		},
		Spec: kbatch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  dummySite.Name + "dummysite-job",
							Image: dummySite.Spec.Image,
							Args:  []string{"app", dummySite.Spec.Website_url},
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}
	return &job
}
func (r *DummySiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stabledwkv1.DummySite{}).
		Complete(r)
}
