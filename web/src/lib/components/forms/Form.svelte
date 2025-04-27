<script lang="ts">
    import type { Snippet } from "svelte";
    import Button from "../Button.svelte";
    import Dropdown from "./Dropdown.svelte";
    import Input from "./input.svelte";

    // This component converts a form description into a rendered form that is two-way bindable
    function onFormSubmit(event: SubmitEvent) {
        event.preventDefault();
        event.stopPropagation();
        console.log("Form submitted");
        dispatchEvent(
            new SubmitEvent("submit", JSON.parse(JSON.stringify(formState))),
        );
        onsubmit?.(JSON.parse(JSON.stringify(formState)));
        return false;
    }

    let {
        state: formState = $bindable({}),
        descriptor = {} as FormDescriptor,
        onsubmit = () => {},
        submit = <Snippet|string>'Submit'
    } = $props();

    let descriptors: ExplicitFormFieldDescriptor[] = $state([]);
    //let formState: { [key: string]: any } = $state({});
    let keylist = $state<string[]>([])

    let doRender = $state(false)
    // hydrate state descriptors
    function hydrate(descriptor:FormDescriptor){
        formState = {...formState}
        keylist = Object.keys(descriptor);
        for (const [key, fieldDescriptor] of Object.entries(descriptor)) {
            let explicitDescriptor: ExplicitFormFieldDescriptor;
            if (fieldDescriptor === String) {
                formState[key] = formState[key] || "";
                explicitDescriptor = {
                    type: "string",
                    required: true,
                    label: key,
                    readonly:false
                };
            } else if (fieldDescriptor === Number) {
                formState[key] = formState[key] || 0;
                explicitDescriptor = {
                    type: "number",
                    required: true,
                    label: key,
                    readonly:false
                };
            } else if (fieldDescriptor === Boolean) {
                formState[key] = !!formState[key]
                explicitDescriptor = {
                    type: "boolean",
                    required: true,
                    label: key,
                    readonly:false
                };
            } else if (Array.isArray(fieldDescriptor)) {
                formState[key] = formState[key] || fieldDescriptor[0];
                explicitDescriptor = {
                    type: "select:single",
    
                    options: fieldDescriptor,
                    default: fieldDescriptor[0],
    
                    required: true,
                    label: key,
                    readonly:false
                };
            } else if (fieldDescriptor === Date) {
                formState[key] = formState[key] || new Date().toISOString()
                explicitDescriptor = {
                    type: "date",
                    required: true,
                    label: key,
                    readonly:false
                };
            }else {
                explicitDescriptor = <ExplicitFormFieldDescriptor>fieldDescriptor;
                formState[key] =
                    formState[key] ||
                    explicitDescriptor.default ||
                    {
                        date: new Date().toISOString(),
                        boolean: false,
                        number: 0,
                        string: "",
                        //@ts-expect-error This path is only hit for select:single in which case this value exitsts.
                        "select:single": explicitDescriptor?.options?.at(0),
                    }[explicitDescriptor.type];
            }
            descriptors.push(explicitDescriptor);
        }
    }

    hydrate(descriptor)
    $effect(()=>{
    })
</script>

<pre>{JSON.stringify(formState,null,4)}</pre>
<form onsubmit={onFormSubmit}>
    {#each descriptors as descriptor, descriptorIndex}
        {#if ['string','date'].includes(descriptor.type)}
            <Input
                readonly={descriptor.readonly}
                required={descriptor.required}
                label={descriptor.label}
                type={{string:'text',date:'datetime-local'}[descriptor.type]}
                bind:value={formState[keylist[descriptorIndex]]}
                class="mb-4"
            />
        {:else if descriptor.type == "select:single"}
            <Dropdown
                readonly={descriptor.readonly}
                required={descriptor.required}
                label={descriptor.label}
                options={descriptor.options}
                bind:value={formState[keylist[descriptorIndex]]}
                labeler={descriptor.labeler}
                class="mb-4"
            />
        {:else}
            <span class="text-rose-600 mb-4"
                >Unsuported input type: {descriptor.type}</span
            >
        {/if}
    {/each}

    {#if submit}
        {#if typeof submit === 'string'}
            <Button type="submit">{submit}</Button>
        {:else}
            {@render submit()}
        {/if}
    {/if}
</form>