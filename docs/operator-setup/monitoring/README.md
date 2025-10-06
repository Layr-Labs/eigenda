## Setup monitoring using Docker
<img width="1024" alt="image" src="https://github.com/Layr-Labs/eigenda-operator-setup/assets/354473/c7c3da8d-a488-441c-a70b-00ebf13e5028">

### Quickstart
EigenDA provides a quickstart guide to run the Prometheus, Grafana, and Node exporter stack.
Follow this section for more details for more details. If you want to manually set this up, follow the steps in [Manual Setup](./README.md#manual-setup) section.

In the folder

* Copy the [.env.example](./.env.example) file to `.env` file:
```bash
cp .env.example .env
```
* Make sure your Prometheus config [file](./prometheus.yml) is updated with the metrics port (`NODE_METRICS_PORT`) of the EigenDA node, Default value is set to `9092`.
* Make sure the EigenDA container name is also set correctly in the Prometheus config file. 
You can find that in EigenDA `.env`(`../<network>/.env.example`) file (`MAIN_SERVICE_NAME`)
* Move prometheus.yml file to path: ${HOME}/.eigenlayer/config/ 
* Make sure the location of prometheus file is correct in [.env](./.env.example) file
 
Once correct config is set up, run the following command to start the monitoring stack
```bash
docker compose up -d
```

Your setup should ensure Prometheus is run in the same Docker network as EigenDA. Run the following command for this purpose:
```bash
docker network connect eigenda-network prometheus
```
Note: `eigenda-network` is the name of the network in which EigenDA is running. You can check the network name in EigenDA `.env`(`../<network>/.env.example`) file (`NETWORK_NAME`).

This will make sure `Prometheus` can scrape the metrics from `EigenDA` node.


### Manual Setup
#### Metrics
To check if the metrics are being emitted, run the following command:
```bash
curl http://localhost:<NODE_METRICS_PORT>/metrics
```

You should see something like
```
# HELP eigen_performance_score The performance metric is a score between 0 and 100 and each developer can define their own way of calculating the score. The score is calculated based on the performance of the Node and the performance of the backing services.
# TYPE eigen_performance_score gauge
eigen_performance_score{avs_name="da-node"} 100
# HELP eigen_registered_stakes Operator stake in <quorum> of <avs_name>'s StakeRegistry contract
# TYPE eigen_registered_stakes gauge
eigen_registered_stakes{avs_name="da-node",quorum_name="eth_quorum",quorum_number="0"} 2.654867142483745e+19
# HELP eigen_rpc_request_duration_seconds Duration of json-rpc <method> in second
...
```
#### Prometheus
[Prometheus](https://prometheus.io/download) is being used to scrape the metrics from the EigenDA node.

Create the following file in `$HOME/.eigenlayer/config/prometheus.yml`
```yaml
global:
  scrape_interval: 15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: "codelab-monitor"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      # Point to the same endpoint that EigenDA is publishing on.
      # If using the sample docker-compose.yml for EigenDA, use the name of the
      # container instead of localhost (e.g. da-node:9092)
      - targets: ["localhost:<NODE_METRICS_PORT>"]
```

Start Prometheus
```bash
prometheus --config.file="$HOME/.eigenlayer/config/prometheus.yml"
```

If you want to use Docker, follow [this](https://prometheus.io/docs/prometheus/latest/installation/#volumes-bind-mount) link.
```bash
docker run -d \
    -p 9090:9090 \
    -v ~/.eigenlayer/config/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus
```

#### Grafana
Grafana is used to visualize the metrics from the EigenDA node.

You can use [OSS Grafana](https://grafana.com/oss/grafana/) for it or any other Dashboard provider.

Start the Grafana server
```bash
grafana server
```
You can also use [Docker](https://grafana.com/docs/grafana/latest/setup-grafana/installation/docker/)
```bash
docker run -d -p 3000:3000 --name=grafana grafana/grafana-enterprise
```

You should be able to navigate to `http://localhost:3000` and login with `admin`/`admin`.
You will need to add a datasource to Grafana. You can do this by navigating to `http://localhost:3000/datasources` and adding a Prometheus datasource. By default, the Prometheus server is running on `http://localhost:9090`. You can use `http://prometheus:9090` as the server URL for the datasource.


#### Node exporter
EigenDA emits DA specific metrics but, it's also important to keep track of the node's health. For this, we will use [Node Exporter](https://prometheus.io/docs/guides/node-exporter/) which is a Prometheus exporter for hardware and OS metrics exposed by *NIX kernels, written in Go with pluggable metric collectors.
Install the binary or use Docker to [run](https://hub.docker.com/r/prom/node-exporter) it.

```bash
docker pull prom/node-exporter
docker run -d -p 9100:9100 --name node-exporter prom/node-exporter
```

### Useful Dashboards
EigenDA provides a set of Grafana dashboards that provide insights into key performance indicators and health metrics of an EigenDA node. These dashboards can be accessed [here](./dashboards).
Once you have Grafana setup, they should be automatically imported.
