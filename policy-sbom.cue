// The predicateType field must match this string
PredicateType: "https://cyclonedx.org/bom"

Predicate: {
  metadata: {
    component: {
      // Enforce the SBOM root component to be the container image digest
      "bom-ref": =~"^pkg:oci/ghcr\\.io/natrontech/pmg-exporter@sha256:[a-f0-9]{64}$"
    }
  }
}
