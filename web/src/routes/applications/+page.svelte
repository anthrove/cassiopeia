<script lang="ts">
    import SearchableTable from "$lib/components/searchableTable.svelte";
    import Modal from "$lib/components/modal.svelte";
    import Form from "$lib/components/forms/Form.svelte";
    import Button from "$lib/components/Button.svelte";

    import { context } from "$lib/stores/context.svelte";

    import { readAll as readAllApplications, create as createApplication } from "$lib/logic/application.svelete";
    import { readAll as readAllCertificates } from "$lib/logic/certificate.svelte";

    let applications: Application[] = $state([]);
    let certificates: Certificate[] = $state([]);
    let createForm_descriptor: FormDescriptor = $state({});
    async function reloadData() {
        applications = await readAllApplications();
        certificates = await readAllCertificates();
        createForm_descriptor = {
            display_name: String,

            certificate_id: {
                type: "select:single",
                labeler: (id: string) =>
                    `${certificates.find((c: Certificate) => c.id == id)?.display_name} (${id})`,
                label: "Signing key",
                required: true,

                options: certificates.map((c) => c.id),
            },
        };
    }

    context.subscribe(reloadData);

    let createModal_open = $state(false);
    function openCreateModal() {
        createModal_open = true;
    }

    type PartialDescriptor = {
        display_name: string
        certificate_id: string
    }

    async function create(descriptior:PartialDescriptor) {
        await createApplication({
            ...descriptior,

            logo:'',

            terms_url:'',
            sign_up_url:'',
            sign_in_url:'',
            forget_url:'',

            redirect_urls:[],
        })

        reloadData()
    }
</script>

<div class="container mx-auto mt-16">
    <Modal bind:open={createModal_open}>
        <div class="prose">
            <h2>Create new Application</h2>
            <Form
                submit="Create"
                descriptor={createForm_descriptor}
                onsubmit={create}
            />
        </div>
    </Modal>

    <SearchableTable rows={applications}>
        {#snippet header()}
            <th class="text-start p-4">Application</th>
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
                <Button
                    href="/tenant/"
                    onclick={() => context.set(item.id)}
                    class="inline-block"
                >
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
                <i class="block">No applications yet ...</i>
                <Button onclick={openCreateModal}
                    >Create your first Application</Button
                >
            </div>
        {/snippet}
    </SearchableTable>

    {#if applications?.length}
        <Button onclick={openCreateModal} class="mt-4"
            >Create an Application</Button
        >
    {/if}
</div>
