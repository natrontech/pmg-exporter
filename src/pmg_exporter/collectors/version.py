from prometheus_client.core import CounterMetricFamily
from prometheus_client.registry import Collector

from proxmoxer import ProxmoxAPI  # pyright: ignore[reportMissingTypeStubs]

import logging

logging.getLogger("pmg_exporter")


class VersionInfoCollector(Collector):
    def __init__(self, proxmox: ProxmoxAPI) -> None:
        self.proxmox = proxmox

    def collect(self):
        logging.debug("Collecting Proxmox Mail Gateway version metrics...")

        version_info = self.proxmox.version.get()  # type: ignore

        release_metric = CounterMetricFamily(
            "pmg_release_info",
            "Proxmox Mail Gateway release information",
            labels=["release"],
        )
        release = str(version_info.get("release", "unknown"))  # type: ignore
        release_metric.add_metric([release], 1)
        yield release_metric

        repo_metric = CounterMetricFamily(
            "pmg_repository_info",
            "Git commit hash from which Proxmox Mail Gateway was built",
            labels=["repo"],
        )
        repo = str(version_info.get("repoid", "unknown"))  # type: ignore
        repo_metric.add_metric([repo], 1)
        yield repo_metric

        version_metric = CounterMetricFamily(
            "pmg_version_info",
            "Currently installed Proxmox Mail Gateway API package version",
            labels=["version"],
        )
        version = str(version_info.get("version", "unknown"))  # type: ignore
        version_metric.add_metric([version], 1)

        yield version_metric
