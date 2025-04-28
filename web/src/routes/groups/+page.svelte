<script lang="ts">
    import Crud from "$lib/pages/CRUD.svelte";

    import * as apiControllerUntyped from "$lib/logic/group.svelte";
    const apiController = apiControllerUntyped as unknown as CRUDAPI;

    const OBJECT_TYPE = "Group";
    const TABLE_COLUMNS = ["Name", "parent_group_id"];

    const build_createForm_descriptor = <FormDescriptor>{
        display_name: String,
        parent_group_id: {
            type:'string',
            label:'Parent group ID',
            required:false
        },
    };

    const build_editForm_descriptor = {
        id: {
            type: "string",
            label: "Group ID",
            readonly: true,
        },
        display_name: String,
        parent_group_id: {
            type:'string',
            label:'Parent group ID',
            required:false
        },
    };
</script>

<Crud
    {apiController}
    {OBJECT_TYPE}
    {TABLE_COLUMNS}
    {build_createForm_descriptor}
    {build_editForm_descriptor}
    searchKey='displayName'
>
    {#snippet name(item: Group)}
        {item.displayName}
    {/snippet}
    {#snippet parent_group_id(item: Group)}
        {JSON.stringify(item.parent_group_id)}
    {/snippet}
</Crud>
