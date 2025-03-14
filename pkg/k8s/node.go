import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
@	"github.com/gruntwork-io/terratest/modules/testing"
)

// AssertAnyNodeHasLabelE will return an error if no nodes are found with the given label.
func AssertAnyNodeHasLabelE(t testing.TestingT, options *KubectlOptions, label string) error {
	nodes, err := GetNodesByFilterE(t, options, metav1.ListOptions{
		LabelSelector: label,
	})

	require.NoError(t, err)
	if len(nodes) == 0 {
		return fmt.Errorf("No nodes found with label %s", label)
	}
	return nil
}

// AssertAnyNodeHasLabel will fail the test if no nodes are found with the given label.
func AssertAnyNodeHasLabel(t testing.TestingT, options *KubectlOptions, label string) {
	err := AssertAnyNodeHasLabelE(t, options, label)
	require.NoError(t, err)
}

func AssertAnyNodeHasTaintE(t testing.TestingT, options *KubectlOptions, taintKey string, taintValue string, taintEffect corev1.TaintEffect) error {
	nodes, err := GetNodesE(t, options)
	require.NoError(t, err)

	for _, node := range nodes {
		for _, taint := range node.Spec.Taints {
			if taint.Key == taintKey && taint.Value == taintValue && taint.Effect == taintEffect {
				return nil
			}
		}
	}

	return fmt.Errorf("No nodes found with taint %s/%s/%s", taintKey, taintValue, taintEffect)
}

// AssertAnyNodeHasTaint will fail the test if no nodes are found with the given taint.
func AssertAnyNodeHasTaint(t testing.TestingT, options *KubectlOptions, taintKey string, taintValue string, taintEffect corev1.TaintEffect) {
	err := AssertAnyNodeHasTaintE(t, options, taintKey, taintValue, taintEffect)
	require.NoError(t, err)
}
