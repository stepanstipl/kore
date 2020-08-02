const HOURLY_CURR_FORMATTER = new Intl.NumberFormat(undefined, { style: 'currency', currency: 'USD', currencyDisplay: 'narrowSymbol', maximumFractionDigits: 3 }).format
const DAILY_CURR_FORMATTER = new Intl.NumberFormat(undefined, { style: 'currency', currency: 'USD', currencyDisplay: 'narrowSymbol' }).format
const MONTHLY_CURR_FORMATTER = new Intl.NumberFormat(undefined, { style: 'currency', currency: 'USD', currencyDisplay: 'narrowSymbol', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format
const MICRODOLLARS_IN_DOLLAR = 1000000
const HOURS_IN_DAY = 24
const HOURS_IN_MONTH = 730

export function formatHourlyCost(c) {
  return `${HOURLY_CURR_FORMATTER((c)/MICRODOLLARS_IN_DOLLAR)}/hr`
}
export function formatDailyCost(c) {
  return `${DAILY_CURR_FORMATTER((c*HOURS_IN_DAY)/MICRODOLLARS_IN_DOLLAR)}/day`
}
export function formatMonthlyCost(c) {
  return `${MONTHLY_CURR_FORMATTER((c*HOURS_IN_MONTH)/MICRODOLLARS_IN_DOLLAR)}/mo`
}
export function formatCost(c) {
  return DAILY_CURR_FORMATTER(c/MICRODOLLARS_IN_DOLLAR)
}