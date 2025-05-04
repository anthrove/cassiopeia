<script lang="ts">
    import Modal from "$lib/components/modal.svelte";
    import { onMount } from "svelte";

    import {
        readCategories,
        readTypes,
        readConfig,
        readAll as readAllProviders,
        create as createProvider,
        update as updateProvider,
        kill as deleteProvider,
    } from "$lib/logic/provider.svelte";
    import Form from "$lib/components/forms/Form.svelte";
    import Button from "$lib/components/Button.svelte";
    import SearchableTable from "$lib/components/searchableTable.svelte";
    let providerCategories = $state<string[]>([]);

    onMount(async () => {
        providerCategories = await readCategories();
    });

    let providers = $state<Provider[]>([]);

    async function reloadData() {
        providers = await readAllProviders();
    }

    import { create as createContext } from "$lib/stores/urlcontext.svelte";
    let selected = createContext("selected");
    import { context } from "$lib/stores/context.svelte";
    context.subscribe(reloadData);
    onMount(reloadData);

    let createModal_state = $state({
        category: "",
        type: "",
    });
    let createModal_open = $state(false);
    function openCreateModal() {
        createModal_state = {
            category: "",
            type: "",
        };
        createModal_open = true;
    }
    let providerTypes = $state<string[]>([]);
    async function selectCategory(category: string) {
        if (createModal_state.category == category) {
            createModal_state.category = "";
        } else {
            providerTypes = await readTypes(category);
            createModal_state.category = category;
        }
        createModal_state.type = "";
    }
    let providerConfig = $state({});
    async function selectType(type: string) {
        if (createModal_state.type == type) {
            createModal_state.type = "";
        } else {
            providerConfig = Object.fromEntries([
                [
                    "display_name",
                    {
                        type: "string",
                        required: true,
                        label: "Display name",
                    },
                ],
                ...(await readConfig(createModal_state.category, type)).map(
                    (field) => [
                        field.field_key,
                        <FormFieldDescriptor>{
                            type:
                                {
                                    text: "string",
                                    bool: "boolean",
                                    int: "number",
                                    secret: "password",
                                }[field.field_type] || field.field_type,
                            required: true,
                            label: field.field_key,
                        },
                    ],
                ),
            ]);
            createModal_state.type = type;
        }
    }
    async function onCreateSubmit(configuration: any) {
        const provider = await createProvider({
            display_name: configuration.display_name,
            provider_type: createModal_state.type,
            category: createModal_state.category,
            parameter: configuration,
        });
        createModal_open = false;
    }

    async function onEditSubmit(configuration: any) {
        let parameter = <any>{};
        for (const key of Object.keys(providerConfig)) {
            parameter[key] = configuration[key];
        }
        const provider = await updateProvider({
            id: $selected,
            display_name: configuration.display_name,
            parameter,
        });
        editModal_open = false;
        reloadData();
    }

    let editModal_open = $state(false);
    async function openEditModal(id: string) {
        const provider = providers.find((p) => p.id == id);
        if (!provider) {
            return;
        }
        selected.set(provider.id);
        createModal_state.category = provider.category;
        await selectType(provider.provider_type);
        editModal_open = true;
    }

    function closeEditModal() {
        selected.set("");
        editModal_open = false;
    }

    function flatProvider(id: string) {
        const provider = providers.find((p) => p.id == id);
        if (!provider) {
            return null;
        }
        return {
            ...provider.parameter,
            ...provider,
        };
    }

    const TABLE_COLUMNS = ["Display Name", "Category"];


    let confirmDeleteModal_open = $state(false);
    function askDelete() {
        confirmDeleteModal_open = true;
    }

    async function confirmDelete() {
        await deleteProvider($selected);
        reloadData();

        confirmDeleteModal_open = false
        editModal_open = false
    }
</script>

<Modal bind:open={createModal_open}>
    <div class="prose">
        <h2>Create new Provider</h2>
        <h3>Step 1: Choose a category</h3>
        <div class="flex gap-4">
            {#each providerCategories as category}
                <button
                    class="px-4 py-2 cursor-pointer border-[1px] border-primary-500 {createModal_state.category ==
                    category
                        ? 'text-white bg-primary-500'
                        : 'text-primary-500'}"
                    onclick={() => selectCategory(category)}>{category}</button
                >
            {/each}
        </div>
        {#if createModal_state.category}
            <h3>Step 2: Choose the type</h3>
            <div class="flex gap-4">
                {#each providerTypes as type}
                    <button
                        class="px-4 py-2 cursor-pointer border-[1px] border-primary-500 {createModal_state.type ==
                        type
                            ? 'text-white bg-primary-500'
                            : 'text-primary-500'}"
                        onclick={() => selectType(type)}>{type}</button
                    >
                {/each}
            </div>
        {/if}
        {#if createModal_state.type}
            <h3>Step 3: Configure Provider</h3>
            {#key createModal_state.type}
                <Form
                    descriptor={providerConfig}
                    onsubmit={onCreateSubmit}
                    submit="Create provider"
                ></Form>
            {/key}
        {/if}
    </div>
</Modal>

<Modal bind:open={editModal_open} onclose={closeEditModal}>
    <div class="prose">
        <h2>Edit Provider</h2>
        {#key createModal_state.type}
            <Form
                descriptor={providerConfig}
                state={flatProvider($selected)}
                onsubmit={onEditSubmit}
            >
                {#snippet submit()}
                    <div class="flex gap-4 items-center">
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
                        <Button
                            onclick={askDelete}
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
                            <span> Delete Provider</span>
                        </Button>
                    </div>
                {/snippet}
            </Form>
        {/key}
    </div>
</Modal>

<Modal bind:open={confirmDeleteModal_open}>
    <div class="prose">
        <h2 class="text-rose-600">
            Delete Provider?
        </h2>
        <p>
            This action cannot be undone. Are you sure you want to do this?
        </p>
        <div class="flex items-center gap-4 flex-wrap">
            <Button variant="danger" onclick={confirmDelete}
                >Confirm permanent deletion</Button
            >
            <Button
                variant="secondary"
                onclick={() => {
                    confirmDeleteModal_open = false;
                }}>Cancel</Button
            >
        </div>
    </div>
</Modal>

<div class="container mx-auto mt-16">
    <SearchableTable rows={providers} columns={TABLE_COLUMNS}>
        {#snippet header()}
            {#each TABLE_COLUMNS as label, index}
                <th class="text-start py-4 {index ? '' : 'px-4'}">{label}</th>
            {/each}
        {/snippet}

        {#snippet row(item: Provider)}
            <td class="p-4">
                {item.display_name}
            </td>
            <td>
                {item.category}-{item.provider_type}
            </td>
            <td class="py-1">
                <Button
                    onclick={() => openEditModal(item.id)}
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
                <i class="block">No providers yet ...</i>
                <Button onclick={openCreateModal}
                    >Create your first provider</Button
                >
            </div>
        {/snippet}
    </SearchableTable>

    {#if providers.length}
        <Button onclick={openCreateModal} class="mt-4">Create Provider</Button>
    {/if}
</div>
