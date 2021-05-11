package cri

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/containerd/containerd/pkg/cri/server"
	gocni "github.com/containerd/go-cni"
	"google.golang.org/grpc"

	// v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/klog/v2"
)

const (
	tokenFile                 = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	rootCAFile                = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	DefautContainerdSock      = "unix:///run/containerd/containerd.sock"
	DefautImageContainerdSock = "unix:///run/containerd/containerd.sock"
	CriSocketPath             = "unix:///run/containerd/containerd.sock"
)

var criclient v1.RuntimeServiceClient
var crimageclient v1.ImageServiceClient
var sockAddr string
var imageSockAddr string

// Initialize the CRI APIs required
func GetClientAPIs(criSocketPath string) (v1.RuntimeServiceClient, v1.ImageServiceClient, error) {
	// Set up a connection to the server.
	conn, err := GetClientConnection(criSocketPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %v", err)
	}
	rc := v1.NewRuntimeServiceClient(conn)
	if rc == nil {
		return nil, nil, fmt.Errorf("failed to create runtime service client")
	}
	ic := v1.NewImageServiceClient(conn)
	if ic == nil {
		return nil, nil, fmt.Errorf("failed to create image service client")
	}
	return rc, ic, err
}

// Initialize CRI client connection
func GetClientConnection(criSocketPath string) (*grpc.ClientConn, error) {
	//conn, err := grpc.Dial(criSocketPath, grpc.WithInsecure(), grpc.WithTimeout(10*time.Second), grpc.WithDialer(UnixDialer))
	conn, err := grpc.Dial(criSocketPath, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	return conn, nil
}

func UnixDialer(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, timeout)
}

func GetCriClient(target string) (v1.RuntimeServiceClient, error) {
	if criclient != nil {
		return criclient, nil
	}

	if len(target) == 0 {
		target = DefautContainerdSock
	}

	sockAddr = target
	gc, err := grpc.Dial(sockAddr, grpc.WithInsecure())
	if err != nil {
		// klog.Errorln(err)
		return nil, err
	}
	klog.Infoln(gc.Target())

	client := v1.NewRuntimeServiceClient(gc)
	criclient = client
	return criclient, nil
}

func GetCriImageClient(target string) (v1.ImageServiceClient, error) {
	if criclient != nil {
		return crimageclient, nil
	}

	if len(target) == 0 {
		target = DefautImageContainerdSock
	}

	imageSockAddr = target
	gc, err := grpc.Dial(imageSockAddr, grpc.WithInsecure())
	if err != nil {
		// klog.Errorln(err)
		return nil, err
	}
	klog.Infoln(gc.Target())

	client := v1.NewImageServiceClient(gc)
	crimageclient = client
	return crimageclient, nil
}

func GetCniConfig() (*gocni.ConfigResult, error) {
	status, err := GetCriStatus()
	if err != nil {
		return nil, err
	}
	var cniConfig gocni.ConfigResult
	err = json.Unmarshal([]byte(status.Info["cniconfig"]), &cniConfig)
	if err != nil {
		return nil, err
	}

	return &cniConfig, nil
}

func GetCriStatus() (*v1.StatusResponse, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}

	req := v1.StatusRequest{Verbose: true}
	status, err := client.Status(context.Background(), &req, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}

	return status, nil
}

func GetSandboxInfo(id string) (*server.SandboxInfo, error) {
	status, err := GetCriPod(id)
	if err != nil {
		return nil, err
	}
	var sandboxInfo server.SandboxInfo
	err = json.Unmarshal([]byte(status.Info["info"]), &sandboxInfo)
	if err != nil {
		return nil, err
	}

	return &sandboxInfo, nil
}

//删掉沙箱
func RemovePodSandBox(id string) (*v1.RemovePodSandboxResponse, error) {

	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}
	request := &v1.RemovePodSandboxRequest{PodSandboxId: id}
	//生成沙箱
	pod, err := client.RemovePodSandbox(context.Background(), request, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

//停止沙箱
func StopPodSandBox(id string) (*v1.StopPodSandboxResponse, error) {

	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}
	request := &v1.StopPodSandboxRequest{PodSandboxId: id}
	pod, err := client.StopPodSandbox(context.Background(), request, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

//拉镜像
func PullImage(image string) (*v1.PullImageResponse, error) {

	client, err := GetCriImageClient(sockAddr)
	if err != nil {
		return nil, err
	}

	request := &v1.PullImageRequest{
		Image: &v1.ImageSpec{
			Image: image,
		},
	}
	pod, err := client.PullImage(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

//创建容器
func CreateContainer(pId string) (*v1.CreateContainerResponse, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}
	podConfig, _ := GeneratePodSandboxConfig()
	config, _ := GenerateContainerConfig()
	request := &v1.CreateContainerRequest{
		PodSandboxId:  pId,
		Config:        config,
		SandboxConfig: podConfig,
	}
	pod, err := client.CreateContainer(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

//运行容器
func StartContainer(cId string) (*v1.StartContainerResponse, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}
	request := &v1.StartContainerRequest{
		ContainerId: cId,
	}
	pod, err := client.StartContainer(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

//建立沙箱
func RunPodSandBox() (*v1.RunPodSandboxResponse, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}
	//生成沙箱
	config, _ := GeneratePodSandboxConfig()
	request := &v1.RunPodSandboxRequest{Config: config}
	pod, err := client.RunPodSandbox(context.Background(), request, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func GetCriPod(id string) (*v1.PodSandboxStatusResponse, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}
	// SandboxInfo
	req := v1.PodSandboxStatusRequest{PodSandboxId: id, Verbose: true}
	pod, err := client.PodSandboxStatus(context.Background(), &req, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func GetCriPodByName(namespace, podname string) (*v1.PodSandbox, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}

	filter := v1.PodSandboxFilter{}
	filter.LabelSelector = make(map[string]string)
	filter.LabelSelector["io.kubernetes.pod.namespace"] = namespace
	filter.LabelSelector["io.kubernetes.pod.name"] = podname
	filter.State = &v1.PodSandboxStateValue{State: v1.PodSandboxState_SANDBOX_READY}
	req := v1.ListPodSandboxRequest{Filter: &filter}
	pods, err := client.ListPodSandbox(context.Background(), &req, grpc.EmptyCallOption{})
	// klog.Infoln(pods)
	if err != nil {
		klog.Errorln(err)
		return nil, err
	}
	for _, pod := range pods.Items {
		if pod.State == v1.PodSandboxState_SANDBOX_READY {
			return pod, nil
		}
	}

	podlist, err := ListCriPod()
	if err != nil {
		klog.Errorln(err)
		return nil, err
	}
	for _, pod := range podlist.Items {
		if pod.Metadata.Namespace == namespace && pod.Metadata.Name == podname {
			return pod, nil
		}
	}

	return nil, fmt.Errorf("pod not found, %v/%v", namespace, podname)
}

func GetCriPodByLabels(namespace string, labels map[string]string) (*v1.PodSandbox, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}

	filter := v1.PodSandboxFilter{}
	filter.LabelSelector = make(map[string]string)
	filter.LabelSelector["io.kubernetes.pod.namespace"] = namespace
	for k, v := range labels {
		filter.LabelSelector[k] = v
	}
	filter.State = &v1.PodSandboxStateValue{State: v1.PodSandboxState_SANDBOX_READY}
	req := v1.ListPodSandboxRequest{Filter: &filter}
	pods, err := client.ListPodSandbox(context.Background(), &req, grpc.EmptyCallOption{})
	// klog.Infoln(pods)
	if err != nil {
		klog.Errorln(err)
		return nil, err
	}
	for _, pod := range pods.Items {
		if pod.State == v1.PodSandboxState_SANDBOX_READY {
			return pod, nil
		}
	}

	return nil, fmt.Errorf("pod not found, %v %v", namespace, labels)
}

func ListCriPodByNamespace(namespace string) (*v1.ListPodSandboxResponse, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}

	filter := v1.PodSandboxFilter{}
	filter.LabelSelector = make(map[string]string)
	filter.LabelSelector["io.kubernetes.pod.namespace"] = namespace
	filter.State = &v1.PodSandboxStateValue{State: v1.PodSandboxState_SANDBOX_READY}
	req := v1.ListPodSandboxRequest{Filter: &filter}

	return client.ListPodSandbox(context.Background(), &req, grpc.EmptyCallOption{})
}

func ListCriPod() (*v1.ListPodSandboxResponse, error) {
	client, err := GetCriClient(sockAddr)
	if err != nil {
		return nil, err
	}

	filter := v1.PodSandboxFilter{}
	filter.LabelSelector = make(map[string]string)
	filter.State = &v1.PodSandboxStateValue{State: v1.PodSandboxState_SANDBOX_READY}
	req := v1.ListPodSandboxRequest{Filter: &filter}

	return client.ListPodSandbox(context.Background(), &req, grpc.EmptyCallOption{})
}

// Greate CRI PodSandboxConfig from the Pod spec
// TODO: This is probably incomplete
func GeneratePodSandboxConfig() (*v1.PodSandboxConfig, error) {

	podUID := "hdishd83djaidwnduwk28bcsb"
	linux := &v1.LinuxPodSandboxConfig{}
	labels := make(map[string]string)
	labels["lys/cluster"] = "cluster-001"
	labels["lys/node"] = "node-001"
	config := &v1.PodSandboxConfig{
		Metadata: &v1.PodSandboxMetadata{
			Name:      "nginx-sandbox",
			Namespace: "default",
			Uid:       podUID,
			Attempt:   1,
		},
		Labels:       labels,
		LogDirectory: "/tmp",
		Linux:        linux,
	}
	return config, nil
}

// Generate the CRI ContainerConfig from the Pod and container specs
// TODO: Probably incomplete
func GenerateContainerConfig() (*v1.ContainerConfig, error) {
	// TODO: Probably incomplete
	config := &v1.ContainerConfig{
		Metadata: &v1.ContainerMetadata{
			Name:    "busybox",
			Attempt: 1,
		},
		Image:   &v1.ImageSpec{Image: "busybox"},
		Command: []string{"top"},
		LogPath: "busybox/0.log",
	}
	return config, nil
}
