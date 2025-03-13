package controllers

import (
	"fmt"
	"strings"

	cguv1alpha1 "github.com/openshift-kni/cluster-group-upgrades-operator/pkg/api/clustergroupupgrades/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

type CGUEventReasonType string

// CGU Event Reasons
const (
	CGUEventReasonCreated CGUEventReasonType = "CguCreated"
	CGUEventReasonSuccess CGUEventReasonType = "CguSuccess"
	CGUEventReasonTimeout CGUEventReasonType = "CguTimeout"

	CGUEventReasonValidationFailure CGUEventReasonType = "CguValidationFailure"

	CGUEventReasonUpgradeStarted  CGUEventReasonType = "CguUpgradeStarted"
	CGUEventReasonUpgradeSuccess  CGUEventReasonType = "CguUpgradeSuccess"
	CGUEventReasonUpgradeTimedout CGUEventReasonType = "CguUpgradeTimeout"
)

// CGU Event Messages
const (
	CGUEventMsgFmtCreated           = "New ClusterGroupUpgrade found: %s"
	CGUEventMsgFmtUpgradeSuccess    = "ClusterGroupUpgrade %s succeeded remediating policies"
	CGUEventMsgFmtUpgradeTimeout    = "ClusterGroupUpgrade %s timed-out remediating policies"
	CGUEventMsgFmtValidationFailure = "ClusterGroupUpgrade %s: validation failure (%s): %s"

	CGUEventMsgFmtBatchUpgradeStarted  = "ClusterGroupUpgrade %s: batch index %d upgrade started"
	CGUEventMsgFmtBatchUpgradeSuccess  = "ClusterGroupUpgrade %s: all clusters in the batch index %d are compliant with managed policies"
	CGUEventMsgFmtBatchUpgradeTimedout = "ClusterGroupUpgrade %s: some clusters in the batch index %d timed out remediating policies"

	CGUEventMsgFmtClusterUpgradeSuccess = "ClusterGroupUpgrade %s: cluster %s upgrade finished successfully"
	CGUEventMsgFmtClusterUpgradeStarted = "ClusterGroupUpgrade %s: cluster %s upgrade started"
	// CGUEventMsgFmtClusterUpgradeTimedout = "ClusterGroupUpgrade %s: cluster %s timed out remediating policies"
)

// CGU Validation Failure literals for the event's message.
const (
	CGUValidationFailureMissingClusters   = "missing clusters"
	CGUValidationFailureMissingPolicies   = "missing policies"
	CGUValidationFailureInvalidPolicies   = "invalid policies"
	CGUValidationFailureAmbiguousPolicies = "ambiguous polcies"
)

// Event annotation keys
const (
	CGUEventAnnotationKeyPrefix = "cgu.openshift.io"

	CGUEventAnnotationKeyEvType            = CGUEventAnnotationKeyPrefix + "/event-type"
	CGUEventAnnotationKeyBatchClustersList = CGUEventAnnotationKeyPrefix + "/batch-clusters"
	CGUEventAnnotationKeyClusterName       = CGUEventAnnotationKeyPrefix + "/cluster"

	// Validation failures
	CGUEventAnnotationKeyMissingClustersList   = CGUEventAnnotationKeyPrefix + "/missing-clusters"
	CGUEventAnnotationKeyMissingPoliciesList   = CGUEventAnnotationKeyPrefix + "/missing-policies"
	CGUEventAnnotationKeyInvalidPoliciesList   = CGUEventAnnotationKeyPrefix + "/invalid-policies"
	CGUEventAnnotationKeyAmbiguousPoliciesList = CGUEventAnnotationKeyPrefix + "/ambiguous-policies"
)

// Values for the CGUEventAnnotationKeyEvType key
const (
	CGUAnnEventBatchUpgrade   = "batch"
	CGUAnnEventClusterUpgrade = "cluster"
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
	r.Recorder.Event(cgu, corev1.EventTypeNormal, string(CGUEventReasonSuccess), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUTimedout(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtUpgradeTimeout, cgu.Name)
	r.Recorder.Event(cgu, corev1.EventTypeWarning, string(CGUEventReasonTimeout), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeStarted(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	batchClusters := []string{}
	for clusterName := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClusters = append(batchClusters, clusterName)
	}

	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeStarted, cgu.Name, cgu.Status.Status.CurrentBatch)

	r.Recorder.AnnotatedEventf(cgu,
		map[string]string{
			CGUEventAnnotationKeyEvType:            CGUAnnEventBatchUpgrade,
			CGUEventAnnotationKeyBatchClustersList: strings.Join(batchClusters, ","),
		},

		corev1.EventTypeNormal,
		string(CGUEventReasonUpgradeStarted), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeSuccess(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	batchClusters := []string{}
	for clusterName := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClusters = append(batchClusters, clusterName)
	}

	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeSuccess, cgu.Name, cgu.Status.Status.CurrentBatch)

	r.Recorder.AnnotatedEventf(cgu,
		map[string]string{
			CGUEventAnnotationKeyEvType:            CGUAnnEventBatchUpgrade,
			CGUEventAnnotationKeyBatchClustersList: strings.Join(batchClusters, ","),
		},
		corev1.EventTypeNormal,
		string(CGUEventReasonUpgradeSuccess), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeTimedout(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	batchClusters := []string{}
	for clusterName := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClusters = append(batchClusters, clusterName)
	}

	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeTimedout, cgu.Name, cgu.Status.Status.CurrentBatch)

	r.Recorder.AnnotatedEventf(cgu,
		map[string]string{
			CGUEventAnnotationKeyEvType:            CGUAnnEventBatchUpgrade,
			CGUEventAnnotationKeyBatchClustersList: strings.Join(batchClusters, ","),
		},
		corev1.EventTypeWarning,
		string(CGUEventReasonUpgradeTimedout), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUClusterUpgradeStarted(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterName string) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtClusterUpgradeStarted, cgu.Name, clusterName)

	r.Recorder.AnnotatedEventf(cgu,
		map[string]string{
			CGUEventAnnotationKeyEvType:      CGUAnnEventClusterUpgrade,
			CGUEventAnnotationKeyClusterName: clusterName,
		},
		corev1.EventTypeNormal,
		string(CGUEventReasonUpgradeStarted), evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUClusterUpgradeSuccess(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterName string) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtClusterUpgradeSuccess, cgu.Name, clusterName)

	r.Recorder.AnnotatedEventf(cgu,
		map[string]string{
			CGUEventAnnotationKeyEvType:      CGUAnnEventClusterUpgrade,
			CGUEventAnnotationKeyClusterName: clusterName,
		},
		corev1.EventTypeNormal,
		string(CGUEventReasonUpgradeSuccess), evMsg)
}

// func (r *ClusterGroupUpgradeReconciler) sendEventCGUClusterUpgradeTimedout(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterName string) {
// 	evMsg := fmt.Sprintf(CGUEventMsgFmtClusterUpgradeTimedout, cgu.Name, clusterName)

// 	r.Recorder.AnnotatedEventf(cgu,
// 		map[string]string{
// 			CGUEventAnnotationKeyEvType:      CGUAnnEventClusterUpgrade,
// 			CGUEventAnnotationKeyClusterName: clusterName,
// 		},
// 		corev1.EventTypeNormal,
// 		string(CGUEventReasonUpgradeTimedout), evMsg)
// }

func (r *ClusterGroupUpgradeReconciler) sendEventCGUValidationFailureMissingClusters(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterNames []string) {
	clusterNamesStr := strings.Join(clusterNames, ",")
	evMsg := fmt.Sprintf(CGUEventMsgFmtValidationFailure, cgu.Name, CGUValidationErrorMsgMissingCluster, clusterNamesStr)

	r.Recorder.AnnotatedEventf(cgu,
		map[string]string{
			CGUEventAnnotationKeyMissingClustersList: clusterNamesStr,
		},
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
		anns[CGUEventAnnotationKeyMissingPoliciesList] = missingPoliciesStr
	case CGUValidationErrorMsgInvalidPolicies:
		invalidPoliciesStr := strings.Join(info.invalidPolicies, ",")
		evMsg = fmt.Sprintf(CGUEventMsgFmtValidationFailure, cgu.Name, failureType, invalidPoliciesStr)
		anns[CGUEventAnnotationKeyInvalidPoliciesList] = invalidPoliciesStr
	case CGUValidationFailureAmbiguousPolicies:
		ambiguousPolicies := []string{}
		for policy := range info.duplicatedPoliciesNs {
			ambiguousPolicies = append(ambiguousPolicies, policy)
		}

		ambiguousPoliciesStr := strings.Join(ambiguousPolicies, ",")

		evMsg = fmt.Sprintf(CGUEventMsgFmtValidationFailure, cgu.Name, failureType, ambiguousPoliciesStr)
		anns[CGUEventAnnotationKeyAmbiguousPoliciesList] = ambiguousPoliciesStr
	}

	r.Recorder.AnnotatedEventf(cgu, anns, corev1.EventTypeWarning, string(CGUEventReasonValidationFailure), evMsg)
}
