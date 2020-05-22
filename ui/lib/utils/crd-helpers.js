
export function isReadOnlyCRD(crd) {
  return crd && crd.metadata && crd.metadata.annotations && crd.metadata.annotations['kore.appvia.io/readonly']
}
