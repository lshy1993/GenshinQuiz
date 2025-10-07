import useSWR from 'swr'

// 定义 fetcher 函数
const fetcher = async (url: string) => {
    const response = await fetch(url)
    if (!response.ok) {
        throw new Error('请求失败')
    }
    return response.json()
}

// 带认证的 fetcher
const authenticatedFetcher = async (url: string) => {
    const token = localStorage.getItem('token') // 或者从其他地方获取 token
    const response = await fetch(url, {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    })
    if (!response.ok) {
        throw new Error('请求失败')
    }
    return response.json()
}

export { fetcher, authenticatedFetcher }
export default useSWR