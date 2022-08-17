<p align="center">
  <a href="https://www.parseable.io" target="_blank"><img src="https://raw.githubusercontent.com/parseablehq/.github/main/images/collector-logo.svg" alt="Parseable" width="300" height="100" /></a>
</p>

<p align="center">
  <img src="https://img.shields.io/website?down_message=red&up_color=green&up_message=online&url=https%3A%2F%2Fwww.parseable.io" alt="website status">
  <img src="https://img.shields.io/github/contributors/parseablehq/collector" alt="contributors">
  <img src="https://img.shields.io/github/commit-activity/m/parseablehq/collector" alt="commits activity monthly">
  <a href="https://github.com/parseablehq/collector/stargazers" target="_blank"><img src="https://img.shields.io/github/stars/parseablehq/collector" alt="Github stars"></a>
  <img src="https://img.shields.io/github/license/parseablehq/collector" alt="License">  
  <a href="https://twitter.com/parseableio" target="_blank"><img src="https://img.shields.io/twitter/follow/parseableio" alt="Twitter"></a>
</p>

Parseable Collector is an automatic log collection system for [Parseable](https://github.com/parseablehq/parseable). Collector is designed to _pull_ application logs from specific containers (using Kubernetes API). Collector sends collected logs to Parseable server for storage, querying and further analysis. 

<p align="center">
  <img src="https://raw.githubusercontent.com/parseablehq/.github/main/images/collector-overview.svg#gh-light-mode-only" alt="Parseable Overview" width="800" height="650" />
  <img src="https://raw.githubusercontent.com/parseablehq/.github/main/images/collector-overview-dark.svg#gh-dark-mode-only" alt="Parseable Overview" width="800" height="650" />
</p>

<h1></h1>

## Why another logging agent?

Log collection has been traditionally done with agents installed on each Kubernetes node. This is too much of installation and compute overhead for most of the use cases. In reality, developers simply want to plug their applications to logging and move on. Collector is designed keeping this in mind.

With Collector we set out to achieve the shortest path between application generating logs and developer analysing those logs.

## Get Started 

Collector is designed to run on Kubernetes only. We recommend installing Collector via the [official helm chart](./helm/) available in this repository. Before deploying the 
collector, make sure you understand the configuration.

### Configuration

Collector takes configuration input in yaml format. 

```yaml
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
