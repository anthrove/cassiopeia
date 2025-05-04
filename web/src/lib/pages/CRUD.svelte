<script lang="ts">
    // Specialization configuration
    /*import {
        readAll as readAllItems,
        create as createItem,
        update as updateItem,
        kill as deleteItem,
    } from "$lib/logic/user.svelte";*/

    // in production these will be basically anything but the idea stands.
    type PartialCreateDescriptor = User__create;
    type PartialUpdateDescriptor = User__update;
    type ItemType = User;

    const {
        apiController = {} as CRUDAPI<
            ItemType,
            PartialCreateDescriptor,
            PartialUpdateDescriptor
        >,
        OBJECT_TYPE = "User",
        TABLE_COLUMNS = ["Name", "Email", "Groups"],

        build_createForm_descriptor,
        build_editForm_descriptor,

        searchKey = "display_name",

        ...otherProps
    } = $props();

    const {
        read: readItem,
        readAll: readAllItems,
        create: createItem,
        update: updateItem,
        kill: deleteItem,
    } = apiController;

    type _ItemType = ItemType | {};

    // Object edit page boilerplate code
    import SearchableTable from "$lib/components/searchableTable.svelte";
    import Modal from "$lib/components/modal.svelte";
    import Form from "$lib/components/forms/Form.svelte";
    import Button from "$lib/components/Button.svelte";

    import { context } from "$lib/stores/context.svelte";

    let items: ItemType[] = $state([]);

    let createForm_descriptor: FormDescriptor = $state({});
    let editForm_descriptor: FormDescriptor = $state({});

    function isFunction(x:any) {
        return typeof x === "function";
    }

    async function reloadData() {
        items = await readAllItems();
        // the form used to populate the PartialDescriptor used in creation
        if(isFunction(build_createForm_descriptor)){
            createForm_descriptor = await build_createForm_descriptor()
        }else{
            createForm_descriptor = build_createForm_descriptor
        }

        if (!Object.keys(current||{}).length) {
            closeEditModal();
        }
        setTimeout(()=>{
            _current = $state.snapshot(current);
        },0)
    }

    context.subscribe(reloadData);
    onMount(()=>{
        reloadData()
        if($selected){
            openEditModal($selected)
        }
    })

    let _current = $state<_ItemType>(<ItemType>{});
    let createModal_open = $state(false);
    function openCreateModal() {
        createModal_open = true;
    }

    async function create(descriptior: PartialCreateDescriptor) {
        await createItem(descriptior);
        reloadData();
        createModal_open = false
    }

    async function update(descriptior: PartialUpdateDescriptor) {
        await updateItem(descriptior);
        reloadData();
    }

    import { create as createContext } from "$lib/stores/urlcontext.svelte";
    import { onMount } from "svelte";
    let selected = createContext("selected");

    let current = $state<ItemType>(<any>null)
    async function openEditModal(id: string) {
        selected.set(id);

        current = await readItem(id);
        _current = $state.snapshot(current || {})

        if(isFunction(build_editForm_descriptor)){
            editForm_descriptor = {}
            editForm_descriptor = await build_editForm_descriptor(id)
        }else{
            editForm_descriptor = build_editForm_descriptor
        }
    }
    function closeEditModal() {
        selected.set("");
    }

    let hasChanges = $derived(
        JSON.stringify(current) != JSON.stringify(_current),
    );

    let confirmDeleteModal_open = $state(false);
    function askDelete() {
        confirmDeleteModal_open = true;
    }

    async function confirmDelete() {
        await deleteItem(current.id);

        confirmDeleteModal_open = false;
        createModal_open = false;

        reloadData();
    }
</script>

<div class="container mx-auto mt-16">
    <Modal bind:open={createModal_open}>
        <div class="prose">
            <h2>Create new {OBJECT_TYPE}</h2>
            <Form
                submit="Create"
                descriptor={createForm_descriptor}
                onsubmit={create}
            />
        </div>
    </Modal>

    <Modal open={!!$selected} onclose={closeEditModal}>
        <div class="prose">
            <h2>Edit {OBJECT_TYPE}</h2>
            {#if Object.keys(editForm_descriptor).length > 0}
                <Form
                    onsubmit={update}
                    bind:state={current}
                    descriptor={editForm_descriptor}
                >
                    {#snippet submit()}
                        <div class="flex gap-4 items-center">
                            <Button
                                disabled={!hasChanges}
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
                                <span> Delete {OBJECT_TYPE}</span>
                            </Button>
                        </div>
                    {/snippet}
                </Form>
            {/if}
        </div>
    </Modal>

    <Modal bind:open={confirmDeleteModal_open}>
        <div class="prose">
            <h2 class="text-rose-600">
                Delete "{current.display_name}"
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

    <SearchableTable rows={items} {searchKey}>
        {#snippet header()}
            {#each TABLE_COLUMNS as label, index}
                <th class="text-start py-4 {index ? '' : 'px-4'}">{label}</th>
            {/each}
        {/snippet}
        {#snippet row(item: ItemType)}
            {#each TABLE_COLUMNS as column, index}
                <td class={index ? "" : "p-4"}>
                    {#if otherProps[column.toLowerCase().replaceAll(' ','_')]}
                        {@render otherProps[column.toLowerCase().replaceAll(' ','_')](item)}
                    {:else}
                        <span>
                            {item[column.toLowerCase().replaceAll(' ','_')]}
                        </span>
                    {/if}
                </td>
            {/each}
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
                <i class="block">No {`${OBJECT_TYPE.toLowerCase()}s`} yet ...</i
                >
                <Button onclick={openCreateModal}
                    >Create your first {OBJECT_TYPE.toLowerCase()}</Button
                >
            </div>
        {/snippet}
    </SearchableTable>

    {#if items?.length}
        <Button onclick={openCreateModal} class="mt-4"
            >Create a {OBJECT_TYPE.toLowerCase()}</Button
        >
    {/if}
</div>
