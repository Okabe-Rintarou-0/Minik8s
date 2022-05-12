package api

//
//import (
//	"fmt"
//	"api/batch/v1"
//	"/api/core/v1"
//	metav1 "/pkg/apis/meta/v1"
//	clientsetscheme "client-go/kubernetes/scheme"
//	"time"
//)
//
//func isJobFinished(j *batchv1.Job) bool {
//	for _, c := range j.Status.Conditions {
//		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == v1.ConditionTrue {
//			return true
//		}
//	}
//	return false
//}
//
//func isJobCompleted(j *batchv1.Job) bool {
//	if len(j.Status.Conditions) != 1 {
//		return false
//	}
//
//	c := j.Status.Conditions[0]
//	if c.Type == batchv1.JobComplete && c.Status == v1.ConditionTrue {
//		return true
//	}
//	return false
//}
//
//func isPodCompleted(p *v1.Pod) bool {
//	if len(p.Status.Conditions) != 1 {
//		return false
//	}
//
//	if p.Status.Phase == v1.PodSucceeded || p.Status.Phase == v1.PodFailed {
//		return true
//	}
//
//	return false
//}
//
//func decodeJob(value []byte) (*batchv1.Job, error) {
//	decoder := clientsetscheme.Codecs.UniversalDeserializer()
//	obj, _, err := decoder.Decode(value, nil, nil)
//	if err != nil {
//		return nil, fmt.Errorf("Decode failed:%v", err)
//	}
//	job := obj.(*batchv1.Job)
//	return job, nil
//}
//
//func decodePod(value []byte) (*v1.Pod, error) {
//	decoder := clientsetscheme.Codecs.UniversalDeserializer()
//	obj, _, err := decoder.Decode(value, nil, nil)
//	if err != nil {
//		return nil, fmt.Errorf("Decode failed:%v", err)
//	}
//	pod := obj.(*v1.Pod)
//	return pod, nil
//}
//
//func jobStartedAtMinutesAgo(job *batchv1.Job, minutes int) bool {
//	if job.Status.StartTime == nil {
//		return false
//	}
//
//	t := job.Status.StartTime.Add(time.Duration(minutes) * time.Minute)
//	now := time.Now()
//	return t.Before(now)
//}
//
//func podStartedAtMinutesAgo(pod *v1.Pod, minutes int) bool {
//	if pod.Status.StartTime == nil {
//		return false
//	}
//
//	t := pod.Status.StartTime.Add(time.Duration(minutes) * time.Minute)
//	now := time.Now()
//	return t.Before(now)
//}
//
//func StringTime(t *metav1.Time) string {
//	if t == nil {
//		return "null"
//	}
//}
