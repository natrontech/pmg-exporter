# Proxmox Mail Gateway Exporter for Prometheus

[![license](https://img.shields.io/github/license/natrontech/pmg-exporter)](https://github.com/natrontech/pmg-exporter/blob/main/LICENSE)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/natrontech/pmg-exporter/badge)](https://scorecard.dev/viewer/?uri=github.com/natrontech/pmg-exporter)
[![release](https://img.shields.io/github/v/release/natrontech/pmg-exporter)](https://github.com/natrontech/pmg-exporter/releases)
![Python version](https://img.shields.io/badge/python-3.10%2B-blue)
[![SLSA 3](https://slsa.dev/images/gh-badge-level3.svg)](https://slsa.dev)

PMG Exporter is a [Prometheus](https://prometheus.io/) exporter designed to monitor [Proxmox Mail Gateway](https://www.proxmox.com/en/products/proxmox-mail-gateway/overview) (PMG) instances and clusters. It collects various metrics related to the health and performance of PMG systems, making it easier to integrate PMG monitoring into your existing Prometheus setup.

## Exported Metrics

| Metric                                       | Description                                                                                    | Labels                                      |
| -------------------------------------------- | ---------------------------------------------------------------------------------------------- | ------------------------------------------- |
| `pmg_exporter_up`                            | PMG Exporter up status (1 if up)                                                               |                                             |
| `pmg_cluster_node_status`                    | Cluster node status (1 if online)                                                              |                                             |
| `pmg_subscription_info`                      | Node subscription info (always 1, labeled)                                                     | `level`, `productname`                      |
| `pmg_subscription_nextdue_timestamp_seconds` | Subscription next due timestamp                                                                | `level`, `productname`                      |
| `pmg_subscription_status`                    | Subscription state per status (`new`, `notfound`, `active`, `invalid`, `expired`, `suspended`) | `status`                                    |
| `pmg_node_insync`                            | Node config in-sync state (1 if in sync)                                                       | `node`                                      |
| `pmg_node_uptime_seconds`                    | Node uptime in seconds                                                                         | `node`                                      |
| `pmg_postfix_queue_size`                     | Postfix mail queue size                                                                        | `age_bucket`, `domain`                      |
| `pmg_cluster_nodes_total`                    | Total number of nodes in the cluster                                                           |                                             |
| `pmg_cluster_node_info`                      | Cluster node info (1 if present)                                                               | `name`                                      |
| `pmg_cluster_domains_total`                  | Total number of domains in the cluster                                                         |                                             |
| `pmg_cluster_domain_info`                    | Cluster domain info (1 if present)                                                             | `domain`                                    |
| `pmg_cluster_backups_remotes_total`          | Total backup remotes in cluster                                                                |                                             |
| `pmg_cluster_backup_remote_info`             | Cluster backup remote info (1 if present, labeled)                                             | `datastore`, `remote`, `server`, `disabled` |
| `pmg_quarantine_spam_count_total`            | Quarantine spam count                                                                          |                                             |
| `pmg_quarantine_spam_average_size_bytes`     | Quarantine spam average size (bytes)                                                           |                                             |
| `pmg_quarantine_spam_average_level`          | Quarantine spam average spam level                                                             |                                             |
| `pmg_quarantine_spam_disk_usage_megabytes`   | Quarantine spam disk usage (MB)                                                                |                                             |
| `pmg_quarantine_virus_count_total`           | Quarantine virus count                                                                         |                                             |
| `pmg_quarantine_virus_average_size_bytes`    | Quarantine virus average size (bytes)                                                          |                                             |
| `pmg_quarantine_virus_disk_usage_megabytes`  | Quarantine virus disk usage (MB)                                                               |                                             |
| `pmg_postfix_messages_total`                 | Total messages processed today                                                                 |                                             |
| `pmg_postfix_messages_in_total`              | Total inbound messages today                                                                   |                                             |
| `pmg_postfix_messages_out_total`             | Total outbound messages today                                                                  |                                             |
| `pmg_postfix_bounces_in_total`               | Inbound bounces today                                                                          |                                             |
| `pmg_postfix_bounces_out_total`              | Outbound bounces today                                                                         |                                             |
| `pmg_postfix_spam_in_total`                  | Inbound spam today                                                                             |                                             |
| `pmg_postfix_spam_out_total`                 | Outbound spam today                                                                            |                                             |
| `pmg_postfix_virus_in_total`                 | Inbound virus messages today                                                                   |                                             |
| `pmg_postfix_virus_out_total`                | Outbound virus messages today                                                                  |                                             |
| `pmg_postfix_rbl_rejects_total`              | Messages rejected by RBL today                                                                 |                                             |
| `pmg_postfix_pregreet_rejects_total`         | Messages rejected during pregreet today                                                        |                                             |
| `pmg_release_info_total`                     | Release information                                                                            | `release`                                   |
| `pmg_repository_info_total`                  | Git commit hash the build was made from                                                        | `repo`                                      |
| `pmg_version_info_total`                     | Installed PMG API package version                                                              | `version`                                   |

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

To run the exporter, install the package and execute the following commands:

```bash
pip install .
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
