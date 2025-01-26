import axios from 'axios'
import type { ApiResponse } from './Types'
import type { AxiosInstance, AxiosResponse } from 'axios'

class APIService {

  private readonly api: AxiosInstance

  public constructor() {
    const axiosConfig = {
      baseURL: import.meta.env.VITE_API_BASE_URL,
      timeout: 3000,
      withCredentials: true,
    }
    console.log('creating apiService with axiosConfig:', axiosConfig)

    const api = axios.create(axiosConfig)
    api.interceptors.request.use(
      config => {
        // add authorization header with jwt token if available
        const url = config.url || ''
        if (!url.startsWith('/auth/login') && !url.startsWith('/auth/register')) {
          config.headers.Authorization = 'Bearer ' + localStorage.getItem('token')
        }
        return config
      },
      error => {
        console.log('REQ HTTP_ERROR: ' + error)
        return Promise.reject(error)
      }
    )

    api.interceptors.response.use(
      (response: AxiosResponse<ApiResponse<any>>) => {
        return response
      },
      error => {
        return Promise.reject(error)
      }
    )
    this.api = api
  }

  public async POST<T> (url: string, body: any): Promise<T> {
    return this.api
      .post<ApiResponse<T>>(url, body)
      .then((response) => {
        const { data } = response
        if (data.code == 0) {
          return response.data.data
        }
        return Promise.reject("Buzz failed[POST], code:" + data.code + ", msg:" + data.msg)
      })
      .catch((error) => {
        throw error;
      });
  }

  public async GET<T> (url: string, urlParams?: Record<string, any>): Promise<T> {
    if (urlParams) {
      const queryString = new URLSearchParams(urlParams).toString();
      url += `?${queryString}`;
    }
    return this.api
      .get<ApiResponse<T>>(url)
      .then((response) => {
        const { data } = response
        if (data.code == 0) {
          return response.data.data
        }
        return Promise.reject("Buzz failed[GET], code:" + data.code + ", msg:" + data.msg)
      })
      .catch((error) => {
        throw error;
      });
  }
}

const API = new APIService()
export default API