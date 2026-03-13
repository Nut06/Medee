import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import { http, HttpResponse } from 'msw'
import { RegisterForm } from '@/components/RegisterForm'
import { server } from '@/test/setup'

const renderRegisterForm = () => {
  return render(
    <MemoryRouter initialEntries={['/register']}>
      <RegisterForm />
    </MemoryRouter>
  )
}

describe('RegisterForm Integration Tests', () => {
  it('should render all required fields', () => {
    renderRegisterForm()
    expect(screen.getByPlaceholderText('สมชาย')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('ใจดี')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('you@example.com')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('********')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /ลงทะเบียน/i })).toBeInTheDocument()
  })

  it('should show validation errors when all fields are empty', async () => {
    renderRegisterForm()
    fireEvent.click(screen.getByRole('button', { name: /ลงทะเบียน/i }))

    await waitFor(() => {
      expect(screen.getByText(/กรุณากรอกชื่อ/i)).toBeInTheDocument()
      expect(screen.getByText(/กรุณากรอกนามสกุล/i)).toBeInTheDocument()
    })
  })

  it('should call register API with correct payload on valid submit', async () => {
    const registerSpy = vi.fn()

    server.use(
      http.post('*/auth/register', async ({ request }) => {
        const body = await request.json()
        registerSpy(body)
        return HttpResponse.json({
          user: { id: '1', firstName: 'สมชาย', email: 'somchai@test.com' },
          companies: [],
        })
      })
    )

    renderRegisterForm()

    fireEvent.input(screen.getByPlaceholderText('สมชาย'), { target: { value: 'สมชาย' } })
    fireEvent.input(screen.getByPlaceholderText('ใจดี'), { target: { value: 'ใจดี' } })
    fireEvent.input(screen.getByPlaceholderText('you@example.com'), {
      target: { value: 'somchai@test.com' },
    })
    fireEvent.input(screen.getByPlaceholderText('********'), {
      target: { value: 'securepass123' },
    })

    fireEvent.click(screen.getByRole('button', { name: /ลงทะเบียน/i }))

    await waitFor(() => {
      expect(registerSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          email: 'somchai@test.com',
          firstName: 'สมชาย',
          lastName: 'ใจดี',
        })
      )
    })
  })
})
