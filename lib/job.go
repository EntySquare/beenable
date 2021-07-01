package lib

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetJob(jobName string, jobParallelism int32, deleteJobAfterFinishSec int32, nodeAffinity corev1.NodeAffinity,
	limitList corev1.ResourceList, requestList corev1.ResourceList,
	imageName string) *batchv1.Job {
	entySwarmImage := imageName

	sectorDataDirHostType := corev1.HostPathDirectoryOrCreate
	jobRestartPolicy := corev1.RestartPolicyNever

	//Dont Restart a failed job pod!!!
	zeroBackoffLimitIsRetryTimeForNeverRestartFailedJobPod := int32(3)

	jobLabelMaps := map[string]string{
		"app":   "enty-bee",
		"phase": "test",
	}

	priorityClassName := ""

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   jobName,
			Labels: jobLabelMaps,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            &zeroBackoffLimitIsRetryTimeForNeverRestartFailedJobPod,
			TTLSecondsAfterFinished: &deleteJobAfterFinishSec,
			Parallelism:             &jobParallelism,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: jobLabelMaps,
				},
				Spec: corev1.PodSpec{
					PriorityClassName: priorityClassName,
					Affinity: &corev1.Affinity{
						NodeAffinity: &nodeAffinity,
					},
					RestartPolicy: jobRestartPolicy,
					Volumes: []corev1.Volume{
						{
							Name: "beedatadir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/root/swarm/",
									Type: &sectorDataDirHostType,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "enty-swarm",
							Image: entySwarmImage,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "beedatadir",
									MountPath: "/etc/bee/",
								},
							},
							Command: []string{"/bin/sh", "-c"},
							//Args: []string{"/entyctl client report -i " + reportIp + " -p " + reportPort + " && /Plotter create -F " +
							//	farmerKey + " -P " + poolKey + " -d /root/rplots/" + userDir + " -t /root/rplots/" + userDir +
							//	" -k " + k + " --rp " + reportIp + " --po " + reportPort + " -b 10000"},
							Args: []string{"/entyctl client -c" + "/root/swarm/"},
							Resources: corev1.ResourceRequirements{
								Limits:   limitList,
								Requests: requestList,
							},
							Env: []corev1.EnvVar{
								{
									Name: "USER_DIR",
									// Value: userDir,
								},
								{
									Name: "JOB_NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								{
									Name: "JOB_POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return job
}
