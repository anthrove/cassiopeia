// This file contains the logic to interact with the Tenant-Application API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(): Promise<Enforcer[]> {
    const certificates = await api<Enforcer[]>(`v1/tenant/${get(context)}/enforcer`)
    return certificates.data!.map(e=>{
        return {
            //@ts-expect-error polyfil broken api name
            display_name: e?.name,
            ...e
        }
    })
}

export async function read(id:string): Promise<Enforcer> {
    const enforcer = await api<Enforcer>(`v1/tenant/${get(context)}/enforcer/${id}`)
    
    return {
        //@ts-expect-error polyfil broken api name
        display_name: enforcer.data?.name,
        ...enforcer.data!,
    }
}

export async function create(descriptor:Enforcer__create) {
    const application = await api<Enforcer>(`v1/tenant/${get(context)}/enforcer`,{
        method:'POST',
        body: descriptor
    })
    console.log(application.data);
    
}

export async function update(descriptor: Enforcer__update){
    const response = await api(`v1/tenant/${get(context)}/enforcer/${descriptor.id}`,{
        method: 'PUT',
        body: {
            name: descriptor.display_name,
            ...descriptor
        },
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${get(context)}/enforcer/${id}`,{
        method:'DELETE'
    })
}