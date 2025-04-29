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

export async function readAll():Promise<Tenant[]> {
    const raw = await api<Tenant[]>("v1/tenant");
    return raw.data!
}

export async function read(id:string) {
    const raw = await api<Tenant>(`v1/tenant/${id}`);
    return raw.data!
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