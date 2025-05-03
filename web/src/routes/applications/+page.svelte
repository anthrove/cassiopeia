<script lang="ts">
    import Crud from "$lib/pages/CRUD.svelte";

    import * as apiControllerUntyped from "$lib/logic/application.svelete";
    import { onMount } from "svelte";
    import { readonly } from "svelte/store";
    const apiController = apiControllerUntyped as unknown as CRUDAPI;

    const OBJECT_TYPE = "Application";
    const TABLE_COLUMNS = ["Name"];

    const build_createForm_descriptor = <FormDescriptor>{
        display_name: String,
        logo: String,
        forget_url: String,
        redirect_urls: {
            type: "string",
            label: "Redirect URLs",
            required: false,
        },
        sign_in_url: String,
        sign_up_url: String,
        terms_url: String,
    };

    const build_editForm_descriptor = {
        id: {
            type: "string",
            label: "Application ID",
            readonly: true,
        },
        ...build_createForm_descriptor,
        client_secret:{
            type: "string",
            label:"Client Secret",
            readonly:true
        }
    };
</script>

<Crud
    {apiController}
    {OBJECT_TYPE}
    {TABLE_COLUMNS}
    {build_createForm_descriptor}
    {build_editForm_descriptor}
    searchKey="display_name"
>
    {#snippet name(item: Application)}
        {item.display_name}
    {/snippet}
</Crud>
