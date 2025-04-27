// This file contains the logic to interact with the Tenant API

import { useApiFetch } from "$lib/composables/apifetch";

const api = useApiFetch()

// Creates a new Tennant
export async function create(definition:Tenant__create) {
    const response = await api('v1/tenant',{
        method: 'POST',
        body: definition
    })
    console.log(response);   
}

function normalizeTenantFields(tenant:Tenant):Tenant{
    //@ts-expect-error in the API this key is named worng. we force normalize it here
    tenant.signing_certificate_id = tenant.signing_key_id
    //@ts-expect-error this field was renamed we do not want orphans
    delete tenant.signing_key_id
    return tenant
}

export async function readAll():Promise<Tenant[]> {
    const raw = await api<Tenant[]>("v1/tenant");
    return raw.data!.map(normalizeTenantFields) || []
}

export async function read(id:string) {
    const raw = await api<Tenant>(`v1/tenant/${id}`);
    return normalizeTenantFields(raw.data!)
}

export async function update(definition: Tenant__update){
    const response = await api(`v1/tenant/${definition.id}`,{
        method: 'PUT',
        body: definition
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${id}`,{
        method:'DELETE'
    })
}