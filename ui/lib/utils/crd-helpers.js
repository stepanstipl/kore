
export function isReadOnlyCRD(crd) {
  return crd && crd.metadata && crd.metadata.annotations && crd.metadata.annotations['kore.appvia.io/readonly']
}

export function getKoreLabel(obj, name) {
  if (obj && obj.metadata && obj.metadata.labels) {
    return obj.metadata.labels[`kore.appvia.io/${name}`] || ''
  }
  return ''
}
