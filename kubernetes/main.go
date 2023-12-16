package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"
)

type CreateDeploymentDto struct {
	Image  string `json:"image"`
	Config string `json:"config"`
	Name   string `json:"name"`
}

type WhanosResources struct {
	Limits   *apiv1.ResourceList `json:"limits,omitempty"`
	Requests *apiv1.ResourceList `json:"requests,omitempty"`
}

type WhanosDeployment struct {
	Replicas  *int32           `yaml:"replicas,omitempty"`
	Resources *WhanosResources `yaml:"resources,omitempty"`
	Ports     *[]int32         `yaml:"ports,omitempty"`
}

type WhanosConfig struct {
	Deployment WhanosDeployment `yaml:"deployment"`
}

func getKubeconfig() string {
	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}
	return filepath.Join("/", ".kube", "config")
}

func getClientset(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func parseConfig(config string) (cfg WhanosConfig) {
	if err := yaml.Unmarshal([]byte(config), &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func main() {
	var kubeconfig = getKubeconfig()
	var clientset = getClientset(kubeconfig)
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	r := gin.Default()
	r.POST("/deployments", func(c *gin.Context) {
		var payload CreateDeploymentDto
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "malformed payload",
				"description": err.Error(),
			})
			return
		}
		value, err := base64.StdEncoding.DecodeString(payload.Config)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "malformed payload",
				"description": err.Error(),
			})
		}
		var config = parseConfig(string(value))
		var ports = []apiv1.ContainerPort{}
		for idx := range *config.Deployment.Ports {
			ports = append(ports, apiv1.ContainerPort{
				ContainerPort: (*config.Deployment.Ports)[idx],
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})

		println("Starting to pull image " + payload.Image)
		cmd := exec.Command("minikube", "cache", "add", payload.Image)
		cmd.Run()
		println("Finished pulling " + payload.Image)

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: payload.Name,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: config.Deployment.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": payload.Name,
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": payload.Name,
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  payload.Name,
								Image: payload.Image,
								Ports: ports,
								Resources: apiv1.ResourceRequirements{
									Limits:   *config.Deployment.Resources.Limits,
									Requests: *config.Deployment.Resources.Requests,
								},
								ImagePullPolicy: apiv1.PullNever,
							},
						},
					},
				},
			},
		}
		result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	})
	r.Run(":3030")
}
