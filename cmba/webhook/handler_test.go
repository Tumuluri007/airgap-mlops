package main

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newPodWithAnnotation(annotation, value string, initContainers, containers []string, image string) corev1.Pod {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-pod",
			Namespace:   "ml-serving",
			Annotations: map[string]string{},
		},
	}
	if annotation != "" {
		pod.Annotations[annotation] = value
	}
	for _, name := range initContainers {
		pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
			Name:  name,
			Image: "registry.airgap.local/cmba/verify:0.1.0",
		})
	}
	for _, name := range containers {
		pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
			Name:  name,
			Image: image,
		})
	}
	return pod
}

func TestPodImageMatches(t *testing.T) {
	pod := newPodWithAnnotation("", "", nil,
		[]string{"main"},
		"registry.airgap.local/ml-models/iris@sha256:abc")
	if !podImageMatches(pod, "registry.airgap.local/ml-models/iris@sha256:abc") {
		t.Fatalf("expected image match")
	}
	if podImageMatches(pod, "registry.airgap.local/ml-models/iris@sha256:other") {
		t.Fatalf("expected image mismatch")
	}
}

func TestHasInitContainerAndSidecar(t *testing.T) {
	pod := newPodWithAnnotation("", "",
		[]string{"cmba-verify"},
		[]string{"main", "cmba-sentinel"},
		"registry.airgap.local/ml-models/iris@sha256:abc")
	if !hasInitContainer(pod, "cmba-verify") {
		t.Fatalf("expected cmba-verify init container present")
	}
	if !hasContainer(pod, "cmba-sentinel") {
		t.Fatalf("expected cmba-sentinel sidecar present")
	}
}

func TestDenyHelper(t *testing.T) {
	res := deny("foo: %s", "bar")
	if res.Allowed {
		t.Fatalf("deny() should produce Allowed=false")
	}
	if res.Reason != "foo: bar" {
		t.Fatalf("unexpected reason: %s", res.Reason)
	}
}
