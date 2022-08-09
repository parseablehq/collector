<p align="center">
  <a href="https://www.parseable.io" target="_blank"><img src="https://raw.githubusercontent.com/parseablehq/docs/main/static/img/logo.svg" alt="Parseable" width="300" height="100" /> <h1 align="center">Collector</h1></a>
</p>

<p align="center">
  <img src="https://img.shields.io/website?down_message=red&up_color=green&up_message=online&url=https%3A%2F%2Fwww.parseable.io" alt="website status">
  <img src="https://img.shields.io/github/contributors/parseablehq/parseable" alt="contributors">
  <img src="https://img.shields.io/github/commit-activity/m/parseablehq/parseable" alt="commits activity monthly">
  <a href="https://github.com/parseablehq/parseable/stargazers" target="_blank"><img src="https://img.shields.io/github/stars/parseablehq/parseable" alt="Github stars"></a>
  <img src="https://img.shields.io/github/license/parseablehq/parseable" alt="License">  
  <a href="https://twitter.com/parseableio" target="_blank"><img src="https://img.shields.io/twitter/follow/parseableio" alt="Twitter"></a>
</p>

<h4 align="center">
  <a href="https://www.parseable.io" target="_blank">Website</a>
</h4>

Parseable Collector is an automatic log collection system for [Parseable](https://github.com/parseablehq/parseable). Collector is designed to _pull_ application logs from specific containers (using Kubernetes API). Collector sends collected logs to Parseable server for storage, querying and further analysis. 

## Why

Log collection has been traditionally done with agents installed on each Kubernetes node. This is too much of installation and compute overhead for most of the use cases. In reality, developers simply want to plug their applications to logging and move on. Collector is designed keeping this in mind.

Our goal with Collector is - Shortest path between application generating logs and developer analysing those logs.

## Get Started 

Collector is designed to run on Kubernetes only. We recommend installing Collector via the [official helm chart](./helm/) available in this repository. Before deploying the 
collector, make sure to 

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
      labels: 
        app: kafka
        namespace: streaming
```
