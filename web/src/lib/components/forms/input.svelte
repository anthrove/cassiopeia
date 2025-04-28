<script lang="ts">
    import DateInput from "./DateInput.svelte";

    let {
        value = $bindable(<any>""),
        label = "",
        type: inputType = "text",
        placeholder = "",
        required = false,
        name = "",
        class: classList = "",
        readonly = false,
    } = $props();
</script>

<label class="flex flex-col {classList}">
    {#if label}
        <div class="relative flex items-center mb-1">
            <span class="px-2 z-10 ml-8 dark:bg-zinc-800 bg-white "
                >{label}{#if required}<span class="text-rose-600">*</span
                    >{/if}</span
            >
            <div class="absolute w-full h-[1px] bg-zinc-200 dark:bg-zinc-600"></div>
        </div>
    {/if}
    {#if ["text"].includes(inputType)}
        <input
            {readonly}
            type={inputType}
            {required}
            {name}
            placeholder={placeholder || label}
            bind:value
            class="bg-inherit {readonly
                ? 'cursor-not-allowed text-zinc-400'
                : ''} p-2 border-zinc-200"
        />
    {:else if ["date", "datetime-local"].includes(inputType)}
        <DateInput
            {readonly}
            type={inputType}
            {required}
            {name}
            placeholder={placeholder || label}
            bind:value
            class="{readonly
                ? 'cursor-not-allowed text-zinc-400'
                : ''} p-2 border-zinc-200"
        />
    {:else}
        <div class="text-lg text-rose-600 italic">
            Input type not yet implemented :(
        </div>
    {/if}
</label>
