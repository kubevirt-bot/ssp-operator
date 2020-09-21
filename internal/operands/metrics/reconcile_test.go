package metrics

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"kubevirt.io/ssp-operator/internal/common"
	"kubevirt.io/ssp-operator/pkg/apis"
	v1 "kubevirt.io/ssp-operator/pkg/apis/ssp/v1"
)

var log = logf.Log.WithName("metrics_operand")

var _ = Describe("Metrics operand", func() {
	const (
		namespace = "kubevirt"
		name      = "test-ssp"
	)

	var (
		request common.Request
	)

	It("should create metrics resources", func() {
		s := scheme.Scheme
		Expect(apis.AddToScheme(s)).ToNot(HaveOccurred())
		Expect(AddWatchTypesToScheme(s)).ToNot(HaveOccurred())

		client := fake.NewFakeClientWithScheme(s)
		request = common.Request{
			Request: reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				},
			},
			Client:  client,
			Scheme:  s,
			Context: context.Background(),
			Instance: &v1.SSP{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
			},
			Logger: log,
		}

		Expect(Reconcile(&request)).ToNot(HaveOccurred())
		expectResourceExists(newPrometheusRule(namespace), request)
	})
})

func expectResourceExists(resource common.Resource, request common.Request) {
	key, err := client.ObjectKeyFromObject(resource)
	Expect(err).ToNot(HaveOccurred())
	Expect(request.Client.Get(request.Context, key, resource)).ToNot(HaveOccurred())
}

func TestMetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrics Suite")
}