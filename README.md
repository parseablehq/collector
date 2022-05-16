# Kubec

Kubec (pronounced "cube - ek") is a log collection system for [Parseable server](https://github.com/parseablehq/parseable). Kubec is designed to _pull_ application logs from specific containers (using Kubernetes API). Kubec sends collected logs to Parseable server for storage, querying and further analysis. 

# Why

Log collection has been traditionally done with agents installed on each Kubernetes node. This is too much of installation and compute overhead for most of the use cases. In reality, developers simply want to plug their applications to logging and move on. Kubec is designed keeping this in mind.

Our goal with Kubec is - Shortest path between application generating logs and developer analysing those logs.

## Get Started 
Kubec is meant to be deployed only on Kubernetes. We recommend installing Kubec via the [official helm chart](./helm/) available in this repository.

### Configuration

```yaml
parseable:
  logStreams:
    - name: backend
      collectionInterval: 3s
      collectFrom: 
        namespace: streaming
        podSelector: 
          app: kafka
      tags: 
        app: kafka
        namespace: streaming
```
