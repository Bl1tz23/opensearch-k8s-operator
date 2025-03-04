package builders

import (
	"fmt"

	sts "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"os-operator.io/pkg/helpers"

	//v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"strconv"

	//v1 "k8s.io/client-go/applyconfigurations/core/v1"
	opsterv1 "os-operator.io/api/v1"
)

/// package that declare and build all the resources that related to the OpenSearch cluster ///

func NewMasterSTSForCR(cr *opsterv1.Os) *sts.StatefulSet {
	disk := fmt.Sprint(cr.Spec.Masters.DiskSize)

	pvt := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: "pvc"},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(disk),
				},
			},
		},
	}

	cluster_init_node := helpers.CreateInitmasters(cr)
	//var vendor string
	labels := map[string]string{
		"app": cr.Name,
	}

	i, err := strconv.ParseInt("420", 10, 32)
	if err != nil {
		fmt.Println("here panic")
		panic(err)
	}
	mode := int32(i)
	//storageClass := "gp2"
	runas := int64(0)

	// if cr.Spec.General.Vendor == "Op" || cr.Spec.General.Vendor == "OP" ||
	// 	cr.Spec.General.Vendor == "Opensearch" ||
	// 	cr.Spec.General.Vendor == "opensearch" {
	//	vendor = "opensearchproject/opensearch"
	// } else {
	//vendor ="elasticsearch"
	// }

	//var jvm string
	////if cr.Spec.Masters.Jvm == "" {
	////	jvm = "-Xmx512M -Xms512M"
	////}

	return &sts.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.General.ClusterName + "-master",
			Namespace: cr.Spec.General.ClusterName,
			Labels:    labels,
		},
		Spec: sts.StatefulSetSpec{
			Replicas: &cr.Spec.Masters.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: nil,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:      "cluster.initial_master_nodes",
									Value:     cluster_init_node,
									ValueFrom: nil,
								},
								{
									Name:      "discovery.seed_hosts",
									Value:     cr.Spec.General.ServiceName + "-headleass-service",
									ValueFrom: nil,
								},
								{
									Name:      "cluster.name",
									Value:     cr.Spec.General.ClusterName,
									ValueFrom: nil,
								},
								{
									Name:      "network.host",
									Value:     "0.0.0.0",
									ValueFrom: nil,
								},
								{
									Name:      "OPENSEARCH_JAVA_OPTS",
									Value:     "-Xmx512M -Xms512M",
									ValueFrom: nil,
								},
								{
									Name:      "node.data",
									Value:     "true",
									ValueFrom: nil,
								},
								{
									Name:      "node.master",
									Value:     "true",
									ValueFrom: nil,
								},
								{
									Name:      "node.ingest",
									Value:     "false",
									ValueFrom: nil,
								},
								{
									Name:      "node.remote_cluster_client",
									Value:     "true",
									ValueFrom: nil,
								},
							},

							Name:  cr.Name,
							Image: "opensearchproject/opensearch:1.0.0",
							Ports: []corev1.ContainerPort{
								{
									Name:          cr.Spec.General.ServiceName + "-port",
									ContainerPort: cr.Spec.General.OsPort,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "pvc",
									MountPath: "/usr/share/opensearch/data",
								},
								{
									Name:      "opensearch-yml",
									MountPath: "/usr/share/opensearch/config/opensearch.yml",
									SubPath:   "opensearch.yml",
								},
							},
						},
					},
					InitContainers: []corev1.Container{{
						Name:    "init",
						Image:   "busybox",
						Command: []string{"sh", "-c"},
						Args:    []string{"chown -R 1000:1000 /usr/share/opensearch/data"},
						SecurityContext: &corev1.SecurityContext{
							RunAsUser: &runas,
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "pvc",
								MountPath: "/usr/share/opensearch/data",
							},
						},
					},
					},
					Volumes: []corev1.Volume{
						{Name: "opensearch-yml",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: "opensearch-yml"},
									DefaultMode:          &mode,
								},
							},
						},
					},
					//NodeSelector:       nil,
					ServiceAccountName: cr.Spec.General.ServiceAccount,
					//	Affinity:           nil,
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{pvt},
			ServiceName:          cr.Spec.General.ServiceName + "-svc",
		},
	}
}

func NewNodeSTSForCR(cr *opsterv1.Os) *sts.StatefulSet {

	/////// Declerations //////

	var jvm string
	disk := fmt.Sprint(cr.Spec.Nodes.DiskSize)
	pvt := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: "pvc"},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(disk),
				},
			},
		},
	}

	cluster_init_node := helpers.CreateInitmasters(cr)
	//var vendor string
	labels := map[string]string{
		"app": cr.Name,
	}

	i, err := strconv.ParseInt("420", 10, 32)
	if err != nil {
		panic(err)
	}
	mode := int32(i)
	//storageClass := "gp2"
	runas := int64(0)

	if cr.Spec.Nodes.Jvm == "" {
		jvm = "-Xmx512M -Xms512M"
	}

	// if cr.Spec.General.Vendor == "Op" || cr.Spec.General.Vendor == "OP" ||
	// 	cr.Spec.General.Vendor == "Opensearch" ||
	// 	cr.Spec.General.Vendor == "opensearch" {
	// 	//	vendor = "opensearchproject/opensearch"
	// } else {
	//vendor ="elasticsearch"
	// }

	//var ingest string

	//if cr.Spec.Nodes.Ingest != "" {
	//	ingest = "true"
	//}
	return &sts.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.General.ClusterName + "-node",
			Namespace: cr.Spec.General.ClusterName,
			Labels:    labels,
		},
		Spec: sts.StatefulSetSpec{
			Replicas: &cr.Spec.Nodes.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: nil,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:      "cluster.initial_master_nodes",
									Value:     cluster_init_node,
									ValueFrom: nil,
								},
								{
									Name:      "discovery.seed_hosts",
									Value:     cr.Spec.General.ServiceName + "-headleass-service",
									ValueFrom: nil,
								},
								{
									Name:      "cluster.name",
									Value:     cr.Spec.General.ClusterName,
									ValueFrom: nil,
								},
								{
									Name:      "network.host",
									Value:     "0.0.0.0",
									ValueFrom: nil,
								},
								{
									Name:      "OPENSEARCH_JAVA_OPTS",
									Value:     jvm,
									ValueFrom: nil,
								},
								{
									Name:      "node.data",
									Value:     "true",
									ValueFrom: nil,
								},
								{
									Name:      "node.master",
									Value:     "false",
									ValueFrom: nil,
								},
								{
									Name:      "node.ingest",
									Value:     "true",
									ValueFrom: nil,
								},
								{
									Name:      "node.remote_cluster_client",
									Value:     "true",
									ValueFrom: nil,
								},
							},

							Name:  cr.Name,
							Image: "opensearchproject/opensearch:1.0.0",
							Ports: []corev1.ContainerPort{
								{
									Name:          cr.Spec.General.ServiceName,
									ContainerPort: cr.Spec.General.OsPort,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "pvc",
									MountPath: "/usr/share/opensearch/data",
								},
								{
									Name:      "opensearch-yml",
									MountPath: "/usr/share/opensearch/config/opensearch.yml",
									SubPath:   "opensearch.yml",
								},
							},
						},
					},
					InitContainers: []corev1.Container{{
						Name:    "init",
						Image:   "busybox",
						Command: []string{"sh", "-c"},
						Args:    []string{"chown -R 1000:1000 /usr/share/opensearch/data"},
						SecurityContext: &corev1.SecurityContext{
							RunAsUser: &runas,
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "pvc",
								MountPath: "/usr/share/opensearch/data",
							},
						},
					},
					},
					Volumes: []corev1.Volume{
						{Name: "opensearch-yml",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: "opensearch-yml"},
									DefaultMode:          &mode,
								},
							},
						},
					},
					//NodeSelector:       nil,
					ServiceAccountName: cr.Spec.General.ServiceAccount,
					//	Affinity:           nil,
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{pvt},
			ServiceName:          cr.Spec.General.ServiceName + "-svc",
		},
	}
}

//func NewPVCForCR(cr *opsterv1alpha1.Es) *corev1.PersistentVolumeClaimTemplate {
//
//
//	return &corev1.PersistentVolumeClaimTemplate{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      cr.Name+"-pvc",
//			Namespace: cr.Spec.ClusterNS,
//		},
//		Spec:       corev1.PersistentVolumeClaimSpec{
//			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
//			Resources:   corev1.ResourceRequirements{
//				Requests: corev1.ResourceList{"storage": x},
//			},
//			VolumeName:  cr.Name+"-pv",
//		},
//	}
//}

func NewHeadlessServiceForCR(cr *opsterv1.Os) *corev1.Service {

	labels := map[string]string{
		"app": cr.Name,
	}

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.General.ServiceName + "-headleass-service",
			Namespace: cr.Spec.General.ClusterName,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: "TCP",
					Port:     cr.Spec.General.OsPort,
					TargetPort: intstr.IntOrString{
						IntVal: cr.Spec.General.OsPort,
					},
				},
				{
					Name:     "transport",
					Protocol: "TCP",
					Port:     9300,
					TargetPort: intstr.IntOrString{
						IntVal: 9300,
						StrVal: "9300",
					},
				},
			},
			Selector: labels,
			Type:     "",
		},
	}
}

func NewServiceForCR(cr *opsterv1.Os) *corev1.Service {

	labels := map[string]string{
		"app": cr.Name,
	}

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.General.ServiceName + "-svc",
			Namespace: cr.Spec.General.ClusterName,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: "TCP",
					Port:     cr.Spec.General.OsPort,
					TargetPort: intstr.IntOrString{
						IntVal: cr.Spec.General.OsPort,
					},
				},
				{
					Name:     "transport",
					Protocol: "TCP",
					Port:     9300,
					TargetPort: intstr.IntOrString{
						IntVal: 9300,
						StrVal: "9300",
					},
				},
				{
					Name:     "metrics",
					Protocol: "TCP",
					Port:     9600,
					TargetPort: intstr.IntOrString{
						IntVal: 9600,
						StrVal: "9600",
					},
				},
				{
					Name:     "rca",
					Protocol: "TCP",
					Port:     9650,
					TargetPort: intstr.IntOrString{
						IntVal: 9650,
						StrVal: "9650",
					},
				},
			},
			Selector: labels,
			Type:     "",
		},
	}
}

func NewNsForCR(cr *opsterv1.Os) *corev1.Namespace {

	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Spec.General.ClusterName,
		},
	}
}

func NewCmForCR(cr *opsterv1.Os) *corev1.ConfigMap {

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "opensearch-yml",
			Namespace: cr.Spec.General.ClusterName,
		},
		Data: map[string]string{
			"opensearch.yml": "      cluster.name: opensearch-cluster\n      network.host: 0.0.0.0\n      plugins:\n        security:\n          allow_default_init_securityindex: true\n          allow_unsafe_democertificates: true\n          audit.type: internal_opensearch\n          authcz:\n            admin_dn:\n            - CN=kirk,OU=client,O=client,L=test, C=de\n          check_snapshot_restore_write_privileges: true\n          enable_snapshot_restore_privilege: true\n          restapi:\n            roles_enabled:\n            - all_access\n            - security_rest_api_access\n          ssl:\n            http:\n              enabled: true\n              pemcert_filepath: esnode.pem\n              pemkey_filepath: esnode-key.pem\n              pemtrustedcas_filepath: root-ca.pem\n            transport:\n              enforce_hostname_verification: false\n              pemcert_filepath: esnode.pem\n              pemkey_filepath: esnode-key.pem\n              pemtrustedcas_filepath: root-ca.pem\n          system_indices:\n            enabled: true\n            indices:\n            - .opendistro-alerting-config\n            - .opendistro-alerting-alert*\n            - .opendistro-anomaly-results*\n            - .opendistro-anomaly-detector*\n            - .opendistro-anomaly-checkpoints\n            - .opendistro-anomaly-detection-state\n            - .opendistro-reports-*\n            - .opendistro-notifications-*\n            - .opendistro-notebooks\n            - .opendistro-asynchronous-search-response*",
		},
	}
}
