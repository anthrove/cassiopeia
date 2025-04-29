// This file contains the logic to interact with the Tenant-Application API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(): Promise<Adapter[]> {
    const certificates = await api<Adapter[]>(`v1/tenant/${get(context)}/adapter`)
    return certificates.data!
}

export async function read(id:string): Promise<Adapter> {
    const certificates = await api<Adapter>(`v1/tenant/${get(context)}/adapter/${id}`)
    return certificates.data!
}

export async function create(descriptor:Adapter__create): Promise<Adapter> {
    const application = await api<Adapter>(`v1/tenant/${get(context)}/adapter`,{
        method:'POST',
        body: descriptor
    })
    console.log(application.data);
    return application.data!
    
}

export async function update(definition: Adapter__update){
    const response = await api(`v1/tenant/${get(context)}/adapter/${definition.id}`,{
        method: 'PUT',
        body: definition
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${get(context)}/adapter/${id}`,{
        method:'DELETE'
    })
}