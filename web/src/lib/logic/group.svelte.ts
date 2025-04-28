// This file contains the logic to interact with the Tenant-Application API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(): Promise<Group[]> {
    const groups = await api<Group[]>(`v1/tenant/${get(context)}/group`)
    return groups.data!.map(
        g=>{
            // TODO: See issue https://github.com/anthrove/identity/issues/29
            //@ts-expect-error temporary overwite until API is fixed
            g.display_name = g.displayName
            return g
        }
    )
}

export async function read(id:string): Promise<Group> {
    const groups = await api<Group>(`v1/tenant/${get(context)}/group/${id}`)
    // TODO: See issue https://github.com/anthrove/identity/issues/29
    //@ts-expect-error temporary overwite until API is fixed
    groups.data!.display_name = groups.data.displayName

    return groups.data!
}

export async function create(descriptor:Group__create) {
    const application = await api<Group>(`v1/tenant/${get(context)}/group`,{
        method:'POST',
        body: descriptor
    })
    console.log(application.data);
    
}

export async function update(definition: Group__update){
    const response = await api(`v1/tenant/${get(context)}/group/${definition.id}`,{
        method: 'PUT',
        body: definition
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${get(context)}/group/${id}`,{
        method:'DELETE'
    })
}