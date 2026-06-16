package controllers

import (
	"context"
	"fmt"
	"strings"

	utils "github.com/openshift-kni/cluster-group-upgrades-operator/controllers/utils"
	cguv1alpha1 "github.com/openshift-kni/cluster-group-upgrades-operator/pkg/api/clustergroupupgrades/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apiValidation "k8s.io/apimachinery/pkg/api/validation"
)

// CGU Event Reasons
const (
	CGUEventReasonCreated  = "CguCreated"
	CGUEventReasonStarted  = "CguStarted"
	CGUEventReasonSuccess  = "CguSuccess"
	CGUEventReasonTimedout = "CguTimedout"

	CGUEventReasonValidationFailure = "CguValidationFailure"
)

// CGU Event Messages
const (
	CGUEventMsgFmtCreated           = "New ClusterGroupUpgrade found: %s"
	CGUEventMsgFmtStarted           = "ClusterGroupUpgrade %s started remediating policies"
	CGUEventMsgFmtUpgradeSuccess    = "ClusterGroupUpgrade %s succeeded remediating policies"
	CGUEventMsgFmtUpgradeTimedout   = "ClusterGroupUpgrade %s timed-out remediating policies"
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

	CGUEventAnnotationKeyEvType                = CGUEventAnnotationKeyPrefix + "/event-type"
	CGUEventAnnotationKeyBatchClustersList     = CGUEventAnnotationKeyPrefix + "/batch-clusters"
	CGUEventAnnotationKeyBatchClustersCount    = CGUEventAnnotationKeyPrefix + "/batch-clusters-count"
	CGUEventAnnotationKeyClusterName           = CGUEventAnnotationKeyPrefix + "/cluster-name"
	CGUEventAnnotationKeyTimedoutClustersList  = CGUEventAnnotationKeyPrefix + "/timedout-clusters"
	CGUEventAnnotationKeyTimedoutClustersCount = CGUEventAnnotationKeyPrefix + "/timedout-clusters-count"
	CGUEventAnnotationKeyTotalBatchesCount     = CGUEventAnnotationKeyPrefix + "/total-batches-count"
	CGUEventAnnotationKeyTotalClustersCount    = CGUEventAnnotationKeyPrefix + "/total-clusters-count"

	// Validation failures
	CGUEventAnnotationKeyMissingClustersList   = CGUEventAnnotationKeyPrefix + "/missing-clusters"
	CGUEventAnnotationKeyMissingClustersCount  = CGUEventAnnotationKeyPrefix + "/missing-clusters-count"
	CGUEventAnnotationKeyMissingPoliciesList   = CGUEventAnnotationKeyPrefix + "/missing-policies"
	CGUEventAnnotationKeyInvalidPoliciesList   = CGUEventAnnotationKeyPrefix + "/invalid-policies"
	CGUEventAnnotationKeyAmbiguousPoliciesList = CGUEventAnnotationKeyPrefix + "/ambiguous-policies"
)

// Values for the CGUEventAnnotationKeyEvType key
const (
	CGUAnnEventGlobalUpgrade  = "global"
	CGUAnnEventBatchUpgrade   = "batch"
	CGUAnnEventClusterUpgrade = "cluster"
)

const (
	maxEventAnnsSize = apiValidation.TotalAnnotationSizeLimitB
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

// emitEvent is a helper that calls the Emitter and logs failures.
// Events are best-effort; errors are logged but never propagated to the caller.
func (r *ClusterGroupUpgradeReconciler) emitEvent(cgu *cguv1alpha1.ClusterGroupUpgrade,
	annotations map[string]string, eventType, reason, note string) {

	if err := r.Emitter.Emit(context.TODO(), cgu, annotations, eventType, reason, "Reconcile", note); err != nil {
		r.Log.Error(err, "failed to emit event", "reason", reason, "cgu", cgu.Namespace+"/"+cgu.Name)
	}
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUCreated(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtCreated, cgu.Name)

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType: CGUAnnEventGlobalUpgrade,
	}

	r.emitEvent(cgu, evAnns, corev1.EventTypeNormal, CGUEventReasonCreated, evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUStarted(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtStarted, cgu.Name)

	clustersCount := getTotalClustersNum(cgu)
	batchesCount := len(cgu.Status.RemediationPlan)

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType:             CGUAnnEventGlobalUpgrade,
		CGUEventAnnotationKeyTotalClustersCount: fmt.Sprint(clustersCount),
		CGUEventAnnotationKeyTotalBatchesCount:  fmt.Sprint(batchesCount),
	}

	r.emitEvent(cgu, evAnns, corev1.EventTypeNormal, CGUEventReasonStarted, evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUSuccess(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtUpgradeSuccess, cgu.Name)

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType:             CGUAnnEventGlobalUpgrade,
		CGUEventAnnotationKeyTotalClustersCount: fmt.Sprint(getTotalClustersNum(cgu)),
	}

	r.emitEvent(cgu, evAnns, corev1.EventTypeNormal, CGUEventReasonSuccess, evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUTimedout(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtUpgradeTimedout, cgu.Name)

	timedoutClusters := []string{}
	for _, clusterState := range cgu.Status.Clusters {
		if clusterState.State == utils.ClusterRemediationTimedout {
			timedoutClusters = append(timedoutClusters, clusterState.Name)
		}
	}

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType:                CGUAnnEventGlobalUpgrade,
		CGUEventAnnotationKeyTimedoutClustersCount: fmt.Sprint(len(timedoutClusters)),
		CGUEventAnnotationKeyTimedoutClustersList:  strings.Join(timedoutClusters, ","),
		CGUEventAnnotationKeyTotalClustersCount:    fmt.Sprint(getTotalClustersNum(cgu)),
	}

	truncateAnnotations(evAnns, maxEventAnnsSize)

	r.emitEvent(cgu, evAnns, corev1.EventTypeWarning, CGUEventReasonTimedout, evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeStarted(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	batchClusters := []string{}
	for clusterName := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClusters = append(batchClusters, clusterName)
	}

	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeStarted, cgu.Name, cgu.Status.Status.CurrentBatch)

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType:             CGUAnnEventBatchUpgrade,
		CGUEventAnnotationKeyBatchClustersCount: fmt.Sprint(len(batchClusters)),
		CGUEventAnnotationKeyTotalClustersCount: fmt.Sprint(getTotalClustersNum(cgu)),
		CGUEventAnnotationKeyBatchClustersList:  strings.Join(batchClusters, ","),
	}

	truncateAnnotations(evAnns, maxEventAnnsSize)

	r.emitEvent(cgu, evAnns, corev1.EventTypeNormal, CGUEventReasonStarted, evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeSuccess(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeSuccess, cgu.Name, cgu.Status.Status.CurrentBatch)

	batchClusters := []string{}
	for clusterName := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClusters = append(batchClusters, clusterName)
	}

	batchClustersCount := len(batchClusters)
	totalClustersCount := getTotalClustersNum(cgu)

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType:             CGUAnnEventBatchUpgrade,
		CGUEventAnnotationKeyBatchClustersCount: fmt.Sprint(batchClustersCount),
		CGUEventAnnotationKeyTotalClustersCount: fmt.Sprint(totalClustersCount),
		CGUEventAnnotationKeyBatchClustersList:  strings.Join(batchClusters, ","),
	}

	truncateAnnotations(evAnns, maxEventAnnsSize)

	r.emitEvent(cgu, evAnns, corev1.EventTypeNormal, CGUEventReasonSuccess, evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUBatchUpgradeTimedout(cgu *cguv1alpha1.ClusterGroupUpgrade) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtBatchUpgradeTimedout, cgu.Name, cgu.Status.Status.CurrentBatch)

	batchClustersCount := 0
	timedoutClusters := []string{}
	for clusterName, clusterProgress := range cgu.Status.Status.CurrentBatchRemediationProgress {
		batchClustersCount++
		if clusterProgress.State == cguv1alpha1.InProgress {
			timedoutClusters = append(timedoutClusters, clusterName)
		}
	}

	timedoutClustersCount := len(timedoutClusters)
	totalClustersCount := getTotalClustersNum(cgu)

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType:                CGUAnnEventBatchUpgrade,
		CGUEventAnnotationKeyTimedoutClustersCount: fmt.Sprint(timedoutClustersCount),
		CGUEventAnnotationKeyTimedoutClustersList:  strings.Join(timedoutClusters, ","),
		CGUEventAnnotationKeyBatchClustersCount:    fmt.Sprint(batchClustersCount),
		CGUEventAnnotationKeyTotalClustersCount:    fmt.Sprint(totalClustersCount),
	}

	truncateAnnotations(evAnns, maxEventAnnsSize)

	r.emitEvent(cgu, evAnns, corev1.EventTypeWarning, CGUEventReasonTimedout, evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUClusterUpgradeStarted(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterName string) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtClusterUpgradeStarted, cgu.Name, clusterName)

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType:      CGUAnnEventClusterUpgrade,
		CGUEventAnnotationKeyClusterName: clusterName,
	}

	truncateAnnotations(evAnns, maxEventAnnsSize)

	r.emitEvent(cgu, evAnns, corev1.EventTypeNormal, CGUEventReasonStarted, evMsg)
}

func (r *ClusterGroupUpgradeReconciler) sendEventCGUClusterUpgradeSuccess(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterName string) {
	evMsg := fmt.Sprintf(CGUEventMsgFmtClusterUpgradeSuccess, cgu.Name, clusterName)

	evAnns := map[string]string{
		CGUEventAnnotationKeyEvType:      CGUAnnEventClusterUpgrade,
		CGUEventAnnotationKeyClusterName: clusterName,
	}

	truncateAnnotations(evAnns, maxEventAnnsSize)

	r.emitEvent(cgu, evAnns, corev1.EventTypeNormal, CGUEventReasonSuccess, evMsg)
}

// func (r *ClusterGroupUpgradeReconciler) sendEventCGUClusterUpgradeTimedout(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterName string) {
// 	evMsg := fmt.Sprintf(CGUEventMsgFmtClusterUpgradeTimedout, cgu.Name, clusterName)

// 	evAnns := map[string]string{
// 		CGUEventAnnotationKeyEvType:      CGUAnnEventClusterUpgrade,
// 		CGUEventAnnotationKeyClusterName: clusterName,
// 	}

// 	truncateAnnotations(evAnns, maxEventAnnsSize)

// 	r.Recorder.AnnotatedEventf(cgu,
// 		evAnns,
// 		corev1.EventTypeNormal,
// 		CGUEventReasonTimedout, evMsg)
// }

func (r *ClusterGroupUpgradeReconciler) sendEventCGUValidationFailureMissingClusters(cgu *cguv1alpha1.ClusterGroupUpgrade, clusterNames []string) {
	clusterNamesStr := strings.Join(clusterNames, ",")
	evMsg := fmt.Sprintf(CGUEventMsgFmtValidationFailure, cgu.Name, CGUValidationErrorMsgMissingCluster, clusterNamesStr)

	evAnns := map[string]string{
		CGUEventAnnotationKeyMissingClustersCount: fmt.Sprint(len(clusterNames)),
		CGUEventAnnotationKeyMissingClustersList:  clusterNamesStr,
	}

	truncateAnnotations(evAnns, maxEventAnnsSize)

	r.emitEvent(cgu, evAnns, corev1.EventTypeNormal, CGUEventReasonValidationFailure, evMsg)
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

	truncateAnnotations(anns, maxEventAnnsSize)

	r.emitEvent(cgu, anns, corev1.EventTypeWarning, CGUEventReasonValidationFailure, evMsg)
}

// truncateAnnotations shrinks annotations whose values can grow unbounded
// (batch clusters, timedout clusters, etc.) to stay within the Kubernetes
// annotation size limit of 64 KiB.
//
// No truncation is made if maxSize is 0.
// nolint: unparam
func truncateAnnotations(anns map[string]string, maxSize int) {
	canBeTruncatedAnnKeys := map[string]bool{
		CGUEventAnnotationKeyBatchClustersList:    true,
		CGUEventAnnotationKeyTimedoutClustersList: true,
		CGUEventAnnotationKeyMissingClustersList:  true,
	}

	var totalAnnsSize int64
	for k, v := range anns {
		totalAnnsSize += int64(len(k)) + int64(len(v))
	}

	if totalAnnsSize <= int64(maxSize) {
		return
	}

	sizeToShrink := totalAnnsSize - int64(maxSize)

	for k, v := range anns {
		if !canBeTruncatedAnnKeys[k] {
			continue
		}

		maxListStrLen := int64(len(v)) - sizeToShrink
		anns[k] = truncateListString(v, maxListStrLen)

		// Only one truncatable annotation per event, so we're done.
		break
	}
}

// truncateListString keeps only the comma-separated elements that fit in maxSize.
func truncateListString(listStr string, maxSize int64) string {
	newElems := []string{}

	elems := strings.Split(listStr, ",")
	for _, elem := range elems {
		newElems = append(newElems, elem)

		newElemsStr := strings.Join(newElems, ",")
		if int64(len(newElemsStr)) > maxSize {
			newElems = newElems[:len(newElems)-1]
			return strings.Join(newElems, ",")
		}
	}

	return strings.Join(newElems, ",")
}

func getTotalClustersNum(cgu *cguv1alpha1.ClusterGroupUpgrade) int {
	total := 0
	for _, clusters := range cgu.Status.RemediationPlan {
		total += len(clusters)
	}

	return total
}
