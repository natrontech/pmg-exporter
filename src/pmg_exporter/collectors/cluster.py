from typing import cast

from prometheus_client.core import GaugeMetricFamily
from prometheus_client.registry import Collector

from proxmoxer import ProxmoxAPI  # pyright: ignore[reportMissingTypeStubs]

import logging

logging.getLogger("pmg_exporter")


class BaseClusterCollector(Collector):
    def __init__(self, proxmox: ProxmoxAPI) -> None:
        self.proxmox = proxmox

    def _get_cluster_status_entries(self) -> list[dict[str, str]]:
        logging.debug("Fetching cluster status entries from ProxmoxAPI")
        raw = (  # type: ignore[assignment]
            self.proxmox.config.cluster.status.get() or []  # type: ignore
        )
        return cast(list[dict[str, str]], raw)

    def _get_cluster_node_entries(self) -> list[dict[str, str]]:
        logging.debug("Fetching cluster node entries from ProxmoxAPI")
        entries = self._get_cluster_status_entries()
        return [e for e in entries if e.get("type") == "node"]

    def _get_cluster_domain_entries(self) -> list[dict[str, str]]:
        logging.debug("Fetching cluster domain entries from ProxmoxAPI")
        raw = (  # type: ignore[assignment]
            self.proxmox.config.domains.get() or []  # pyright: ignore
        )
        return cast(list[dict[str, str]], raw)

    def _get_cluster_backup_remote_entries(self) -> list[dict[str, str]]:
        logging.debug("Fetching cluster backup remote entries from ProxmoxAPI")
        raw = (  # type: ignore[assignment]
            self.proxmox.config.pbs.get() or []  # pyright: ignore
        )
        return cast(list[dict[str, str]], raw)


class ClusterStatusCollector(BaseClusterCollector):
    def collect(self):
        logging.debug("Collecting cluster status metrics")
        status_metric = GaugeMetricFamily(
            "pmg_cluster_node_status",
            "Proxmox Mail Gateway cluster node status (1 if online)",
            labels=["name", "status"],
        )

        for entry in self._get_cluster_node_entries():
            name = entry.get("name", "unknown")
            status = entry.get("status", "unknown")
            value = 1 if status == "online" else 0

            status_metric.add_metric(
                [name, status],
                value,
            )

        yield status_metric


class ClusterNodesCollector(BaseClusterCollector):
    def collect(self):
        logging.debug("Collecting cluster nodes metrics")
        nodes_total_metric = GaugeMetricFamily(
            "pmg_cluster_nodes_total",
            "Total number of nodes in the Proxmox Mail Gateway cluster",
        )

        nodes_info_metric = GaugeMetricFamily(
            "pmg_cluster_node_info",
            "Proxmox Mail Gateway cluster node info (1 if present)",
            labels=["name"],
        )

        node_entries = self._get_cluster_node_entries()

        for entry in node_entries:
            name = entry.get("name", "unknown")
            nodes_info_metric.add_metric(
                [name],
                1,
            )

        total_nodes = len(node_entries)
        nodes_total_metric.add_metric([], total_nodes)

        yield nodes_total_metric
        yield nodes_info_metric


class ClusterDomainsCollector(BaseClusterCollector):
    def collect(self):
        logging.debug("Collecting cluster domains metrics")
        domains_total_metric = GaugeMetricFamily(
            "pmg_cluster_domains_total",
            "Total number of domains in the Proxmox Mail Gateway cluster",
        )

        domains_info_metric = GaugeMetricFamily(
            "pmg_cluster_domain_info",
            "Proxmox Mail Gateway cluster domain info (1 if present)",
            labels=["domain"],
        )

        domain_entries = self._get_cluster_domain_entries()

        total_domains = len(domain_entries)
        domains_total_metric.add_metric([], total_domains)

        for entry in domain_entries:
            domain = entry.get("domain", "unknown")
            domains_info_metric.add_metric(
                [domain],
                1,
            )

        yield domains_total_metric
        yield domains_info_metric


class ClusterBackupCollector(BaseClusterCollector):
    def collect(self):
        logging.debug("Collecting cluster backup metrics")
        backups_remotes_total_metric = GaugeMetricFamily(
            "pmg_cluster_backups_remotes_total",
            "Total number of backup remotes in the Proxmox Mail Gateway cluster",
        )

        backups_remotes_info_metric = GaugeMetricFamily(
            "pmg_cluster_backup_remote_info",
            "Proxmox Mail Gateway cluster backup remote info (1 if present)",
            labels=["datastore", "remote", "server", "disabled"],
        )

        remote_entries = self._get_cluster_backup_remote_entries()

        total_remotes = len(remote_entries)
        backups_remotes_total_metric.add_metric([], total_remotes)

        for entry in remote_entries:
            datastore = entry.get("datastore", "unknown")
            remote = entry.get("remote", "unknown")
            server = entry.get("server", "unknown")
            disabled = entry.get("disabled", "0")

            backups_remotes_info_metric.add_metric(
                [datastore, remote, server, disabled],
                1,
            )

        yield backups_remotes_total_metric
        yield backups_remotes_info_metric
