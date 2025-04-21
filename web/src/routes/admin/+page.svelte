<script lang="ts">
    import { useApiFetch } from "$lib/composables/apifetch";
    const api = useApiFetch();
    import { onMount } from "svelte";

    import SearchableTable from "$lib/components/searchableTable.svelte";

    let tenants: Tenant[] = $state([]);
    onMount(async () => {
        tenants = await api<Tenant[]>("v1/tenant");
        console.log("table loaded");
        console.log(tenants);
    });

    import Modal from "$lib/components/modal.svelte";
    import Input from "$lib/components/forms/input.svelte";
    import Form from "$lib/components/forms/Form.svelte";
    import Button from "$lib/components/Button.svelte";

    let createTenantModal_open = $state(false);
    function openCreateTenantModal() {
        createTenantModal_open = true;
    }

    const createTenantModal_descriptor: FormDescriptor = {
        display_name: String,
        password_type: ["bcrypt"],
    };

    import { create as createTenant } from "$lib/logic/tenant.svelte";
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
                <Button href="/tennant/{item.id}">
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
</div>
