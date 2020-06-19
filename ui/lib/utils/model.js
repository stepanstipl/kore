import V1ObjectMeta from '../kore-api/model/V1ObjectMeta'
import V1Ownership from '../kore-api/model/V1Ownership'

const NewV1ObjectMeta = (name, namespace) => {
  const meta = new V1ObjectMeta()
  meta.setName(name)
  if (namespace) {
    meta.setNamespace(namespace)
  }
  return meta
}

const NewV1Ownership = ({ group, version, kind, name, namespace }) => {
  const ownership = new V1Ownership()
  ownership.setGroup(group)
  ownership.setVersion(version)
  ownership.setKind(kind)
  ownership.setName(name)
  ownership.setNamespace(namespace)
  return ownership
}

module.exports = {
  NewV1ObjectMeta,
  NewV1Ownership
}
