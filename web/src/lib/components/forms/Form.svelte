<script lang="ts">
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
        submit = 'Submit'
    } = $props();

    // hydrate state descriptors
    let descriptors: ExplicitFormFieldDescriptor[] = [];
    let initialState: { [key: string]: any } = {};
    for (const [key, fieldDescriptor] of Object.entries(descriptor)) {
        let explicitDescriptor: ExplicitFormFieldDescriptor;
        if (fieldDescriptor === String) {
            explicitDescriptor = {
                type: "string",
                required: true,
                label: key,
            };
            initialState[key] = "";
        } else if (fieldDescriptor === Number) {
            explicitDescriptor = {
                type: "number",
                required: true,
                label: key,
            };
            initialState[key] = 0;
        } else if (fieldDescriptor === Boolean) {
            explicitDescriptor = {
                type: "boolean",
                required: true,
                label: key,
            };
        } else if (Array.isArray(fieldDescriptor)) {
            explicitDescriptor = {
                type: "select:single",

                options: fieldDescriptor,
                default: fieldDescriptor[0],

                required: true,
                label: key,
            };
            initialState[key] = explicitDescriptor.default;
        } else {
            explicitDescriptor = <ExplicitFormFieldDescriptor>fieldDescriptor;
            initialState[key] =
                explicitDescriptor.default ||
                {
                    boolean: false,
                    number: 0,
                    string: "",
                    //@ts-expect-error This path is only hit for select:single in which case this value exitsts.
                    "select:single": explicitDescriptor?.options?.at(0),
                }[explicitDescriptor.type];
        }
        descriptors.push(explicitDescriptor);
    }
    formState = initialState;
    const keylist = Object.keys(initialState);
</script>

<form onsubmit={onFormSubmit}>
    {#each descriptors as descriptor, descriptorIndex}
        {#if descriptor.type == "string"}
            <Input
                required={descriptor.required}
                label={descriptor.label}
                type="text"
                bind:value={formState[keylist[descriptorIndex]]}
                class="mb-4"
            />
        {:else if descriptor.type == "select:single"}
            <Dropdown
                required={descriptor.required}
                label={descriptor.label}
                options={descriptor.options}
                bind:value={formState[keylist[descriptorIndex]]}
                class="mb-4"
            />
        {:else}
            <span class="text-rose-600 mb-4"
                >Unsuported input type: {descriptor.type}</span
            >
        {/if}
    {/each}

    <Button>{submit}</Button>
</form>
