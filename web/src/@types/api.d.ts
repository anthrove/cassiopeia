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
  updated_at:string
}