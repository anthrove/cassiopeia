<script lang="ts">
    import { onMount } from "svelte";

    import Button from "$lib/components/Button.svelte";
    import Dropdown from "$lib/components/forms/Dropdown.svelte";
    import Form from "$lib/components/forms/Form.svelte";
    import Modal from "$lib/components/modal.svelte";

    import apiDescriptor from "../../data/apidescriptor.json";

    let createModal_open = $state(false);
    function openCreateModal() {
        createModal_open = true;
    }

    const createForm_descriptor = <FormDescriptor>{
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

    const editEnforcer_descriptor = <FormDescriptor>{
        display_name: String,

        description: String,
    };
    const editAdapter_descriptor = <FormDescriptor>{
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

    const editModel_descriptor = <FormDescriptor>{
        model:{
            type:'textarea',
            label:'Model definition',
            required:true
        }
    }

    // reload logic
    let enforcers: Enforcer[] = $state([]);
    async function reloadData() {
        enforcers = await readAllEnforcers();
    }
    
    // context loader logic
    import { context } from "$lib/stores/context.svelte";
    context.subscribe(reloadData);
    onMount(()=>{
        reloadData()
    })

    let enforcerIds = $derived(enforcers.map(e=>e.id))
    function enforcerLabeler(id:string){
        const e = enforcers.find(e=>e.id == id)
        return `${e?.display_name} (${e?.id})`
    }

    // Selection logic
    import { create as createContext } from "$lib/stores/urlcontext.svelte";
    let selected = createContext("selected");

    

    // API logic
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

    async function create(descriptor: any): Promise<void> {
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
    }

    
    // load the linked objects
    let currentEnforcer = $state<Enforcer>({});//$derived(enforcers.find(e=>e.id == $selected))
    let currentModel = $state<Model>({})
    let currentAdaper = $state<Adapter>({})

    let _referenceState = $state({
        currentEnforcer:'null',
        currentModel:'',
        currentAdaper:''
    });

    $effect(()=>{
        (async ()=>{
            const curr = JSON.parse(JSON.stringify(enforcers.find(e=>e.id == $selected)||{}))
            if(!$selected || !curr.model_id){return}
            currentEnforcer = curr
            currentModel = await readModel(curr.model_id)
            currentAdaper = await readAdapter(curr.adapter_id)

            // create cached reference
            _referenceState = {
                currentEnforcer: JSON.stringify(currentEnforcer),
                currentModel: JSON.stringify(currentModel),
                currentAdaper: JSON.stringify(currentAdaper),
            }
        })()
    })

    async function onSaveChanges(){
        if(JSON.stringify(currentEnforcer) != _referenceState.currentEnforcer){
            await updateEnforcer(currentEnforcer)
        }
        if(JSON.stringify(currentModel) != _referenceState.currentModel){
            await updateModel(currentModel)
        }
        if(JSON.stringify(currentAdaper) != _referenceState.currentAdaper){
            await updateAdapter(currentAdaper)
        }
    }
</script>

<Modal bind:open={createModal_open}>
    <div class="prose">
        <h2>Create new Enforcer</h2>
        <Form
            submit="Create"
            descriptor={createForm_descriptor}
            onsubmit={create}
        />
    </div>
</Modal>

<div class="max-w-3xl">
    <div class="flex gap-4 mb-4">
        {#if enforcerIds.length}
        <div class="flex-1">
            <Dropdown class="flex-1 w-full" bind:value={$selected} options={enforcerIds} labeler={enforcerLabeler}/>
        </div>
        {/if}
        <Button class="flex-1 text-center!" onclick={openCreateModal}
            >Create new Enforcer</Button
        >
    </div>
    {#if !$selected}
    <div class="border-[1px] py-8 border-gray-300 flex items-center justify-center">
        <span>No enforcer selected</span>
    </div>
    {:else}
    {#if Object.keys(currentEnforcer||{}).length > 0}
        <Form descriptor={editEnforcer_descriptor} state={currentEnforcer} class="mt-12" submit=''/>
    {/if}
    {#if Object.keys(currentAdaper||{}).length > 0}
        <Form descriptor={editAdapter_descriptor} state={currentAdaper} class="mt-12" submit=''/>
    {/if}
    {#if Object.keys(currentModel||{}).length > 0}
        <Form descriptor={editModel_descriptor} state={currentModel} class="mt-12" submit='Save changes' onsubmit={onSaveChanges}/>
    {/if}
    {/if}
</div>
