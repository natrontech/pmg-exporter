FROM python:3.15.0-slim@sha256:f5b92f2a669a4c347044f29367f31983280b648250ca94dea66f98c53bc87e58 AS runtime

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

RUN apt-get update -qq && apt-get install -y --no-install-recommends \
    tzdata \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY pyproject.toml README.md /app/
COPY src /app/src

ARG VERSION=0.0.dev0
ENV SETUPTOOLS_SCM_PRETEND_VERSION=${VERSION}
ENV EXPORTER_VERSION=${VERSION}

LABEL org.opencontainers.image.title="pmg-exporter" \
      org.opencontainers.image.description="Prometheus exporter for Proxmox Mail Gateway" \
      org.opencontainers.image.version="${EXPORTER_VERSION}" \
      org.opencontainers.image.authors="NatronTech <info@natron.io>, Nicolo Lüscher <nicolo.luescher@natron.io>"

RUN pip install --no-cache-dir --upgrade pip \
    && pip install --no-cache-dir hatchling hatch-vcs \
    && pip install --no-cache-dir .

RUN useradd -u 1000 -r -s /usr/sbin/nologin pmg-exporter
USER pmg-exporter


EXPOSE 10069

CMD ["pmg-exporter"]
