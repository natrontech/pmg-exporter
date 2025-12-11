from prometheus_client import start_http_server, REGISTRY

from typing import Any
import asyncio

from proxmoxer import ProxmoxAPI  # pyright: ignore[reportMissingTypeStubs]

from pmg_exporter.collectors import (
    ExporterCollector,
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

from pmg_exporter.config import load_config

import logging

logging.getLogger("pmg_exporter")


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

        sanitized_kwargs: list[str] = []
        for k, arg in kwargs.items():
            sanitized_kwargs.append(str({k: arg} if k != "password" else {k: "***"}))

        logging.debug(f"ProxmoxAPI initialization parameters: {sanitized_kwargs}")
        self.proxmox = ProxmoxAPI(**kwargs)
        logging.info("ProxmoxAPI client initialized.")

    def register_collectors(self) -> None:
        logging.info("Registering collectors...")
        mapping: dict[str, Any] = {
            "exporter_status": ExporterCollector,
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
            REGISTRY.register(collector(self.proxmox))
            logging.debug(f"Collector '{name}' registered successfully.")

    async def run(self) -> None:
        logging.info("Starting HTTP server for Prometheus metrics...")
        port = int(self.config.get("exporter_port", 10069))
        default_addr = "127.0.0.1"
        addr = str(self.config.get("exporter_address", default_addr))
        if addr == "0.0.0.0":
            logging.warning(
                (
                    "Configured to bind to all interfaces (0.0.0.0);"
                    " consider restricting to a specific address for security."
                )
            )
        start_http_server(port=port, addr=addr)
        logging.info(f"HTTP server started on {addr}:{port}.")
        await asyncio.Event().wait()
