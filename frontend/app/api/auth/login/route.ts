import { NextRequest, NextResponse } from 'next/server'

const AUTH0_DOMAIN = process.env.AUTH0_DOMAIN!
const AUTH0_CLIENT_ID = process.env.AUTH0_CLIENT_ID!
const AUTH0_CLIENT_SECRET = process.env.AUTH0_CLIENT_SECRET!

export async function POST(req: NextRequest) {
  try {
    const { email, password } = await req.json()

    if (!email || !password) {
      return NextResponse.json(
        { error: 'Email and password are required' },
        { status: 400 }
      )
    }

    const auth0Res = await fetch(`https://${AUTH0_DOMAIN}/oauth/token`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        grant_type: 'password',
        username: email,
        password,
        client_id: AUTH0_CLIENT_ID,
        client_secret: AUTH0_CLIENT_SECRET,
        connection: 'Username-Password-Authentication',
        scope: 'openid profile email',
      }),
    })

    const data = await auth0Res.json()

    if (!auth0Res.ok) {
      return NextResponse.json(
        {
          error: data.error || 'authentication_failed',
          error_description:
            data.error_description || 'Authentication failed',
        },
        { status: auth0Res.status }
      )
    }

    return NextResponse.json({
      access_token: data.access_token,
      id_token: data.id_token,
      token_type: data.token_type,
      expires_in: data.expires_in,
    })
  } catch {
    return NextResponse.json(
      { error: 'server_error', error_description: 'Internal server error' },
      { status: 500 }
    )
  }
}
