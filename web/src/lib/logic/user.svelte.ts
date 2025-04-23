// This file contains the logic to interact with the Tenant-Application API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(): Promise<User[]> {
    const certificates = await api<User[]>(`v1/tenant/${get(context)}/user`)
    return certificates.data!
}

export async function read(id:string): Promise<User> {
    const certificates = await api<User>(`v1/tenant/${get(context)}/user/${id}`)
    return certificates.data!
}

export async function create(descriptor:User__create) {
    const application = await api<User>(`v1/tenant/${get(context)}/user`,{
        method:'POST',
        body: descriptor
    })
    console.log(application.data);
    
}