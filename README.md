# Proxmox Mail Gateway Exporter for Prometheus

PMG Exporter is a Prometheus exporter designed to monitor Proxmox Mail Gateway (PMG) instances and clusters. It collects various metrics related to the health and performance of PMG systems, making it easier to integrate PMG monitoring into your existing Prometheus setup.

## Configuration
PMG Exporter can be configured using a `.env` file or environment variables. The configuration options include:
```ini .env
PMG_HOST=your_pmg_host
PMG_USER=your_username (Auditor role recommended)
PMG_PASSWORD=your_password
PMG_VERIFY_SSL=true
PMG_BACKEND=https
PMG_SERVICE=pmg
PMG_LOG_LEVEL=INFO
```

## Usage
### CLI
To run the exporter, use the following command:
```bash
pmg-exporter --config /path/to/your/.env
```

### Docker
You can also run PMG Exporter using Docker:
```bash
docker run -d \
  -p 10069:10069 \
  --env-file /path/to/your/.env \
  ghcr.io/natrontech/pmg-exporter:latest
```

## Metrics
```
# HELP pmg_cluster_node_status Proxmox Mail Gateway cluster node status (1 if online)
# TYPE pmg_cluster_node_status gauge
# HELP pmg_subscription_info Proxmox Mail Gateway node subscription info (always 1, labeled)
# TYPE pmg_subscription_info gauge
# HELP pmg_subscription_nextdue_timestamp_seconds Proxmox Mail Gateway node subscription next due timestamp
# TYPE pmg_subscription_nextdue_timestamp_seconds gauge
# HELP pmg_subscription_status Proxmox Mail Gateway node subscription status (1 if matching status)
# TYPE pmg_subscription_status gauge
# HELP pmg_node_insync Proxmox Mail Gateway node configuration in sync status (1 if in sync)
# TYPE pmg_node_insync gauge
# HELP pmg_node_uptime_seconds Proxmox Mail Gateway node uptime in seconds
# TYPE pmg_node_uptime_seconds gauge
# HELP pmg_cluster_nodes_total Total number of nodes in the Proxmox Mail Gateway cluster
# TYPE pmg_cluster_nodes_total gauge
# HELP pmg_cluster_node_info Proxmox Mail Gateway cluster node info (1 if present)
# TYPE pmg_cluster_node_info gauge
# HELP pmg_cluster_domains_total Total number of domains in the Proxmox Mail Gateway cluster
# TYPE pmg_cluster_domains_total gauge
# HELP pmg_cluster_domain_info Proxmox Mail Gateway cluster domain info (1 if present)
# TYPE pmg_cluster_domain_info gauge
# HELP pmg_cluster_backups_remotes_total Total number of backup remotes in the Proxmox Mail Gateway cluster
# TYPE pmg_cluster_backups_remotes_total gauge
# HELP pmg_cluster_backup_remote_info Proxmox Mail Gateway cluster backup remote info (1 if present)
# TYPE pmg_cluster_backup_remote_info gauge
```

