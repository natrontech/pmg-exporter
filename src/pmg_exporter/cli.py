from argparse import ArgumentParser
import os
import pathlib
import asyncio

from pmg_exporter import PMGExporter


def cli() -> None:
    parser = ArgumentParser(
        description=(
            "Proxmox Mail Gateway Exporter CLI"
            " - Available collectors: all, cluster_status, subscription, "
            "node_status, cluster_nodes, cluster_domains, cluster_backups"
        )
    )
    exporterFlags = parser.add_argument_group("Exporter Flags")
    exporterFlags.add_argument(
        "--config-file",
        type=pathlib.Path,
        default=pathlib.Path(os.getenv("PMG_EXPORTER_CONFIG_FILE", ".env")),
        help=(
            "Path to the configuration file"
            "( default: .env or PMG_EXPORTER_CONFIG_FILE env variable)"
        ),
    )
    exporterFlags.add_argument(
        "--collectors",
        type=str,
        default="all",
        dest="collectors",
        help="Comma-separated list of collectors to enable (default: all)",
    )

    parameters = parser.parse_args()

    exporter = PMGExporter(
        config_file=str(parameters.config_file),
        duration=10,
        collectors=parameters.collectors.split(","),
    )
    exporter.start()
    exporter.register_collectors()
    asyncio.run(exporter.run())


if __name__ == "__main__":
    cli()
