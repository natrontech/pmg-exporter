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
pmg_subscription_info{level, productname}
pmg_subscription_nextdue_timestamp_seconds{level, productname} 
pmg_subscription_status{status="new, notfound, active, invalid, expired, suspended"}
pmg_node_insync{node}
pmg_node_uptime_seconds{node}
pmg_postfix_queue_size{age_bucket="5m | 10m | 20m| 40m | 80n | 160m | 320m | 640m | 1280m | 1280m+ | total", domain}
pmg_cluster_nodes_total 
pmg_cluster_domains_total 
pmg_cluster_backups_remotes_total 
pmg_quarantine_spam_count_total 
pmg_quarantine_spam_average_size_bytes 
pmg_quarantine_spam_average_level 
pmg_quarantine_spam_disk_usage_megabytes 
pmg_quarantine_virus_count_total 
pmg_quarantine_virus_average_size_bytes 
pmg_quarantine_virus_disk_usage_megabytes 
pmg_postfix_messages_total
pmg_postfix_messages_in_total 
pmg_postfix_messages_out_total
pmg_postfix_bounces_in_total 
pmg_postfix_bounces_out_total 
pmg_postfix_spam_in_total 
pmg_postfix_spam_out_total 
pmg_postfix_virus_in_total 
pmg_postfix_virus_out_total 
pmg_postfix_rbl_rejects_total 
pmg_postfix_pregreet_rejects_total 
pmg_release_info_total{release=}
pmg_repository_info_total{repo}
pmg_version_info_total{version}
```

