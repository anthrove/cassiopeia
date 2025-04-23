type Tenant = {
  id: string
  display_name: string
  created_at: string
  updated_at: string
  password_type: string
  signing_certificate_id: string
}

type Tenant__create = Pick<Tenant, 'display_name' | 'password_type'>
type Tenant__update = Pick<Tenant, 'display_name' | 'password_type' | 'signing_certificate_id'>

// --- Certificates ---

type Certificate = {
  id: string
  display_name: string
  algorithm: string
  bit_size: number
  certificate: string
  created_at: string,
  display_name: string
  expired_at: string
  private_key: string
  tenant_id: string
  updated_at: string
}

// --- Applications ---

type Application = {
  certificate_id: string
  client_secret: string
  createdAt: string
  display_name: string
  forget_url: string
  id: string
  logo: string
  redirect_urls: string[]
  sign_in_url: string
  sign_up_url: string
  tenant_id: string
  terms_url: string
  updatedAt: string
}

type Application__create = Pick<Application, 'display_name' | 'certificate_id' | 'forget_url' | 'logo' | 'redirect_urls' | 'sign_in_url' | 'sign_up_url' | 'terms_url'>

// --- Groups ---
type Group = {
  createdAt: string,
  displayName: string,
  id: string,
  parent_group_id: string,
  tenant_id: string,
  updatedAt: string,
  users: [
    string
  ]
}
// --- Users ---

type User = {
  created_at: string,
  deleted_at: {
    time: string,
    valid: boolean
  },
  display_name: string,
  email: string,
  email_verified: true,
  groups: Group[],
  id: string,
  tenant_id: string,
  updated_at: string,
  username: string
}

type User__create = Pick<User,'display_name' | 'email' | 'username'> & {password: string}