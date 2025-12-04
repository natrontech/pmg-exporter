from typing import Optional, Any, cast

from prometheus_client.core import GaugeMetricFamily
from prometheus_client.registry import Collector

from proxmoxer import ProxmoxAPI  # pyright: ignore[reportMissingTypeStubs]

import logging

logging.getLogger("pmg_exporter")


POSSIBLE_SUBSCRIPTION_STATUSES: tuple[str, ...] = (
    "new",
    "notfound",
    "active",
    "invalid",
    "expired",
    "suspended",
)


def get_single_node_name(proxmox: ProxmoxAPI) -> Optional[str]:
    """
    Return the first node name from proxmox, or None if there are no nodes.
    """
    nodes = proxmox.nodes.get() or []  # type: ignore[assignment]
    for entry in cast(list[dict[str, Any]], nodes):
        return cast(str, entry["node"])
    return None


class NodeStatusCollector(Collector):
    def __init__(self, proxmox: ProxmoxAPI) -> None:
        self.proxmox = proxmox

    def collect(self):
        logging.debug("Collecting node status metrics...")
        node = get_single_node_name(self.proxmox)
        if node is None:
            return

        node_insync_metric = GaugeMetricFamily(
            "pmg_node_insync",
            "Proxmox Mail Gateway node configuration in sync status (1 if in sync)",
            labels=["node"],
        )

        node_uptime_metric = GaugeMetricFamily(
            "pmg_node_uptime_seconds",
            "Proxmox Mail Gateway node uptime in seconds",
            labels=["node"],
        )

        status = cast(dict[str, Any], self.proxmox.nodes(node).status.get())  # type: ignore

        node_insync_metric.add_metric(
            [node],
            1 if status.get("insync", 0) == 1 else 0,
        )

        node_uptime_metric.add_metric(
            [node],
            int(status.get("uptime", 0)),
        )

        yield node_insync_metric
        yield node_uptime_metric


class NodeSubscriptionCollector(Collector):
    def __init__(self, proxmox: ProxmoxAPI) -> None:
        self.proxmox = proxmox

    def collect(self):
        logging.debug("Collecting node subscription metrics...")
        node = get_single_node_name(self.proxmox)
        if node is None:
            return

        subscription_info_metric = GaugeMetricFamily(
            "pmg_subscription_info",
            "Proxmox Mail Gateway node subscription info (always 1, labeled)",
            labels=["level", "productname"],
        )

        subscription_nextdue_timestamp_metric = GaugeMetricFamily(
            "pmg_subscription_nextdue_timestamp_seconds",
            "Proxmox Mail Gateway node subscription next due timestamp",
            labels=["level", "productname"],
        )

        subscription_status_metric = GaugeMetricFamily(
            "pmg_subscription_status",
            "Proxmox Mail Gateway node subscription status (1 if matching status)",
            labels=["status"],
        )

        subscription_info = cast(
            dict[str, Any],
            self.proxmox.nodes(node).subscription.get(),  # type: ignore
        )

        level = subscription_info.get("level", "unknown")
        productname = subscription_info.get("productname", "unknown")
        nextdue_raw = subscription_info.get("nextdue", 0)

        try:
            nextdue = int(nextdue_raw)
        except (TypeError, ValueError):
            nextdue = 0

        subscription_info_metric.add_metric([level, productname], 1)
        subscription_nextdue_timestamp_metric.add_metric(
            [level, productname],
            nextdue,
        )

        status = subscription_info.get("status", "unknown")

        for possible_status in POSSIBLE_SUBSCRIPTION_STATUSES:
            value = 1 if status == possible_status else 0
            subscription_status_metric.add_metric([possible_status], value)

        yield subscription_info_metric
        yield subscription_nextdue_timestamp_metric
        yield subscription_status_metric
