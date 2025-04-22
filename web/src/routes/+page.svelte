<script lang="ts">
    import { onMount } from "svelte";

    import SearchableTable from "$lib/components/searchableTable.svelte";
    import { create as createTenant, readAll as readAllTenants } from "$lib/logic/tenant.svelte";

    let tenants: Tenant[] = $state([]);
    onMount(async () => {
        tenants = await readAllTenants();
    });

    import Modal from "$lib/components/modal.svelte";
    import Form from "$lib/components/forms/Form.svelte";
    import Button from "$lib/components/Button.svelte";

    let createTenantModal_open = $state(false);
    function openCreateTenantModal() {
        createTenantModal_open = true;
    }

    import apiconfig from '../data/apidescriptor.json'

    const createTenantModal_descriptor: FormDescriptor = {
        display_name: String,
        password_type: apiconfig.tenant.password_types,
    };

    import { context } from "$lib/stores/context.svelte";
</script>

<div class="container mx-auto mt-16">
    <Modal bind:open={createTenantModal_open}>
        <div class="prose">
            <h2>Create new Tenant</h2>
            <Form
                submit="Create tennant"
                descriptor={createTenantModal_descriptor}
                onsubmit={createTenant}
            />
        </div>
    </Modal>

    <SearchableTable rows={tenants}>
        {#snippet header()}
            <th class="text-start p-4">Tennant</th>
            <th class="text-start">Password Type</th>
        {/snippet}
        {#snippet row(item: Tenant)}
            <td class="p-4">
                {item.display_name}
            </td>
            <td>
                {item.password_type}
            </td>
            <td>
                <Button href="/tenant/" onclick={()=>context.set(item.id)} class="inline-block">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        ><path
                            fill="currentColor"
                            d="M3 21v-4.25L16.2 3.575q.3-.275.663-.425t.762-.15t.775.15t.65.45L20.425 5q.3.275.438.65T21 6.4q0 .4-.137.763t-.438.662L7.25 21zM17.6 7.8L19 6.4L17.6 5l-1.4 1.4z"
                        /></svg
                    >
                </Button>
            </td>
        {/snippet}
        {#snippet empty()}
            <div
                class="w-full p-8 flex items-center justify-center flex-col gap-4"
            >
                <i class="block">No Tenants yet ...</i>
                <Button onclick={openCreateTenantModal}
                    >Create your first Tenant</Button
                >
            </div>
        {/snippet}
    </SearchableTable>

    {#if tenants?.length}
        <Button onclick={openCreateTenantModal} class="mt-4"
            >Create a Tenant</Button
        >
    {/if}
</div>
