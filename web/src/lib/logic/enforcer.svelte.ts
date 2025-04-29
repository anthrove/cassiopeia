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
    const certificates = await api<Enforcer>(`v1/tenant/${get(context)}/enforcer/${id}`)
    console.log(certificates);
    
    return {
        //@ts-expect-error polyfil broken api name
        display_name: certificates.data?.name,
        ...certificates.data!,
    }
}

export async function create(descriptor:Enforcer__create) {
    const application = await api<Enforcer>(`v1/tenant/${get(context)}/enforcer`,{
        method:'POST',
        body: {
            ...descriptor,
            name: descriptor.display_name,
        }
    })
    console.log(application.data);
    
}

export async function update(descriptor: Enforcer__update){
    const response = await api(`v1/tenant/${get(context)}/enforcer/${descriptor.id}`,{
        method: 'PUT',
        body: {
            ...descriptor,
            name: descriptor.display_name,
        }
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${get(context)}/enforcer/${id}`,{
        method:'DELETE'
    })
}