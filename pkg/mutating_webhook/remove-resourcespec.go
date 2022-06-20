package mutating_webhook

import (
	"confidential-containers/peer-pods-webhook/pkg/utils"
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

const (
	VM_DEFAULT_CPU_MILLIVALUE = 1000     // millicpu
	VM_DEFAULT_MEM            = 16777216 // bytes
	RUNTIME_CLASS_NAME        = "kata-remote-cc"
	VM_ANNOTATION_CPU         = "kata.peerpods.io/vmcpu"
	VM_ANNOTATION_MEM         = "kata.peerpods.io/vmmem"
)

// remove the POD resource spec
func removePodResourceSpec(pod *corev1.Pod) (*corev1.Pod, error) {
	mpod := pod.DeepCopy()

	// Mutate only if the POD is using specific runtimeClass
	if mpod.Spec.RuntimeClassName == nil || *mpod.Spec.RuntimeClassName != RUNTIME_CLASS_NAME {
		return mpod, nil
	}

	vmCpuTotal := utils.GetResourceRequest(mpod, corev1.ResourceCPU)
	vmMemTotal := utils.GetResourceRequest(mpod, corev1.ResourceMemory)

	if vmCpuTotal < VM_DEFAULT_CPU_MILLIVALUE {
		vmCpuTotal = VM_DEFAULT_CPU_MILLIVALUE
	}
	if vmMemTotal < VM_DEFAULT_MEM {
		vmMemTotal = VM_DEFAULT_MEM
	}

	// Add vmMemTotal and vmCpuTotal as POD annotations

	if mpod.Annotations == nil {
		mpod.Annotations = map[string]string{}
	}

	mpod.Annotations[VM_ANNOTATION_CPU] = strconv.FormatInt(vmCpuTotal, 10)
	mpod.Annotations[VM_ANNOTATION_MEM] = strconv.FormatInt(vmMemTotal, 10)

	for idx, _ := range mpod.Spec.Containers {
		mpod.Spec.Containers[idx].Resources = corev1.ResourceRequirements{}
	}
	return mpod, nil
}
