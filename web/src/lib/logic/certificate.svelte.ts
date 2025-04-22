// This file contains the logic to interact with the Certificate API

import { useApiFetch } from "$lib/composables/apifetch";
import { context } from "$lib/stores/context.svelte";
import { get } from "svelte/store";

const api = useApiFetch()

export async function readAll(): Promise<Certificate[]> {
    const certificates = await api<Certificate[]>(`v1/tenant/${get(context)}/certificate`)
    return certificates.data!
}