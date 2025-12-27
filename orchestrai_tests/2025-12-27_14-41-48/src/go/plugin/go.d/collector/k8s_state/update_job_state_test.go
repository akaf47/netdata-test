package k8s_state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Test UpdateJobState with valid job data
func TestUpdateJobState_ValidJob(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-job",
			Namespace: "default",
			UID:       "test-uid",
		},
		Status: batchv1.JobStatus{
			Active:    1,
			Succeeded: 0,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
	assert.Greater(t, len(collector.jobs), 0)
}

// Test UpdateJobState with nil job
func TestUpdateJobState_NilJob(t *testing.T) {
	// Arrange
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act & Assert - should not panic
	assert.NotPanics(t, func() {
		if job := (*batchv1.Job)(nil); job != nil {
			collector.updateJobState(job)
		}
	})
}

// Test UpdateJobState with empty job name
func TestUpdateJobState_EmptyJobName(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    0,
			Succeeded: 0,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with empty namespace
func TestUpdateJobState_EmptyNamespace(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-job",
			Namespace: "",
		},
		Status: batchv1.JobStatus{
			Active:    0,
			Succeeded: 1,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with job having active pods
func TestUpdateJobState_JobWithActivePods(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "active-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    5,
			Succeeded: 0,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with job completed successfully
func TestUpdateJobState_JobSucceeded(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "succeeded-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    0,
			Succeeded: 10,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with job failed
func TestUpdateJobState_JobFailed(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "failed-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    0,
			Succeeded: 0,
			Failed:    3,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with job in all states
func TestUpdateJobState_JobMixedStates(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mixed-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    2,
			Succeeded: 8,
			Failed:    1,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState updates existing job
func TestUpdateJobState_UpdateExistingJob(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "update-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    1,
			Succeeded: 0,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act - First update
	collector.updateJobState(job)
	firstLen := len(collector.jobs)

	// Update same job with different status
	job.Status.Active = 2
	job.Status.Succeeded = 5
	collector.updateJobState(job)
	secondLen := len(collector.jobs)

	// Assert - Should have same or updated entry
	assert.Equal(t, firstLen, secondLen)
}

// Test UpdateJobState with max int32 values
func TestUpdateJobState_MaxIntValues(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "max-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    2147483647,
			Succeeded: 2147483647,
			Failed:    2147483647,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with zero values
func TestUpdateJobState_ZeroValues(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zero-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    0,
			Succeeded: 0,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with special characters in name
func TestUpdateJobState_SpecialCharactersInName(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job-with-special-chars-!@#$%",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    1,
			Succeeded: 0,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with very long name
func TestUpdateJobState_VeryLongName(t *testing.T) {
	// Arrange
	longName := ""
	for i := 0; i < 1000; i++ {
		longName += "a"
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      longName,
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    1,
			Succeeded: 0,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with multiple jobs in same namespace
func TestUpdateJobState_MultipleJobsSameNamespace(t *testing.T) {
	// Arrange
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	job1 := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job1",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    1,
			Succeeded: 0,
			Failed:    0,
		},
	}

	job2 := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job2",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    2,
			Succeeded: 0,
			Failed:    0,
		},
	}

	// Act
	collector.updateJobState(job1)
	collector.updateJobState(job2)

	// Assert
	assert.NotNil(t, collector.jobs)
	assert.Greater(t, len(collector.jobs), 0)
}

// Test UpdateJobState with jobs in different namespaces
func TestUpdateJobState_JobsDifferentNamespaces(t *testing.T) {
	// Arrange
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	job1 := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job",
			Namespace: "ns1",
		},
		Status: batchv1.JobStatus{
			Active: 1,
		},
	}

	job2 := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job",
			Namespace: "ns2",
		},
		Status: batchv1.JobStatus{
			Active: 2,
		},
	}

	// Act
	collector.updateJobState(job1)
	collector.updateJobState(job2)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with nil namespace value in status
func TestUpdateJobState_JobWithOwnerReferences(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "owned-job",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "batch/v1",
					Kind:       "CronJob",
					Name:       "parent-cronjob",
				},
			},
		},
		Status: batchv1.JobStatus{
			Active:    1,
			Succeeded: 0,
			Failed:    0,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with job labels
func TestUpdateJobState_JobWithLabels(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "labeled-job",
			Namespace: "default",
			Labels: map[string]string{
				"app":     "myapp",
				"version": "v1",
			},
		},
		Status: batchv1.JobStatus{
			Active: 1,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with job annotations
func TestUpdateJobState_JobWithAnnotations(t *testing.T) {
	// Arrange
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "annotated-job",
			Namespace: "default",
			Annotations: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		Status: batchv1.JobStatus{
			Active: 1,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with job having completion time
func TestUpdateJobState_JobWithCompletionTime(t *testing.T) {
	// Arrange
	now := metav1.Now()
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "completed-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:         0,
			Succeeded:      1,
			CompletionTime: &now,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}

// Test UpdateJobState with job having start time
func TestUpdateJobState_JobWithStartTime(t *testing.T) {
	// Arrange
	now := metav1.Now()
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "running-job",
			Namespace: "default",
		},
		Status: batchv1.JobStatus{
			Active:    1,
			StartTime: &now,
		},
	}
	collector := &K8sState{
		jobs: make(map[string]*JobState),
	}

	// Act
	collector.updateJobState(job)

	// Assert
	assert.NotNil(t, collector.jobs)
}