package cri

import (
	"flag"
	"log"

	"study/running"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/klog/v2"
)

func TestGetCriStatus(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := running.Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")

	flag.Parse()
	client, err := GetCriClient("unix:///run/containerd/containerd.sock")
	assert.Nil(t, err)
	klog.Infoln(client)
	status, err := GetCriStatus()
	assert.Nil(t, err)
	klog.Infof("%#v", status.Info["cniconfig"])
	cniConfig, err := GetCniConfig()
	assert.Nil(t, err)
	klog.Infof("%#v", cniConfig)
}

func TestGetCriPodID(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := running.Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")

	flag.Parse()
	client, err := GetCriClient("unix:///run/containerd/containerd.sock")
	assert.Nil(t, err)
	klog.Infoln(client)
	pod, err := GetCriPod("2165591e88119126246475446cf4fabc53f33489b002431a1b9fe29b03d991cd")
	assert.Nil(t, err)
	klog.Infof("%#v", pod.GetStatus())
	// klog.Infof("%#v", pod.Info["info"])
	info, err := GetSandboxInfo("2165591e88119126246475446cf4fabc53f33489b002431a1b9fe29b03d991cd")
	assert.Nil(t, err)
	klog.Infof("%#v", info)
	klog.Infof("%v", info)
}

func TestRunPodSandBox(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := running.Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")

	flag.Parse()
	pod, err := RunPodSandBox()
	assert.Nil(t, err)
	klog.Infof("%#v", pod.GetPodSandboxId())
	klog.Infof("%#v", pod)
	klog.Infof("%v", pod)
}

func TestPullImage(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := running.Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")

	flag.Parse()
	pod, err := PullImage("nginx:alpine")
	assert.Nil(t, err)
	klog.Infof("%#v", pod.ImageRef)
	klog.Infof("%#v", pod)
	klog.Infof("%v", pod)
}
func TestCreateContainer(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := running.Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")

	flag.Parse()
	pod, err := CreateContainer("a377fa3d10f1dfc427b3ffe60fb5e7ed6d5e99d25c6707741c2dec4da4d20df4")
	assert.Nil(t, err)
	klog.Infof("%#v", pod.ContainerId)
	klog.Infof("%#v", pod)
	klog.Infof("%v", pod)
}
func TestStartContainer(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := running.Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")

	flag.Parse()
	pod, err := StartContainer("48cb180a51f4783fb23927080c2a23b2208296d083e4181bc69240d7c1537556")
	assert.Nil(t, err)
	klog.Infof("%#v", pod.String())
	klog.Infof("%#v", pod)
	klog.Infof("%v", pod)
}

func TestStopPodSandBox(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := running.Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")

	flag.Parse()
	pod, err := StopPodSandBox("7db571f6c81692a2e1f086c78c14c73d6a7d5e170abec03a928eac0c64a328a5")
	assert.Nil(t, err)
	klog.Infof("%#v", pod.String())
	klog.Infof("%#v", pod)
	klog.Infof("%v", pod)
}

func TestRemovePodSandBox(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := running.Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")

	flag.Parse()
	pod, err := RemovePodSandBox("7db571f6c81692a2e1f086c78c14c73d6a7d5e170abec03a928eac0c64a328a5")
	assert.Nil(t, err)
	klog.Infof("%#v", pod.String())
	klog.Infof("%#v", pod)
	klog.Infof("%v", pod)
}
