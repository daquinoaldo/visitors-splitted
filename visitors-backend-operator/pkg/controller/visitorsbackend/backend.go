package visitorsbackend

import (
	"context"
	"time"

	examplev1 "git.extrasys.it/aldo.daquino/visitors-backend-operator/pkg/apis/example/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// == Const and names ==========
const backendPort = 8000
const backendServicePort = 30685
const backendImage = "jdob/visitors-service:1.0.0"

func backendDeploymentName(instance *examplev1.VisitorsBackend) string {
	return instance.Name + "-backend"
}

func backendServiceName(instance *examplev1.VisitorsBackend) string {
	return instance.Name + "-backend-service"
}

// == Resources functions ==========
func (r *ReconcileVisitorsBackend) backendDeployment(instance *examplev1.VisitorsBackend) *appsv1.Deployment {
	labels := labels(instance, "backend")
	size := instance.Spec.BackendSize

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName(instance)},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName(instance)},
			Key:                  "password",
		},
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendDeploymentName(instance),
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
					Containers: []corev1.Container{{
						Image:           backendImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "visitors-service",
						Ports: []corev1.ContainerPort{{
							ContainerPort: backendPort,
							Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_DATABASE",
								Value: "visitors",
							},
							{
								Name:  "MYSQL_SERVICE_HOST",
								Value: mysqlServiceName(instance),
							},
							{
								Name:      "MYSQL_USERNAME",
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

func (r *ReconcileVisitorsBackend) backendService(instance *examplev1.VisitorsBackend) *corev1.Service {
	labels := labels(instance, "backend")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendServiceName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       backendPort,
				TargetPort: intstr.FromInt(backendPort),
				NodePort:   30685,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(instance, service, r.scheme)
	return service
}

// == Reconcile functions ==========
func (r *ReconcileVisitorsBackend) updateBackendStatus(instance *examplev1.VisitorsBackend) error {
	instance.Status.BackendImage = backendImage
	err := r.client.Status().Update(context.TODO(), instance)
	return err
}

func (r *ReconcileVisitorsBackend) handleBackendChanges(instance *examplev1.VisitorsBackend) (*reconcile.Result, error) {
	// Get deployment
	found, err := r.getClient(instance, backendDeploymentName(instance))
	if err != nil {
		// The deployment may not have been created yet, so requeue
		return &reconcile.Result{RequeueAfter: 5 * time.Second}, err
	}

	// Ensure replicas
	size := instance.Spec.BackendSize
	if size != *found.Spec.Replicas {
		found.Spec.Replicas = &size
		return r.updateClient(found)
	}

	return nil, nil
}
