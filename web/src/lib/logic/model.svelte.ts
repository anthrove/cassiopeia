// This file contains the logic to interact with the Tenant-Application API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(): Promise<Model[]> {
    const certificates = await api<Model[]>(`v1/tenant/${get(context)}/model`)
    return certificates.data!
}

export async function read(id:string): Promise<Model> {
    const certificates = await api<Model>(`v1/tenant/${get(context)}/model/${id}`)
    return certificates.data!
}

export async function create(descriptor:Model__create):Promise<Model> {
    const application = await api<Model>(`v1/tenant/${get(context)}/model`,{
        method:'POST',
        body: {
            ...descriptor,
            name:descriptor.display_name
        }
    })
    console.log(application.data);
    return application.data!
}

export async function update(definition: Model__update){
    const response = await api(`v1/tenant/${get(context)}/model/${definition.id}`,{
        method: 'PUT',
        body: definition
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${get(context)}/model/${id}`,{
        method:'DELETE'
    })
}