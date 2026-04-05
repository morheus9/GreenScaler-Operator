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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appv1alpha1 "github.com/morheus9/GreenScaler-Operator/api/v1alpha1"
)

var _ = Describe("GreenScalerService Controller", func() {
	Context("When reconciling a resource", func() {
		const (
			resourceName   = "test-resource"
			deploymentName = "test-workload"
			testNamespace  = "default"
		)
		initialReplicas := int32(1)
		scheduleReplicas := int32(2)

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: testNamespace,
		}
		deployNamespacedName := types.NamespacedName{
			Name:      deploymentName,
			Namespace: testNamespace,
		}
		greenscalerservice := &appv1alpha1.GreenScalerService{}

		BeforeEach(func() {
			By("creating the target Deployment")
			deploy := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName,
					Namespace: testNamespace,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &initialReplicas,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": deploymentName},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"app": deploymentName},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{Name: "pause", Image: "registry.k8s.io/pause:3.9"},
							},
						},
					},
				},
			}
			err := k8sClient.Get(ctx, deployNamespacedName, deploy)
			if err != nil && errors.IsNotFound(err) {
				Expect(k8sClient.Create(ctx, deploy)).To(Succeed())
			}

			By("creating the custom resource for the Kind GreenScalerService")
			err = k8sClient.Get(ctx, typeNamespacedName, greenscalerservice)
			if err != nil && errors.IsNotFound(err) {
				resource := &appv1alpha1.GreenScalerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: testNamespace,
					},
					Spec: appv1alpha1.GreenScalerServiceSpec{
						Targets: []appv1alpha1.ScaleTarget{
							{Kind: "Deployment", Name: deploymentName},
						},
						// from == to means the window matches the full day (see desiredReplicas).
						Schedule: []appv1alpha1.ScaleWindow{
							{From: "00:00", To: "00:00", Replicas: scheduleReplicas},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &appv1alpha1.GreenScalerService{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				By("cleaning up the GreenScalerService instance")
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}

			deploy := &appsv1.Deployment{}
			err = k8sClient.Get(ctx, deployNamespacedName, deploy)
			if err == nil {
				By("cleaning up the Deployment")
				Expect(k8sClient.Delete(ctx, deploy)).To(Succeed())
			}
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &GreenScalerServiceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
