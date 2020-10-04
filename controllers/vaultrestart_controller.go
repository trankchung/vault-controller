/*
Copyright 2020 Chung Tran <chung.k.tran@gmail.com>.

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
	"time"

	"github.com/go-logr/logr"
	vapi "github.com/hashicorp/vault/api"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	restarterv1 "github.com/trankchung/vault-controller/api/v1"
)

// defaultPollingInterval defines polling interval if not defined in custom resource.
const defaultPollingInterval = "60s"

var (
	vault, _ = vapi.NewClient(vapi.DefaultConfig())
)

// VaultRestartReconciler reconciles a VaultRestart object
type VaultRestartReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=restarter.tran-scending.net,resources=vaultrestarts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=restarter.tran-scending.net,resources=vaultrestarts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=list;deletecollection

// Reconcile
func (r *VaultRestartReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("vaultrestart", req.NamespacedName)

	var rst restarterv1.VaultRestart

	// Fetch the VaultRestart custom resource and stop reconciling if previous
	// custom resource no longer exists.
	if err := r.Get(ctx, req.NamespacedName, &rst); err != nil {
		log.Info(fmt.Sprintf("resource %s no longer exists, skip until next update", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Custom resource does not define any matching labels so skip
	// processing it until it's updated.
	if len(rst.Spec.MatchingLabels) <= 0 {
		log.Info("no matching labels defined, skip until next update")
		return ctrl.Result{}, nil
	}

	s, err := vault.Logical().List("/p2/secret1")
	if err != nil {
		log.Error(err, "unable to fetch vault data")
		return ctrl.Result{}, err
	}

	log.Info(fmt.Sprintf("SECRET: %+v", s))

	// Get next run duration. If not specified or error, use default polling interval.
	dur, err := time.ParseDuration(rst.Spec.PollingInterval)
	if err != nil {
		dur, _ = time.ParseDuration(defaultPollingInterval)
	}

	// Re-queue the next reconciliation.
	return ctrl.Result{Requeue: true, RequeueAfter: dur}, nil
}

func (r *VaultRestartReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&restarterv1.VaultRestart{}).
		Complete(r)
}
