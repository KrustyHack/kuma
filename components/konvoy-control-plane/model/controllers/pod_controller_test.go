package controllers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	konvoy_mesh "github.com/Kong/konvoy/components/konvoy-control-plane/model/api/v1alpha1"
	. "github.com/Kong/konvoy/components/konvoy-control-plane/model/controllers"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	client_fake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("PodController", func() {
	Describe("ProxyTemplateToPodsMapper", func() {
		It("should use Client to list Pods in the same namespace", func() {
			// given
			template := &konvoy_mesh.ProxyTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "custom-proxy-template",
					Namespace: "example",
				},
			}

			pod1 := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "app1",
					Namespace: "example",
					Annotations: map[string]string{
						konvoy_mesh.ProxyTemplateAnnotation: "custom-proxy-template",
					},
				},
			}

			pod2 := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "app2",
					Namespace: "example",
				},
			}

			pod3 := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "app3",
					Namespace: "example",
					Annotations: map[string]string{
						konvoy_mesh.ProxyTemplateAnnotation: "another-template",
					},
				},
			}

			pod4 := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "app4",
					Namespace: "demo",
					Annotations: map[string]string{
						konvoy_mesh.ProxyTemplateAnnotation: "custom-proxy-template",
					},
				},
			}

			expected := []reconcile.Request{
				{NamespacedName: types.NamespacedName{Namespace: "example", Name: "app1"}},
			}

			// setup
			scheme := runtime.NewScheme()
			corev1.AddToScheme(scheme)
			mapper := &ProxyTemplateToPodsMapper{
				Client: client_fake.NewFakeClientWithScheme(scheme, pod1, pod2, pod3, pod4),
			}

			// when
			actual := mapper.Map(handler.MapObject{Meta: template})

			// then
			Expect(actual).To(ConsistOf(expected))
		})

		It("should handle Client errors gracefully", func() {
			// given
			template := &konvoy_mesh.ProxyTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "custom-proxy-template",
					Namespace: "example",
				},
			}

			// setup
			mapper := &ProxyTemplateToPodsMapper{
				// List operation will fail since Pod type hasn't been registered with the scheme
				Client: client_fake.NewFakeClientWithScheme(runtime.NewScheme()),
			}

			// when
			actual := mapper.Map(handler.MapObject{Meta: template})

			// then
			Expect(actual).To(BeNil())
		})
	})
})
