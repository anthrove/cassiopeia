import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readCategories() {
    const categories = await api<string[]>(`v1/tenant/${get(context)}/provider/category`)
    return categories.data!
}

export async function readTypes(category:string) {
    const types = await api<string[]>(`v1/tenant/${get(context)}/provider/category/${category}`)
    return types.data!
}

export async function readConfig(category:string,type:string) {
    const config = await api<[{field_key:string,field_type:'text'|'int'|'bool'|'secret'}]>(`v1/tenant/${get(context)}/provider/category/${category}/${type}`)
    return config.data!
}

export async function create(descriptor:Provider__create):Promise<Provider> {
    const application = await api<Model>(`v1/tenant/${get(context)}/provider`,{
        method:'POST',
        body: descriptor
    })
    console.log(application.data);
    return application.data!
}

export async function readAll():Promise<Provider[]> {
    const raw = await api<Tenant[]>(`v1/tenant/${get(context)}/provider`);
    return raw.data!
}

export async function update(definition: Provider_update){
    const response = await api(`v1/tenant/${get(context)}/provider/${definition.id}`,{
        method: 'PUT',
        body: definition
    })
    console.log(response);
}

export async function kill(id:string) {
    await api(`v1/tenant/${get(context)}/provider/${id}`,{
        method:'DELETE'
    })
}