// This file contains the logic to interact with the Tenant-Application API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(): Promise<Enforcer[]> {
    const certificates = await api<Enforcer[]>(`v1/tenant/${get(context)}/enforcer`)
    return certificates.data!
}

export async function read(id:string): Promise<Enforcer> {
    const certificates = await api<Enforcer>(`v1/tenant/${get(context)}/enforcer/${id}`)
    return certificates.data!
}

export async function create(descriptor:Enforcer__create) {
    const application = await api<Enforcer>(`v1/tenant/${get(context)}/enforcer`,{
        method:'POST',
        body: descriptor
    })
    console.log(application.data);
    
}

export async function update(definition: Enforcer__update){
    const response = await api(`v1/tenant/${get(context)}/enforcer/${definition.id}`,{
        method: 'PUT',
        body: definition
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${get(context)}/enforcer/${id}`,{
        method:'DELETE'
    })
}