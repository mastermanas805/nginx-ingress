package pods

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ExposePod godoc
// @Summary Expose a pod
// @Description Expose a pod with the given pod name in Kubernetes
// @ID expose-pod
// @Accept  json
// @Produce  json
// @Param podName path string true "Pod Name"
// @Param port path string true "Port Number"
// @Success 200 {object} string "Pod exposed successfully"
// @Failure 400 {object} string "Invalid pod name"
// @Failure 500 {object} string "Failed to expose pod"
// @Router /exposepod/{podName}/{port} [post]
func ExposePod(c *gin.Context) {
	podname := c.Param("podName")
	port := c.Param("port")
	ingressclass := "nginx"

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pathtype := networking.PathTypePrefix

	ingress := &networking.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:          podname,
			ManagedFields: []metav1.ManagedFieldsEntry{metav1.ManagedFieldsEntry{Manager: "nginx-ingress-controller"}},
		},
		Spec: networking.IngressSpec{
			IngressClassName: &ingressclass,
			Rules: []networking.IngressRule{
				{
					Host: podname + "-" + port + ".r.localdev.me",
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathtype,
									Backend: networking.IngressBackend{
										Service: &networking.IngressServiceBackend{
											Name: podname,
											Port: networking.ServiceBackendPort{
												Number: 80,
											},
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

	// Create the ingress resource in the Kubernetes cluster
	_, err = clientset.NetworkingV1().Ingresses("default").Create(c, ingress, metav1.CreateOptions{})
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("error creating the ingress resource: %s", err))
		return
	}
	c.JSON(http.StatusBadRequest, fmt.Errorf("bad request"))
}
