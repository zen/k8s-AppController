
AppController [![Build Status](https://travis-ci.org/Mirantis/k8s-AppController.svg?branch=master)](https://travis-ci.org/Mirantis/k8s-AppController) [![Stories in Progress](https://badge.waffle.io/Mirantis/k8s-AppController.png?label=in%20progress&title=In%20Progress)](http://waffle.io/Mirantis/k8s-AppController)
=============
AppController is a pod that you can spawn in your Kubernetes cluster which will take care of your complex deployments for you.

## Basic concepts

AppController uses three basic concepts:

### K8s Objects

AppController interacts with bare Kubernetes objects by creating them (if they are needed by deployment and do not exist yet) and reading their state. The state is used by AppController to ensure that dependencies for other objects are met.

### Dependencies

Dependencies are objects that represent vertices in your deployment graph. You can define them and easily create them with kubectl. Dependencies are ThirdPartyResource which is API extension provided by AppController. It's worth mentioning, that Dependencies can represent dependency between pre-existing K8s object (not orchestrated by AppController) and Resource Definitions, so parts of your deployment graph can depend on objects that were created in your cluster before you even started AppController-aided-deployment. Dependency could have metadata which can contain additional informations about how to determine if it's fulfilled. Current implementation assings this metadata into parent resource and merges it with metadatas from other dependencies definitions so resulting check for resource status uses same values for each dependency. Current implementation does not permit a different values for same metadata key in different dependencies.

Dependency on Replica Set accepts `success_factor` key with stringified percentage integer value of how many replicas should be ready to fulfill the status check.

### Resource Definitions

Resource Definitions are objects that represent Kubernetes Objects that are not yet created, but are part of deployment graph. They store manifests of underlying objects. Objects currently supported by Resource Definitions: (the list is growing steadily)
* Jobs
* Pods
* Services
* Replica Sets
* Pet Sets

Resource Definitions are (the same as Dependencies) ThirdPartyResource API extension.

# Demo
[![asciicast](https://asciinema.org/a/c4ujuq2f8mv1cl16h0u5x0sl1.png)](https://asciinema.org/a/c4ujuq2f8mv1cl16h0u5x0sl1)

[Voice demo from sig-apps meeting](https://youtu.be/BXRToNV4Rdw?t=178)

[Voice demo from kubernetes community meeting](https://youtu.be/NzkoocVeFMQ?t=31)

# Usage

Clone repo:

`git clone https://github.com/Mirantis/k8s-AppController.git`
`cd k8s-AppController`

Run AppController pod:

`kubectl create -f manifests/appcontroller.yaml`

Suppose you have some yaml files with single k8s object definitions (pod and jobs are supported right now). Create AppController ResourceDefintions for them:

`cat path_to_your_pod.yaml | kubectl exec -i k8s-appcontroller wrap <resource_name> | kubectl create -f -`

Create file with dependencies:
```yaml
apiVersion: appcontroller.k8s1/v1alpha1
kind: Dependency
metadata:
  name: dependency-1
parent: pod/<pod_resource_name_1>
child: job/<job_resource_name_2>
---
apiVersion: appcontroller.k8s1/v1alpha1
kind: Dependency
metadata:
  name: dependency-2
parent: pod/<pod_resource_name_2>
child: pod/<pod_resource_name_3>
---
apiVersion: appcontroller.k8s1/v1alpha1
kind: Dependency
metadata:
  name: dependency-3
parent: job/<job_resource_name_1>
child: job/<job_resource_name_3>
---
apiVersion: appcontroller.k8s1/v1alpha1
kind: Dependency
metadata:
  name: dependency-4
parent: replicaset/<replicaset_resource_name_1>
child: job/<job_resource_name_1>
meta:
  success_factor: "80"
```
Load it to k8s:

`kubectl create -f dependencies_file.yaml`

Start appcontroller process:

`kubectl exec k8s-appcontroller ac-run`

You can stop appcontroller process by:

`kubectl exec k8s-appcontroller ac-stop`

# Development

We use github to manage AppController project, so please follow usual github workflow: clone, branch, code, request pull.

We also use [Waffle.io](http://waffle.io/Mirantis/k8s-AppController) to track our progress. Feel free to create issues for functionality we are missing

[![Throughput Graph](https://graphs.waffle.io/Mirantis/k8s-AppController/throughput.svg)](https://waffle.io/Mirantis/k8s-AppController/metrics/throughput)

##Roadmap

We are currently working on support for all Kubernetes objects. Next steps will include:

* better UX
* chaining dependency graphs
* cooperation with [Helm](https://github.com/kubernetes/helm) project

Detailed roadmap is being still being discussed.

##Release schedule

AppController is being released every 3 weeks based on readiness.

### AppController 0.1 alpha

Release date (planned): **14.10.2016**

Planned features:

* Basic functionality completed
* Support for Jobs, Pods, Services, Replica Sets, Pet Sets, Daemon Sets
* Support for dependencies on already existing objects
* Basic CLI
* Reporting on deployment progress
* Basic control over deployment process (start, stop, resume)
* Compliance with Kubernetes incubator repo structure
* Basic graphical graph builder

### AppController 0.2

Release date (planned): 04.11.2016

Planned features:

* Improvements of graph visualisation
* Better failure handling (ability for user to define action of failure)
* Extended logic for handling cases like automatic slave promotion after failure