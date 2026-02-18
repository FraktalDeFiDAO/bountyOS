export const ErrorTypes = {
  NETWORK: 'network',
  TIMEOUT: 'timeout',
  SERVER: 'server',
  VALIDATION: 'validation',
  UNKNOWN: 'unknown'
}

export function classifyError(error) {
  if (!error) {
    return { type: ErrorTypes.UNKNOWN, message: 'Unknown error', retryable: false, retryAfter: null }
  }

  if (error.name === 'AbortError' || error.code === 'UND_ERR_CONNECT_TIMEOUT') {
    return { type: ErrorTypes.TIMEOUT, message: 'Request timed out', retryable: true, retryAfter: 1000 }
  }

  if (error.name === 'TypeError' && error.message.includes('fetch')) {
    return { type: ErrorTypes.NETWORK, message: 'Network error - check your connection', retryable: true, retryAfter: 1000 }
  }

  if (error.status || error.statusCode) {
    const status = error.status || error.statusCode
    if (status >= 500) {
      return { type: ErrorTypes.SERVER, message: `Server error (${status})`, retryable: true, retryAfter: 2000 }
    }
    if (status === 429) {
      const retryAfter = error.headers?.get?.('retry-after') || 5
      return { type: ErrorTypes.SERVER, message: 'Rate limited - please wait', retryable: true, retryAfter: parseInt(retryAfter) * 1000 }
    }
    if (status >= 400 && status < 500) {
      return { type: ErrorTypes.VALIDATION, message: `Request error (${status})`, retryable: false, retryAfter: null }
    }
  }

  if (error.message) {
    const msg = error.message.toLowerCase()
    if (msg.includes('network') || msg.includes('connection') || msg.includes('offline')) {
      return { type: ErrorTypes.NETWORK, message: error.message, retryable: true, retryAfter: 1000 }
    }
    if (msg.includes('timeout')) {
      return { type: ErrorTypes.TIMEOUT, message: error.message, retryable: true, retryAfter: 1000 }
    }
  }

  return { type: ErrorTypes.UNKNOWN, message: error.message || 'Unknown error', retryable: false, retryAfter: null }
}

export function formatError(error) {
  const classified = classifyError(error)
  const messages = {
    [ErrorTypes.NETWORK]: 'Unable to connect. Please check your internet connection.',
    [ErrorTypes.TIMEOUT]: 'The request took too long. Please try again.',
    [ErrorTypes.SERVER]: 'Server is having issues. Please try again later.',
    [ErrorTypes.VALIDATION]: 'There was a problem with your request.',
    [ErrorTypes.UNKNOWN]: 'Something went wrong. Please try again.'
  }
  return messages[classified.type] || classified.message
}

export function isRetryable(error) {
  const classified = classifyError(error)
  return classified.retryable
}
