from prometheus_client.core import GaugeMetricFamily
from prometheus_client.registry import Collector

from proxmoxer import ProxmoxAPI  # pyright: ignore[reportMissingTypeStubs]

import logging

logging.getLogger("pmg_exporter")


class QuarantineSpamCollector(Collector):
    def __init__(self, proxmox: ProxmoxAPI) -> None:
        self.proxmox = proxmox

    def collect(self):
        logging.debug("Collecting quarantine spam metrics...")
        spam_count_metric = GaugeMetricFamily(
            "pmg_quarantine_spam_count_total",
            "Proxmox Mail Gateway quarantine spam count",
        )
        spam_average_size_metric = GaugeMetricFamily(
            "pmg_quarantine_spam_average_size_bytes",
            "Proxmox Mail Gateway quarantine spam average size in bytes",
        )
        spam_average_level_metric = GaugeMetricFamily(
            "pmg_quarantine_spam_average_level",
            "Proxmox Mail Gateway quarantine spam average spam level",
        )
        spam_disk_usage_metric = GaugeMetricFamily(
            "pmg_quarantine_spam_disk_usage_megabytes",
            "Proxmox Mail Gateway quarantine estimated spam disk usage in megabytes",
        )

        quarantine_info = self.proxmox.quarantine.spamstatus.get()  # type: ignore

        spam_count = int(quarantine_info.get("count", 0))  # type: ignore
        spam_average_size = float(quarantine_info.get("avgbytes", 0.0))  # type: ignore
        spam_average_level = float(quarantine_info.get("avgspam", 0.0))  # type: ignore
        spam_disk_usage = float(quarantine_info.get("mbytes", 0.0))  # type: ignore

        spam_count_metric.add_metric([], spam_count)
        spam_average_size_metric.add_metric([], spam_average_size)
        spam_average_level_metric.add_metric([], spam_average_level)
        spam_disk_usage_metric.add_metric([], spam_disk_usage)

        yield spam_count_metric
        yield spam_average_size_metric
        yield spam_average_level_metric
        yield spam_disk_usage_metric


class QuarantineVirusCollector(Collector):
    def __init__(self, proxmox: ProxmoxAPI) -> None:
        self.proxmox = proxmox

    def collect(self):
        logging.debug("Collecting quarantine virus metrics...")
        virus_metric_count = GaugeMetricFamily(
            "pmg_quarantine_virus_count_total",
            "Proxmox Mail Gateway quarantine virus count",
        )
        virus_average_size_metric = GaugeMetricFamily(
            "pmg_quarantine_virus_average_size_bytes",
            "Proxmox Mail Gateway quarantine virus average size in bytes",
        )
        virus_disk_usage_metric = GaugeMetricFamily(
            "pmg_quarantine_virus_disk_usage_megabytes",
            "Proxmox Mail Gateway quarantine estimated virus disk usage in megabytes",
        )

        quarantine_info = self.proxmox.quarantine.virusstatus.get()  # type: ignore

        virus_count = int(quarantine_info.get("count", 0))  # type: ignore
        virus_average_size = float(quarantine_info.get("avgbytes", 0.0))  # type: ignore
        virus_disk_usage = float(quarantine_info.get("mbytes", 0.0))  # type: ignore

        virus_metric_count.add_metric([], virus_count)
        virus_average_size_metric.add_metric([], virus_average_size)
        virus_disk_usage_metric.add_metric([], virus_disk_usage)

        yield virus_metric_count
        yield virus_average_size_metric
        yield virus_disk_usage_metric
