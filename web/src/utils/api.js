import { classifyError, ErrorTypes } from './errors'

const DEFAULT_TIMEOUT = 10000
const MAX_RETRIES = 3
const RETRY_DELAY = 1000

export class ApiClient {
  constructor(baseUrl = '', options = {}) {
    this.baseUrl = baseUrl
    this.timeout = options.timeout || DEFAULT_TIMEOUT
    this.maxRetries = options.maxRetries || MAX_RETRIES
    this.retryDelay = options.retryDelay || RETRY_DELAY
    this.headers = options.headers || {}
  }

  withTimeout(promise, ms) {
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        const error = new Error('Request timeout')
        error.name = 'AbortError'
        reject(error)
      }, ms)

      promise
        .then((value) => {
          clearTimeout(timer)
          resolve(value)
        })
        .catch((error) => {
          clearTimeout(timer)
          reject(error)
        })
    })
  }

  getRetryDelay(attempt) {
    const baseDelay = this.retryDelay
    const exponentialDelay = baseDelay * Math.pow(2, attempt)
    const jitter = Math.random() * 500
    return Math.min(exponentialDelay + jitter, 30000)
  }

  async fetch(path, options = {}) {
    const url = this.baseUrl + path
    const timeout = options.timeout || this.timeout

    const fetchOptions = {
      ...options,
      headers: {
        ...this.headers,
        ...options.headers
      }
    }

    if (options.body && typeof options.body === 'object') {
      fetchOptions.body = JSON.stringify(options.body)
      fetchOptions.headers['Content-Type'] = 'application/json'
    }

    const response = await this.withTimeout(fetch(url, fetchOptions), timeout)

    if (!response.ok) {
      const error = new Error(`HTTP ${response.status}: ${response.statusText}`)
      error.status = response.status
      error.headers = response.headers
      throw error
    }

    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      return response.json()
    }
    return response.text()
  }

  async fetchWithRetry(path, options = {}, retries = this.maxRetries) {
    let lastError

    for (let attempt = 0; attempt <= retries; attempt++) {
      try {
        return await this.fetch(path, options)
      } catch (error) {
        lastError = error
        const classified = classifyError(error)

        if (!classified.retryable || attempt === retries) {
          throw error
        }

        const delay = classified.retryAfter || this.getRetryDelay(attempt)
        await this.sleep(delay)
      }
    }

    throw lastError
  }

  sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms))
  }

  get(path, options = {}) {
    return this.fetchWithRetry(path, { ...options, method: 'GET' })
  }

  post(path, body, options = {}) {
    return this.fetchWithRetry(path, { ...options, method: 'POST', body })
  }

  put(path, body, options = {}) {
    return this.fetchWithRetry(path, { ...options, method: 'PUT', body })
  }

  delete(path, options = {}) {
    return this.fetchWithRetry(path, { ...options, method: 'DELETE' })
  }
}

export const api = new ApiClient()
