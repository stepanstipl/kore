import moment from 'moment'

export function startOfMonth(momentlike) {
  if (!momentlike) {
    momentlike = moment().utc(false)
  }
  return moment(momentlike).utc(false).startOf('month')
}

export function endOfMonth(momentlike) {
  if (!momentlike) {
    momentlike = moment().utc(false)
  }
  return moment(momentlike).utc(false).endOf('month')
}

export function apiDateTime(momentlike) {
  if (!momentlike) {
    momentlike = moment().utc(false)
  }
  return moment(momentlike).utc(false).format('YYYY-MM-DDTHH:mm:ssZ')
}