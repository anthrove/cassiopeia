// This file contains the logic to interact with the Tenant-Application API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(): Promise<Application[]> {
    const certificates = await api<Application[]>(`v1/tenant/${get(context)}/application`)
    return certificates.data!
}

export async function read(id:string): Promise<Application> {
    const certificates = await api<Application>(`v1/tenant/${get(context)}/application/${id}`)
    return certificates.data!
}

export async function create(descriptor:Application__create) {
    const application = await api<Application>(`v1/tenant/${get(context)}/application`,{
        method:'POST',
        body: descriptor
    })
    console.log(application.data);
    
}