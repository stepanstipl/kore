/* eslint-disable no-useless-escape */
export const patterns = {
  uriCompatible10CharMax: {
    pattern: '^[a-z]?[a-z0-9-]{0,8}[a-z0-9]$',
    message: 'Must consist of lower case alphanumeric characters or "-", it must start with a letter and end with an alphanumeric and must be no longer than 10 characters'
  },
  uriCompatible40CharMax: {
    pattern: '^[a-z][a-z0-9-]{0,38}[a-z0-9]$',
    message: 'Must consist of lower case alphanumeric characters or "-", it must start with a letter and end with an alphanumeric and must be no longer than 40 characters'
  },
  uriCompatible63CharMax: {
    pattern: '^[a-z][a-z0-9-]{0,61}[a-z0-9]$',
    message: 'Must consist of lower case alphanumeric characters or "-", it must start with a letter and end with an alphanumeric and must be no longer than 63 characters'
  },
  email: {
    // double backslash as it seems to be escaped otherwise
    pattern: '^[^\\s@]+@[^\\s@]+\.[^\\s@]+$',
    message: 'Must be a valid email address'
  },
  amazonIamRoleArn: {
    pattern: '^arn:aws:iam::[0-9]{12}:role\\/[\\S]+$',
    message: 'Must be an Amazon IAM role ARN in the format "arn:aws:iam::[account]:role/[role]"'
  }
}
