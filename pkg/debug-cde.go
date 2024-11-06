/*
Copyright 2024 Mario Loriedo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"context"
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	debugIDEExample = `
	# Create a copy of the Pod <pod-name> with an extra sidecar container running an IDE and including the <repository-url> source code
	%[1]s debug-ide <pod-name> --image <debug-image> --git-repository <repository-url>`

	errNoContext = fmt.Errorf("no context is currently set, use %q to select a new one", "kubectl config use-context <context>")
)

const (
	defaultIdeReference = "https://eclipse-che.github.io/che-plugin-registry/main/v3/plugins/che-incubator/che-code/latest/devfile.yaml"
	defaultDebugImage   = "quay.io/devfile/universal-developer-image:ubi8-latest"
)

// DebugIDEOptions provides information required to update
// the current context on a user's KUBECONFIG
type DebugIDEOptions struct {
	configFlags *genericclioptions.ConfigFlags

	resultingContext     *api.Context
	resultingContextName string

	userSpecifiedCluster   string
	userSpecifiedContext   string
	userSpecifiedAuthInfo  string
	userSpecifiedNamespace string

	targetPodName       string
	targetPodContainers []ContainerInfo

	debugImage     string
	copyToPodName  string
	shareProcesses bool
	gitRepository  string
	ideReference   string
	rawConfig      api.Config
	args           []string

	genericiooptions.IOStreams
}

// NewDebugIDEOptions provides an instance of DebugIDEOptions with default values
func NewDebugIDEOptions(streams genericiooptions.IOStreams) *DebugIDEOptions {
	return &DebugIDEOptions{
		configFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// NewCmdDebugIDE provides a cobra command wrapping DebugIDEOptions
func NewCmdDebugIDE(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewDebugIDEOptions(streams)

	cmd := &cobra.Command{
		Use:          "debug-cde [pod] [flags]",
		Short:        "Create a copy of a Pod and add a Cloud Development Environment to debug it.",
		Example:      fmt.Sprintf(debugIDEExample, "kubectl"),
		SilenceUsage: true,
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "kubectl debug-ide",
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	//cmd.Flags().BoolVar(&o.listNamespaces, "list", o.listNamespaces, "if true, print the list of all namespaces in the current KUBECONFIG")
	cmd.Flags().StringVar(&o.ideReference, "ide", defaultIdeReference, "URI to the devfile with the IDE definition")
	cmd.Flags().StringVar(&o.debugImage, "image", defaultDebugImage, "Image of the debug sidecar container")
	cmd.Flags().StringVar(&o.gitRepository, "git-repository", o.gitRepository, "URL of the git repository with the source code of the application we want to debug")
	cmd.Flags().StringVar(&o.copyToPodName, "copy-to", o.copyToPodName, "Name of the new Pod, copy of the target Pod")
	cmd.Flags().BoolVar(&o.shareProcesses, "share-processes", o.shareProcesses, "If true, enable process namespace sharing in the copy")
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete sets all information required for creating a DevWorkspace
func (o *DebugIDEOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("cannot omit the target pod to debug")
	}

	if len(args) > 1 {
		return fmt.Errorf("cannot specify more than one pod (args number is %d)", len(args))
	}

	o.targetPodName = args[0]

	if len(args) > 1 {
		return fmt.Errorf("cannot specify more than one pod")
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return fmt.Errorf("ClientConfig error: %v", err)
	}

	namespace, _, _ := kubeConfig.Namespace()
	clientset, err := kubernetes.NewForConfig(config)
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), o.targetPodName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return fmt.Errorf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		return fmt.Errorf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("🎯 found the target pod %s in namespace %s.\n", pod.Name, namespace)
	}

	podContainersNum := len(pod.Spec.Containers)
	o.targetPodContainers = make([]ContainerInfo, podContainersNum)
	for i := 0; i < podContainersNum; i++ {
		o.targetPodContainers[i].name = pod.Spec.Containers[i].Name
		o.targetPodContainers[i].image = pod.Spec.Containers[i].Image
		o.targetPodContainers[i].memoryLimit = pod.Spec.Containers[i].Resources.Limits.Memory().String()
		o.targetPodContainers[i].cpuLimit = pod.Spec.Containers[i].Resources.Limits.Cpu().String()
		containerPorts := len(pod.Spec.Containers[i].Ports)
		o.targetPodContainers[i].endpoints = make([]ContainerEndpoint, podContainersNum)
		for j := 0; j < containerPorts; j++ {
			portName := pod.Spec.Containers[i].Ports[j].Name
			portNumber := int(pod.Spec.Containers[i].Ports[j].ContainerPort)
			if portName == "" {
				portName = "port" + strconv.Itoa(portNumber)
			}
			o.targetPodContainers[i].endpoints[i].name = portName
			o.targetPodContainers[i].endpoints[i].targetPort = portNumber
		}
	}

	o.userSpecifiedNamespace, err = cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}

	o.userSpecifiedContext, err = cmd.Flags().GetString("context")
	if err != nil {
		return err
	}

	o.userSpecifiedCluster, err = cmd.Flags().GetString("cluster")
	if err != nil {
		return err
	}

	o.userSpecifiedAuthInfo, err = cmd.Flags().GetString("user")
	if err != nil {
		return err
	}

	currentContext, exists := o.rawConfig.Contexts[o.rawConfig.CurrentContext]
	if !exists {
		return errNoContext
	}

	o.resultingContext = api.NewContext()
	o.resultingContext.Cluster = currentContext.Cluster
	o.resultingContext.AuthInfo = currentContext.AuthInfo

	// if a target context is explicitly provided by the user,
	// use that as our reference for the final, resulting context
	if len(o.userSpecifiedContext) > 0 {
		o.resultingContextName = o.userSpecifiedContext
		if userCtx, exists := o.rawConfig.Contexts[o.userSpecifiedContext]; exists {
			o.resultingContext = userCtx.DeepCopy()
		}
	}

	// override context info with user provided values
	o.resultingContext.Namespace = o.userSpecifiedNamespace

	if len(o.userSpecifiedCluster) > 0 {
		o.resultingContext.Cluster = o.userSpecifiedCluster
	}
	if len(o.userSpecifiedAuthInfo) > 0 {
		o.resultingContext.AuthInfo = o.userSpecifiedAuthInfo
	}

	// generate a unique context name based on its new values if
	// user did not explicitly request a context by name
	if len(o.userSpecifiedContext) == 0 {
		o.resultingContextName = generateContextName(o.resultingContext)
	}

	return nil
}

func generateContextName(fromContext *api.Context) string {
	name := fromContext.Namespace
	if len(fromContext.Cluster) > 0 {
		name = fmt.Sprintf("%s/%s", name, fromContext.Cluster)
	}
	if len(fromContext.AuthInfo) > 0 {
		cleanAuthInfo := strings.Split(fromContext.AuthInfo, "/")[0]
		name = fmt.Sprintf("%s/%s", name, cleanAuthInfo)
	}

	return name
}

// Validate ensures that all required arguments and flag values are provided
func (o *DebugIDEOptions) Validate() error {
	if len(o.rawConfig.CurrentContext) == 0 {
		return errNoContext
	}
	return nil
}

// Run lists all available namespaces on a user's KUBECONFIG or updates the
// current context based on a provided namespace.
// Apply a DevWorkspace object
func (o *DebugIDEOptions) Run() error {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return fmt.Errorf("ClientConfig error: %v", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("dynamic client creation failed: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("client creation failed: %v", err)
	}

	// Generate the DevWorkspace
	dw, err := generate(*o)
	if err != nil {
		return fmt.Errorf("Error generating devworkspace: %v", err)
	}

	// Convert the DevWorkspace to an Unstructured
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&dw)
	if err != nil {
		return fmt.Errorf("Found error while coverting resource to unstructured err - %v", err)
	}
	unstructuredResource := &unstructured.Unstructured{Object: obj}

	// Automatically get the GroupVersionResouce for the DevWorkspace
	gvk := dw.GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return fmt.Errorf("restmapper error: %v", err)
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)
	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return fmt.Errorf("RESTMapping error: %v", err)
	}

	// Check if the DevWorksapce already exist
	namespace, _, _ := kubeConfig.Namespace()
	result, err := dynClient.Resource(mapping.Resource).Namespace(namespace).Get(
		context.TODO(),
		dw.Name,
		metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("Error checking for the existance of a DevWorkspace named %s: %v", dw.Name, err)
	}
	if result != nil && result.Object != nil {
		return fmt.Errorf("A DevWorkspace named %s already exist. Delete it to create a new one.", dw.Name)
	}

	// Create the DevWorkspace
	result, err = dynClient.Resource(mapping.Resource).Namespace(namespace).Create(
		context.TODO(),
		unstructuredResource,
		metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Error creating custom resource: %v\n", err)
	}
	dwName := result.Object["metadata"].(map[string]interface{})["name"].(string)
	fmt.Printf("⌨️ created devworkspace %s in namespace %s.\n", dwName, namespace)

	// Get the deployment name
	dwUnstruct, err := dynClient.Resource(mapping.Resource).Namespace(namespace).Get(
		context.TODO(),
		dwName,
		metav1.GetOptions{})
	deploymentName := dwUnstruct.Object["status"].(map[string]interface{})["devworkspaceId"].(string)

	// Wait for deployment status condition available == true
	fmt.Printf("⏳ waiting for the deployment %s to be available...", deploymentName)
	var d *appv1.Deployment
	available := false
	timeout := 30
	for i := 0; i < timeout; i++ {
		d, err = clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		for _, condition := range d.Status.Conditions {
			if condition.Type == appv1.DeploymentAvailable {
				if condition.Status == "True" {
					available = true
					break
				} else if condition.Status == "False" {
					fmt.Print(".")
				}
			}
		}
		if available {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !available {
		return fmt.Errorf("the deployment %s is not available after %d seconds (%v)", deploymentName, timeout, d.Status.Conditions)
	}

	fmt.Printf("done\n")

	// Get the Pod name
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(),
		metav1.ListOptions{
			LabelSelector: "controller.devfile.io/devworkspace_name=" + dwName,
		})
	if err != nil {
		return fmt.Errorf("Error listing devworkspace pods: %v", err)
	}
	if len(pods.Items) == 0 {
		return fmt.Errorf("No devworkspace pods found for DevWorkspace: %s", dwName)
	}
	if len(pods.Items) > 1 {
		return fmt.Errorf("More than one devworkspace pods found for DevWorkspace: %s", dwName)
	}
	podName := pods.Items[0].Name

	// Wait for pod status condition ready == true
	fmt.Printf("🥑 waiting for the pod %s to be ready...", podName)
	var p *corev1.Pod
	ready := false
	podReadinessTimeout := 30
	for i := 0; i < podReadinessTimeout; i++ {
		p, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		for _, condition := range p.Status.Conditions {
			if condition.Type == corev1.PodReady {
				if condition.Status == "True" {
					ready = true
					break
				} else if condition.Status == "False" {
					fmt.Print(".")
				}
			}
		}
		if ready {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !ready {
		return fmt.Errorf("the pod %s is not ready after %d seconds (%v)", podName, podReadinessTimeout, p.Status.Conditions)
	}

	fmt.Printf("done\n")

	// Retrieve IDE URL
	dwUnstruct, err = dynClient.Resource(mapping.Resource).Namespace(namespace).Get(
		context.TODO(),
		dwName,
		metav1.GetOptions{})
	dwMainURL := dwUnstruct.Object["status"].(map[string]interface{})["mainUrl"].(string)
	fmt.Printf("🐞 click on the following link ⬇️ and start debugging\n\n")
	fmt.Printf("%s\n", dwMainURL)
	return nil
}
