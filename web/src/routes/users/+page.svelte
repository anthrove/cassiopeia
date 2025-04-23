<script lang="ts">
    // Specialization configuration
    import {
        readAll as readAllItems,
        create as createItem,
    } from "$lib/logic/user.svelte";

    const OBJECT_TYPE = "User"
    const tableColumns = ['Name','Email','Groups']

    type PartialDescriptor = {
        display_name: string;
        username: string;
        email: string,
        password: string
    };

    // Object edit page boilerplate code
    import SearchableTable from "$lib/components/searchableTable.svelte";
    import Modal from "$lib/components/modal.svelte";
    import Form from "$lib/components/forms/Form.svelte";
    import Button from "$lib/components/Button.svelte";

    import { context } from "$lib/stores/context.svelte";

    let items: User[] = $state([]);

    let createForm_descriptor: FormDescriptor = $state({});
    let editForm_descriptor:FormDescriptor = $state({})

    async function reloadData() {
        items = await readAllItems();
        // the form used to populate the PartialDescriptor used in creation
        createForm_descriptor = {
            display_name: String,

            username: String,
            email: String,

            password: String,
        };

        editForm_descriptor = {
            id:{
                type:'string',
                label:'User ID',
                readonly:true
            },
            display_name: String,

            username: String,
            email: String,

            password: String,
        }
    }

    context.subscribe(reloadData);

    let createModal_open = $state(false);
    function openCreateModal() {
        createModal_open = true;
    }

    async function create(descriptior: PartialDescriptor) {
        await createItem({
            ...descriptior,
            // optional overwrites
        });

        reloadData();
    }

    import { create as createContext } from "$lib/stores/urlcontext.svelte";
    let selected = createContext('selected')

    function openEditModal(id:string){
        selected.set(id)
    }
    function closeEditModal(){
        selected.set('')
    }
    let current = $derived(items.find(u=>u.id==$selected)||{})
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
            <pre>{JSON.stringify(current,null,4)}</pre>
            <Form
                submit="Update"
                bind:state={current}
                descriptor={editForm_descriptor}
            />
        </div>
    </Modal>

    <SearchableTable rows={items}>
        {#snippet header()}
            {#each tableColumns as label,index}
                <th class="text-start py-4 {index?'':'px-4'}">{label}</th>
            {/each}
        {/snippet}
        {#snippet row(item: User)}
            <td class="p-4">
                {item.display_name}
            </td>
            <td class="flex gap-2 my-4 items-center">
                {#if item.email_verified}
                    <svg
                        class="text-primary-400"
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        ><path
                            fill="currentColor"
                            d="m8.6 22.5l-1.9-3.2l-3.6-.8l.35-3.7L1 12l2.45-2.8l-.35-3.7l3.6-.8l1.9-3.2L12 2.95l3.4-1.45l1.9 3.2l3.6.8l-.35 3.7L23 12l-2.45 2.8l.35 3.7l-3.6.8l-1.9 3.2l-3.4-1.45zm2.35-6.95L16.6 9.9l-1.4-1.45l-4.25 4.25l-2.15-2.1L7.4 12z"
                        /></svg
                    >
                {/if}
                {item.email}
            </td>
            <td>
                {item.groups}
            </td>
            <td>
                <Button
                    onclick={()=>openEditModal(item.id)}
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
                <i class="block">No users yet ...</i>
                <Button onclick={openCreateModal}>Create your first user</Button
                >
            </div>
        {/snippet}
    </SearchableTable>

    {#if items?.length}
        <Button onclick={openCreateModal} class="mt-4">Create a user</Button>
    {/if}
</div>
