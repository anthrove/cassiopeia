// --- Meta-Types ---
type APIResponse<T> = {
  data: T,
  error: null
} | {
  data: null,
  error: string
}

interface CRUDAPI<Base = any, Create = any, Update = any> {
  create(descriptor: Create): Promise<Base>
  readAll(): Promise<Base[]>
  read(id:string): Promise<Base>
  update(descriptor: Update): Promise<void>
  kill(id: string): Promise<void>
}

// --- Tenant ---
type Tenant = {
  id: string
  display_name: string
  created_at: string
  updated_at: string
  password_type: string
  signing_certificate_id: string
}

type Tenant__create = Pick<Tenant, 'display_name' | 'password_type'>
type Tenant__update = Pick<Tenant, 'id' | 'display_name' | 'password_type' | 'signing_certificate_id'>

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

type Certificate__create = Pick<Certificate, 'algorithm' | 'bit_size' | 'display_name' | 'expired_at'>
type Certificate__update = Pick<Certificate, 'id' | 'display_name'>
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
  id: string,
  display_name: string,
  parent_group_id: string,
  tenant_id: string,
  created_at: string,
  updated_at: string,
}

type Group__create = Pick<Group, 'display_name' | 'parent_group_id'>
type Group__update = Pick<Group, 'id' | 'display_name' | 'parent_group_id'>
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

type User__create = Pick<User, 'display_name' | 'email' | 'username'> & { password: string }

type User__update = Pick<User, 'id' | 'display_name'>

// --- Enforcers ---

type Enforcer = {
  id: string,
  adapter_id: string,
  description: string,
  model_id: string,
  display_name: string,
  tenant_id: string
}

type Enforcer__create = {
  adapter_id: string,
  description: string,
  model_id: string,
  display_name: string,
}

type Enforcer__update = {
  id: string,
  adapter_id: string,
  description: string,
  model_id: string,
  display_name: string,
}

// --- Adapters ---

type Adapter = {
  id: string,
  database_name: string,
  driver: string,
  external_db: boolean,
  host: string,
  display_name: string,
  password: string,
  port: string,
  table_name: string,
  tenant_id: string,
  username: string
}

type Adapter__create = {
  database_name: string,
  driver: string,
  external_db: boolean,
  host: string,
  display_name: string,
  password: password,
  port: number,
  table_name: string,
  username: string
}

type Adapter__update = {
  id: string
}

// --- Models ---

type Model = {
  id: string,
  description: string,
  model: string,
  display_name: string,
  tenant_id: string
}

type Model__create = {
  description: string,
  model: string,
  display_name: string,
}

type Model__update = {
  id: string,
  description: string,
  model: string,
  display_name: string,
}