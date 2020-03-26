module.exports = string => string.trim().replace(/[^a-zA-Z0-9-_\s]/g, '').replace(/\W+/g, '-').toLowerCase()
