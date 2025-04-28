<script lang="ts">
    import Crud from "$lib/pages/CRUD.svelte";

    import * as apiControllerUntyped from "$lib/logic/enforcer.svelte";
    const apiController = apiControllerUntyped as unknown as CRUDAPI;

    const OBJECT_TYPE = "Enforcer";
    const TABLE_COLUMNS = ["Name", "Description"];

    const build_createForm_descriptor = {
        display_name: String,

        username: String,
        email: String,

        password: String,
    };

    const build_editForm_descriptor = {
        id: {
            type: "string",
            label: "Enforcer ID",
            readonly: true,
        },
        display_name: String
    };
</script>

<Crud
    {apiController}
    {OBJECT_TYPE}
    {TABLE_COLUMNS}
    {build_createForm_descriptor}
    {build_editForm_descriptor}
>
    {#snippet name(item: User)}
        {item.display_name}
    {/snippet}
    {#snippet groups(item: User)}
        {JSON.stringify(item.groups)}
    {/snippet}
</Crud>
