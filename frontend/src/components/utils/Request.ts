import axios from 'axios'

const server = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 3000,
  // withCredentials: true,
})

server.interceptors.request.use(
  config => {
    // add authorization header with jwt token if available
    const url = config.url || ''
    if (!url.startsWith('/auth/login') && !url.startsWith('/auth/register')) {
      config.headers.Authorization = 'Bearer ' + localStorage.getItem('token')
    }
    return config
  },
  error => {
    console.log('ERROR: ' + error)
    return Promise.reject(error)
  }
)

server.interceptors.response.use(
  response => {
    console.log('Resp:', response.data)
    return response.data
  },
  error => {
    console.log('ERROR: ' + error)
    return Promise.reject(error)
  }
)

export default server
