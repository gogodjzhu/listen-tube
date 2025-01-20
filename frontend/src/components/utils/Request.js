import axios from 'axios'

const server = axios.create({
  baseURL: 'http://127.0.0.1:8080/',
  timeout: 3000
})

server.interceptors.request.use(
  config => {
    // add authorization header with jwt token if available
    const url = config.url
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
    console.log('ERROR: ' + (error.response ? error.response.data : error.message))
    return Promise.reject(new Error(error.response ? error.response.data.message : error.message))
  }
)

export default server
