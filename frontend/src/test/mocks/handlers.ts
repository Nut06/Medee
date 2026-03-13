import { http, HttpResponse } from 'msw'

const mockUser = {
  id: 'test-user-id-123',
  firstName: 'Integration',
  lastName: 'Tester',
  email: 'test@example.com',
  avatarURL: null,
  bio: null,
  linkedInURL: null,
  githubURL: null,
  websiteURL: null,
}

export const handlers = [
  // --- Auth Endpoints ---
  http.post('*/auth/login', () => {
    return HttpResponse.json({
      user: mockUser,
      companies: [],
      accessToken: 'fake-access-token-xyz',
    })
  }),

  http.post('*/auth/register', () => {
    return HttpResponse.json({
      user: mockUser,
      companies: [],
    })
  }),

  http.post('*/auth/logout', () => {
    return HttpResponse.json({ message: 'Logged out' })
  }),

  http.post('*/auth/refresh', () => {
    return HttpResponse.json({ accessToken: 'refreshed-fake-token' })
  }),

  // --- User Endpoints ---
  http.get('*/user', () => {
    return HttpResponse.json(mockUser)
  }),

  http.put('*/user', async ({ request }) => {
    const body = await request.json() as Record<string, string>
    return HttpResponse.json({ ...mockUser, ...body })
  }),
]
