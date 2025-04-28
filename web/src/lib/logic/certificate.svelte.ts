// This file contains the logic to interact with the Certificate API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(id:string=''): Promise<Certificate[]> {
    const certificates = await api<Certificate[]>(`v1/tenant/${id||get(context)}/certificate`)
    return certificates.data!
}


export async function read(id:string): Promise<Certificate> {
    const certificates = await api<Certificate>(`v1/tenant/${get(context)}/certificate/${id}`)
    return certificates.data!
}

export async function create(descriptor:Certificate__create) {
    const application = await api<Certificate>(`v1/tenant/${get(context)}/certificate`,{
        method:'POST',
        body: descriptor
    })
    console.log(application.data);
    
}

export async function update(definition: Certificate__update){
    const response = await api(`v1/tenant/${get(context)}/certificate/${definition.id}`,{
        method: 'PUT',
        body: definition
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${get(context)}/certificate/${id}`,{
        method:'DELETE'
    })
}