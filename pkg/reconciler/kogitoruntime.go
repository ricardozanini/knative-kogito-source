package reconciler

import (
	"context"
	"fmt"
	"github.com/kiegroup/kogito-operator/api/v1beta1"
	kogitoclient "github.com/kiegroup/kogito-operator/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"

	// knative.dev imports
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"

	sourcesv1 "knative.dev/eventing/pkg/apis/sources/v1"

	"go.uber.org/zap"
)

// newKogitoRuntimeCreated makes a new reconciler event with event type Normal, and
// reason KogitoRuntimeCreated.
func newKogitoRuntimeCreated(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, "KogitoRuntimeCreated", "created kogitoruntime: \"%s/%s\"", namespace, name)
}

// newKogitoRuntimeFailed makes a new reconciler event with event type Warning, and
// reason KogitoRuntimeFailed.
func newKogitoRuntimeFailed(namespace, name string, err error) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeWarning, "KogitoRuntimeFailed", "failed to create kogitoruntime: \"%s/%s\", %w", namespace, name, err)
}

// newKogitoRuntimeUpdated makes a new reconciler event with event type Normal, and
// reason KogitoRuntimeUpdated.
func newKogitoRuntimeUpdated(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, "KogitoRuntimeUpdated", "updated kogitoruntime: \"%s/%s\"", namespace, name)
}

// newKogitoRuntimeNotReady makes a new reconciler event with type Normal, and
// reason KogitoRuntimeNotReady
func newKogitoRuntimeNotReady(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, "KogitoRuntimeNotReady", "waiting for kogitoruntime: \"%s/%s\"", namespace, name)
}

type KogitoRuntimeReconciler struct {
	KubeClientSet   kubernetes.Interface
	KogitoClientSet kogitoclient.Interface
}

// ReconcileKogitoRuntime reconciles kogitoruntime resource for KogitoSource
func (r *KogitoRuntimeReconciler) ReconcileKogitoRuntime(
	ctx context.Context,
	owner kmeta.OwnerRefable,
	binder *sourcesv1.SinkBinding,
	expected *v1beta1.KogitoRuntime,
) (*appsv1.Deployment, *sourcesv1.SinkBinding, pkgreconciler.Event) {
	namespace := owner.GetObjectMeta().GetNamespace()
	ra, err := r.KogitoClientSet.AppV1beta1().KogitoRuntimes(namespace).Get(ctx, expected.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		ra, err = r.KogitoClientSet.AppV1beta1().KogitoRuntimes(namespace).Create(ctx, expected, metav1.CreateOptions{})
		if err != nil {
			return nil, binder, newKogitoRuntimeFailed(expected.Namespace, expected.Name, err)
		}
		return nil, binder, newKogitoRuntimeCreated(ra.Namespace, ra.Name)
	} else if err != nil {
		return nil, binder, fmt.Errorf("error getting receive adapter %q: %v", expected.Name, err)
	} else if !metav1.IsControlledBy(ra, owner.GetObjectMeta()) {
		return nil, binder, fmt.Errorf("kogitoruntime %q is not owned by %s %q",
			ra.Name, owner.GetGroupVersionKind().Kind, owner.GetObjectMeta().GetName())
	}

	deployment, err := r.getDeployment(ctx, ra)
	if err != nil {
		return nil, binder, fmt.Errorf("error getting kogitoruntime %q deployment", ra.Name)
	} else if deployment == nil {
		// KogitoRuntime doesn't have a deployment created yet, reconcile
		return nil, binder, newKogitoRuntimeNotReady(ra.Namespace, ra.Name)
	}
	if kogitoRuntimeSpecSync(ctx, binder, expected.Spec, ra.Spec, deployment.Spec.Template.Spec) {
		if ra, err = r.KogitoClientSet.AppV1beta1().KogitoRuntimes(namespace).Update(ctx, ra, metav1.UpdateOptions{}); err != nil {
			return deployment, binder, err
		}
		return deployment, binder, newKogitoRuntimeUpdated(ra.Namespace, ra.Name)
	}
	logging.FromContext(ctx).Debugw("Reusing existing receive adapter", zap.Any("receiveAdapter", ra))
	return deployment, binder, nil
}

// getDeployment gets the associated Deployment of the given Receive Adapter
func (r *KogitoRuntimeReconciler) getDeployment(ctx context.Context, ra *v1beta1.KogitoRuntime) (*appsv1.Deployment, error) {
	// KogitoRuntime has a associated Deployment with same name and namespace
	deployment, err := r.KubeClientSet.AppsV1().Deployments(ra.Namespace).Get(ctx, ra.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (r *KogitoRuntimeReconciler) FindOwned(ctx context.Context, owner kmeta.OwnerRefable, selector labels.Selector) (*v1beta1.KogitoRuntime, error) {
	kogitoRuntimeList, err := r.KogitoClientSet.AppV1beta1().KogitoRuntimes(owner.GetObjectMeta().GetNamespace()).List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		logging.FromContext(ctx).Error("Unable to list kogitoruntime: %v", zap.Error(err))
		return nil, err
	}
	for _, kogitoruntime := range kogitoRuntimeList.Items {
		if metav1.IsControlledBy(&kogitoruntime, owner.GetObjectMeta()) {
			return &kogitoruntime, nil
		}
	}
	return nil, apierrors.NewNotFound(schema.GroupResource{}, "")
}

// Returns true if an update is needed.
func kogitoRuntimeSpecSync(ctx context.Context, binder *sourcesv1.SinkBinding, expected v1beta1.KogitoRuntimeSpec, now v1beta1.KogitoRuntimeSpec, nowPodSpec corev1.PodSpec) bool {
	old := *now.DeepCopy()
	syncSink(ctx, binder, nowPodSpec)
	// TODO: should we guarantee labels and envs added in expected?
	return !equality.Semantic.DeepEqual(old, now)
}

func syncSink(ctx context.Context, binder *sourcesv1.SinkBinding, now corev1.PodSpec) {
	// call Do() to project sink information.
	ps := &duckv1.WithPod{}
	ps.Spec.Template.Spec = now

	binder.Do(ctx, ps)
}
