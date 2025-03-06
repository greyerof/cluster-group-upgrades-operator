package controllers

import (
	"fmt"
	"strings"

	cguv1alpha1 "github.com/openshift-kni/cluster-group-upgrades-operator/pkg/api/clustergroupupgrades/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

type CGUEventReasonType string

const (
	CGUValidationFailureMissingClusters   = "missing clusters"
	CGUValidationFailureMissingPolicies   = "missing policies"
	CGUValidationFailureInvalidPolicies   = "invalid policies"
	CGUValidationFailureAmbiguousPolicies = "ambiguous polcies"
)

// CGU Event Reasons
const (
	CGUEventReasonCreated               CGUEventReasonType = "CguCreated"
	CGUEventReasonUpgradeSuccess        CGUEventReasonType = "CguUpgradeSuccess"
	CGUEventReasonUpgradeTimeout        CGUEventReasonType = "cguUpgradeTimeout"
	CGUEventReasonBatchUpgradeStarted   CGUEventReasonType = "CguBatchUpgradeStarted"
	CGUEventReasonBatchUpgradeSuccess   CGUEventReasonType = "CguBatchUpgradeSuccess"
	CGUEventReasonBatchUpgradeTimedout  CGUEventReasonType = "CguBatchUpgradeTimeout"
	CGUEventReasonClusterUpgradeSuccess CGUEventReasonType = "CguClusterUpgradeSuccess"
	CGUEventReasonValidationFailure     CGUEventReasonType = "CguValidationFailure"
)

// CGU Event Messages
const (
	CGUEventMsgFmtCreated               = "New ClusterGroupUpgrade found: %s"
	CGUEventMsgFmtUpgradeSuccess        = "ClusterGroupUpgrade %s succeeded remediating policies"
	CGUEventMsgFmtUpgradeTimeout        = "ClusterGroupUpgrade %s timed-out remediating policies"
	CGUEventMsgFmtBatchUpgradeStarted   = "ClusterGroupUpgrade %s: batch index %d upgrade started"
	CGUEventMsgFmtBatchUpgradeSuccess   = "ClusterGroupUpgrade %s: all clusters in the batch index %d are compliant with managed policies"
	CGUEventMsgFmtBatchUpgradeTimedout  = "ClusterGroupUpgrade %s: some clusters in the batch index %d timed-out remediating policies"
	CGUEventMsgFmtClusterUpgradeSuccess = "ClusterGroupUpgrade %s: cluster %s upgrade finished successfully"
	CGUEventMsgFmtValidationFailure     = "ClusterGroupUpgrade %s: validation failure (%s): %s"
)

// CGU Validation errors
type PoliciesValidationFailureType string

const (
	CGUValidationErrorMsgMissingCluster = "missing clusters"

	CGUValidationErrorMsgNone              PoliciesValidationFailureType = "none"
	CGUValidationErrorMsgMissingPolicies   PoliciesValidationFailureType = "missing policies"
	CGUValidationErrorMsgAmbiguousPolicies PoliciesValidationFailureType = "ambiguous policies"
	CGUValidationErrorMsgInvalidPolicies   PoliciesValidationFailureType = "invalid policies"
)

func (r *ClusterGroupUpgradeReconciler) sendEventCGUCreated(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtCreated, cgu.Name)
	r.Recorder.Event(cgu, corev1.EventTypeNormal, string(CGUEventReasonCreated), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUSuccess(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtUpgradeSuccess, cgu.Name)
	r.Recorder.Event(cgu, corev1.EventTypeNormal, string(CGUEventReasonUpgradeSuccess), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUTimedout(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtUpgradeTimeout, cgu.Name)
	r.Recorder.Event(cgu, corev1.EventTypeWarning, string(CGUEventReasonUpgradeTimeout), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeStarted(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	batchClusters := []string{}
	for clusterName := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClusters = append(batchClusters, clusterName)
	}

	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeStarted, cgu.Name, cgu.Status.Status.CurrentBatch)

	r.Recorder.AnnotatedEventf(cgu, map[string]string{"cgu.openshift.io/batch-clusters": strings.Join(batchClusters, ",")},
		corev1.EventTypeNormal,
		string(CGUEventReasonBatchUpgradeStarted), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeSuccess(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	batchClusters := []string{}
	for clusterName := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClusters = append(batchClusters, clusterName)
	}

	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeSuccess, cgu.Name, cgu.Status.Status.CurrentBatch)

	r.Recorder.AnnotatedEventf(cgu, map[string]string{"cgu.openshift.io/batch-clusters": strings.Join(batchClusters, ",")},
		corev1.EventTypeNormal,
		string(CGUEventReasonBatchUpgradeSuccess), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeTimedout(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	batchClusters := []string{}
	for clusterName := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClusters = append(batchClusters, clusterName)
	}

	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeTimedout, cgu.Name, cgu.Status.Status.CurrentBatch)

	r.Recorder.AnnotatedEventf(cgu, map[string]string{"cgu.openshift.io/batch-clusters": strings.Join(batchClusters, ",")},
		corev1.EventTypeWarning,
		string(CGUEventReasonBatchUpgradeTimedout), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUClusterUpgradeSuccess(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterName string) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtClusterUpgradeSuccess, cgu.Name, clusterName)

	r.Recorder.AnnotatedEventf(cgu, map[string]string{"cgu.openshift.io/cluster": clusterName},
		corev1.EventTypeNormal,
		string(CGUEventReasonClusterUpgradeSuccess), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUValidationFailureMissingCluster(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterNames []string) {
	clusterNamesStr := strings.Join(clusterNames, ",")
	evMsg := fmt.Sprintf(CGUEventMsgFmtValidationFailure, cgu.Name, CGUValidationErrorMsgMissingCluster, clusterNamesStr)

	r.Recorder.AnnotatedEventf(cgu, map[string]string{"cgu.openshift.io/cluster": clusterNamesStr},
		corev1.EventTypeNormal,
		string(CGUEventReasonValidationFailure), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUVPoliciesValidationFailure(cgu *cguv1alpha1.ClusterGroupUpgrade, failureType PoliciesValidationFailureType, info policiesInfo) {
	r.Log.Info("Sending policies validation failure event",
		"cgu", cgu.Namespace+"/"+cgu.Name,
		"failureType", string(failureType),
	)

	var evMsg string
	anns := map[string]string{}

	switch failureType {
	case CGUValidationErrorMsgMissingPolicies:
		missingPoliciesStr := strings.Join(info.missingPolicies, ",")
		evMsg = fmt.Sprintf(CGUEventMsgFmtValidationFailure, cgu.Name, failureType, missingPoliciesStr)
		anns["cgu.openshift.io/missing-policies"] = missingPoliciesStr
	case CGUValidationErrorMsgInvalidPolicies:
		invalidPoliciesStr := strings.Join(info.invalidPolicies, ",")
		evMsg = fmt.Sprintf(CGUEventMsgFmtValidationFailure, cgu.Name, failureType, invalidPoliciesStr)
		anns["cgu.openshift.io/invalid-policies"] = invalidPoliciesStr
	case CGUValidationFailureAmbiguousPolicies:
		ambiguousPolicies := []string{}
		for policy := range info.duplicatedPoliciesNs {
			ambiguousPolicies = append(ambiguousPolicies, policy)
		}

		ambiguousPoliciesStr := strings.Join(ambiguousPolicies, ",")

		evMsg = fmt.Sprintf(CGUEventMsgFmtValidationFailure, cgu.Name, failureType, ambiguousPoliciesStr)
		anns["cgu.openshift.io/ambiguous-policies"] = ambiguousPoliciesStr
	}

	r.Recorder.AnnotatedEventf(cgu, anns, corev1.EventTypeWarning, string(CGUEventReasonValidationFailure), evMsg)
}
