[![Build Status](https://travis-ci.org/containerum/containerum.svg?branch=master)](https://travis-ci.org/containerum/containerum) [![HitCount](http://hits.dwyl.com/containerum/containerum.svg)](http://hits.dwyl.com/containerum/containerum) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Coverage](coverage_badge.png)]('https://github.com/jpoles1/gopherbadger)

# Embark to Containerum!

![Containerum logo](../logo.svg)

A tool to quickly install [Containerum](https://containerum.com/software/) on a Kubernetes cluster.

## Getting Started

### Prerequisites
Embark requires installed and configured Kubernetes. If you need to install a Kubernetes cluster, check out these articles:
  + [4 ways to bootstrap a Kubernetes cluster](https://medium.com/containerum/4-ways-to-bootstrap-a-kubernetes-cluster-de0d5150a1e4)
  + [How to deploy Kubernetes and Containerum on Digital Ocean](https://medium.com/containerum/how-to-deploy-kubernetes-and-containerum-on-digital-ocean-eca93e6b4d26)
  + [Installing Kubernetes from packages](https://docs.kdc.containerum.com/installation/packages/)

### Installation
Just run 
```bash 
kubect create -f -
```
and then paste 

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: containerum/embark
  labels:
    app: containerum
    name: embark
spec:
  serviceAccountName: containerum-embark
  selector: {}
  template:
    metadata:
      labels:
        app: containerum
        name: embark
    spec:
      containers:
        - name: embark
          image: containerum/embark
          imagePullPolicy: Always
          volumeMounts:
            - name: kube
              mountPath: /etc/kube
      volumes:
        - name: kube
          configMap:
            defaultMode: 420
            name: kube-config
      restartPolicy: Never
```
This will launch Containerum installation job. Then wait until installation is complete:
```bash
kubectl wait --for=condition=complete containerum/embark
```

## Running tests

__Run unit tests__
```bash
go test ./...
```

__Run integration tests__
You need Kubernetes cluster up and running and a kubectl config located at `~/.kube/config`.

**⛔️ !!DANGEROUS!!! ⛔️**

**DUE TO THE FACT THAT TESTS CAN CREATE AND DELETE OBJECTS IN THE KUBERNETES CLUSTER, YOU MUST BE SURE THAT THE CONFIGURATION FILE GIVES ACCESS ONLY TO THE TEST CLUSTER, SO NO PRODUCTION SERVICES WILL BE DESTROYED OR CORRUPTED!**

```bash
go test -tags="IntegrationTests" ./...  
```

## Contributing

We welcome any help from the open source community. To submit your contributions, fork the project you want to contribute to (e.g. permissions, auth, etc.), commit changes and create a pull request to the develop branch. We will review the changes and include them to the project. Read more about contributing in this document.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/containerum/containerum/tags).

## Authors

* **Paul Petrukhin** - *Initial work* - [ninedraft](https://github.com/ninedraft)

See also the list of [contributors](https://github.com/containerum/containerum/contributors) who participated in this project.

## License

This project is licensed under the Apache License - see the [LICENSE](../LICENSE) file for details

## Acknowledgments

* Thanks to Kubernetes team for providing well documented code for such components as `apimachinery` and `client-go`
* Thanks to Helm gang for excelent quality of their product
