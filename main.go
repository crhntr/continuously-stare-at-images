package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	k8sConfigPath := os.Getenv("KUBERNETES_CONFIG_PATH")
	if k8sConfigPath == "" {
		home, _ := os.UserHomeDir()
		k8sConfigPath = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", k8sConfigPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	cluster, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	nsWatcher, err := cluster.CoreV1().Namespaces().Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

nsEventLoop:
	for nsEvent := range nsWatcher.ResultChan() {
		namespace, ok := nsEvent.Object.(*corev1.Namespace)
		if !ok {
			continue nsEventLoop
		}

		switch nsEvent.Type {
		default: // ignored
			// fmt.Println("ignored namespace event", nsEvent)

		case watch.Added, watch.Modified:
		}

		podWatcher, err := cluster.CoreV1().Pods(namespace.ObjectMeta.Name).Watch(metav1.ListOptions{})

		if err != nil {
			fmt.Fprintf(os.Stderr, "could not start pod watcher in namespace %s: %s", namespace.ObjectMeta.Name, err.Error())
			continue nsEventLoop
		}

	podEventLoop:
		for podEvent := range podWatcher.ResultChan() {
			switch podEvent.Type {
			default:
				continue podEventLoop
			case watch.Modified, watch.Added:
			}

			pod, ok := podEvent.Object.(*corev1.Pod)
			if !ok {
				continue podEventLoop
			}

			for _, container := range pod.Spec.InitContainers {
				fmt.Println(container.Image)
			}
			for _, container := range pod.Spec.Containers {
				fmt.Println(container.Image)
			}
			for _, container := range pod.Spec.EphemeralContainers {
				fmt.Println(container.Image)
			}
		}
	}
}
