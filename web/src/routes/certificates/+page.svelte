<script lang="ts">
    import Crud from "$lib/pages/CRUD.svelte";
    import apiDescriptor from "../../data/apidescriptor.json";

    import * as apiControllerUntyped from "$lib/logic/certificate.svelte";
    const apiController = apiControllerUntyped as unknown as CRUDAPI;

    const OBJECT_TYPE = "Certificate";
    const TABLE_COLUMNS = ["Name", "Bit size", "Created at", "Expired at"];

    const build_createForm_descriptor = <FormDescriptor>{
        display_name:String,
        algorithm: {
            type: "select:single",
            label: "Signing key",
            required: true,

            options: apiDescriptor.certificate.algorithms,
        },
        bit_size:{
            type: "select:single",
            label:"Bit-Size",
            required:true,

            options: apiDescriptor.certificate.bit_sizes
        },
        expired_at:Date
        //'algorithm'|'bit_size'|'display_name'|'expired_at'
    };

    const build_editForm_descriptor = {
        id: {
            type: "string",
            label: "Certificate ID",
            readonly: true,
        },
        display_name: String,
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
