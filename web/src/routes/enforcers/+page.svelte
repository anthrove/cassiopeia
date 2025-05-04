<script lang="ts">
    import Crud from "$lib/pages/CRUD.svelte";
    import apiDescriptor from "../../data/apidescriptor.json";

    import * as apiControllerUntyped from "$lib/logic/enforcer.svelte";
    import {
        read as readAdapter,
        create as createAdapter,
        update as updateAdapter,
    } from "$lib/logic/adapter.svelte";

    import {
        create as createModel,
        read as readModel,
        update as updateModel,
    } from "$lib/logic/model.svelte";

    import {
        create as createEnforcer,
        readAll as readAllEnforcers,
        read as readEnforcer,
        update as updateEnforcer,
    } from "$lib/logic/enforcer.svelte";
    const apiController = <CRUDAPI>{
        readAll: readAllEnforcers,

        async read(id: string): Promise<any> {
            const enforcer = await readEnforcer(id);
            const adapter = await readAdapter(enforcer.adapter_id);

            const model = await readModel(enforcer.model_id);

            return {
                ...model,
                modelId: model.id,
                ...adapter,
                adapterId: adapter.id,
                ...enforcer,
            };
        },

        async create(descriptor: any): Promise<void> {
            // 1) Create Dependencies
            // 1.1) Create new Adapter
            const newAdapter = await createAdapter({
                database_name: descriptor.database_name,
                driver: descriptor.driver,
                external_db: descriptor.external_db,
                host: descriptor.host,
                display_name: "__internal_adapter",
                password: descriptor.password,
                port: descriptor.port + "",
                table_name: descriptor.table_name || descriptor.database_name,
                username: descriptor.username,
            });

            // 1.2) Create new Model
            const newModel = await createModel({
                display_name: "__internal__model",
                description: "__internal",
                model: apiDescriptor.models.default,
            });

            // 2) Create Enforcer with link

            const newEnforcer = await createEnforcer({
                display_name: descriptor.display_name,
                adapter_id: newAdapter.id,
                description: descriptor.description,
                model_id: newModel.id,
            });
        },

        async update(descriptor: any) {
            // we have to reconstruct the 3 base objects from the merged form state
            // 1) Model
            const model = <Model__update>{
                id: descriptor.modelId,
                description: "__internal__model",
                model: descriptor.model,
                display_name: "__internal__model",
            };
            updateModel(model)
            // 2) Adapter
            const adapter = <Adapter__update>{
                id: descriptor.adapterId,
                database_name: descriptor.database_name,
                driver: descriptor.driver,
                external_db: descriptor.external_db,
                host: descriptor.host,
                display_name: "__internal",
                password: descriptor.password,
                port: descriptor.port,
                table_name: descriptor.table_name,
                //tenant_id: descriptor.tenant_id,
                username: descriptor.username,
            };
            updateAdapter(adapter)
            // 3) Enforcer
            const enforcer = <Enforcer__update>{
                id: descriptor.id,
                adapter_id: descriptor.adapterId,
                description: descriptor.description,
                model_id: descriptor.modelId,
                display_name: descriptor.display_name,
            };
            updateEnforcer(enforcer)
        },

        async kill(id: string) {
            // cleanup all related objects to avoid dangling references
        },
    };

    const OBJECT_TYPE = "Enforcer";
    const TABLE_COLUMNS = ["Name", "Description"];

    const build_createForm_descriptor = <FormDescriptor>{
        display_name: String,

        description: String,

        table_name: String,
        external_db: {
            type: "boolean",
            default: false,
            label: "Use external database",
        },
        driver: {
            type: "select:single",
            label: "Database type",
            options: apiDescriptor.adapters.databases,
            required: true,
            conditional: (state: any) => state.external_db,
        },
        host: {
            type: "string",
            required: true,
            label: "Database Host",
            conditional: (state: any) => state.external_db,
        },
        password: {
            type: "string",
            required: true,
            label: "Database Password",
            conditional: (state: any) => state.external_db,
        },
        port: {
            type: "number",
            max: 65535,
            min: 1,
            required: true,
            label: "Database port number",
            conditional: (state: any) => state.external_db,
        },
        database_name: {
            type: "string",
            required: true,
            label: "Table name",
            conditional: (state: any) => state.external_db,
        },
        username: {
            type: "string",
            required: true,
            label: "Database username",
            conditional: (state: any) => state.external_db,
        },
    };

    const build_editForm_descriptor = {
        id: {
            type: "string",
            label: "Enforcer ID",
            readonly: true,
        },
        display_name: String,

        description: String,

        table_name: String,
        external_db: {
            type: "boolean",
            default: false,
            label: "Use external database",
        },
        driver: {
            type: "select:single",
            label: "Database type",
            options: apiDescriptor.adapters.databases,
            required: true,
            conditional: (state: any) => state.external_db,
        },
        host: {
            type: "string",
            required: true,
            label: "Database Host",
            conditional: (state: any) => state.external_db,
        },
        password: {
            type: "string",
            required: true,
            label: "Database Password",
            conditional: (state: any) => state.external_db,
        },
        port: {
            type: "number",
            max: 65535,
            min: 1,
            required: true,
            label: "Database port number",
            conditional: (state: any) => state.external_db,
        },
        database_name: {
            type: "string",
            required: true,
            label: "Table name",
            conditional: (state: any) => state.external_db,
        },
        username: {
            type: "string",
            required: true,
            label: "Database username",
            conditional: (state: any) => state.external_db,
        },

        model: {
            type: "textarea",
            label: "Model definition",
            required: true,
        },
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
