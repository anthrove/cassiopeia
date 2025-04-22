<script lang="ts">
    import Form from "$lib/components/forms/Form.svelte";
    import { context } from "$lib/stores/context.svelte";
    import { onMount } from "svelte";

    import {
        update as updateTenant,
        read as readTenant,
        kill as deleteTenant
    } from "$lib/logic/tenant.svelte";
    import { readAll as readAllCertificates } from "$lib/logic/certificate.svelte";

    let tenantData = $state<Tenant>(<any>null);
    let _tenantData: Tenant;

    let certificates = $state<Certificate[]>([]);
    import apidescriptor from "../../data/apidescriptor.json";
    import Button from "$lib/components/Button.svelte";
    import Modal from "$lib/components/modal.svelte";
    import { goto } from "$app/navigation";

    let tenantData_descriptor: FormDescriptor = $state(<any>{});
    async function loadData(ctx: string) {
        certificates = await readAllCertificates();

        tenantData_descriptor = <FormDescriptor>{
            id: {
                type: "string",
                label: "Tenant ID",
                readonly: true,
            },
            display_name: {
                type: "string",
                label: "Display name",
                required: true,
            },
            password_type: {
                type: "select:single",
                label: "Password type",
                required: true,

                options: apidescriptor.tenant.password_types,
            },
            signing_certificate_id: {
                type: "select:single",
                labeler: (id: string) =>
                    `${certificates.find((c: Certificate) => c.id == id)?.display_name} (${id})`,
                label: "Signing key",
                required: true,

                options: certificates.map((c) => c.id),
            },

            created_at: {
                type: "string",
                label: "Created at",
                readonly: true,
            },
            updated_at: {
                type: "string",
                label: "Last update",
                readonly: true,
            },
        };

        tenantData = await readTenant(ctx);
        _tenantData = $state.snapshot(tenantData);
    }

    onMount(() => {
        context.subscribe(loadData);
    });

    async function onSaveChanges(tenantConfig: Tenant) {
        updateTenant(tenantConfig);
        console.log(tenantConfig);
    }

    let hasChanges = $derived(
        JSON.stringify(tenantData) != JSON.stringify(_tenantData),
    );

    let showConfirmDeleteModal = $state(false);
    function askDeleteTenant(){
        showConfirmDeleteModal = true
    }

    async function confirmDeleteTenant() {
        await deleteTenant(tenantData.id)
        goto('/')
    }
</script>

<Modal bind:open={showConfirmDeleteModal}>
    <div class="prose">
        <h2 class="text-rose-600">
            Delete "{tenantData.display_name}""
        </h2>
        <p>This action cannot be undone. Are you sure you want to do this?</p>
        <Button variant="danger" onclick={confirmDeleteTenant}>Delete tenant</Button>
    </div>
</Modal>

<div class="prose">
    {#if tenantData_descriptor.display_name}
        {#key tenantData_descriptor}
            <Form
                onsubmit={onSaveChanges}
                bind:state={tenantData}
                descriptor={tenantData_descriptor}
            >
                {#snippet submit()}
                    <div class="flex gap-4 items-center">
                        {#if hasChanges}
                            <Button
                                type="submit"
                                variant="primary"
                                class="flex gap-2 items-center"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    width="24"
                                    height="24"
                                    viewBox="0 0 24 24"
                                    ><path
                                        fill="currentColor"
                                        d="M21 7v14H3V3h14zm-9 11q1.25 0 2.125-.875T15 15t-.875-2.125T12 12t-2.125.875T9 15t.875 2.125T12 18m-6-8h9V6H6z"
                                    /></svg
                                >
                                <span> Save changes</span>
                            </Button>
                        {/if}
                        <Button
                            onclick={askDeleteTenant}
                            variant="danger"
                            class="flex gap-2 items-center"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="24"
                                height="24"
                                viewBox="0 0 24 24"
                                ><path
                                    fill="currentColor"
                                    d="M5 21V6H4V4h5V3h6v1h5v2h-1v15zm2-2h10V6H7zm2-2h2V8H9zm4 0h2V8h-2zM7 6v13z"
                                /></svg
                            >
                            <span> Delete tennant </span>
                        </Button>
                    </div>
                {/snippet}
            </Form>
        {/key}
    {/if}
</div>
