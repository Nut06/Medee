import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import { http, HttpResponse } from 'msw'
import { LoginForm } from '@/components/LoginForm'
import { server } from '@/test/setup'

// Helper to render LoginForm inside a Router (required for useNavigate and Link)
const renderLoginForm = () => {
  return render(
    <MemoryRouter initialEntries={['/login']}>
      <LoginForm />
    </MemoryRouter>
  )
}

describe('LoginForm Integration Tests', () => {
  it('should render email and password fields', () => {
    renderLoginForm()
    expect(screen.getByPlaceholderText('you@example.com')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('********')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /เข้าสู่ระบบ/i })).toBeInTheDocument()
  })

  it('should show validation error if email is invalid', async () => {
    renderLoginForm()

    fireEvent.input(screen.getByPlaceholderText('you@example.com'), {
      target: { value: 'not-an-email' },
    })
    fireEvent.click(screen.getByRole('button', { name: /เข้าสู่ระบบ/i }))

    await waitFor(() => {
      expect(screen.getByText(/please enter a valid email/i)).toBeInTheDocument()
    })
  })

  it('should show validation error if password is too short', async () => {
    renderLoginForm()

    fireEvent.input(screen.getByPlaceholderText('you@example.com'), {
      target: { value: 'user@test.com' },
    })
    fireEvent.input(screen.getByPlaceholderText('********'), {
      target: { value: '123' },
    })
    fireEvent.click(screen.getByRole('button', { name: /เข้าสู่ระบบ/i }))

    await waitFor(() => {
      expect(screen.getByText(/at least 6 characters/i)).toBeInTheDocument()
    })
  })

  it('should display error message when login API returns 401', async () => {
    // Override handler for this test only - simulate wrong credentials
    server.use(
      http.post('*/auth/login', () => {
        return HttpResponse.json({ message: 'Invalid credentials' }, { status: 401 })
      })
    )

    renderLoginForm()

    fireEvent.input(screen.getByPlaceholderText('you@example.com'), {
      target: { value: 'wrong@test.com' },
    })
    fireEvent.input(screen.getByPlaceholderText('********'), {
      target: { value: 'wrongpass' },
    })
    fireEvent.click(screen.getByRole('button', { name: /เข้าสู่ระบบ/i }))

    await waitFor(() => {
      expect(screen.getByText(/invalid credentials/i)).toBeInTheDocument()
    })
  })

  it('should call login API with correct payload on valid submit', async () => {
    const loginSpy = vi.fn()

    server.use(
      http.post('*/auth/login', async ({ request }) => {
        const body = await request.json()
        loginSpy(body)
        return HttpResponse.json({
          user: { id: '1', firstName: 'Test', email: 'user@test.com' },
          companies: [],
          accessToken: 'fake-token',
        })
      })
    )

    renderLoginForm()

    fireEvent.input(screen.getByPlaceholderText('you@example.com'), {
      target: { value: 'user@test.com' },
    })
    fireEvent.input(screen.getByPlaceholderText('********'), {
      target: { value: 'password123' },
    })
    fireEvent.click(screen.getByRole('button', { name: /เข้าสู่ระบบ/i }))

    await waitFor(() => {
      expect(loginSpy).toHaveBeenCalledWith({
        email: 'user@test.com',
        password: 'password123',
      })
    })
  })
})
