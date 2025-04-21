// This file contains the logic to interact with the Tenant API

import { useApiFetch } from "$lib/composables/apifetch";

const api = useApiFetch()

// Creates a new Tennant
export async function create(definition:Tenant__create) {
    const response = await api('v1/tenant/',{
        method: 'POST',
        body: definition
    })

    console.log(response);
    
}