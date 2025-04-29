<script lang="ts">
    import Crud from "$lib/pages/CRUD.svelte";

    import * as apiControllerUntyped from "$lib/logic/enforcer.svelte";
    const apiController = apiControllerUntyped as unknown as CRUDAPI;

    const OBJECT_TYPE = "Enforcer";
    const TABLE_COLUMNS = ["Name", "Description"];

    import {
        readAll as readAllAdapters,
        create as createAdapter,
    } from "$lib/logic/adapter.svelte";

    import {
        create as createModel,
    } from "$lib/logic/model.svelte";

    import {
        create as createEnforcer,
        readAll as readAllEnforcers,
        read as readEnforcer,
        update as updateEnforcer,
    } from "$lib/logic/enforcer.svelte"

    import apiDescriptor from '../../data/apidescriptor.json'

    const MultiCRUD = <CRUDAPI>{
        async create(descriptor: any): Promise<void> {
            // 1) Create Dependencies
            // 1.1) Create new Adapter
            const newAdapter = await createAdapter({
                database_name: descriptor.database_name,
                driver: descriptor.driver,
                external_db: descriptor.external_db,
                host: descriptor.host,
                display_name: '__internal_adapter',
                password: descriptor.password,
                port: descriptor.port+'',
                table_name: descriptor.table_name || descriptor.database_name,
                username: descriptor.username,
            });

            // 1.2) Create new Model
            const newModel = await createModel({
                display_name:'__internal__model',
                description:'__internal',
                model: apiDescriptor.models.default
            })

            // 2) Create Enforcer with link

            const newEnforcer = await createEnforcer({
                display_name: descriptor.display_name,
                adapter_id: newAdapter.id,
                description: descriptor.description,
                model_id: newModel.id
            })
        },
        update: updateEnforcer,
        readAll: readAllEnforcers,
        read: readEnforcer,
        async kill(id: string): Promise<void> {},
    };

    const build_createForm_descriptor = <FormDescriptor>{
        display_name: String,

        description: String,

        database_name: String,
        driver: {
            type:'select:single',
            label:'Database type',
            options:apiDescriptor.adapters.databases,
            required:true
        },
        external_db: {
            type:'boolean',
            default: false,
            label:'Use external database'
        },
        host: {
            type: 'string',
            required:true,
            label:'Database Host',
            conditional:(state:any)=>state.external_db
        },
        password: {
            type: 'string',
            required:true,
            label:'Database Password',
            conditional:(state:any)=>state.external_db
        },
        port: {
            type:'number',
            max:65535,
            min:1,
            required:true,
            label:'Database port number',
            conditional:(state:any)=>state.external_db
        },
        table_name: {
            type: 'string',
            required:true,
            label:'Table name',
            conditional:(state:any)=>state.external_db
        },
        username: {
            type: 'string',
            required:true,
            label:'Database username',
            conditional:(state:any)=>state.external_db
        },
    };

    const build_editForm_descriptor = <FormDescriptor>{
        id: {
            type: "string",
            label: "Enforcer ID",
            readonly: true,
        },
        adapter_id:{
            type:'string',
            readonly:true,
            label:'Adapter ID (internal)'
        },
        model_id:{
            type:'string',
            readonly:true,
            label:'Model ID (internal)'
        },
        display_name: String,
    };
</script>

<Crud
    apiController={MultiCRUD}
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
