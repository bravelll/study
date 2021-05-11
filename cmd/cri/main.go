// package main1

// import (
// 	"context"
// 	"log"

// 	"google.golang.org/grpc"

// 	// v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
// 	v1 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
// 	"k8s.io/klog/v2"
// )

// const (
// 	tokenFile            = "/var/run/secrets/kubernetes.io/serviceaccount/token"
// 	rootCAFile           = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
// 	DefautContainerdSock = "unix:///run/containerd/containerd.sock"
// )

// var criclient v1.RuntimeServiceClient
// var sockAddr string

// func main() {
// 	if pod, err := CreateCriPod("nginx-sandbox"); err != nil {
// 		log.Fatal(err)
// 	} else {
// 		klog.Infof("%#v", pod.GetPodSandboxId())
// 	}

// }
// func CreateCriPod(id string) (*v1.RunPodSandboxResponse, error) {
// 	client, err := GetCriClient("unix:///run/containerd/containerd.sock")
// 	klog.Infof("%#v", client)
// 	//client.RunPodSandbox()
// 	if err != nil {
// 		return nil, err
// 	}
// 	// SandboxInfo

// 	//req := v1.PodSandboxStatusRequest{PodSandboxId: id, Verbose: true}
// 	//req := v1.RunPodSandboxRequest{}
// 	//labels := make(map[string]string)
// 	//labels["lys/cluster"] = "cluster-001"
// 	//labels["lys/node"] = "node-001"
// 	// req.Config.Metadata.Name = "nginx-sandbox"
// 	// req.Config.Metadata.Namespace = "default"
// 	// req.Config.Metadata.Uid = "hdishd83djaidwnduwk28bcsb"
// 	// req.Config.Metadata.Attempt = 1
// 	// req.Config.LogDirectory = "/tmp"
// 	//req.Config.Labels = labels
// 	//req.Config.Linux = nil

// 	//pod, err := client.PodSandboxStatus(context.Background(), &req, grpc.EmptyCallOption{})
// 	config, err := generatePodSandboxConfig()
// 	request := &v1.RunPodSandboxRequest{Config: config}
// 	pod, err := client.RunPodSandbox(context.Background(), request, grpc.EmptyCallOption{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return pod, nil
// }
// func GetCriClient(target string) (v1.RuntimeServiceClient, error) {
// 	if criclient != nil {
// 		return criclient, nil
// 	}

// 	if len(target) == 0 {
// 		target = DefautContainerdSock
// 	}

// 	sockAddr = target
// 	gc, err := grpc.Dial(target, grpc.WithInsecure())
// 	if err != nil {
// 		// klog.Errorln(err)
// 		return nil, err
// 	}
// 	klog.Infoln(gc.Target())

// 	client := v1.NewRuntimeServiceClient(gc)
// 	criclient = client
// 	return criclient, nil
// }

// // Greate CRI PodSandboxConfig from the Pod spec
// // TODO: This is probably incomplete
// func generatePodSandboxConfig() (*v1.PodSandboxConfig, error) {

// 	podUID := "hdishd83djaidwnduwk28bcsb"
// 	linux := &v1.LinuxPodSandboxConfig{}
// 	labels := make(map[string]string)
// 	labels["lys/cluster"] = "cluster-001"
// 	labels["lys/node"] = "node-001"
// 	config := &v1.PodSandboxConfig{
// 		Metadata: &v1.PodSandboxMetadata{
// 			Name:      "nginx-sandbox",
// 			Namespace: "default",
// 			Uid:       podUID,
// 			Attempt:   1,
// 		},
// 		Labels:       labels,
// 		LogDirectory: "/tmp",
// 		Linux:        linux,
// 	}
// 	return config, nil
// }
