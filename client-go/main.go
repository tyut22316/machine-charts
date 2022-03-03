/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	informers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

//const metaCRD = `
//apiVersion: "machart.machart.rd/v1"
//kind: Mcmachines
//metadata:
//  name: mcmachines-sample
//spec:
//  # Add fields here
//  cpu: "4"
//`

//func GetGVRdyClient(config *rest.Config, gvk *schema.GroupVersionKind, namespace string) (dr dynamic.ResourceInterface, err error) {
//
//	//gvk
//	//discover
//	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
//	if err != nil {
//		return
//	}
//	// 获取GVK GVR 映射
//	mapperGVRGVK := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))
//
//	// 根据资源GVK 获取资源的GVR GVK映射
//	resourceMapper, err := mapperGVRGVK.RESTMapping(gvk.GroupKind(), gvk.Version)
//	if err != nil {
//		return
//	}
//	// 创建动态客户端
//	dynamicClient, err := dynamic.NewForConfig(config)
//	if err != nil {
//		return
//	}
//
//	if resourceMapper.Scope.Name() == meta.RESTScopeNameNamespace {
//		// 获取gvr对应的动态客户端
//		dr = dynamicClient.Resource(resourceMapper.Resource).Namespace(namespace)
//	} else {
//		// 获取gvr对应的动态客户端
//		dr = dynamicClient.Resource(resourceMapper.Resource)
//	}
//	return
//}

type ThingSpec struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// getGVR :- gets GroupVersionResource for dynamic client
func getGVR(group, version, resource string) schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
}

func main() {

	// 配置 k8s 集群外 kubeconfig 配置文件
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	//在 kubeconfig 中使用当前上下文环境，config 获取支持 url 和 path 方式
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 根据指定的 config 创建一个新的 clientset
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//同步时间间隔
	sharedInformerFactory := informers.NewSharedInformerFactory(clientset, time.Minute)
	stopCh := make(chan struct{})
	defer close(stopCh)
	//sharedInformerFactory.Start(stopCh)

	nodeinformer := sharedInformerFactory.Core().V1().Nodes().Informer()
	nodeinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//pod资源对象创建的时候出发的回调方法
		AddFunc: func(obj interface{}) {
			//	obja := obj.(v1.Object)
			//fmt.Println(obja)
			//fmt.Println("来了老弟")
		},
		//更新回调
		UpdateFunc: func(oldObj, newObj interface{}) {

			obja := oldObj.(*v1.Node)
			//fmt.Println("旧的", obja)
			//obja.String()
			//fmt.Println("输出结果", obja.GetLabels())
			status := obja.Status.Allocatable.Cpu()
			fmt.Println("旧的取这个", status)

			newobja := oldObj.(*v1.Node)
			newstatus := newobja.Status.Allocatable.Cpu()
			fmt.Println("新的取这个", newstatus)
			//
			//things := make([]ThingSpec, 1)
			//things[0].Op = "replace"
			//things[0].Path = "/spec/cpu"
			//things[0].Value = "5"
			//
			//patchBytes, err := json.Marshal(things)
			if err != nil {
				panic(err.Error())
			}
			newclient, err := dynamic.NewForConfig(config)
			if err != nil {
				panic(err)
			}
			//background := context.
			fmt.Println("来了老弟")
			//gvr := getGVR("machart.machart.rd", "v1", "mcmachines")
			var gvr = schema.GroupVersionResource{
				Group:    "machart.machart.rd",
				Version:  "v1",
				Resource: "mcmachines",
			}

			patchData := []byte(`{"spec": {"cpu": "5"}}`)
			//打字打错了
			if _, err := newclient.Resource(gvr).Namespace("default").Patch(context.Background(), "mcmachines-sample", types.MergePatchType, patchData, metav1.PatchOptions{}); err != nil {
				panic(err)
			}

			//myclient := clientset.RESTClient()
			//
			//result, err := myclient.Patch(api.JSONPatchType).
			//	Namespace(api.NamespaceDefault).
			//	Resource("mcmachines").
			//	Name("mcmachines-sample").
			//	Body(patchBytes).
			//	Do().
			//	Get()

			if err != nil {
				panic(err.Error())
			}

			name := types.NamespacedName{Namespace: "default", Name: "mcmachines-sample"}
			fmt.Println("新的取这个", name.String())

			//crdobj := &unstructured.Unstructured{}
			//var gvk *schema.GroupVersionKind
			//
			//var dr dynamic.ResourceInterface
			//
			//_, gvk, err = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode([]byte(metaCRD), nil, crdobj)
			//if err != nil {
			//	panic(fmt.Errorf("failed to get GVK: %v", err))
			//}
			//
			//dr, err = GetGVRdyClient(config, gvk, crdobj.GetNamespace())
			//
			//if err != nil {
			//	panic(fmt.Errorf("failed to get dr: %v", err))
			//}
			//
			////创建
			//
			//if err != nil {
			//	//panic(fmt.Errorf("Create resource ERROR: %v", err))
			//	log.Print(err)
			//}
			//
			//// 查询
			//
			////缺少上下文  少了return
			//background := context.Background()
			//crdobjGET, err := dr.Get(background, crdobj.GetName(), metav1.GetOptions{})
			//if err != nil {
			//	panic(fmt.Errorf("select resource ERROR: %v", err))
			//}
			//fmt.Printf("GET: ", crdobjGET)

			//newja := newObj.(v2.Node)
			//newLabels := newja.GetLabels()

			//obja := oldObj.(v1.Object)
			////fmt.Println("旧的", obja)
			//oldLabels := obja.GetLabels()
			//newja := newObj.(v1.Object)
			//newLabels := newja.GetLabels()

			//if reflect.DeepEqual(oldLabels, newLabels) {
			//	fmt.Println("相同没变化", obja.GetLabels())
			//} else {
			//	fmt.Println("变化了", newja.GetLabels()["test"])
			//}

			//fmt.Println("新的", new)
			//fmt.Println("检测到edit")
			//	podinfomer.()

		},
		//删除回调
		DeleteFunc: func(obj interface{}) {

		},
	})
	nodeinformer.Run(stopCh)

	/*
		//通过sharedinfomers 我们获取到pod的informer
		podinfomer := sharedInformerFactory.Core().V1().Pods().Informer()
		//为pod informer添加 controller的handlerfunc  触发回调函数之后 会通过addch 传给nextCh 管道然后调用controller的对应的handler来做处理
		podinfomer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			//pod资源对象创建的时候出发的回调方法
			AddFunc: func(obj interface{}) {
				//	obja := obj.(v1.Object)
				//fmt.Println(obja)
				//fmt.Println("来了老弟")
			},
			//更新回调
			UpdateFunc: func(oldObj, newObj interface{}) {

				obja := oldObj.(v1.Object)
				//fmt.Println("旧的", obja)
				oldLabels := obja.GetLabels()
				newja := newObj.(v1.Object)
				newLabels := newja.GetLabels()

				if reflect.DeepEqual(oldLabels, newLabels) {
					fmt.Println("相同没变化", obja.GetLabels())
				} else {
					fmt.Println("变化了", newja.GetLabels()["test"])
				}

				//fmt.Println("新的", new)
				//fmt.Println("检测到edit")
				//	podinfomer.()

			},
			//删除回调
			DeleteFunc: func(obj interface{}) {

			},
		})
		podinfomer.Run(stopCh)
	*/

	//lister := sharedInformerFactory.Core().V1().Nodes().Lister()
	//返回值是list
	//list, err := lister.List(labels.Everything())
	//list, err := lister.Get("master")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(list)
	/*
		for {
			// 通过实现 clientset 的 CoreV1Interface 接口列表中的 PodsGetter 接口方法 Pods(namespace string)返回 PodInterface
			// PodInterface 接口拥有操作 Pod 资源的方法，例如 Create、Update、Get、List 等方法
			// 注意：Pods() 方法中 namespace 不指定则获取 Cluster 所有 Pod 列表

			nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("There are %d pods in the cluster\n", len(nodes.Items))

			fmt.Println(nodes.String())

	*/

	/*
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// 获取指定 namespace 中的 Pod 列表信息

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		//	namespace := "default"


				namespace := "kube-system"

			pod := "etcd-master"
			//pod := "example-xxxxx"
			_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				fmt.Printf("Error getting pod %s in namespace %s: %v\n",
					pod, namespace, statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
			}
		}
	*/

	time.Sleep(10 * time.Second)
}
