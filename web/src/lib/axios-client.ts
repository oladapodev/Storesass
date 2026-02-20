import axios, { type AxiosRequestConfig } from 'axios'

const getBaseUrl = () => {
  try {
    return (import.meta.env.VITE_API_URL || 'http://localhost:8080') + '/api/v1'
  } catch (e) {
    return 'http://localhost:8080/api/v1'
  }
}

const BASE_URL = getBaseUrl()

export const axiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export type ErrorType<Error> = Error

export const customInstance = <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig,
): Promise<T> => {
  const source = axios.CancelToken.source()
  const promise = axiosInstance({
    ...config,
    ...options,
    cancelToken: source.token,
  }).then(({ data }) => data)

  // @ts-expect-error cancel is not in Promise type
  promise.cancel = () => {
    source.cancel('Query was cancelled')
  }

  return promise
}
