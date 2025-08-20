import { useState } from 'react'
import './App.css'
import axios from 'axios'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'

const api = axios.create({ baseURL: import.meta.env.VITE_API_BASE || 'http://localhost:8080/api' })

const signupSchema = z.object({
  name: z.string().min(1).max(100),
  agent_id: z.string().min(3).max(64),
  password: z.string().min(8).max(128),
})

type SignupForm = z.infer<typeof signupSchema>

const loginSchema = z.object({
  agent_id: z.string().min(3).max(64),
  password: z.string().min(8).max(128),
})

type LoginForm = z.infer<typeof loginSchema>

function Signup() {
  const { register, handleSubmit, formState: { errors, isSubmitting }, reset } = useForm<SignupForm>({ resolver: zodResolver(signupSchema) })
  const [message, setMessage] = useState<string | null>(null)
  const onSubmit = async (data: SignupForm) => {
    setMessage(null)
    try {
      await api.post('/signup', data)
      setMessage('Signup successful')
      reset()
    } catch (err: any) {
      setMessage(err?.response?.data?.error || 'Signup failed')
    }
  }
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="card">
      <h2>Create Account</h2>
      <label>Name</label>
      <input placeholder="Your name" {...register('name')} />
      {errors.name && <span className="error">{errors.name.message}</span>}

      <label>Agent ID</label>
      <input placeholder="agent id" {...register('agent_id')} />
      {errors.agent_id && <span className="error">{errors.agent_id.message}</span>}

      <label>Password</label>
      <input type="password" placeholder="password" {...register('password')} />
      {errors.password && <span className="error">{errors.password.message}</span>}

      <button disabled={isSubmitting}>{isSubmitting ? 'Submitting...' : 'Sign up'}</button>
      {message && <div className="notice">{message}</div>}
    </form>
  )
}

function Login() {
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<LoginForm>({ resolver: zodResolver(loginSchema) })
  const [message, setMessage] = useState<string | null>(null)
  const [user, setUser] = useState<any>(null)
  const onSubmit = async (data: LoginForm) => {
    setMessage(null)
    setUser(null)
    try {
      const res = await api.post('/login', data)
      setUser(res.data)
      setMessage('Login successful')
    } catch (err: any) {
      setMessage(err?.response?.data?.error || 'Login failed')
    }
  }
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="card">
      <h2>Login</h2>
      <label>Agent ID</label>
      <input placeholder="agent id" {...register('agent_id')} />
      {errors.agent_id && <span className="error">{errors.agent_id.message}</span>}

      <label>Password</label>
      <input type="password" placeholder="password" {...register('password')} />
      {errors.password && <span className="error">{errors.password.message}</span>}

      <button disabled={isSubmitting}>{isSubmitting ? 'Submitting...' : 'Login'}</button>
      {message && <div className="notice">{message}</div>}
      {user && (
        <pre className="result">{JSON.stringify(user, null, 2)}</pre>
      )}
    </form>
  )
}

function App() {
  const [tab, setTab] = useState<'login' | 'signup'>('login')
  return (
    <div className="container">
      <header>
        <h1>RD Website</h1>
        <p className="subtitle">Golang + PostgreSQL starter with secure auth</p>
      </header>
      <div className="tabs">
        <button className={tab === 'login' ? 'active' : ''} onClick={() => setTab('login')}>Login</button>
        <button className={tab === 'signup' ? 'active' : ''} onClick={() => setTab('signup')}>Sign up</button>
      </div>
      {tab === 'login' ? <Login /> : <Signup />}
      <footer>
        <small>Configure API base via VITE_API_BASE. Default: http://localhost:8080/api</small>
      </footer>
    </div>
  )
}

export default App
