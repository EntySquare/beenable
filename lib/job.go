package lib

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func GetJob(jobName string, jobParallelism int32, deleteJobAfterFinishSec int32, nodeAffinity corev1.NodeAffinity,
	limitList corev1.ResourceList, requestList corev1.ResourceList, swapEndpoint, swapEnable, swapGas, swapInitDeposit,
	debugApiEnable, networkId, mainnet, fullNode, verbosity, clefEnable, imageName, password, dataDir string,
	port1, port2, port3 int) *batchv1.Job {
	entyBeeImage := imageName

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
							Name: "bee-datadir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/root/swarm/bee/bee-docker/file",
									Type: &sectorDataDirHostType,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "enty-bee",
							Image: entyBeeImage,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "bee-datadir",
									MountPath: "/home/bee/bee/file",
								},
							},
							Command: []string{"/bin/sh", "-c"},
							Args: []string{"bee start --swap-endpoint=" + swapEndpoint + " --swap-enable=" + swapEnable + " --debug-api-enable=" +
								debugApiEnable + " --swap-initial-deposit=" + swapInitDeposit + " --network-id=" + networkId + " --full-node=" + fullNode +
								" --verbosity=" + verbosity + " --clef-signer-enable=" + clefEnable + " --swap-deployment-gas-price " + swapGas +
								" --password=" + password + " --data-dir=" + dataDir + " --mainnet=" + mainnet + "--api-addr=" + strconv.Itoa(port1) +
								"--p2p-addr=" + strconv.Itoa(port2) + "debug-api-addr=" + strconv.Itoa(port3)},
							Ports: []corev1.ContainerPort{
								{
									Name:          "api-addr",
									HostPort:      int32(port1),
									ContainerPort: 1633,
								},
								{
									Name:          "p2p-addr",
									HostPort:      int32(port2),
									ContainerPort: 1634,
								},
								{
									Name:          "debug-api-addr",
									HostPort:      int32(port3),
									ContainerPort: 1635,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits:   limitList,
								Requests: requestList,
							},
							Env: []corev1.EnvVar{
								//{
								//	Name: "USER_DIR",
								//	// Value: userDir,
								//},
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
