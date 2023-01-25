---
theme: gaia
_class: lead
paginate: true
backgroundColor: #fff
backgroundImage: url('https://marp.app/assets/hero-background.svg')

---

<style>
img[alt~="center"] {
	display: block;
	margin: 0 auto;
}

p,li { font-size: 24px }
</style>

# Demo App

<!-- Hi everyone, I'm Francesco Ilario and today I'm going to show a simple microservices application we developed to use as demo application -->
[https://github.com/primaza/demo-app](https://github.com/primaza/demo-app)

---

## What & Why

The demo app is a simple microservice application we developed to use in demos for showing the value added by projects like Service Binding Operator and Primaza when working with complex architectures

The application is a simple e-commerce with no authentication for users' requests. That means that whoever can access the application can perform operations on it

The use-cases implemented allow users to
* navigate the catalog of products
* buy products (one at a time)
* see how many times a product has been bought
* see placed orders

---

## Architecture

<br/>
<div style="columns:2">
<div>

![final-architecture center width:680px](../docs/architecture/architecture-final.png)

</div>
<div/>
<div style=" width: 450px; margin-left:140px">
<br/>

The architecture is composed of the following 4 microservices:
- **front-end**: an NGINX server that provides the React application
- **catalog**: handles the catalog of products
- **orders**: manages orders
- **orders-events-consumer**: reacts to the events fired by the Orders' database and places messages in a Message Brokers

</div>
</div>

---

## Target platform and services

It is possible to run the application on different targets. Also, services can be AWS Cloud services or local development instances of them

In the following the tested configurations:
* Docker compose with local instances
* Minikube with local instances
* Minikube with AWS Services
* Local OpenShift (crc) with local instances
* Local OpenShift (crc) with AWS Services

---

## Use AWS services from kubernetes

For using AWS services from kubernetes the application relies on [Amazon Controllers for Kubernetes (ACK)](https://aws-controllers-k8s.github.io/community/docs/community/overview/)

ACK lets you define and use AWS service resources directly from Kubernetes

With ACK, you can take advantage of AWS-managed services for your Kubernetes applications without needing to define resources outside of the cluster or run services that provide supporting capabilities like databases or message queues within the cluster

---

## Bind to AWS services

We can not use Service Binding Operator (SBO) to bind microservices to AWS services, because their Custom Resource Definitions (CRDs) do not implement the [Service Binding specification](https://github.com/servicebinding/spec)

To enable binding for ACK's resources, we use the [Service Mapper Operator (SMO)](https://github.com/openshift-app-service-poc/service-mapper)

* SMO allows binding of services with SBO without needing changes to CRDs
* You define a set of rules for generating the Service Endpoint Definition (SED)
* SMO Creates a SED and a Service Proxy implementing the Service Binding Specification
* You can now use the Service Proxy to bind to the service

---

# Demo

* Create a Local OpenShift cluster (crc)
* Install Service Binding Operator
* Install AWS Controllers for Kubernetes
* Install Service Mapper Operator
* Install the demo-app
