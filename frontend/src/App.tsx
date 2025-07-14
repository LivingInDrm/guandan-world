import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)
  const [healthStatus, setHealthStatus] = useState<string>('检查中...')

  // 检查后端健康状态
  useEffect(() => {
    const checkHealth = async () => {
      try {
        const response = await fetch('http://localhost:8080/healthz')
        const data = await response.json()
        if (data.status === 'pong') {
          setHealthStatus('服务正常运行 ✅')
        } else {
          setHealthStatus('服务异常 ❌')
        }
      } catch (error) {
        setHealthStatus('服务连接失败 ❌')
      }
    }

    checkHealth()
  }, [])

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>掼蛋世界</h1>
      
      {/* 健康状态检查显示 */}
      <div className="card">
        <p><strong>后端服务状态：</strong>{healthStatus}</p>
      </div>
      
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
