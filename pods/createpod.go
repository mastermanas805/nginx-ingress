package pods

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// CreatePod godoc
// @Summary Create a pod
// @Description Create a pod with the given image name in Kubernetes
// @Accept  json
// @Produce  json
// @Title POD
// @Param podName path string true "Pod Name"
// @Param image path string true "Image Name"
// @Success 200 {object} string "Pod created successfully"
// @Failure 400 {object} string "Invalid image name"
// @Failure 500 {object} string "Failed to create pod"
// @Router /createpod/{podName}/{image} [post]
func CreatePod(c *gin.Context) {
	podName := c.Param("podName")
	image := c.Param("image")

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//create deployment
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": podName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": podName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  podName,
							Image: image,
							Ports: []v1.ContainerPort{
								v1.ContainerPort{Name: "http", HostPort: 80, ContainerPort: 80},
							},
						},
					},
				},
			},
		},
	}
	// Create the deployment resource in the Kubernetes cluster
	_, err = clientset.AppsV1().Deployments("default").Create(c, deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	// Define the service resource
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": podName,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// Create the service resource in the Kubernetes cluster
	_, err = clientset.CoreV1().Services("default").Create(c, service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "pod, deployment and service created successfully"})
}
