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
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"
)

type NewDeploymentInput struct {
	Image  string `json:"image"`
	Config string `json:"config"`
	Name   string `json:"name"`
}

type CustomResources struct {
	Limits   *apiv1.ResourceList `json:"limits,omitempty"`
	Requests *apiv1.ResourceList `json:"requests,omitempty"`
}

type CustomService struct {
	Type       string `yaml:"type,omitempty"`
	TargetPort int32  `yaml:"targetPort,omitempty"`
}

type CustomDeployment struct {
	Replicas  *int32           `yaml:"replicas,omitempty"`
	Resources *CustomResources `yaml:"resources,omitempty"`
	Ports     *[]int32         `yaml:"ports,omitempty"`
	Service   *CustomService   `yaml:"service,omitempty"`
}

type CustomConfig struct {
	Deployment CustomDeployment `yaml:"deployment"`
}

func GetconfigPath() string {
	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}
	return filepath.Join("/", ".kube", "config")
}

func NewKubernetesClient(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	newKube, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return newKube
}

func getConf(config string) (cfg CustomConfig) {
	if err := yaml.Unmarshal([]byte(config), &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func main() {
	var kubeconfig = GetconfigPath()
	var newKube = NewKubernetesClient(kubeconfig)
	ClientsDep := newKube.AppsV1().Deployments(apiv1.NamespaceDefault)
	r := gin.Default()
	r.POST("/deployments", func(c *gin.Context) {
		var payload NewDeploymentInput
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
		var config = getConf(string(value))
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

		// Create the deployment
		result, err := ClientsDep.Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

		// Always create the service
		svc := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: payload.Name + "-service",
			},
			Spec: apiv1.ServiceSpec{
				Selector: map[string]string{
					"app": payload.Name,
				},
				Ports: []apiv1.ServicePort{
					{
						Protocol:   apiv1.ProtocolTCP,
						Port:       80,
						TargetPort: intstr.FromInt(int(3000)), // Assuming your app runs on port 3000
					},
				},
				Type: apiv1.ServiceTypeClusterIP,
			},
		}

		// Create the service
		_, err = newKube.CoreV1().Services(apiv1.NamespaceDefault).Create(context.TODO(), svc, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created service %q.\n", svc.GetObjectMeta().GetName())
	})
	r.Run(":3030")
}
