package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Printf("Writing a k8s custom controller")
	kubeconfig := flag.String("kubeconfig", "/home/ubuntu/.kube/config", "location of kubeconfig")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		fmt.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}
	//fmt.Println(config)

	//set of clients used to insteract with the resources from different API versions .
	//If you want to List set of pods you can call coreV1 API version .
	//for deployment resources from apps/v1 from clientset to execute CRUD operations .
	clinetset, err := kubernetes.NewForConfig(config)
	if err != nil {
		//handle error
	}
	//fmt.Println(clinetset)

	ctx := context.Background()

	//list expects a context and ListOptions
	pods, err := clinetset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		//handle error
	}
	//fmt.Println(pods)

	//the return value from pods contains a list of pods
	for _, pod := range pods.Items {
		fmt.Printf("%s", pod.Name)
	}

	deploymets, err := clinetset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	for _, d := range deploymets.Items {
		fmt.Printf("%s", d.Name)
	}

	/** Informers can reduce the load on the API server which can be seen while using watch verb
	because watch verb is going to query the API server a lot of time which might slower down the processes
	on the other hand Informers have a where we can store all the informers after a certain interval of time .
	- If you are experiencing a internet outage informers can handle it with any human intervention where as watch verb requires
	- we have to create informers for every group, version or resource
	  If we want to watch 3 gvr we have to create 3 informers
	- sharedInformerFactory to avoid such heavy load on the API server can get all the resources from all the namespaces
	- once the informer starts we have to initialize the inmemory cache that informer maintains
	  It makes a list call to the API server and then store it in memory
	- subsequent calls are watch insted of List
	- once the initializtion process over the in mermory store provides lister using which we can get and list objects

	- might be a watch request to API server fails informer is going to make another request if the request is no longer available
	- NewSharedInformetFactory after the dedicated memory Informer is going to resync with the cluster state to avoid above situation
	- After the provided time passes the update function is going to call again and the in memory cache was resynced
	- k edit pod -n kube-system kube-scheduler
	  labels:
	    lastResourceVersion: 580
	- Dont update the inmemory cache if you want you can create a deepcopy of the cacahe
	- NewFilterSharedInformerFactory is going to get the resources from test namespace as well as listoptions
	  like lableselector, API version etc .
	- we have registerd some function like addFunc, deleteFun what if the logic written failed in the function defination
	- we can use a queue and enque and add the item into the queue and in another process we have the business Logic
	  It is going to get the element from the work queue and process the things
	- If the process is not completed we can enque the items into the queue because the process is not completed
	  but if the process is completed we can expose a function on the queue Done()

	*/

	ch := make(chan struct{})

	informers := informers.NewSharedInformerFactory(clinetset, 10*time.Minute)

	c := newController(clinetset, informers.Apps().V1().Deployments())
	informers.Start(ch)
	c.run(ch)

	fmt.Println(informers)

}
