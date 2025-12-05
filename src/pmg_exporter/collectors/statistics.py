from typing import Any, cast

from prometheus_client.core import CounterMetricFamily
from prometheus_client.registry import Collector

from proxmoxer import ProxmoxAPI  # pyright: ignore[reportMissingTypeStubs]

import logging

logging.getLogger("pmg_exporter")


class StatisticsMailcountCollector(Collector):
    def __init__(self, proxmox: ProxmoxAPI) -> None:
        self.proxmox = proxmox

    def collect(self):
        logging.debug("Collecting node Postfix daily stats metrics...")

        stats = cast(
            list[dict[str, Any]],
            self.proxmox.statistics.mailcount.get(),  # type: ignore
        )

        if not stats:
            return

        totals: dict[str, float] = {
            "count": 0.0,
            "count_in": 0.0,
            "count_out": 0.0,
            "bounces_in": 0.0,
            "bounces_out": 0.0,
            "spamcount_in": 0.0,
            "spamcount_out": 0.0,
            "viruscount_in": 0.0,
            "viruscount_out": 0.0,
            "rbl_rejects": 0.0,
            "pregreet_rejects": 0.0,
        }

        for row in stats:
            for key in totals.keys():
                raw = row.get(key, 0)
                try:
                    value = float(raw)
                except (TypeError, ValueError):
                    value = 0.0
                totals[key] += value

        m_total = CounterMetricFamily(
            "pmg_postfix_messages_total",
            "Total messages processed today (in + out). Resets daily.",
        )
        m_total.add_metric([], totals["count"])
        yield m_total

        m_in = CounterMetricFamily(
            "pmg_postfix_messages_in_total",
            "Total inbound messages today. Resets daily.",
        )
        m_in.add_metric([], totals["count_in"])
        yield m_in

        m_out = CounterMetricFamily(
            "pmg_postfix_messages_out_total",
            "Total outbound messages today. Resets daily.",
        )
        m_out.add_metric([], totals["count_out"])
        yield m_out

        b_in = CounterMetricFamily(
            "pmg_postfix_bounces_in_total",
            "Total inbound bounces today. Resets daily.",
        )
        b_in.add_metric([], totals["bounces_in"])
        yield b_in

        b_out = CounterMetricFamily(
            "pmg_postfix_bounces_out_total",
            "Total outbound bounces today. Resets daily.",
        )
        b_out.add_metric([], totals["bounces_out"])
        yield b_out

        s_in = CounterMetricFamily(
            "pmg_postfix_spam_in_total",
            "Inbound messages classified as spam today. Resets daily.",
        )
        s_in.add_metric([], totals["spamcount_in"])
        yield s_in

        s_out = CounterMetricFamily(
            "pmg_postfix_spam_out_total",
            "Outbound messages classified as spam today. Resets daily.",
        )
        s_out.add_metric([], totals["spamcount_out"])
        yield s_out

        v_in = CounterMetricFamily(
            "pmg_postfix_virus_in_total",
            "Inbound messages with viruses today. Resets daily.",
        )
        v_in.add_metric([], totals["viruscount_in"])
        yield v_in

        v_out = CounterMetricFamily(
            "pmg_postfix_virus_out_total",
            "Outbound messages with viruses today. Resets daily.",
        )
        v_out.add_metric([], totals["viruscount_out"])
        yield v_out

        rbl = CounterMetricFamily(
            "pmg_postfix_rbl_rejects_total",
            "Messages rejected by RBL today. Resets daily.",
        )
        rbl.add_metric([], totals["rbl_rejects"])
        yield rbl

        pregreet = CounterMetricFamily(
            "pmg_postfix_pregreet_rejects_total",
            "Messages rejected during pregreeting today. Resets daily.",
        )
        pregreet.add_metric([], totals["pregreet_rejects"])
        yield pregreet
