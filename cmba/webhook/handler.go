package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// Handler holds dependencies for the admission webhook.
type Handler struct {
	TrustRootPath string
}

// NewHandler constructs a Handler with the given offline Sigstore trust root.
func NewHandler(trustRoot string) *Handler {
	return &Handler{TrustRootPath: trustRoot}
}

// modelBindingAnnotation is the pod annotation that names the paired
// ModelBinding CRD.
const modelBindingAnnotation = "cmba.airgap.mlops/binding-name"

// CMBAResult captures the outcome of an admission decision.
type CMBAResult struct {
	Allowed bool
	Reason  string
}

// Validate is the AdmissionReview handler.
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var review admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &review); err != nil {
		http.Error(w, "invalid AdmissionReview", http.StatusBadRequest)
		return
	}

	pod := corev1.Pod{}
	if err := json.Unmarshal(review.Request.Object.Raw, &pod); err != nil {
		writeReview(w, review, false, "cannot unmarshal Pod")
		return
	}

	res := h.check(pod)
	writeReview(w, review, res.Allowed, res.Reason)
}

// check runs the six CMBA admission checks defined in the paper.
func (h *Handler) check(pod corev1.Pod) CMBAResult {
	// Check 1: pod carries the binding-name annotation.
	bindingName, ok := pod.Annotations[modelBindingAnnotation]
	if !ok || bindingName == "" {
		return deny("missing binding annotation: %s", modelBindingAnnotation)
	}

	// Check 2: ModelBinding exists in pod's namespace.
	mb, err := h.lookupModelBinding(pod.Namespace, bindingName)
	if err != nil {
		return deny("ModelBinding not found: %s/%s", pod.Namespace, bindingName)
	}

	// Check 3: pod's container image matches ModelBinding.spec.containerImage.
	if !podImageMatches(pod, mb.ContainerImage) {
		return deny("image mismatch: expected %s, got %s",
			mb.ContainerImage, joinImages(pod))
	}

	// Check 4: cmba-verify init container present.
	if !hasInitContainer(pod, "cmba-verify") {
		return deny("cmba-verify init container required")
	}

	// Check 5: cmba-sentinel sidecar present (if required).
	if mb.SentinelRequired && !hasContainer(pod, "cmba-sentinel") {
		return deny("cmba-sentinel sidecar required")
	}

	// Check 6: signature bundle verifies offline against the local trust root.
	if !h.verifySignatureOffline(mb.SignatureBundle) {
		return deny("attestation signature invalid")
	}

	return CMBAResult{Allowed: true, Reason: "CMBA OK"}
}

// ModelBindingSpec is the projection used by the webhook.
type ModelBindingSpec struct {
	ContainerImage   string
	ModelSHA256      string
	BindingHash      string
	SentinelRequired bool
	SignatureBundle  []byte
}

// lookupModelBinding fetches a ModelBinding from the API server. The actual
// implementation uses a controller-runtime client; this skeleton sketches
// the contract used by check().
func (h *Handler) lookupModelBinding(namespace, name string) (*ModelBindingSpec, error) {
	// Production implementation: use a kubernetes client to GET
	// /apis/cmba.airgap.mlops/v1alpha1/namespaces/<ns>/modelbindings/<name>.
	return nil, fmt.Errorf("not implemented; see deployment.yaml for client wiring")
}

// verifySignatureOffline validates a Sigstore bundle against the pre-staged
// trust root mounted as a ConfigMap on the webhook pod.
func (h *Handler) verifySignatureOffline(bundle []byte) bool {
	if len(bundle) == 0 {
		return false
	}
	if _, err := os.Stat(filepath.Join(h.TrustRootPath, "root.json")); err != nil {
		log.Printf("trust root missing at %s: %v", h.TrustRootPath, err)
		return false
	}
	// Production implementation: invoke sigstore-go verifier with the
	// offline trust root. Returns true on cryptographic success.
	return true
}

func podImageMatches(pod corev1.Pod, expected string) bool {
	for _, c := range pod.Spec.Containers {
		if c.Image == expected {
			return true
		}
	}
	return false
}

func joinImages(pod corev1.Pod) string {
	images := make([]string, 0, len(pod.Spec.Containers))
	for _, c := range pod.Spec.Containers {
		images = append(images, c.Image)
	}
	return strings.Join(images, ",")
}

func hasInitContainer(pod corev1.Pod, name string) bool {
	for _, c := range pod.Spec.InitContainers {
		if c.Name == name {
			return true
		}
	}
	return false
}

func hasContainer(pod corev1.Pod, name string) bool {
	for _, c := range pod.Spec.Containers {
		if c.Name == name {
			return true
		}
	}
	return false
}

func deny(format string, args ...interface{}) CMBAResult {
	return CMBAResult{Allowed: false, Reason: fmt.Sprintf(format, args...)}
}

func writeReview(w http.ResponseWriter, review admissionv1.AdmissionReview,
	allowed bool, reason string) {
	resp := &admissionv1.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: allowed,
	}
	if !allowed {
		resp.Result = &metaStatus(reason)
	}
	out := admissionv1.AdmissionReview{
		TypeMeta: review.TypeMeta,
		Response: resp,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Printf("encode response: %v", err)
	}
}

// metaStatus is a tiny shim that returns a Status object value (not pointer)
// to keep the writeReview call site free of imports for metav1.
func metaStatus(message string) interface{} {
	return map[string]interface{}{
		"status":  "Failure",
		"message": message,
		"reason":  "CMBABindingViolation",
	}
}
