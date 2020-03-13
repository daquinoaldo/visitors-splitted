package visitorsbackend

import (
	examplev1 "git.extrasys.it/aldo.daquino/visitors-backend-operator/pkg/apis/example/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// == Const and names ==========
const mysqlPort = 3306
const mysqlImage = "mysql:5.7" // TODO: handle updates
const mysqlAuthUsername = "visitors-user"
const mysqlAuthPassword = "visitors-pass"

func mysqlDeploymentName(instance *examplev1.VisitorsBackend) string {
	return instance.Name + "-mysql"
}

func mysqlServiceName(instance *examplev1.VisitorsBackend) string {
	return instance.Name + "-mysql-service"
}

func mysqlAuthName(instance *examplev1.VisitorsBackend) string {
	return instance.Name + "-mysql-auth"
}

func mysqlVolumeClaimName(instance *examplev1.VisitorsBackend) string {
	return instance.Name + "-mysql-pv-claim"
}

// == Resources functions ==========
func (r *ReconcileVisitorsBackend) mysqlAuthSecret(instance *examplev1.VisitorsBackend) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlAuthName(instance),
			Namespace: instance.Namespace,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"username": "visitors-user",
			"password": "visitors-pass",
		},
	}
	controllerutil.SetControllerReference(instance, secret, r.scheme)
	return secret
}

func (r *ReconcileVisitorsBackend) mysqlVolume(instance *examplev1.VisitorsBackend) *corev1.PersistentVolumeClaim {
	volume := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlVolumeClaimName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}
	controllerutil.SetControllerReference(instance, volume, r.scheme)
	return volume
}

func (r *ReconcileVisitorsBackend) mysqlDeployment(instance *examplev1.VisitorsBackend) *appsv1.Deployment {
	labels := labels(instance, "mysql")
	size := instance.Spec.DatabaseSize // TODO: But how to be consistent between replicas?

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: mysqlAuthName(instance),
			},
			Key: "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: mysqlAuthName(instance),
			},
			Key: "password",
		},
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlDeploymentName(instance),
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{
						Name: "mysql-persistent-storage",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: mysqlVolumeClaimName(instance),
							},
						},
					}},
					Containers: []corev1.Container{{
						Name:  "visitors-mysql",
						Image: mysqlImage,
						Ports: []corev1.ContainerPort{{
							ContainerPort: mysqlPort,
							Name:          "mysql",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "mysql-persistent-storage",
							MountPath: "/var/lib/mysql",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: "password",
							},
							{
								Name:  "MYSQL_DATABASE",
								Value: "visitors",
							},
							{
								Name:      "MYSQL_USER",
								ValueFrom: userSecret,
							},
							{
								Name:      "MYSQL_PASSWORD",
								ValueFrom: passwordSecret,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, deployment, r.scheme)
	return deployment
}

func (r *ReconcileVisitorsBackend) mysqlService(instance *examplev1.VisitorsBackend) *corev1.Service {
	labels := labels(instance, "mysql")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlServiceName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port: mysqlPort,
			}},
			ClusterIP: "None",
		},
	}

	controllerutil.SetControllerReference(instance, service, r.scheme)
	return service
}

// == Reconcile functions ==========
// Returns whether or not the MySQL deployment is running
func (r *ReconcileVisitorsBackend) isMysqlUp(instance *examplev1.VisitorsBackend) bool {
	// Get deployment
	found, err := r.getClient(instance, mysqlDeploymentName(instance))
	if err != nil {
		log.Error(err, "Deployment mysql not found")
		return false
	}

	// TODO: scaling is probably not so easy
	if found.Status.ReadyReplicas == instance.Spec.DatabaseSize {
		return true
	}

	return false
}
