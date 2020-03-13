package visitorsbackend

import (
	"context"

	examplev1 "git.extrasys.it/aldo.daquino/visitors-backend-operator/pkg/apis/example/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func labels(v *examplev1.VisitorsBackend, tier string) map[string]string {
	return map[string]string{
		"app":             "visitors",
		"visitorssite_cr": v.Name,
		"tier":            tier,
	}
}

func (r *ReconcileVisitorsBackend) ensureDeployment(request reconcile.Request,
	instance *examplev1.VisitorsBackend,
	deployment *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      deployment.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		log.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return &reconcile.Result{}, err
		}
		// Deployment was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileVisitorsBackend) ensureService(request reconcile.Request,
	instance *examplev1.VisitorsBackend,
	service *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      service.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileVisitorsBackend) ensureSecret(request reconcile.Request,
	instance *examplev1.VisitorsBackend,
	secret *corev1.Secret,
) (*reconcile.Result, error) {
	found := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      secret.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the secret
		log.Info("Creating a new secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.client.Create(context.TODO(), secret)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the secret not existing
		log.Error(err, "Failed to get Secret")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileVisitorsBackend) ensureVolume(request reconcile.Request,
	instance *examplev1.VisitorsBackend,
	volume *corev1.PersistentVolumeClaim,
) (*reconcile.Result, error) {
	found := &corev1.PersistentVolumeClaim{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      volume.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the volume
		log.Info("Creating a new volume", "Volume.Namespace", volume.Namespace, "Volume.Name", volume.Name)
		err = r.client.Create(context.TODO(), volume)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Volume", "Volume.Namespace", volume.Namespace, "Volume.Name", volume.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the volume not existing
		log.Error(err, "Failed to get Volume")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileVisitorsBackend) getClient(instance *examplev1.VisitorsBackend, deploymentName string) (*appsv1.Deployment, error) {
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      deploymentName,
		Namespace: instance.Namespace,
	}, found)
	return found, err
}

func (r *ReconcileVisitorsBackend) updateClient(found *appsv1.Deployment) (*reconcile.Result, error) {
	err := r.client.Update(context.TODO(), found)
	if err != nil {
		log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
		return &reconcile.Result{}, err
	}
	// Spec updated - return and requeue
	return &reconcile.Result{Requeue: true}, nil
}
