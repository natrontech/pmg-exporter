from prometheus_client import start_http_server, REGISTRY

from typing import Any
import asyncio

from proxmoxer import ProxmoxAPI  # pyright: ignore[reportMissingTypeStubs]

from pmg_exporter.collectors import (
    ClusterStatusCollector,
    ClusterNodesCollector,
    ClusterDomainsCollector,
    ClusterBackupCollector,
    NodeStatusCollector,
    NodeSubscriptionCollector,
    NodePostfixQueueCollector,
    QuarantineSpamCollector,
    QuarantineVirusCollector,
    StatisticsMailcountCollector,
    VersionInfoCollector,
)

import logging

logging.getLogger("pmg_exporter")

from pmg_exporter.config import load_config


class PMGExporter:
    def __init__(self, config_file: str, duration: int, collectors: list[str]) -> None:
        self.config = load_config(config_file)
        self.duration = duration
        self.collectors = collectors
        self.proxmox: ProxmoxAPI | None = None
        logging.basicConfig(level=self.config.get("log_level", "INFO"))
        logging.info("PMGExporter initialized with config:")
        for key, value in self.config.items():
            logging.info(
                f"  {key}: {value}" if key != "password" else "  password: ****"
            )

    def start(self) -> None:
        logging.info("Starting PMGExporter...")
        self.initialize_proxmox_client()

    def initialize_proxmox_client(self) -> None:
        logging.info("Initializing ProxmoxAPI client...")
        excluded_keys = {"exporter_port", "log_level"}
        raw = dict(self.config)
        for k in excluded_keys:
            raw.pop(k, None)

        kwargs: dict[str, Any] = {
            k: (bool(v) if k == "verify_ssl" else str(v)) for k, v in raw.items()
        }

        logging.debug(f"ProxmoxAPI initialization parameters: {kwargs}")
        self.proxmox = ProxmoxAPI(**kwargs)
        logging.info("ProxmoxAPI client initialized.")

    def register_collectors(self) -> None:
        logging.info("Registering collectors...")
        mapping: dict[str, Any] = {
            "cluster_status": ClusterStatusCollector,
            "subscription": NodeSubscriptionCollector,
            "node_status": NodeStatusCollector,
            "node_postfix_queue": NodePostfixQueueCollector,
            "cluster_nodes": ClusterNodesCollector,
            "cluster_domains": ClusterDomainsCollector,
            "cluster_backups": ClusterBackupCollector,
            "quarantine_spam": QuarantineSpamCollector,
            "quarantine_virus": QuarantineVirusCollector,
            "statistics_mailcount": StatisticsMailcountCollector,
            "version_info": VersionInfoCollector,
        }
        if self.collectors == ["all"]:
            self.collectors = list(mapping.keys())
        for name in self.collectors:
            logging.debug(f"Registering collector: {name}")
            collector = mapping.get(name)
            if collector is None:
                logging.warning(f"Collector '{name}' not found. Skipping registration.")
                continue
            assert self.proxmox is not None
            REGISTRY.register(collector(self.proxmox))
            logging.debug(f"Collector '{name}' registered successfully.")

    async def run(self) -> None:
        logging.info("Starting HTTP server for Prometheus metrics...")
        start_http_server(
            port=int(self.config.get("exporter_port", 10069)),
            addr=str(self.config.get("exporter_address", "0.0.0.0")),
        )
        logging.info(
            f"HTTP server started. Exporter is running on port {self.config.get('exporter_port', 10069)}."
        )
        await asyncio.Event().wait()
