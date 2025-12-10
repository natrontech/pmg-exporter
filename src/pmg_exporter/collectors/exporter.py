from prometheus_client.core import GaugeMetricFamily
from prometheus_client.registry import Collector
from proxmoxer import ProxmoxAPI  # pyright: ignore[reportMissingTypeStubs]

class ExporterCollector(Collector):
    """
    Collector for exporter-level metrics.
    """
    
    def __init__(self, proxmox: ProxmoxAPI) -> None:
        self.proxmox = proxmox
    
    def collect(self):

        exporter_up_metric = GaugeMetricFamily(
            "pmg_exporter_up",
            "PMG Exporter up status (1 if up)",
        )
        exporter_up_metric.add_metric([], 1)
        yield exporter_up_metric